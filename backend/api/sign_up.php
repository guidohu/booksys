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
		echo json_encode(Status::errorStatus("no API action provided"));
		exit;
	}
	
	$response = "";
	switch($_GET['action']){
		case 'sign_up':
			$response = sign_up();
			break;
		default:
			$response = Status::errorStatus("action not known to this API");
			break;
	}
	
	echo json_encode($response);
	return;
  
	// registers a user in the database
	// the user will have 'guest' status and will be locked
	function sign_up(){
		$data = json_decode(file_get_contents('php://input'));
		
		// input validation
		$sanitizer = new Sanitizer();
		if(! isset($data->username) or !$sanitizer->isAlphaNumEmail($data->username)){
			return Status::errorStatus("Your email address seems to have a typo in it.");
		}
		if(! isset($data->password) or !$sanitizer->isAsciiText($data->password)){
			return Status::errorStatus("Password is somehow corrupt.");
		}
		if(! isset($data->first_name)){
			return Status::errorStatus("No first name provided.");
		}
		if(! isset($data->last_name)){
			return Status::errorStatus("No last name provided.");
		}
		if(! isset($data->address)){
			return Status::errorStatus("No address given.");
		}
		if(! isset($data->mobile) or !$sanitizer->isMobileNumber($data->mobile)){
			return Status::errorStatus("The mobile number is not correct.");
		}
		if(! isset($data->plz) or !$sanitizer->isInt($data->plz)){
			return Status::errorStatus("Your Zip Code is wrong");
		}
		if(! isset($data->city)){
			return Status::errorStatus("No city provided");
		}
		if(! isset($data->email) or !$sanitizer->isEmail($data->email)){
			return Status::errorStatus("Please check your Email again.");
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
				return Status::errorStatus("reCAPTCHA token missing, are you a robot? Pleaes click I'm not a robot.");
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
				return Status::errorStatus("reCAPTCHA token not valid.");
			} 
		}
		
		// check if user already exists in the database
		$query = 'SELECT username FROM user WHERE username = ?';
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("The backend cannot connect to the database.");
		}
		$db->prepare($query);
		$db->bind_param('s', $data->username);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(!isset($res) or count($res)>0){
			$db->disconnect();
			return Status::errorStatus("Username/Email already taken.");
		}

		// calculate the password hash
		$password_salt = rand(0, 65635);
		$password_hash = crypt($data->password, '$6$rounds=5000$'.$password_salt);
		
		// add the user to the database
		$query = 'INSERT INTO user 
					(username, password_salt, password_hash, first_name, last_name, 
					address, city, plz, mobile, email, license, 
					status, locked, comment)
								VALUES 
					(?,?,?,?,?,?,?,?,?,?,?,?,1,?);';
		$db->prepare($query);
		$db->bind_param('sssssssissiis',
			$data->username,
			$password_salt,
			$password_hash,
			$data->first_name,
			$data->last_name,
			$data->address,
			$data->city,
			$data->plz,
			$data->mobile,
			$data->email,
			$data->license,
			$configuration->default_guest_user_status_id,
			$data->comment);
		
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/sign_up: Cannot add user to database');
			return Status::errorStatus("User cannot be added due to unknown reason. Please contact the administrator.");
		}

		// get the user-id of the created user
		$query = 'SELECT id FROM user WHERE username = ?';
		$db->prepare($query);
		$db->bind_param('s', $data->username);
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/sign_up: Cannot get ID of newly created user');
			return Status::errorStatus("User was not created and cannot be found in the database. Please contact the administrator.");
		}
		$result = $db->fetch_stmt_hash();
		if(count($result) != 1){
			$db->disconnect();
			error_log('api/sign_up: Cannot get ID of newly created user (non-unique)');
			return Status::errorStatus("User was not created and cannot be found in the database. Please contact the administrator.");
		}
		$db->disconnect();

		// return response
		$resp = array();
		$resp['user_id'] = $result[0]['id'];

		return Status::successDataResponse('success', $resp);
	}  
?>