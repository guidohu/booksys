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
		exit;
	}

	// check if we have an action
	if(!isset($_GET['action'])){
		HttpHeader::setResponseCode(200);
		exit;
	}

	switch($_GET['action']){
		case 'get_booking_day':
			get_booking_day($configuration, $db);
			$db->disconnect();
			exit;
		case 'get_booking_month':
			get_booking_month($configuration, $db);
			$db->disconnect();
			exit;
		case 'get_session':
			get_session($configuration, $db);
			$db->disconnect();
			exit;
		case 'add_session':
			add_session($configuration, $db);
			$db->disconnect();
			exit;
		case 'delete_session':
			delete_session($configuration, $db);
			$db->disconnect();
			exit;
		case 'delete_user':
			delete_rider_from_session($configuration, $db);
			$db->disconnect();
			exit;
		case 'add_users':
			add_rider_to_session($configuration, $db);
			$db->disconnect();
			exit;
		case 'add_self':
			add_myself_to_session($configuration, $db);
			$db->disconnect();
			exit;
	}

	HttpHeader::setResponseCode(400);
	echo "Action not supported: " . $_GET['action'];
	return;

	function get_booking_day($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// input validation
		$error = array();
		$sanitizer = new Sanitizer();
		if(! isset($data->start)){
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No date start parameter given';
			echo json_encode($error);
			return;
		}
		if(! isset($data->end)){
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No date end parameter given';
			echo json_encode($error);
			return;
		}
		if(!$sanitizer->isInt($data->start)){
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No valid date start parameter given';
			echo json_encode($error);
			return;
		}
		if(!$sanitizer->isInt($data->end)){
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No valid date end parameter given';
			echo json_encode($error);
			return;
		}

		$res = get_booking($data->start, $data->end, $configuration, $db);

		echo json_encode($res);
		return;
	}

	function get_booking_month($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// input validation
		$sanitizer = new Sanitizer();
		$year = 0;
		$month = 0;
		if(!isset($data->timeWindows)){
			HttpHeader::setResponseCode(400);
			echo 'No time windows provided';
			return;
		}
		
		// validate all time windows
		for($i = 0; $i<count($data->timeWindows); $i++){
			$start = $data->timeWindows[$i]->start;
			$end   = $data->timeWindows[$i]->end;

			if(!isset($start) or !$sanitizer->isInt($start)){
				HttpHeader::setResponseCode(400);
				echo 'Invalid start time provided for timeWindow ' . $i;
				return;
			}
			if(!isset($end) or !$sanitizer->isInt($end)){
				HttpHeader::setResponseCode(400);
				echo 'Invalid end time provided for timeWindow ' . $i;
				return;
			}
		}

		// get the booking for each timeWindow
		for($i=0; $i<count($data->timeWindows); $i++){
			$start = $data->timeWindows[$i]->start;
			$end   = $data->timeWindows[$i]->end;
			$res[$i] = get_booking($start, $end, $configuration, $db);
		}

		echo json_encode($res);
		return;
	}

	function get_session($configuration, $db){
		$id = $_GET['id'];

		// sanitize input
		$sanitizer = new Sanitizer();
		if(!isset($id) or !$sanitizer->isInt($id)){
			error_log('api/booking.php: Illegal session ID provided: ' . $id);
			HttpHeader::setResponseCode(400);
			echo 'No valid session ID provided';
			return;
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
			$id
		);
		$db->execute();
		$res = $db->fetch_stmt_hash();

		// In case there is no session, return already now
		if(!isset($res) or  $res === FALSE or count($res) != 1){
			error_log('api/booking.php: No session found for id: ' . $id);
			echo json_encode(new stdClass);
			return;
		}

		$res = $res[0];
		
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
			$id
		);
		$db->execute();
		$res["riders"] = $db->fetch_stmt_hash();

		echo json_encode($res);
		return;
	}

	function add_session($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// sanitize input
		$error = array();
		$sanitizer = new Sanitizer();
		if(!isset($data->start) or !$sanitizer->isInt($data->start)){
			error_log('api/booking.php: Illegal start provided: ' . $data->start);
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No valid start provided';
			echo json_encode($error);
			return;
		}
		if(!isset($data->end) or !$sanitizer->isInt($data->end)){
			error_log('api/booking.php: Illegal end provided: ' . $data->end);
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No valid end provided';
			echo json_encode($error);
			return;
		}
		if(!isset($data->max_riders) or !$sanitizer->isInt($data->max_riders)){
			error_log('api/booking.php: Illegal max_riders provided: ' . $data->max_riders);
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No valid number of maximum riders provided';
			echo json_encode($error);
			return;
		}
		if(!isset($data->type) or !$sanitizer->isInt($data->type)){
			error_log('api/booking.php: Illegal session type provided: ' . $data->type);
			HttpHeader::setResponseCode(400);
			$error["error"] = 'No valid session type provided';
			echo json_encode($error);
			return;
		}
		if($data->start >= $data->end){
			error_log('api/booking.php: Session start cannot be after end: ' . $data->start . ' to ' . $data->end);
			HttpHeader::setResponseCode(400);
			$error["error"] = 'Session start is after the end and that does not make sense.';
			echo json_encode($error);
			return;
		}

		// check validity of data
		// check type of session
		$query = sprintf('SELECT id FROM session_type WHERE id = %d', $data->type);
		if(!$db->exists($query)){
			HttpHeader::setResponseCode(400);
			error_log('api/booking.php: Session type does not exist ' . $data->type);
			$error["error"] = 'No valid session type provided';
			echo json_encode($error);
			return;
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
			HttpHeader::setResponseCode(400);
			error_log('api/booking.php: The session is overlapping with an existing one');
			error_log('api/booking.php: Oberlapping with: ' . $res[0]['id']);
			$error["error"] = 'Session overlaps with an existing one';
			echo json_encode($error);
			return;
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
			HttpHeader::setResponseCode(500);
			error_log('Cannot add new session with: ' . $query);
			$error["error"] = 'Cannot add new session. An error appeared.';
			echo json_encode($error);
			return;
		}

		$query = 'SELECT LAST_INSERT_ID() as session_id;';
		$db->prepare($query);
		if(!$db->execute()){
			HttpHeader::setResponseCode(500);
			error_log('Cannot get last inserted session id: ' . $query);
			$error["error"] = 'Cannot add new session. An error appeared.';
			echo json_encode($error);
			return;
		}
		$res = $db->fetch_stmt_hash();

		HttpHeader::setResponseCode(200);
		echo json_encode($res[0]);
		return;
	}

	function get_booking($start, $end, $configuration, $db){
		// get sunrise and sunset for given position
		// position can be configured in the configuration script
		$res['window_start'] = $start;
		$res['window_end']   = $end;

		// get my offset to UTC (Timezone configured in configuration)
		$dateTimeZoneLocation = new DateTimeZone($configuration->location_time_zone);
		$dateTimeZoneUTC      = new DateTimeZone("UTC");
		$dateTimeLocation     = new DateTime("now", $dateTimeZoneLocation);
		$dateTimeLocation->setTimestamp($end);
		$dateTimeUTC          = new DateTime("now", $dateTimeZoneUTC);
		$dateTimeLocation->setTimestamp($end);
		$timezoneOffset       = $dateTimeZoneLocation->getOffset($dateTimeUTC);

		// calculate sunrise/sunset for that day
		// Location is configurable
		$res['sunrise']       = date_sunrise($end+$timezoneOffset, SUNFUNCS_RET_TIMESTAMP, $configuration->location_latitude, $configuration->location_longitude);
		$res['sunset']        = date_sunset($end+$timezoneOffset, SUNFUNCS_RET_TIMESTAMP, $configuration->location_latitude, $configuration->location_longitude);
		
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
			HttpHeader::setResponseCode(500);
			$res['error'] = "Cannot get bookings for time window: " . $start . " - " . $end;
			echo json_encode($res);
			return;
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
				HttpHeader::setResponseCode(500);
				return;
			}
			for($j=0; $j<count($riders); $j++){
				$res['sessions'][$i]['riders'][$j]['id']   = $riders[$j]['uid'];
				$res['sessions'][$i]['riders'][$j]['name'] = $riders[$j]['fn'] . " " . $riders[$j]['ln'];
			}
		}

		return $res;

	}

	// Deletes a session and all it's riders, it also removes all the invitations
	// All registered users are notified about the cancellation
	function delete_session($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// sanitize input
		$sanitizer = new Sanitizer();
		if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
			error_log('api/booking.php: Illegal session id provided: ' . $data->session_id);
			HttpHeader::setResponseCode(400);
			echo 'No valid session id provided';
			return;
		}

		// only an admin is allowed to do so
		$session_data = Login::getSessionData($configuration, $db);
		if($session_data['user_role_id'] != $configuration->admin_user_status_id){
			error_log('api/booking.php: Attempt to delete a session as non-admin ' . $session_data[username]);
			HttpHeader::setResponseCode(400);
			echo "Only an administrator can delete sessions.";
			return;
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
			return;
		}
		if(count($res) > 0){
			error_log('api/booking.php: Attempt to delete session with heats.');
			HttpHeader::setResponseCode(400);
			echo("A session with finished heats cannot be deleted.");
			return;
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
			HttpHeader::setResponseCode(500);
			echo "Cannot remove session.";
			return;
		}

		return;
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
			HttpHeader::setResponseCode(400);
			echo 'No valid user_id provided';
			return;
		}
		if(!isset($data->user_id)){
			$data->user_id = $session_data['user_id'];
		}
		if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
			error_log('api/booking.php: Illegal session id provided: ' . $data->session_id);
			HttpHeader::setResponseCode(400);
			echo 'No valid session id provided';
			return;
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
			HttpHeader::setResponseCode(400);
			echo "User and session do not exist in this combination";
			return;
		}

		// check if the user is still allowed to cancel the session
		$time = strtotime($res[0]['date'] .' ' . $res[0]['start']);
		if($session_data['user_role_id'] != $configuration->admin_user_status_id
		   && $session_data['user_id'] != $data->user_id){
		    error_log('api/booking.php: Attempt to delete another user from a session (logged_in_user/session): '
			          . $session_data['user_id'] . ' / ' . $data->session_id);
			HttpHeader::setResponseCode(403);
			echo "You can only cancel your own session";
			return;
		}
		if(time() > $time - $configuration->session_cancel_graceperiod
		   && $session_data['user_role_id'] != $configuration->admin_user_status_id){
			error_log('api/booking.php: Attempt to delete a user from a session too late (user/session) ' . $data->user_id . '/' . $data->session_id );
			HttpHeader::setResponseCode(400);
			echo "The session cannot be cancelled anymore at the time being.";
			return;
		}

		// Delete the user from the session
		$query = "DELETE FROM user_to_session "
		        ."WHERE session_id = ?"
			    . " AND user_id = ? LIMIT 1;";
		$db->prepare($query);
		$db->bind_param('ii', $data->session_id, $data->user_id);
		if(!$db->execute()){
			error_log("api/booking.php: Cannot delete user from session");
			HttpHeader::setResponseCode(500);
			echo "The session cannot be cancelled due to a technical issue. Please try again and let us know about this error.";
			return;
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

		return;
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
			HttpHeader::setResponseCode(400);
			echo 'No valid session_id provided';
			return;
		}
		if(!isset($data->user_ids) or count($data->user_ids) < 1){
			error_log('api/booking.php: Illegal user_ids provided.');
			HttpHeader::setResponseCode(400);
			echo 'No valid user_ids provided';
			return;
		}
		for($i=0; $i<count($data->user_ids); $i++){
			if(! $sanitizer->isInt($data->user_ids[$i])){
				error_log('api/booking.php: Illegal user_id provided: ' . $data->user_ids[$i]);
				HttpHeader::setResponseCode(400);
				echo 'No valid user_id provided';
				return;
			}
		}

		// all is valid, add all users
		for($i=0; $i<count($data->user_ids); $i++){
			if($session_data['user_role_id'] == $configuration->admin_user_status_id){
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
	}

	// Adds the current user to the session
	function add_myself_to_session($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// Get user information
		$session_data = Login::getSessionData($configuration, $db);

		// sanitize input
		$sanitizer = new Sanitizer();
		if(isset($data->session_id) and !$sanitizer->isInt($data->session_id)){
			error_log('api/booking.php: Illegal session_id provided: ' . $data->session_id);
			HttpHeader::setResponseCode(400);
			echo 'No valid user_id provided';
			return;
		}

		_add_rider_to_session($session_data['user_id'], $data->session_id, $db);
	}

	// Function to add a rider to a session
	// This function does not make any sanitation
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
			return;
		}
		if($res[0]['count'] > 0){
			error_log('api/booking.php: User already part of session (session/user): '
			          . $session_id . '/' . $user_id);
			return;
		}

		// remove invitations since now the user is added to the session
		$query = 'DELETE FROM invitation
		          WHERE user_id = ?
				    AND session_id = ?;';
		$db->prepare($query);
		$db->bind_param('ii', $user_id, $session_id);
		if(!$db->execute()){
			error_log('api/booking.php: Cannot delete inviation for user (session/user): '
			          . $session_id . '/' . $user_id);
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
			return;
		}
		if($res[0]['free'] < 1){
			error_log('api/booking.php: No free space in session (session/user): '
			          . $session_id . '/' . $user_id);
			return;
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
			return;
		}

		// Decrease free
		$query = 'UPDATE session SET free = free - 1
		          WHERE id = ?;';
		$db->prepare($query);
		$db->bind_param('i', $session_id);
		if(!$db->execute()){
			error_log('api/booking.php: Cannot decrement free of session (session): '
			          . $session_id);
		}
		return;
	}




?>
