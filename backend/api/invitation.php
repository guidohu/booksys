<?php
  include_once '../config/config.php';
	
  // automatically load all classes
  spl_autoload_register('invitation_autoloader');
  function invitation_autoloader($class){
	include '../classes/'.$class.'.php';
  }

  // only users which are logged in are allowed to handle invitations
  $lc = new Login($config);
  if(!$lc->isLoggedIn()){
	header("Location:".$configuration->$configuration->login_page);
	HttpHeader::setResponseCode(302);
	exit;
  }

  // check if we have an action
  if(!isset($_GET['action'])){
	HttpHeader::setResponseCode(200);
	exit;
  }
  
  switch($_GET['action']){
	case 'decline':
		decline_invitation($config);
		exit;
	case 'accept':
		accept_invitation($config);
		exit;
	case 'invite':
		invite_user($config);
		exit;
  }
  
  HttpHeader::setResponseCode(400);
  return;
  
  function invite_user($config){
	$data = json_decode(file_get_contents('php://input'));
	
	// input validation
	$sanitizer = new Sanitizer();
	if( !isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
		HttpHeader::setResponseCode(400);
		echo 'No valid session ID provided.';
		return;
	}
	if( !isset($data->user_ids) or count($data->user_ids)<1){
		HttpHeader::setResponseCode(400);
		echo 'No users selected';
		return;
	}else{
		for($i=0; $i<count($data->user_ids); $i++){
			if(! $sanitizer->isInt($data->user_ids[$i])){
				HttpHeader::setResponseCode(400);
				echo 'No valid user_id given';
				return;
			}
		}
	}
	
	// create database connection
	$db = new DBAccess($config);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
		error_log('api/invitation: Cannot connect to database');
		echo 'Cannot connect to database';
		return;
	}
	
	// check if the session exists and it has some free space
	$query = "SELECT free FROM session
	          WHERE id = ?;";
	$db->prepare($query);
	$db->bind_param('i', $data->session_id);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(! isset($res) or $res === FALSE){
		error_log('api/invitation: Cannot get free slots for session: ' . $data->session_id);
		HttpHeader::setResponseCode(400);
		echo "Check if you selected an existing session.";
		$db->disconnect();
		return;
	}
	if(count($res) == 0){
		error_log('api/invitation: Session does not exist: ' . $data->session_id);
		HttpHeader::setResponseCode(400);
		echo "Session was not found";
		$db->disconnect();
		return;
	}
	if($res[0]['free'] < 1){
		HttpHeader::setResponseCode(400);
		echo "Session does not have free slots";
		$db->disconnect();
		return;
	}
	
	// get session_data
	$session_data = Login::getSessionData($config, $db);
	
	// add user_ids to invitation
	for($i = 0; $i<count($data->user_ids); $i++){
		echo "Handle user: " . $data->user_ids[$i];
		
		// check if this user already got an invitation for this session
		$query = 'SELECT count(*) as count FROM
					( SELECT i.id FROM invitation i WHERE i.session_id = ?
					     AND i.user_id = ?
					  UNION ALL
					  SELECT u2s.id FROM user_to_session u2s WHERE u2s.session_id = ?
					     AND u2s.user_id = ?
					) existing';
		$db->prepare($query);
		$db->bind_param('iiii',
		                $data->session_id,
						$data->user_ids[$i],
						$data->session_id,
						$data->user_ids[$i]);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(!isset($res) or $res===FALSE){
			error_log('api/invitation: Cannot get if user was already invited (session_id, user_id): ' 
			          . $data->session_id .' ,' . $data->user_ids[$i]);
		}
		if($res[0]['count'] > 0){
			// user already got an invitation
			continue;
		}
		
		// add an invitation
		$query = 'INSERT INTO invitation
		          (session_id, user_id, inviter_id, status) 
				  VALUES(?,?,?,5);';
		$db->prepare($query);
		$db->bind_param('iii',
                        $data->session_id,
                        $data->user_ids[$i],
						$session_data['user_id']);
		if(!$db->execute()){
			error_log('api/invitation: Cannot add invitation to database');
		}
	}
	
	// send an email notification to all the users that have invitation status 5
	send_invitation_email($data->session_id, $session_data['user_id'], $db, $config);
	
	$db->disconnect();
	HttpHeader::setResponseCode(200);
	return;	
  }
  
  // this function only removes invitations (not already booked sessions -> booking.php)
  function decline_invitation($config){
	$data = json_decode(file_get_contents('php://input'));
	
	// input validation
	$sanitizer = new Sanitizer();
	if( !isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
		HttpHeader::setResponseCode(400);
		echo 'No valid session ID provided.';
		return;
	}
	
	// create database connection
	$db = new DBAccess($config);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
		error_log('api/invitation: Cannot connect to database');
		echo 'Cannot connect to database';
		return;
	}
	
	// get session_data
	$session_data = Login::getSessionData($config, $db);
	
	// decline the invitation
	$query = 'UPDATE invitation SET status = 4 
	          WHERE session_id = ?
                AND user_id = ?';
	$db->prepare($query);
	$db->bind_param('ii', $data->session_id, $session_data['user_id']);
	if(!$db->execute()){
		HttpHeader::setResponseCode(500);
		error_log('api/invitation: Cannot decline invitation (session / user_id): ' 
		          . $data->session_id . '/' . $session_data['user_id']);
		echo 'Error while declining the invitation.';
		$db->disconnect();
		return;
	}
	
	$db->disconnect();
	return;	
  }
  
  // accept an invitation
  function accept_invitation($config){
	$data = json_decode(file_get_contents('php://input'));
	
	// input validation
	$sanitizer = new Sanitizer();
	if( !isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
		HttpHeader::setResponseCode(400);
		echo 'No valid session ID provided.';
		return;
	}
	
	// create database connection
	$db = new DBAccess($config);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
		error_log('api/invitation: Cannot connect to database');
		echo 'Cannot connect to database';
		return;
	}
	
	// get session_data
	$session_data = Login::getSessionData($config, $db);
	
	// check if the user already accepted the session or was
	// added by an admin to the session
	$query = 'SELECT user_id FROM user_to_session
	          WHERE session_id = ? AND user_id = ?';
	$db->prepare($query);
	$db->bind_param('ii', $data->session_id, $session_data['user_id']);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(!isset($res) or $res === FALSE){
		HttpHeader::setResponseCode(500);
		error_log('api/invitation: Cannot check if user already part of session (session / user_id): '
		          . $data->session_id . '/' . $session_data['user_id']);
		echo 'Error while accepting the invitation.';
		$db->disconnect();
		return;
	}
	if(count($res)>0){
		HttpHeader::setResponseCode(400);
		error_log('api/invitation: User tried to accept invitation a second time (session / user_id): '
		         . $data->session_id . '/' . $session_data['user_id']);
		echo 'You are already part of the session.';
		$db->disconnect();
		return;
	}
	
	// accept the invitation
	// update the own invitation status
	$query = 'UPDATE invitation SET status = 3 
	          WHERE session_id = ?
                AND user_id = ?';
	$db->prepare($query);
	$db->bind_param('ii', $data->session_id, $session_data['user_id']);
	if(!$db->execute()){
		HttpHeader::setResponseCode(500);
		error_log('api/invitation: Cannot accept invitation (session / user_id): ' 
		          . $data->session_id . '/' . $session_data['user_id']);
		echo 'Error while accepting the invitation.';
		$db->disconnect();
		return;
	}
	
	// add the user to the session
	$query = 'INSERT INTO user_to_session (user_id, session_id) 
	          VALUES (?, ?);';
	$db->prepare($query);
	$db->bind_param('ii', $session_data['user_id'], $data->session_id);
	if(!$db->execute()){
		HttpHeader::setResponseCode(500);
		error_log('api/invitation: Cannot add user to session (session / user_id): ' 
		          . $data->session_id . '/' . $session_data['user_id']);
		echo 'Error while accepting the invitation.';
		$db->disconnect();
		return;
	}
	
	// Decrement the amount of free space in the session
	$query = 'UPDATE session SET free = free - 1
              WHERE id = ?;';
	$db->prepare($query);
	$db->bind_param('i', $data->session_id);
	if(!$db->execute()){
		error_log('api/invitation: Cannot lower "free" of session: ' 
		          . $data->session_id);
	}
	
	// inform if session is full by now
	$query = 'SELECT free FROM session WHERE id = ?';
	$db->prepare($query);
	$db->bind_param('i', $data->session_id);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(! isset($res) or $res === FALSE){
		error_log('api/invitation: Cannot check for needed boat-full-messages for session: ' . $data->session_id);
	}
	if( count($res) < 1){
		error_log('api/invitation: Cannot check for needed boat-full-messages for session, session not found: ' . $data->session_id);
	}elseif($res[0]['free'] < 1){
		// Inform users that their invitation is gone
		send_session_full_email($data->session_id, $session_data['user_id'], $db, $config);
		
		// Set all the still existing invitations to 'boat-is-full'
		$query = 'UPDATE invitation SET status = 1
		          WHERE session_id = ?
				    AND status = 2';
		$db->prepare($query);
		$db->bind_param('i', $data->session_id);
		if(!$db->execute()){
			error_log('api/invitation: Cannot set inviations to boat-full-status: ' . $data->session_id);
		}
	}
	
	$db->disconnect();
	return;	
  }
  
  // send notification to the users that they got invited to a session
  function send_invitation_email($session_id, $my_user_id, $db, $config){
	$query = 'SELECT u.first_name AS fn, u.last_name AS ln,
                     u.email AS email, 
                     s.date AS date, s.start AS start, s.end AS end,
                     s.title AS title
              FROM invitation inv, session s, user u
              WHERE inv.session_id = ? 
                AND inv.session_id = s.id
                AND inv.user_id    = u.id
                AND inv.user_id    <> ?
                AND inv.status = 5';
	$db->prepare($query);
	$db->bind_param('ii', $session_id, $my_user_id);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(! isset($res) or $res===FALSE){
		error_log('api/invitation: Cannot get users to send invitation email to session: ' . $data->session_id);
		return false;
	}
	
	for($i = 0; $i<count($res); $i++){
		$message = 'Dear ' . $res[$i]['fn'] . ' ' . $res[$i]['ln'] . "\n";
        $message .= "\nThere is a new upcoming session. You just got an invitation:\n";
        $message .= " Title: " . utf8_decode($res[$i]['title']) . "\n";
        
        $start = strtotime($res[$i]['date'] . " " . $res[$i]['start']);
        $end = strtotime($res[$i]['date'] . " " . $res[$i]['end']);
        $message .= " Date : " . date("D d.m.Y", $start) . "\n";
        $message .= "        " . date("H:i", $start) . " to " . date("H:i", $end) . "\n";
        $message .= "\nPlease login to your account on www.wakeandsurf.ch and check your Schedule to accept or decline.\n";
        $message .= "\nSee you on the lake soon\n";
        Email::sendMail($res[$i]['email'], 'Session Invitation', $message, $config);
	}
	
	$query = 'UPDATE invitation 
                 SET status = 2 
                 WHERE session_id = ? AND status = 5';
	$db->prepare($query);
	$db->bind_param('i', $session_id);
	if(!$db->execute()){
		error_log('api/invitation: Cannot update invitation status to 5 for session: ' . $data->session_id);
		return false;
	}
	
	return true;
  }
  
  // send 'session full' to the invitees
  function send_session_full_email($session_id, $my_user_id, $db, $config){
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
	$db->bind_param('ii', $session_id, $my_user_id);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(! isset($res) or $res===FALSE){
		error_log('api/invitation: Cannot get users to send boat-full email to session: ' . $data->session_id);
		return false;
	}
	
	for($i = 0; $i<count($res); $i++){
		$message = 'Dear ' . $res[$i]['fn'] . ' ' . $res[$i]['ln'] . "\n";
        $message .= "\nThis E-mail is to inform you that you have been invited to the following session.\n";
        $message .= "Unfortunately the boat is already full and therefore the invitation is cancelled:\n";
        $message .= "\n Title: " . utf8_decode($res[$i]['title']) . "\n";
        
        $start = strtotime($res[$i]['date'] . " " . $res[$i]['start']);
        $end = strtotime($res[$i]['date'] . " " . $res[$i]['end']);
        $message .= " Date : " . date("D d.m.Y", $start) . "\n";
        $message .= "        " . date("H:i", $start) . " to " . date("H:i", $end) . "\n";
        $message .= "\nSee you on the lake soon\n";
        Email::sendMail($res[$i]['email'], 'Session Invitation', $message, $config);
	}
	
	$query = 'UPDATE invitation 
                 SET status = 2 
                 WHERE session_id = ? AND status = 5';
	$db->prepare($query);
	$db->bind_param('i', $session_id);
	if(!$db->execute()){
		error_log('api/invitation: Cannot update invitation status to 5 for session: ' . $data->session_id);
		return false;
	}
	
	return true;
    
  }
  
?>