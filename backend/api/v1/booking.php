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

?>
