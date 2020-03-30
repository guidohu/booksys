<?php	
	// automatically load all classes
	spl_autoload_register('heat_autoloader');
	function heat_autoloader($class){
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
	
	// only users which are logged in and are admin are allowed to see booking information
	$lc = new Login($configuration);
	if(!$lc->isLoggedInAndRoleRedirect($configuration->admin_user_status_id, $configuration->login_page)){
		exit;
	}
	
	// check if we have an action
	if(!isset($_GET['action'])){
		HttpHeader::setResponseCode(200);
		exit;
	}
	
	switch($_GET['action']){
		case 'add_heat':
			add_heat($configuration, $db);
			$db->disconnect();
			exit;
		case 'add_heats':
			add_heats($configuration, $db);
			$db->disconnect();
			exit;
		case 'get_heats':
			get_heats($configuration, $db);
			$db->disconnect();
			exit;
		case 'get_session_heats':
			get_session_heats($configuration, $db);
			$db->disconnect();
			exit;
	}
	
	$db->disconnect();
	HttpHeader::setResponseCode(400);
	return;
	
	function add_heat($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));
		$status = check_and_add_heat($configuration, $db, $data);
		echo json_encode($status);
		return;
	}

	function add_heats($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		$response = array();
		foreach($data as $time => $entry){
			//echo json_encode($entry);
			//echo json_encode($entry->session_id);
			$response[$time] = check_and_add_heat($configuration, $db, $entry);
		}
		echo json_encode($response);
		return;
	}

	function get_session_heats($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		$response = array();
		$response['ok'] = TRUE;

		// sanitize input
		$sanitizer = new Sanitizer();
		if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
			error_log('api/heat.php: No "session_id" provided: ' . $data->from);
			$response['ok'] = FALSE;
			$response['msg'] = 'No "session_id" provided';
		}
		if(!$response['ok']){
			echo json_encode($response);
			return;
		}

		// query the heats
		$query = 'SELECT 
				h.id as heat_id,
				u.id as user_id,
				u.first_name as first_name,
				u.last_name as last_name,
				h.session_id as session_id,
				UNIX_TIMESTAMP(h.timestamp) as timestamp,
				h.duration_s as duration_s,
				FORMAT(h.cost_chf, 2) as cost,
				h.comment as comment
			FROM user u
			JOIN heat h ON h.user_id = u.id
			WHERE h.session_id = ? 
			ORDER BY h.timestamp DESC';
		$db->prepare($query);
		$db->bind_param('i', 
			$data->session_id
		);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		
		$response = array();
		$response['heats'] = $res;
		$response['currency'] = $configuration->currency;

		echo json_encode($response);
		return;
	}

	function get_heats($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		$response = array();
		$response['ok'] = TRUE;

		// sanitize input
		$sanitizer = new Sanitizer();
		if(!isset($data->from) or !$sanitizer->isInt($data->from)){
			error_log('api/heat.php: No "from" time-range provided: ' . $data->from);
			$response['ok'] = FALSE;
			$response['msg'] = 'No "from" time-range provided';
		}
		if(!isset($data->to) or !$sanitizer->isInt($data->to)){
			error_log('api/heat.php: No "to" time-range provided: ' . $data->from);
			$response['ok'] = FALSE;
			$response['msg'] = 'No "to" time-range provided';
		}
		if(!$response['ok']){
			echo json_encode($response);
			return;
		}

		// query the heats
		$query = 'SELECT 
				h.id as heat_id,
				u.id as user_id,
				u.first_name as first_name,
				u.last_name as last_name,
				h.session_id as session_id,
				UNIX_TIMESTAMP(h.timestamp) as timestamp,
				h.duration_s as duration_s,
				FORMAT(h.cost_chf, 2) as cost,
				h.comment as comment
			FROM user u
			JOIN heat h ON h.user_id = u.id
			WHERE UNIX_TIMESTAMP(h.timestamp) >= ?
			  AND UNIX_TIMESTAMP(h.timestamp) < ? 
			ORDER BY h.timestamp DESC';
		$db->prepare($query);
		$db->bind_param('ii', 
			$data->from,
			$data->to
		);
		$db->execute();
		$res = $db->fetch_stmt_hash();

		$response = array();
		$response['heats'] = $res;
		$response['currency'] = $configuration->currency;

		echo json_encode($response);
		return;
	}

	function check_and_add_heat($configuration, $db, $data){
		$status = array();

		// sanitize input
		$sanitized_msg = validate_heat_input($data);
		if($sanitized_msg != 'ok'){
			$status['ok']  = FALSE;
			$status['msg'] = $sanitized_msg;
			return $status;
		}

		// check if user exists
		$user_exists_msg = validate_user_exists($configuration, $data->user_id);
		if($user_exists_msg != 'ok'){
			$status['ok']  = FALSE;
			$status['msg'] = $user_exists_msg;
			return $status;
		}

		// check if session_id exists / or is Null
		$session_valid_msg = validate_session_id($db, $data->session_id);
		if($session_valid_msg != 'ok'){
			$status['ok']  = FALSE;
			$status['msg'] = $session_valid_msg;
			return $status;
		}

		// get cost_per minute for this user
		$cost_per_minute = get_price($configuration, $db, $data->user_id);
		if(!isset($cost_per_minute)){
			$status['ok']  = FALSE;
			$status['msg'] = 'Cannot get price for user';
			return $status;
		}		
		
		// round up (ceil) costs to the next 0.05
		$cost = floatval($cost_per_minute) * floatval($data->duration_s) / 60;
		$cost = ceil($cost / 0.05) * 0.05;
				
		// enter the heat to the database
		$entry = array();
		$entry['user_id']    = $data->user_id;
		$entry['session_id'] = $data->session_id;
		$entry['duration_s'] = $data->duration_s;
		$entry['cost']       = $cost;
		$entry['comment']    = $data->comment;

		if(db_add_heat($db, $entry)){
			$status['ok']  = TRUE;
			$status['msg'] = "Added heat for user $data->user_id";
			return $status;
		}else{
			$status['ok']  = FALSE;
			$status['msg'] = "Could not add heat for user $data->user_id";
			return $status;
		}
	}

	function validate_heat_input($data){
		$sanitizer = new Sanitizer();
		if(!isset($data->user_id) or !$sanitizer->isInt($data->user_id)){
			error_log('api/heat.php: No user_id provided: ' . $data->user_id);
			return 'No valid user_id provided';
		}
		if(isset($data->session_id) and 
		   !($sanitizer->isInt($data->session_id) or $data->session_id=='' or $data->session_id == NULL)){
			error_log('api/heat.php: No session_id provided: ' . $data->session_id);
			return 'No valid session_id provided';
		}
		if(!isset($data->duration_s) or !$sanitizer->isInt($data->duration_s)){
			error_log('api/heat.php: No duration_s provided: ' . $data->duration_s);
			return 'No valid duration provided';
		}
		// data->comment needs no sanitization as it will be treated as text
		return "ok";
	}

	// checks if a specific user_id exists
	function validate_user_exists($configuration, $user_id){
		// check if data makes sense
		// valid user_id?
		$u = new User($configuration);
		$user = $u->getUserById($user_id);
		if(!isset($user)){
			error_log('api/heat.php: Cannot load user to user_id: ' . $data->user_id);
			return "No user found with this ID";
		}
		return "ok";
	}

	function validate_session_id($db, $session_id){
		// valid session_id?
		$id = NULL;
		if($session_id != NULL && $session_id != ''){
			$id = $session_id;
			$query = 'SELECT id FROM session WHERE id = ?';
			$db->prepare($query);
			$db->bind_param('i', $id);
			$db->execute();
			$res = $db->fetch_stmt_hash();
			if(count($res) != 1){
				error_log('api/heat.php: Session to ID not found: ' . $session_id);
				return "No session found with this ID";
			}
		}
		return 'ok';
	}
	
	// Returns the price for a certain user status
	function get_price($configuration, $db, $user_id){
		$u = new User($configuration);
		$user = $u->getUserById($user_id);

		$query = 'SELECT price_chf_min FROM pricing WHERE user_status_id = ?';
		$db->prepare($query);
		$db->bind_param('i', $user['status']);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(count($res) != 1){
			error_log('api/heat.php: Cannot get price for user'. $user_id);
			return NULL;
		}else{
			return $res[0]['price_chf_min'];
		}
	}

	// add a given entry to the heat table in the database
	function db_add_heat($db, $entry){
		// enter the heat to the database
		error_log(print_r($entry, TRUE));
		$query = 'INSERT INTO heat (user_id, session_id, duration_s, cost_chf, comment)
		          VALUES (?, ?, ?, ?, ?);';
		$db->prepare($query);
		$db->bind_param('iiids',
			$entry['user_id'],
			$entry['session_id'],
			$entry['duration_s'],
			$entry['cost'],
			$entry['comment']);
		if($db->execute()){
			return TRUE;
		}else{
			error_log("api/heat.php: Cannot add heat for user/time:" . $entry->user_id . "/".$entry->duration_s);
			return FALSE;
		}
	}


?>