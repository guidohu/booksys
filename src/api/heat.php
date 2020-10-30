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
	
	$response = '';
	switch($_GET['action']){
		case 'add_heat':
			$response = add_heat($configuration, $db);
			break;
		case 'add_heats':
			$response = add_heats($configuration, $db);
			break;
		case 'get_heats':
			$response = get_heats($configuration, $db);
			break;
		case 'get_heat':
			$response = get_heat($configuration, $db);
			break;
		case 'update_heat':
			$response = update_heat($configuration, $db);
			break;
		case 'delete_heat':
			$response = delete_heat($configuration, $db);
			break;
		case 'get_session_heats':
			$response = get_session_heats($configuration, $db);
			break;
		default:
			HttpHeader::setResponseCode(400);
	}

	$db->disconnect();
	echo json_encode($response);
	return;
	
	// add a single heat
	function add_heat($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));
		$response = check_and_add_heat($configuration, $db, $data);
		return $response;
	}

	// add multiple heats
	function add_heats($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		$response = array();
		foreach($data as $time => $entry){
			$id = $time;
			if(isset($entry->uid)){
				$id = $entry->uid;
			}
			$response[$id] = check_and_add_heat($configuration, $db, $entry);
		}
		return Status::successDataResponse("check individual responses per heat", $response);
	}

	// get the sessions of a specific heat
	function get_session_heats($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// sanitize input
		$sanitizer = new Sanitizer();
		if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
			error_log('api/heat.php: No or invalid "session_id" provided');
			return Status::errorStatus('No "session_id" provided');
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
				p.price_chf_min as price_per_min,
				h.comment as comment
			FROM user u
			JOIN heat h ON h.user_id = u.id
			JOIN user_status us ON u.status = us.id 
			JOIN pricing p ON u.status = p.user_status_id 
			WHERE h.session_id = ? 
			ORDER BY h.timestamp DESC';
		$db->prepare($query);
		$db->bind_param('i', 
			$data->session_id
		);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		
		return Status::successDataResponse("success", $res);
	}

	// get all heats in between from and to
	function get_heats($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		$response = array();
		$response['ok'] = TRUE;

		// sanitize input
		$sanitizer = new Sanitizer();
		if(!isset($data->from) or !$sanitizer->isInt($data->from)){
			error_log('api/heat.php: No "from" time-range provided: ' . $data->from);
			return Status::errorStatus('No "from" time-range provided');
		}
		if(!isset($data->to) or !$sanitizer->isInt($data->to)){
			error_log('api/heat.php: No "to" time-range provided: ' . $data->from);
			return Status::errorStatus('No "to" time-range provided');
		}
		if(!$response['ok']){
			return $response;
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
		return $response;
	}

	// get a specific heat given its id
	function get_heat($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		$response = array();
		$response['ok'] = TRUE;

		// sanitize input
		$sanitizer = new Sanitizer();
		$error = validate_heat_id($data->heat_id, $sanitizer);
		if($error != NULL){
			return Status::errorStatus($error);
		}

		// get heat details
		$query = 'SELECT 
			h.id as heat_id, h.user_id as user_id, 
			h.duration_s as duration_s, h.cost_chf as cost, 
			UNIX_TIMESTAMP(h.timestamp) as timestamp,
			u.first_name as first_name, u.last_name as last_name, 
			p.price_chf_min as price_per_min
			FROM heat h 
			JOIN user u ON h.user_id = u.id 
			JOIN user_status us ON u.status = us.id 
			JOIN pricing p ON u.status = p.user_status_id 
			WHERE h.id = ?;';
		$db->prepare($query);
		$db->bind_param('i',
			$data->heat_id
		);
		$db->execute();
		$response = array();
		$result = $db->fetch_stmt_hash();
		if(count($result) != 1){
			$response['heat'] = new stdClass();
		}else{
			$response['heat'] = $result[0];
		}
		$response['currency'] = $configuration->currency;
		return $response;
	}

	// update a specific heat
	// - only the duration_s of a given heat_id/user_id is changed
	function update_heat($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// sanitize input
		$sanitizer = new Sanitizer();
		$error = validate_heat_id($data->heat_id, $sanitizer);
		if($error != NULL){
			return Status::errorStatus($error);
		}
		$error = validate_user_id($data->user_id, $sanitizer);
		if($error != NULL){
			return Status::errorStatus($error);
		}
		$error = validate_duration($data->duration_s, $sanitizer);
		if($error != NULL){
			return Status::errorStatus($error);
		}

		$comment = NULL;
		if(isset($data->comment)){
			$comment = $data->comment;
		}

		// check for existing user
		$error = validate_user_exists($configuration, $data->user_id);
		if($error != NULL){
			return Status::errorStatus($error);
		}
		// check for existing heat
		$error = validate_heat_exists($db, $data->heat_id);
		if($error != NULL){
			return Status::errorStatus($error);
		}

		// get cost for this user and heat
		$cost = 0;
		try{
			$cost = get_cost_for_heat($configuration, $db, $data->user_id, $data->duration_s);
		}catch (Exception $e){
			return Status::errorStatus($e->getMessage());
		}

		// update heat in database
		$query = 'UPDATE heat SET duration_s = ?, cost_chf = ?, comment = ?
			WHERE id = ?;';
		$db->prepare($query);
		$db->bind_param('idsi',
			$data->duration_s,
			$cost,
			$comment,
			$data->heat_id);
		if($db->execute()){
			return Status::successStatus("heat updated");
		}else{
			error_log("api/heat.php: Cannot update heat for user/time:" . $data->user_id . "/".$data->duration_s);
			return Status::errorStatus("Cannot update heat.");
		}
	}

	// delete a specific heat
	function delete_heat($configuration, $db){
		$data = json_decode(file_get_contents('php://input'));

		// sanitize input
		$sanitizer = new Sanitizer();
		$error = validate_heat_id($data->heat_id, $sanitizer);
		if($error != NULL){
			return Status::errorStatus($error);
		}

		// check for existing heat
		$error = validate_heat_exists($db, $data->heat_id);
		if($error != NULL){
			return Status::errorStatus($error);
		}

		// remove heat from database
		$query = 'DELETE FROM heat
			WHERE id = ?;';
		$db->prepare($query);
		$db->bind_param('i', $data->heat_id);
		if($db->execute()){
			return Status::successStatus("heat deleted");
		}else{
			error_log("api/heat.php: Cannot delete heat with id: " . $data->heat_id);
			return Status::errorStatus("Cannot delete heat.");
		}
	}

	function check_and_add_heat($configuration, $db, $data){
		$status = array();

		// sanitize input
		$sanitized_msg = validate_heat_input($data);
		if($sanitized_msg != 'ok'){
			return Status::errorStatus($sanitized_msg);
		}

		// check if user exists
		$error = validate_user_exists($configuration, $data->user_id);
		if($error != NULL){
			return Status::errorStatus($error);
		}

		// check if session_id exists / or is Null
		if($data->session_id != NULL && $data->session_id != ''){
			$error = validate_session_exists($db, $data->session_id);
			if($error != NULL){
				return Status::errorStatus($error);
			}
		}

		// get cost for this user and heat
		$cost = 0;
		try{
			$cost = get_cost_for_heat($configuration, $db, $data->user_id, $data->duration_s);
		}catch (Exception $e){
			return Status::errorStatus($e->getMessage());
		}

		// get the comment
		$comment = NULL;
		if(isset($data->comment)){
			$comment = $data->comment;
		}
				
		// enter the heat to the database
		$entry = array();
		$entry['user_id']    = $data->user_id;
		$entry['session_id'] = $data->session_id;
		$entry['duration_s'] = $data->duration_s;
		$entry['cost']       = $cost;
		$entry['comment']    = $comment;

		if(db_add_heat($db, $entry)){
			return Status::successStatus("Added heat for user $data->user_id");
		}else{
			return Status::errorStatus("Could not add heat for user $data->user_id");
		}
	}

	function validate_heat_input($data){
		$sanitizer = new Sanitizer();
		$error = validate_user_id($data->user_id, $sanitizer);
		if($error != NULL){
			return $error;
		}
		$error = validate_session_id($data->session_id, $sanitizer);
		if($error != NULL){
			return $error;
		}
		$error = validate_duration($data->duration_s, $sanitizer);
		if($error != NULL){
			return $error;
		}
		// data->comment needs no sanitization as it will be treated as text
		return "ok";
	}

	function validate_user_id($user_id, $sanitizer){
		if(!isset($user_id) or !$sanitizer->isInt($user_id)){
			error_log('api/heat.php: No user_id provided: ' . $user_id);
			return 'No valid user_id provided';
		}
		return NULL;
	}

	function validate_session_id($session_id, $sanitizer){
		if(isset($session_id) and 
		   !($sanitizer->isInt($session_id) or $session_id=='' or $session_id == NULL)){
			error_log('api/heat.php: No session_id provided: ' . $session_id);
			return 'No valid session_id provided';
		}
		return NULL;
	}

	function validate_duration($duration, $sanitizer){
		if(!isset($duration) or !$sanitizer->isInt($duration)){
			error_log('api/heat.php: No duration_s provided: ' . $duration);
			return 'No valid duration provided';
		}
		return NULL;
	}

	function validate_heat_id($heat_id, $sanitizer){
		if(!isset($heat_id) or !$sanitizer->isInt($heat_id)){
			error_log('api/heat.php: No "heat_id" provided: ' . $heat_id);
			return 'No "heat_id" provided';
		}
		return NULL;
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
		return NULL;
	}

	function validate_session_exists($db, $session_id){
		// valid session_id?
		$query = 'SELECT id FROM session WHERE id = ?';
		$db->prepare($query);
		$db->bind_param('i', $session_id);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(count($res) != 1){
			error_log('api/heat.php: Session to ID not found: ' . $session_id);
			return "No session found with this ID";
		}
		return NULL;
	}

	function validate_heat_exists($db, $heat_id){
		// valid heat?
		$query = 'SELECT id FROM heat WHERE id = ?';
		$db->prepare($query);
		$db->bind_param('i', $heat_id);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(count($res) != 1){
			error_log('api/heat.php: Heat to ID not found: ' . $heat_id);
			return "No heat found with this ID";
		}
		return NULL;
	}

	// returns the cost for the given user and duration
	function get_cost_for_heat($configuration, $db, $user_id, $duration_s){
		// get price for this user
		$price_per_min = get_price($configuration, $db, $user_id);
		if($price_per_min == NULL){
			throw new Exception('Cannot get price for the given user.');
		}

		// calculate the cost for the given duration
		$cost = floatval($price_per_min) * floatval($duration_s) / 60;
		$cost = ceil($cost / 0.05) * 0.05;

		return $cost;
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