<?php
    // automatically load all classes
    spl_autoload_register('booking_autoloader');
    function booking_autoloader($class){
        include '../classes/'.$class.'.php';
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
        case 'add_session':
            $response = add_session($configuration, $db);
            break;
        case 'edit_session':
            $response = edit_session($configuration, $db);
            break;
        case 'delete_session':
            $response = delete_session($configuration, $db);
            break;
        case 'delete_user':
            $response = delete_rider_from_session($configuration, $db);
            break;
        case 'add_users':
            $response = add_rider_to_session($configuration, $db);
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

    function add_session($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // sanitize input
        $error = array();
        $sanitizer = new Sanitizer();
        if(!isset($data->start) or !$sanitizer->isInt($data->start)){
            error_log('api/booking.php: Illegal start provided: ' . $data->start);
            return Status::errorStatus('No valid start provided');
        }
        if(!isset($data->end) or !$sanitizer->isInt($data->end)){
            error_log('api/booking.php: Illegal end provided: ' . $data->end);
            return Status::errorStatus('No valid end provided');
        }
        if(!isset($data->max_riders) or !$sanitizer->isInt($data->max_riders)){
            error_log('api/booking.php: Illegal max_riders provided: ' . $data->max_riders);
            return Status::errorStatus('No valid number of maximum riders provided');
        }
        if(!isset($data->type) or !$sanitizer->isInt($data->type)){
            error_log('api/booking.php: Illegal session type provided: ' . $data->type);
            return Status::errorStatus('No valid session type provided');
        }
        if($data->start >= $data->end){
            error_log('api/booking.php: Session start cannot be after end: ' . $data->start . ' to ' . $data->end);
            return Status::errorStatus('Session start is after the end and that does not make sense.');
        }

        // check validity of data
        // check type of session
        $query = sprintf('SELECT id FROM session_type WHERE id = %d', $data->type);
        if(!$db->exists($query)){
            error_log('api/booking.php: Session type does not exist ' . $data->type);
            return Status::errorStatus('No valid session type provided');
        }
        // check if the new session collides with an existing one
        $query = 'SELECT id
            FROM session
            WHERE
                ( start_time <= FROM_UNIXTIME(?)
                    AND end_time >= FROM_UNIXTIME(?)
                )
                OR
                ( start_time <= FROM_UNIXTIME(?)
                    AND end_time >= FROM_UNIXTIME(?)
                )
                OR
                ( start_time <= FROM_UNIXTIME(?)
                    AND end_time >= FROM_UNIXTIME(?)
                )
                OR
                ( start_time >= FROM_UNIXTIME(?)
                    AND end_time <= FROM_UNIXTIME(?)
                )
        ;';
        $db->prepare($query);
        $db->bind_param('ssssssss',
            $data->start,
            $data->end,
            $data->end,
            $data->end,
            $data->start,
            $data->start,
            $data->start,
            $data->end
        );
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(count($res) > 0){
            error_log('api/booking.php: The session is overlapping with an existing one');
            error_log('api/booking.php: Oberlapping with: ' . $res[0]['id']);
            return Status::errorStatus('Session overlaps with an existing one');
        }

        // Everything is fine and we will add the session
        $session = Login::getSessionData($configuration, $db);
        $query = 'INSERT INTO session
            (title, start_time, end_time, comment, free, creator_id, type)
            VALUES
            (?, FROM_UNIXTIME(?), FROM_UNIXTIME(?),?,?,?,?);';
        $db->prepare($query);
        $db->bind_param('siisiii',
                        $data->title,
                        $data->start,
                        $data->end,
                        $data->comment,
                        $data->max_riders,
                        $session['user_id'],
                        $data->type);
        if(!$db->execute()){
            error_log('Cannot add new session with: ' . $query);
            return Status::errorStatus('Cannot add new session. An error appeared.');
        }

        $query = 'SELECT LAST_INSERT_ID() as session_id;';
        $db->prepare($query);
        if(!$db->execute()){
            error_log('Cannot get last inserted session id: ' . $query);
            return Status::errorStatus('Cannot add new session. An error appeared.');
        }
        $res = $db->fetch_stmt_hash();

        return Status::successDataResponse("Session has been created", $res[0]);
    }

    function edit_session($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // sanitize input
        $error = array();
        $sanitizer = new Sanitizer();
        if(!isset($data->id) or !$sanitizer->isInt($data->id)){
            error_log('api/booking.php: Illegal id provided: ' . $data->id);
            return Status::errorStatus('No valid id provided');
        }
        if(!isset($data->start) or !$sanitizer->isInt($data->start)){
            error_log('api/booking.php: Illegal start provided: ' . $data->start);
            return Status::errorStatus('No valid start provided');
        }
        if(!isset($data->end) or !$sanitizer->isInt($data->end)){
            error_log('api/booking.php: Illegal end provided: ' . $data->end);
            return Status::errorStatus('No valid end provided');
        }
        if(!isset($data->max_riders) or !$sanitizer->isInt($data->max_riders)){
            error_log('api/booking.php: Illegal max_riders provided: ' . $data->max_riders);
            return Status::errorStatus('No valid number of maximum riders provided');
        }
        if(!isset($data->type) or !$sanitizer->isInt($data->type)){
            error_log('api/booking.php: Illegal session type provided: ' . $data->type);
            return Status::errorStatus('No valid session type provided');
        }
        if($data->start >= $data->end){
            error_log('api/booking.php: Session start cannot be after end: ' . $data->start . ' to ' . $data->end);
            return Status::errorStatus('Session start is after the end and that does not make sense.');
        }

        // check validity of data
        // check type of session
        $query = sprintf('SELECT id FROM session_type WHERE id = %d', $data->type);
        if(!$db->exists($query)){
            error_log('api/booking.php: Session type does not exist ' . $data->type);
            return Status::errorStatus('No valid session type provided');
        }

        // check that a session with this ID exists indeed
        $query = 'SELECT id
            FROM session
            WHERE id = ?';
        $db->prepare($query);
        $db->bind_param('i', $data->id);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(count($res) != 1){
            error_log('api/booking.php: A session with this ID does not exist.');
            return Status::errorStatus('This session was not found and can thus not be changed.');
        }

        // check if the new session collides with an existing one
        $query = 'SELECT id
            FROM session
            WHERE
            (
                ( start_time <= FROM_UNIXTIME(?)
                    AND end_time >= FROM_UNIXTIME(?)
                )
                OR
                ( start_time <= FROM_UNIXTIME(?)
                    AND end_time >= FROM_UNIXTIME(?)
                )
                OR
                ( start_time <= FROM_UNIXTIME(?)
                    AND end_time >= FROM_UNIXTIME(?)
                )
                OR
                ( start_time >= FROM_UNIXTIME(?)
                    AND end_time <= FROM_UNIXTIME(?)
                )
            ) AND id <> ?
        ;';
        $db->prepare($query);
        $db->bind_param('ssssssssi',
            $data->start,
            $data->end,
            $data->end,
            $data->end,
            $data->start,
            $data->start,
            $data->start,
            $data->end,
            $data->id
        );
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(count($res) > 0){
            error_log('api/booking.php: The session is overlapping with an existing one.');
            error_log('api/booking.php: Oberlapping with: ' . $res[0]['id']);
            return Status::errorStatus('Session would overlap with an existing one. Please change start and/or end time.');
        }

        // Everything is fine and we will update the session
        $session = Login::getSessionData($configuration, $db);
        $query = 'UPDATE session SET
            title = ?,
            start_time = FROM_UNIXTIME(?),
            end_time = FROM_UNIXTIME(?),
            comment = ?,
            free = ?,
            type = ?
            WHERE id = ?;';
        $db->prepare($query);
        $db->bind_param('siisiii',
            $data->title,
            $data->start,
            $data->end,
            $data->comment,
            $data->max_riders,
            $data->type,
            $data->id);
        if(!$db->execute()){
            error_log('Cannot edit session with: ' . $query);
            return Status::errorStatus('Cannot edit session. An error occured.');
        }

        return Status::successStatus("Session has been edited");
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
            id as id,
            UNIX_TIMESTAMP(start_time) as start,
            UNIX_TIMESTAMP(end_time) as end,
            title as title,
            comment as description,
            type as type,
            free as free,
            creator_id as creator_id
            FROM session
            WHERE
            ( UNIX_TIMESTAMP(start_time) >= ? AND UNIX_TIMESTAMP(start_time) < ? )
            OR
            ( UNIX_TIMESTAMP(end_time) >= ? AND UNIX_TIMESTAMP(end_time) < ? )
            OR
            ( UNIX_TIMESTAMP(start_time) < ? AND UNIX_TIMESTAMP(end_time) >= ? )
            ORDER BY start_time;';
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

    // Deletes a session and all it's riders, it also removes all the invitations
    // All registered users are notified about the cancellation
    function delete_session($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // sanitize input
        $sanitizer = new Sanitizer();
        if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
            error_log('api/booking.php: Illegal session id provided: ' . $data->session_id);
            return Status::errorStatus('No valid session id provided');
        }

        // only an admin is allowed to do so
        $session_data = Login::getSessionData($configuration, $db);
        if($session_data['user_role_id'] != $configuration->admin_user_role_id){
            error_log('api/booking.php: Attempt to delete a session as non-admin ' . $session_data[username]);
            return Status::errorStatus('Only an administrator can delete sessions.');
        }

        // check if there exist heats to this session
        $query = 'SELECT id FROM heat WHERE session_id = ?';
        $db->prepare($query);
        $db->bind_param('i', $data->session_id);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(!isset($res)){
            error_log('api/booking.php: Cannot check for heats for this session.');
            HttpHeader::setResponseCode(500);
        }
        if(count($res) > 0){
            error_log('api/booking.php: Attempt to delete session with heats.');
            return Status::errorStatus('This session already has heats. Only sessions without heats can be deleted.');
        }

        // inform invitees
        _inform_invitees($data->session_id, $session_data, $db, $configuration);

        // delete the invitations
        $query = 'DELETE FROM invitation
                  WHERE session_id = ?';
        $db->prepare($query);
        $db->bind_param('i', $data->session_id);
        if(!$db->execute()){
            error_log('api/booking.php: Cannot delete invitations for session: ' . $data->session_id);
        }

        // inform riders of this session
        _inform_registered($data->session_id, $session_data, $db, $configuration);

        // delete registered users
        $query = 'DELETE FROM user_to_session WHERE session_id = ?';
        $db->prepare($query);
        $db->bind_param('i', $data->session_id);
        if(!$db->execute()){
            error_log('api/booking.php: Cannot delete users for session: ' . $data->session_id);
        }

        // Delete the session itself
        $query = 'DELETE FROM session WHERE id = ?';
        $db->prepare($query);
        $db->bind_param('i', $data->session_id);
        if(!$db->execute()){
            error_log('api/booking.php: Cannot delete session: ' . $data->session_id);
            return Status::errorStatus('Session cannot be removed due to unknown reasons.');
        }

        return Status::successStatus("session deleted");
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

    function _inform_registered($session_id, $session_data, $db, $configuration){
        $query = 'SELECT u.first_name AS fn, u.last_name AS ln,
                         u.email AS email,
                         s.date AS date, s.start AS start, s.end AS end,
                         s.title AS title
                  FROM user_to_session u2s, session s, user u
                  WHERE u2s.session_id = ?
                     AND u2s.session_id = s.id
                     AND u2s.user_id    = u.id
                     AND u2s.user_id    <> ?';
        $db->prepare($query);
        $db->bind_param('ii', $session_id, $session_data['user_id']);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(!isset($res) or $res===FALSE){
            error_log('_inform_registered: Cannot get invitations from database');
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

    // Deletes a rider from a session
    // - Removing a general user from a session can only be done by the admin user
    // - Removing the own user can be done by the user itself if it is done
    //   early enough
    function delete_rider_from_session($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // Get user information
        $session_data = Login::getSessionData($configuration, $db);

        // sanitize input
        $sanitizer = new Sanitizer();
        if(isset($data->user_id) and !$sanitizer->isInt($data->user_id)){
            error_log('api/booking.php: Illegal user_id provided: ' . $data->user_id);
            return Status::errorStatus("No valid user_id provided");
        }
        if(!isset($data->user_id)){
            $data->user_id = $session_data['user_id'];
        }
        if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
            error_log('api/booking.php: Illegal session id provided: ' . $data->session_id);
            return Status::errorStatus("No valid session id provided");
        }

        // check if there is such a session
        $query = "SELECT s.date, s.start, u2s.user_id
                  FROM user_to_session u2s, session s
                  WHERE u2s.session_id = ?
                    AND u2s.session_id = s.id
                    AND u2s.user_id = ?;";
        $db->prepare($query);
        $db->bind_param('ii', $data->session_id, $data->user_id);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(!isset($res) or count($res) != 1){
            error_log('api/booking.php: No session found for user/session: ' . $data->user_id . '/' . $data->session_id);
            return Status::errorStatus("User and session do not exist in this combination");
        }

        // check if the user is still allowed to cancel the session
        $time = strtotime($res[0]['date'] .' ' . $res[0]['start']);
        if($session_data['user_role_id'] != $configuration->admin_user_role_id
           && $session_data['user_id'] != $data->user_id){
            error_log('api/booking.php: Attempt to delete another user from a session (logged_in_user/session): '
                      . $session_data['user_id'] . ' / ' . $data->session_id);
            return Status::errorStatus("You can only cancel your own session");

        }
        if(time() > $time - $configuration->session_cancel_graceperiod
           && $session_data['user_role_id'] != $configuration->admin_user_role_id){
            error_log('api/booking.php: Attempt to delete a user from a session too late (user/session) ' . $data->user_id . '/' . $data->session_id );
            return Status::errorStatus("The session cannot be cancelled anymore at the time being.");
        }

        // check if the user already had a heat in that session
        $query = "SELECT count(*) as nbr_of_heats
            FROM heat h 
            WHERE user_id = ? 
                AND session_id = ?;";
        $db->prepare($query);
        $db->bind_param("ii",
            $data->user_id,
            $data->session_id
        );
        if(!$db->execute()){
            error_log('api/booking.php: Attempt to check whether user already had a heat did not work.');
            return Status::errorStatus("Cannot check whether user already had sessions");
        }
        $res = $db->fetch_stmt_hash();
        if($res[0]["nbr_of_heats"] > 0){
            return Status::errorStatus("This user cannot be removed because there are heats from this user in the selected session.");
        }        

        // Delete the user from the session
        $query = "DELETE FROM user_to_session "
                ."WHERE session_id = ?"
                . " AND user_id = ? LIMIT 1;";
        $db->prepare($query);
        $db->bind_param('ii', $data->session_id, $data->user_id);
        if(!$db->execute()){
            error_log("api/booking.php: Cannot delete user from session");
            return Status::errorStatus("The session cannot be cancelled due to a technical issue. Please try again and let us know about this error.");
        }

        // Also delete the invitations regarding this session
        $query = "DELETE FROM invitation "
                ."WHERE session_id = ?"
                . " AND user_id = ? LIMIT 1;";
        $db->prepare($query);
        $db->bind_param('ii', $data->session_id, $data->user_id);
        if(!$db->execute()){
            error_log("api/booking.php: Cannot delete user and session from invitation table (user/session): " .
                      $data->user_id . "/" . $data->session_id);
        }

        // Add a free space to the session
        $query = "UPDATE session SET free=(free+1) WHERE id = ?;";
        $db->prepare($query);
        $db->bind_param('i', $data->session_id);
        if(!$db->execute()){
            error_log("api/booking.php: Cannot add free space to session: " . $data->session_id);
        }

        return Status::successStatus("user has been removed");
    }

    // Adds all given riders to the session
    function add_rider_to_session($configuration, $db){
        $data = json_decode(file_get_contents('php://input'));

        // Get user information
        $session_data = Login::getSessionData($configuration, $db);

        // sanitize input
        $sanitizer = new Sanitizer();
        if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
            error_log('api/booking.php: Illegal session_id provided: ' . $data->session_id);
            return Status::errorStatus("No valid session ID provided");
        }
        if(!isset($data->user_ids) or count($data->user_ids) < 1){
            error_log('api/booking.php: Illegal user_ids provided.');
            return Status::errorStatus('No user ID provided');
        }
        for($i=0; $i<count($data->user_ids); $i++){
            if(! $sanitizer->isInt($data->user_ids[$i])){
                error_log('api/booking.php: Illegal user_id provided: ' . $data->user_ids[$i]);
                return Status::errorStatus('Invalid user ID provided');
            }
        }

        // all is valid, add all users
        for($i=0; $i<count($data->user_ids); $i++){
            if($session_data['user_role_id'] == $configuration->admin_user_role_id){
                _add_rider_to_session($data->user_ids[$i], $data->session_id, $db);
            }else{
                error_log('api/booking.php: User is not admin and cannot add person to session: user '
                          . $data->user_ids[$i]
                          . ' (logged in user: ' . $session_data['user_id'] . ')');
            }

            // Notify user by email
            $query = 'SELECT u.first_name AS fn, u.last_name AS ln,
                             u.email AS email,
                             s.date AS date, s.start AS start, s.end AS end,
                             s.title AS title
                      FROM session s, user u
                      WHERE s.id = ?
                        AND u.id = ?
                        AND u.id <> ?;';
            $db->prepare($query);
            $db->bind_param('iii',
                            $data->session_id,
                            $data->user_ids[$i],
                            $session_data['user_id']);
            $db->execute();
            $res = $db->fetch_stmt_hash();
            if(!isset($res) or $res===FALSE){
                error_log('api/booking.php: Cannot get notification info (session/user): '
                          . $data->session_id . '/'
                          . $data->user_ids[$i]);
            }elseif(count($res) > 1){
                $message = 'Dear ' . $res[0]['fn'] . ' ' . $res[0]['ln'] . "\n";
                $message .= "\nYou have been added to the following session:\n";
                $message .= " Title: " . $res[0]['title'] . "\n";

                $start = strtotime($res[0]['date'] . " " . $res[0]['start']);
                $end = strtotime($res[0]['date'] . " " . $res[0]['end']);
                $message .= " Date : " . date("D d.m.Y", $start) . "\n";
                $message .= "        " . date("H:i", $start) . " to " . date("H:i", $end) . "\n";
                $message .= "\nPlease login to www.wakeandsurf.ch and check your Schedule to decline.\n";
                $message .= "\nSee you on the lake soon\n";
                Email::sendMail($res[0]['email'], 'Session Confirmation', $message, $configuration);
            }
        }

        return Status::successStatus("users have been added to session");
    }

    // Function to add a rider to a session
    // This function does not make any sanitization
    function _add_rider_to_session($user_id, $session_id, $db){

        // check if user is already part of the session
        $query = 'SELECT count(*) as count FROM user_to_session
                  WHERE user_id = ?
                    AND session_id = ?;';
        $db->prepare($query);
        $db->bind_param('ii', $user_id, $session_id);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(!isset($res) or $res === FALSE){
            error_log('api/booking.php: Cannot check if user part of session (session/user): '
                      . $session_id . '/' . $user_id);
            return Status::errorStatus('Cannot check if user is part of session already.');
        }
        if($res[0]['count'] > 0){
            error_log('api/booking.php: User already part of session (session/user): '
                      . $session_id . '/' . $user_id);
            return Status::errorStatus('User already part of session.');
        }

        // remove invitations since now the user is added to the session
        $query = 'DELETE FROM invitation
                  WHERE user_id = ?
                    AND session_id = ?;';
        $db->prepare($query);
        $db->bind_param('ii', $user_id, $session_id);
        if(!$db->execute()){
            error_log('api/booking.php: Cannot delete invitation for user (session/user): '
                      . $session_id . '/' . $user_id);
            return Status::errorStatus('Cannot delete invitation for user.');
        }

        // check if the session has still space for another user
        $query = 'SELECT free FROM session
                  WHERE id = ?;';
        $db->prepare($query);
        $db->bind_param('i', $session_id);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        if(!isset($res) or $res===FALSE){
            error_log('api/booking.php: Cannot check free space for session (session/user): '
                      . $session_id . '/' . $user_id);
            return Status::errorStatus('Cannot check for free space in session.');
        }
        if($res[0]['free'] < 1){
            error_log('api/booking.php: No free space in session (session/user): '
                      . $session_id . '/' . $user_id);
            return Status::errorStatus('No free space in session.');
        }

        // Insert user to session
        $query = 'INSERT INTO user_to_session
                  (user_id, session_id)
                  VALUES (?, ?);';
        $db->prepare($query);
        $db->bind_param('ii', $user_id, $session_id);
        if(!$db->execute()){
            error_log('api/booking.php: Cannot add user to session (session/user): '
                      . $session_id . '/' . $user_id);
            return Status::errorStatus('Cannot add user to session.');
        }

        // Decrease free
        $query = 'UPDATE session SET free = free - 1
                  WHERE id = ?;';
        $db->prepare($query);
        $db->bind_param('i', $session_id);
        if(!$db->execute()){
            error_log('api/booking.php: Cannot decrement free of session (session): '
                      . $session_id);
            return Status::errorStatus('Cannot decrement free space.');
        }
        return Status::errorStatus('User added to session');
    }




?>
