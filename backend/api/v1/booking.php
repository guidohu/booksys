<?php
    // automatically load all classes
    spl_autoload_register('booking_autoloader');
    function booking_autoloader($class){
        include '../../classes/'.$class.'.php';
    }

    // get configuration access
    $configuration = new Configuration();

    // create database connection
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log('Cannot connect to database');
        echo 'Cannot connect to database';
        return;
    }

    // only users which are logged in are allowed to see booking information
    $lc = new Login($configuration);
    if(!$lc->isLoggedIn()){
        $response = array();
        $response["redirect"] = $configuration->login_page;
        echo json_encode($response);
		$db->disconnect();
        exit;
    }

    // check if we have an action
    if(!isset($_GET['action'])){
        HttpHeader::setResponseCode(200);
		$db->disconnect();
        exit;
    }

	$response = null;

    switch($_GET['action']){
        case 'get_booking_day':
            $response = get_booking_day($configuration, $db);
            break;
        case 'get_booking_month':
            $response = get_booking_month($configuration, $db);
            break;
        case 'get_session':
            $response = get_session($configuration, $db);
            break;
        default:
			HttpHeader::setResponseCode(400);
            $response = Status::errorStatus("Action not supported");
            break;
    }

    $db->disconnect();
	echo json_encode($response);
    return;

    function get_booking_day($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // input validation
        $error = array();
        $sanitizer = new Sanitizer();
        if(! isset($data->start)){
            return Status::errorStatus('No date start parameter given');
        }
        if(! isset($data->end)){
            return Status::errorStatus('No date end parameter given');
        }
        if(!$sanitizer->isInt($data->start)){
            return Status::errorStatus('No valid date start parameter given');
        }
        if(!$sanitizer->isInt($data->end)){
            return Status::errorStatus('No valid date end parameter given');
        }

        $res = get_booking($data->start, $data->end, $configuration, $db);
        return $res;
    }

    function get_booking_month($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // input validation
        $sanitizer = new Sanitizer();
        $year = 0;
        $month = 0;
        if(!isset($data->timeWindows)){
            return Status::errorStatus("No time windows provided");
        }
        
        // validate all time windows
        for($i = 0; $i<count($data->timeWindows); $i++){
            $start = $data->timeWindows[$i]->start;
            $end   = $data->timeWindows[$i]->end;

            if(!isset($start) or !$sanitizer->isInt($start)){
                return Status::errorStatus('Invalid start time provided for timeWindow ' . $i);
            }
            if(!isset($end) or !$sanitizer->isInt($end)){
                return Status::errorStatus('Invalid end time provided for timeWindow ' . $i);
            }
        }

        // get the booking for each timeWindow
        for($i=0; $i<count($data->timeWindows); $i++){
            $start = $data->timeWindows[$i]->start;
            $end   = $data->timeWindows[$i]->end;
            $bookings = get_booking($start, $end, $configuration, $db);
            if($bookings['ok'] == TRUE){
                $res[$i] = $bookings['data'];
            }else{
                return Status::errorStatus($bookins['msg']);
            }
        }
        return Status::successDataResponse("sessions retrieved", $res);
    }

    function get_session($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // sanitize input
        $sanitizer = new Sanitizer();
        if(!isset($data->id) or !$sanitizer->isInt($data->id)){
            error_log('api/booking.php: Illegal session ID provided: ' . $data->id);
            return Status::errorStatus("No valid session ID provided");
        }

        // build the query to get session details
        $query = 'SELECT s.id as id, 
                s.date as date,
                s.start as start,
                UNIX_TIMESTAMP(s.start_time) as start_time,
                s.end as end,
                UNIX_TIMESTAMP(s.end_time) as end_time,
                s.title as title,
                s.comment as comment,
                s.type as type,
                st.name as type_name,
                s.free as free,
                s.creator_id as creator_id,
                u.first_name as creator_first_name,
                u.last_name as creator_last_name
            FROM session s 
            JOIN session_type st ON s.type = st.id
            JOIN user u ON u.id = s.creator_id
            WHERE s.id = ?';
        $db->prepare($query);
        $db->bind_param('d',
            $data->id
        );
        $db->execute();
        $res = $db->fetch_stmt_hash();

        // In case there is no session, return already now
        if(!isset($res) or  $res === FALSE or count($res) != 1){
            error_log('api/booking.php: No session found for id: ' . $data->id);
            return Status::errorStatus("No session found with this the provided ID");
        }

        $res = $res[0];

        // get sunset and sunrise
        $sun_info = _get_sunrise_and_sunset($res["start_time"], $res["end_time"], $configuration);
        $res['sunrise'] = $sun_info["sunrise"];
        $res['sunset']  = $sun_info["sunset"];
        
        // build the query to get its riders
        $query = 'SELECT 
                u.id as id,
                u.first_name as first_name,
                u.last_name as last_name
            FROM user_to_session u2s
            JOIN user u ON u2s.user_id = u.id
            WHERE u2s.session_id = ?';
        $db->prepare($query);
        $db->bind_param('d',
            $data->id
        );
        $db->execute();
        $res["riders"] = $db->fetch_stmt_hash();
        $res["riders_max"] = count($res["riders"]) + $res["free"];

        return Status::successDataResponse("success", $res);
    }

    function _get_sunrise_and_sunset($start, $end, $configuration){
        // get my offset to UTC (Timezone configured in configuration)
        $dateTimeZoneLocation = new DateTimeZone($configuration->location_time_zone);
        $dateTimeZoneUTC      = new DateTimeZone("UTC");
        $dateTimeLocation     = new DateTime("@".$end, $dateTimeZoneLocation);
        $dateTimeLocation->setTimestamp($end);
        $dateTimeUTC          = new DateTime("@".$end, $dateTimeZoneUTC);
        $dateTimeLocation->setTimestamp($end);
        $timezoneOffset       = $dateTimeZoneLocation->getOffset($dateTimeUTC);

        // calculate sunrise/sunset for that day
        // Location is configurable
        $res = array();
        $res['sunrise']       = date_sunrise($start+$timezoneOffset, SUNFUNCS_RET_TIMESTAMP, $configuration->location_latitude, $configuration->location_longitude);
        $res['sunset']        = date_sunset($start+$timezoneOffset, SUNFUNCS_RET_TIMESTAMP, $configuration->location_latitude, $configuration->location_longitude);
        return $res;    
    }

    function get_booking($start, $end, $configuration, $db){
        // get sunrise and sunset for given position
        // position can be configured in the configuration script
        $res['window_start'] = intval($start);
        $res['window_end']   = intval($end);

        // get my offset to UTC (Timezone configured in configuration)
        $dateTimeZoneLocation = new DateTimeZone($configuration->location_time_zone);
        $dateTimeZoneUTC      = new DateTimeZone("UTC");
        $dateTimeLocation     = new DateTime("@".$end, $dateTimeZoneLocation);
        $dateTimeLocation->setTimestamp($end);
        $dateTimeUTC          = new DateTime("@".$end, $dateTimeZoneUTC);
        $dateTimeLocation->setTimestamp($end);
        $timezoneOffset       = $dateTimeZoneLocation->getOffset($dateTimeUTC);

        $res['timezone']      = $configuration->location_time_zone;

        // calculate sunrise/sunset for that day
        // Location is configurable
        $res['sunrise']       = date_sunrise($start+$timezoneOffset, SUNFUNCS_RET_TIMESTAMP, $configuration->location_latitude, $configuration->location_longitude);
        $res['sunset']        = date_sunset($start+$timezoneOffset, SUNFUNCS_RET_TIMESTAMP, $configuration->location_latitude, $configuration->location_longitude);
        
        $res['error']               = "";
        $res['business_day_start']  = $configuration->business_day_start;
        $res['business_day_end']    = $configuration->business_day_end;
        $res['sessions']   = Array();

        // get all sessions that happen between start and end
        $query = 'SELECT 
            s.id as id,
            UNIX_TIMESTAMP(s.start_time) as start,
            UNIX_TIMESTAMP(s.end_time) as end,
            s.title as title,
            s.comment as description,
            s.type as type,
            s.free as free,
            s.creator_id as creator_id,
			u.first_name as first_name,
			u.last_name as last_name
            FROM session s 
			LEFT JOIN user u ON u.id = s.creator_id
            WHERE
            ( UNIX_TIMESTAMP(s.start_time) >= ? AND UNIX_TIMESTAMP(s.start_time) < ? )
            OR
            ( UNIX_TIMESTAMP(s.end_time) >= ? AND UNIX_TIMESTAMP(s.end_time) < ? )
            OR
            ( UNIX_TIMESTAMP(s.start_time) < ? AND UNIX_TIMESTAMP(s.end_time) >= ? )
            ORDER BY s.start_time;';
        $db->prepare($query);
        $db->bind_param('iiiiii', 
            $start,
            $end,
            $start,
            $end,
            $start,
            $end
        );
        $db->execute();
        $result = $db->fetch_stmt_hash();
        if(!isset($result) or  $result === FALSE){
            error_log('Cannot get bookings for time window: ' . $start . ' - ' . $end);
            return Status::errorStatus("Cannot get bookings for time window: " . $start . " - " . $end);
        }

        // go through all the sessions in the database
        //-------------------------------------
        for($i=0; $i < count($result); $i++){
            $res['sessions'][$i] = Array (
                'id'         => $result[$i]['id'],
                'start'      => $result[$i]['start'],
                'end'        => $result[$i]['end'],
                'title'      => $result[$i]['title'],
                'comment'    => $result[$i]['description'],
                'free'       => $result[$i]['free'],
                'type'       => $result[$i]['type'],
                'creator_id' => $result[$i]['creator_id'],
				'creator_first_name' => $result[$i]['first_name'],
				'creator_last_name' => $result[$i]['last_name'],
                'duration'   => $result[$i]['end'] - $result[$i]['start']);
        }

        // Get all users of the session
        for($i=0; $i<count($res['sessions']); $i++){
            $id = $res['sessions'][$i]['id'];
            $res['sessions'][$i]['riders'] = Array();
            
            // get all riders of this session
            $query = 'SELECT u.id uid, u.first_name AS fn, u.last_name AS ln
                        FROM user_to_session u2s, user u
                        WHERE
                        u2s.user_id = u.id
                        AND u2s.session_id = ? ORDER BY u.first_name;';
            $db->prepare($query);
            $db->bind_param('i', $id);
            $db->execute();
            $riders = $db->fetch_stmt_hash();
            if(!isset($result) or  $result === FALSE){
                error_log('Cannot get riders for: ' . $id . ' - ' . $query);
                return Status::errorStatus("Cannot get riders for the session");

            }
            for($j=0; $j<count($riders); $j++){
                $res['sessions'][$i]['riders'][$j]['id']   = $riders[$j]['uid'];
                $res['sessions'][$i]['riders'][$j]['name'] = $riders[$j]['fn'] . " " . $riders[$j]['ln'];
            }
        }

        return Status::successDataResponse("sessions retrieved", $res);

    }

    function _inform_invitees($session_id, $session_data, $db, $configuration){
        $query = 'SELECT u.first_name AS fn, u.last_name AS ln,
                         u.email AS email,
                         s.date AS date, s.start AS start, s.end AS end,
                         s.title AS title
                  FROM invitation inv, session s, user u
                  WHERE inv.session_id = ?
                     AND inv.session_id = s.id
                     AND inv.user_id    = u.id
                     AND inv.user_id    <> ?
                     AND inv.status = 2';
        $db->prepare($query);
        $db->bind_param('ii', $session_id, $session_data['user_id']);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(!isset($res) or $res===FALSE){
            error_log('_inform_invitees: Cannot get invitations from database');
            return;
        }

        for($i=0; $i<count($res); $i++){
            $message = 'Dear ' . $res[$i]['fn'] . ' ' . $res[$i]['ln'] . "\n";
            $message .= "\nThe following session has been cancelled:\n";
            $message .= " Title: " . $res[$i]['title'] . "\n";

            $start = strtotime($res[$i]['date'] . " " . $res[$i]['start']);
            $end = strtotime($res[$i]['date'] . " " . $res[$i]['end']);
            $message .= " Date : " . date("D d.m.Y", $start) . "\n";
            $message .= "        " . date("H:i", $start) . " to " . date("H:i", $end) . "\n";
            $message .= "\nSee you soon on the lake\n";
            Email::sendMail($res[$i]['email'], 'Session Cancelled', $message, $configuration);
        }
    }

?>
