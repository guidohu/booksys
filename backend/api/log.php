<?php	
	// automatically load all classes
	spl_autoload_register('log_autoloader');
	function log_autoloader($class){
		include '../classes/'.$class.'.php';
	}

	// get configuration access
	$configuration = new Configuration();
	
	// Check if the user is already logged in and is of type admin
	$lc = new Login($configuration);
	if(!$lc->isAdmin()){
		header("Location:".$configuration->$configuration->login_page);
		HttpHeader::setResponseCode(302);
		exit;
	}
	
	// check if we have an action
	if(!isset($_GET['action'])){
		http_response_code(200);
		exit;
	}
	
	$response = NULL;
	switch($_GET['action']){
		case 'get_logs':
			$response = get_logs($configuration);
			break;
		default:
			$response = Status::errorStatus('action not supported');
			break;
	}
	
	echo json_encode($response);
	return;
	
	// FUNCTIONS
	//---------------------------------------------------	
	
	/* Returns the transactions */
	function get_logs($configuration){
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database");
		}
		
		$query = "SELECT c_log.id as id, c_log.type as type, c_log.time as time, c_log.log as log FROM (
					SELECT h.id as id, \"heat\" as type, h.timestamp as time, 
						   CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, 
						          ' was riding for ', FORMAT(h.duration_s/60,1), 
								  ' min and ', FORMAT(h.cost_chf, 2), ' $configuration->currency') as log 
						FROM user u, heat h 
						WHERE u.id = h.user_id AND h.comment IS NULL
					UNION ALL
					SELECT h.id as id, \"heat\" as type, h.timestamp as time, 
						   CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, 
						          ' was riding for ', FORMAT(h.duration_s/60,1), 
								  ' min and ', FORMAT(h.cost_chf, 2), ' $configuration->currency (Comment: ', h.comment, ')') as log 
						FROM user u, heat h 
						WHERE u.id = h.user_id AND NOT h.comment IS NULL
					UNION ALL
					SELECT p.id as id, \"payment\" as type, p.timestamp as time, 
					       CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, 
						          ' paid for sessions ', FORMAT(p.amount_chf, 2), ' $configuration->currency') as log 
						FROM user u, payment p 
						WHERE u.id = p.user_id AND p.type_id = 4
					UNION ALL
					SELECT bf.id as id, \"boat_fuel\" as type, bf.timestamp as time, 
					       CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, ' added ', 
						          FORMAT(bf.liters, 2), 'L fuel for ', FORMAT(cost_chf, 2), ' $configuration->currency') as log 
						FROM user u, boat_fuel bf 
						WHERE bf.user_id = u.id AND bf.contributes_to_balance = 1
					UNION ALL
					SELECT bf.id as id, \"boat_fuel\" as type, bf.timestamp as time, 
					       CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, ' added ', 
						          FORMAT(bf.liters, 2), 'L fuel for ', FORMAT(cost_chf, 2), ' $configuration->currency (on credit)') as log 
						FROM user u, boat_fuel bf 
						WHERE bf.user_id = u.id AND bf.contributes_to_balance = 0
                  ) as c_log ORDER BY time DESC";

		// limit the number of returned log line for mobile browsers
		if(MobileDevice::isMobileBrowser() or isset($_GET['mobile'])){
			$query .= "  LIMIT 0, 200";
		}
				  
		$res = $db->fetch_data_hash($query, -1);
		$db->disconnect();
		if(!$res){
			return Status::errorStatus("Cannot get logs from database");
		}
		return Status::successDataResponse("success", $res);
	}