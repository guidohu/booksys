<?php
	// Function which handles signing up of a user
	/// It adds a user to the database but sets him
	
	// automatically load all classes
	spl_autoload_register('sign_up_autoloader');
	function sign_up_autoloader($class){
			include '../classes/'.$class.'.php';
	}
	
	// User does not need to be logged in. This is basically an open API
	// therefore Input Validation is crucial
	
	// check if we have an action
	if(!isset($_GET['action'])){
			HttpHeader::setResponseCode(200);
			exit;
	}
	
	switch($_GET['action']){
		case 'sign_up':
			sign_up();
			exit;
	}
	
	HttpHeader::setResponseCode(400);
	return;
  
	// registers a user in the database
	// the user will have 'guest' status and will be locked
	function sign_up(){
		$data = json_decode(file_get_contents('php://input'));

		$status = array();
		$status['ok'] = TRUE;
		
		// input validation
		$sanitizer = new Sanitizer();
		if(! isset($data->username) or !$sanitizer->isAlphaNumEmail($data->username)){
			print_status(FALSE, "Username is not valid, no special characters are allowed.");
			return;
		}
		if(! isset($data->password) or !$sanitizer->isAsciiText($data->password)){
			print_status(FALSE, "Password is somehow corrupt.");
			return;
		}
		if(! isset($data->first_name)){
			print_status(FALSE, "No first name provided.");
			return;
		}
		if(! isset($data->last_name)){
			print_status(FALSE, "No last name provided");
			return;
		}
		if(! isset($data->address)){
			print_status(FALSE, "No address given");
			return;
		}
		if(! isset($data->mobile) or !$sanitizer->isMobileNumber($data->mobile)){
			print_status(FALSE, "The mobile number is not correct.");
			return;
		}
		if(! isset($data->plz) or !$sanitizer->isInt($data->plz)){
			print_status(FALSE, "Your Zip Code is wrong");
			return;
		}
		if(! isset($data->city)){
			print_status(FALSE, "No city provided");
			return;
		}
		if(! isset($data->email) or !$sanitizer->isEmail($data->email)){
			print_status(FALSE, "Please check your Email again.");
			return;
		}
		if(!isset($data->license) or $data->license != 1){
			$data->license = 0;
		}
		if(! isset($data->comment)){
			$data->comment = '';
		}

		// Get configuration access
		$configuration = new Configuration();

		if(isset($configuration->recaptcha_privatekey)){
			if(!isset($data->recaptcha_token)){
				print_status(FALSE, "reCAPTCHA token missing, are you a robot? Pleaes click I'm not a robot.");
				return;
			}
			// check recaptcha
			$recaptcha_url      = 'https://www.google.com/recaptcha/api/siteverify';
			$recaptcha_secret   = $configuration->recaptcha_privatekey;
			$recaptcha_response = $data->recaptcha_token;

			// Make and decode POST request:
			$recaptcha = file_get_contents($recaptcha_url . '?secret=' . $recaptcha_secret . '&response=' . $recaptcha_response);
			$recaptcha = json_decode($recaptcha);

			// Take action based on the result returned
			if($recaptcha->success != 1){
				print_status(FALSE, "reCAPTCHA token not valid");
				return;
			} 
		}
		
		
		// check if user already exists in the database
		$query = 'SELECT username FROM user WHERE username = ?';
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			print_status(FALSE, "The backend cannot connect to the database");
			return;
		}
		$db->prepare($query);
		$db->bind_param('s', $data->username);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(!isset($res) or count($res)>0){
			print_status(FALSE, "Username/Email already taken");
			$db->disconnect();
			return;
		}
		
		// add the user to the database
		$query = 'INSERT INTO user 
					(username, password, first_name, last_name, 
					address, city, plz, mobile, email, license, 
					status, locked, comment)
								VALUES 
					(?,?,?,?,?,?,?,?,?,?,?,1,?);';
		$db->prepare($query);
		$db->bind_param('ssssssissiis',
			$data->username,
			$data->password,
			$data->first_name,
			$data->last_name,
			$data->address,
			$data->city,
			$data->plz,
			$data->mobile,
			$data->email,
			$data->license,
			$configuration->guest_user_status_id,
			$data->comment);
		
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/sign_up: Cannot add user to database');
			print_status(FALSE, "User cannot be added due to unknown reason. Please contact the administrator.");
			return;
		}

		// get the user-id of the created user
		$query = 'SELECT id FROM user WHERE username = ?';
		$db->prepare($query);
		$db->bind_param('s', $data->username);
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/sign_up: Cannot get ID of newly created user');
			print_status(FALSE, "User was not created and cannot be found in the database. Please contact the administrator.");
			return;
		}
		$result = $db->fetch_stmt_hash();
		if(count($result) != 1){
			$db->disconnect();
			error_log('api/sign_up: Cannot get ID of newly created user (non-unique)');
			print_status(FALSE, "User was not created and cannot be found in the database. Please contact the administrator.");
			return;
		}
		$db->disconnect();

		// return response
		$status['user_id'] = $result[0]['id'];
		echo json_encode($status);

		return;
	}

	// helper function to print the status
	function print_status($ok, $msg){
		$status = array();
		$status['ok'] = $ok;
		$status['msg'] = $msg;
		echo json_encode($status);
		return;
	}
  
?>