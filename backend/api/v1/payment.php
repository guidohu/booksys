<?php	
	// automatically load all classes
	spl_autoload_register('payment_autoloader');
	function payment_autoloader($class){
		include '../../classes/'.$class.'.php';
	}

	// get configuration access
	$configuration = new Configuration();
	
	// Check if the user is already logged in and is of type admin
	$lc = new Login($configuration);
	if(!$lc->isAdmin()){
		echo json_encode(Status::errorStatus("Insufficient permissions"));
		exit;
	}
	
	// check if we have an action
	if(!isset($_GET['action'])){
		echo json_encode(Status::errorStatus("API action not known"));
		exit;
	}
	
	$response = '';
	switch($_GET['action']){
		case 'add_expenditure':
		    $response = add_expenditure($configuration);
			break;
		case 'add_payment':
			$response = add_payment($configuration);
			break;
		default:
			$respone = Status::errorStatus("action unknown for this API call");
			break;
	}

	echo json_encode($response);
	return;
	
	// FUNCTIONS
	//---------------------------------------------------	
	
	function add_payment($configuration){
		$post_data = json_decode(file_get_contents('php://input'));
		
		$sanitizer = new Sanitizer();
		
		// general input validation
		if(!isset($post_data->user_id) or !$sanitizer->isInt($post_data->user_id)){
			return Status::errorStatus("No valid user selected, please select a user");
		}
		if(!isset($post_data->type_id) or !$sanitizer->isInt($post_data->type_id)){
			return Status::errorStatus("No valid income type");
		}
		if(!isset($post_data->date) or !$sanitizer->isDate($post_data->date)){
			return Status::errorStatus("No valid date");
		}
		if(!isset($post_data->amount) or !$sanitizer->isFloat($post_data->amount)){
			return Status::errorStatus("No valid amount");
		}
		if($post_data->type_id != 4 and $post_data->type_id != 6
		   and (!isset($post_data->comment) or $post_data->comment == '')){
			return Status::errorStatus("No comment provided");
		}else if((!isset($post_data->comment) or $post_data->comment == '') and $post_data->type_id == 4){
			$post_data->comment = 'Session Payment';
		}else if((!isset($post_data->comment) or $post_data->comment == '') and $post_data->type_id == 6){
			$post_data->comment = 'Membership Fee';
		}
		
		# add a new payment to the database
		$query = "INSERT INTO payment (user_id, type_id, timestamp, amount_chf, comment)
		          VALUES ( ?,?,?,?,? );";
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to the database");
		}
		
		$db->prepare($query);
		$db->bind_param('iisds',
		                $post_data->user_id,
						$post_data->type_id,
						$post_data->date,
						sprintf('%.2f', $post_data->amount),
						$post_data->comment);
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/payment: Cannot add payment');
			return Status::errorStatus("Income transaction cannot be added, due to an unknown error.");
		}
		$db->disconnect();
		
		# in case it was a user who paid for sessions
		# we write an email to the user
		if($post_data->type_id == 4){
			$user = new User($configuration);
			$user_data = $user->getUserById($post_data->user_id);
			$msg = 'Dear ' . $user_data['first_name'] . ' ' . $user_data['last_name'] . "\n\n"
			     . 'Thank you for your payment. We received ' . $post_data->amount . ' ' . $configuration->currency . ' and' . "\n"
				 . 'added the amount to your account balance.'."\n\n"
				 . 'See you soon on the boat' . "\n"
				 . 'wakeandsurf Crew';
			Email::sendMail($user_data['email'], 'Payment received', $msg, $configuration);
		}
		
		return Status::successStatus("Income added");
	}
	
	function add_expenditure($configuration){
		$post_data = json_decode(file_get_contents('php://input'));
		
		$sanitizer = new Sanitizer();
		
		// general input validation
		if(!isset($post_data->user_id) or !$sanitizer->isInt($post_data->user_id)){
			return Status::errorStatus("No valid user selected, please select a user");
		}
		if(!isset($post_data->type_id) or !$sanitizer->isInt($post_data->type_id)){
			return Status::errorStatus("No valid expenditure type");
		}
		if(!isset($post_data->date) or !$sanitizer->isDate($post_data->date)){
			return Status::errorStatus("No valid date");
		}
		if(!isset($post_data->amount) or !$sanitizer->isFloat($post_data->amount)){
			return Status::errorStatus("No valid amount");
		}
		if(!isset($post_data->comment) or $post_data->comment == ''){
			return Status::errorStatus("No comment provided");
		}
		if($post_data->type_id == 0){
			return Status::errorStatus("Use the boat API for fuel expenses");
		}
		
		# add a new expenditure to the database
		$query = "INSERT INTO expenditure (user_id, type_id, timestamp, amount_chf, comment)
		          VALUES (?,?,?,?,?);";
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database");
		}
		
		$db->prepare($query);
		$db->bind_param('iisds',
		                $post_data->user_id,
						$post_data->type_id,
						$post_data->date,
						sprintf('%.2f', $post_data->amount),
						$post_data->comment);
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/payment: Cannot add expenditure');
			return Status::errorStatus("Internal Server Error, cannot add the expenditure.");
		}
		
		$db->disconnect();
		return Status::successStatus("Expense added");
	}

?>