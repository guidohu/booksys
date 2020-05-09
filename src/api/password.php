<?php
  /**
   * The password API provides workflows to change the password, either through getting
   * a password reset token via email or by changing the password by providing the old
   * password too.
   */

  // automatically load all classes
  spl_autoload_register('password_autoloader');
  function password_autoloader($class){
	include '../classes/'.$class.'.php';
  }
  
  // check if we have an action
  if(!isset($_GET['action'])){
	HttpHeader::setResponseCode(200);
	exit;
  }

  $configuration = new Configuration();
  $lc = new Login($configuration);
  $response = array(
	  'ok'	=> FALSE,
	  'msg' => 'unknown API call'
  );
  
  // check which action to perform
  switch($_GET['action']){
	case 'token_request':
		$response = token_request($configuration);
		break;
	case 'change_password_by_token':
		$response = change_password_by_token($configuration);
		break;
	case 'change_password_by_password':
	    // only a user which is logged in is allowed to change a 
		// password
		if(! $lc->isLoggedIn()){
			$response['ok'] = FALSE;
			$response['message'] = "You are not logged in.";
			break;
	    }
		$response = change_password_by_password($configuration);
		break;
  }

  echo json_encode($response);
  return;
  
  /**
   * token_request will generate a password reset token that is sent to the user via email
   *
   * @param $configuration		configuration object
   * @return $response of the form
   * {
   * 	ok		=>	TRUE|FALSE		indicating success or failure
   * 	message => "somestring"		indicating additional information especially for failures
   * }
   * 
   * Note: the function does hide non existing users requests towards the frontend to prevent
   *       user enumeration
   */
  function token_request($configuration){
	$data = json_decode(file_get_contents('php://input'));
  
    # input validation
	$sanitizer = new Sanitizer();
	$response  = array(
		'ok' 		=> TRUE,
		'message' 	=> 'token requested'
	);

	if(!isset($data->email) or !$sanitizer->isEmail($data->email)){
		$response['ok'] = FALSE;
		$response['message'] = "Your provided Email is not a valid email address";
		return $response;
	}

	if(isset($configuration->recaptcha_privatekey)){
		if(!isset($data->recaptcha_token)){
			$response['ok'] = FALSE;
			$response['message'] = "reCAPTCHA token missing, are you a robot? Pleaes click I'm not a robot.";
			return $response;
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
			$response['ok'] = FALSE;
			$response['message'] = "reCAPTCHA token not valid";
			return $response;
		} 
	}
	
	# get user by email
	$db = new DBAccess($configuration);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
		echo "Internal server error.";
		exit;
	}
	
	$query = 'SELECT id as id, username as username, 
	                 first_name as first_name, last_name as last_name
			  FROM user
			  WHERE email = ?';
	$db->prepare($query);
	$db->bind_param('s', $data->email);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(!isset($res) or count($res)<1){
		// do not give any information
		$response['ok'] = TRUE;
		$response['message'] = "token requested, please check your email inbox";
		error_log("api/password: The provided email is not known: $data->email");
		$db->disconnect();
		return $response;
	}elseif(count($res)>1){
		// do not give any information
		$response['ok'] = TRUE;
		$response['message'] = "token requested, please check your email inbox";
		error_log("api/password: The provided email is not unique: $data->email");
		$db->disconnect();
		return $response;
	}
	
	$user_id    = $res[0]['id'];
	$username   = $res[0]['username'];
	$first_name = $res[0]['first_name'];
	$last_name  = $res[0]['last_name'];
	
	# generate random token
	$token = rand(100000, 999999);
	error_log("new token generated: " . $token); // TODO delete
	$token_hash = hash('sha256', $token);
	
	# store token in db
	$query = sprintf(
		'INSERT INTO password_reset (user_id, token, valid) VALUES ( "%d", "%s", 1);',
		$user_id,
		$token_hash
	);
	$res = $db->query($query);
	if($res == FALSE){
		HttpHeader::setResponseCode(500);
		echo "Internal server error";
		$db->disconnect();
		exit;
	}
	$db->disconnect();
	
	# send token by mail
	$mail = new Email();
	$message = <<< _END
Dear $first_name $last_name

You've just requested a password reset. Please use this token when asked, to reset your password:

Token: $token
 
See you on the lake soon
 
_END;
 
	if(!$mail->sendMail($data->email, 'Password Reset Token', $message, $configuration)){
		HttpHeader::setResponseCode(500);
		echo "Password reset token could not be sent, please verify your<br>input or contact us.";
		exit;
	}
	
	$response['message'] = "token requested, please check your email inbox";
	return $response;
  }
  
  /**
   * change_password_by_token changes a user's password given the presented token is valid.
   * It only allows one attempt after that one, all tokens of this user are invalidated and a
   * new one has to be requested.
   *
   * @param $configuration		configuration object
   * @return $response of the form
   * {
   * 	ok		=>	TRUE|FALSE		indicating success or failure
   * 	message => "somestring"		indicating additional information especially for failures
   * }
   */
  function change_password_by_token($configuration){
	$data = json_decode(file_get_contents('php://input'));
	
	# validate input
	$sanitizer = new Sanitizer();
	$response  = array(
		'ok' 		=> TRUE,
		'message' 	=> 'token requested'
	);
	if(!isset($data->token) or !$sanitizer->isCookie($data->token)){
	    $response['ok'] = FALSE;
		$response['message'] = "Your provided token is not in a valid format";
		return $response;
	}	
	if(!isset($data->email) or !$sanitizer->isEmail($data->email)){
		$response['ok'] = FALSE;
		$response['message'] = "Your provided email is not in a valid format";
		return $response;
	}
	if(!isset($data->password) or !$sanitizer->isCookie($data->password)){
	    $response['ok'] = FALSE;
		$response['message'] = "Your provided password is not in a valid format";
		return $response;
	}
	
	# calculate hash of token
	$token_hash = hash('sha256', $data->token);
	
	# check if we have a user_id with this hashed token
	$db = new DBAccess($configuration);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
		echo "Internal server error. Please let us know and try later.";
		exit;
	}	
	$query = 'SELECT pr.id FROM password_reset pr JOIN user u ON pr.user_id = u.id
	          WHERE u.email = ?
				AND pr.token   = ?
				AND pr.valid   = 1';
	$db->prepare($query);
	$db->bind_param(
		'ss', 
		$data->email, 
		$token_hash
	);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(!isset($res) or count($res)<1){
		# we do not know this token code
		$response['ok'] = FALSE;
		$response['message'] = "Your provided token code is not correct, please request a new one";
	}else{
		# calculate the password hash
		$new_salt          = rand(0, 65635);
		$new_password_hash = crypt($data->password, '$6$rounds=5000$'.$new_salt);

		# change password
		$query = 'UPDATE user SET password_salt = ?, password_hash = ? WHERE email = ?';
		$db->prepare($query);
		$db->bind_param(
			'sss', 
			$new_salt,
			$new_password_hash, 
			$data->email
		);
		if(!$db->execute()){
			error_log('Change password did not work: ' . $query);
			HttpHeader::setResponseCode(500);
			echo "Internal error, password could not be changed.<br>Please try again or contact us";
			$db->disconnect();
		}

		$response['ok'] = TRUE;
		$response['message'] = "Password has been reset";
	}
	
	# set all tokens of this user to invalid
	$query = 'UPDATE password_reset pr 
		JOIN user u ON u.id = pr.user_id 
		SET pr.valid = 0 
		WHERE u.email = ?';
	$db->prepare($query);
	$db->bind_param(
		's', 
		$data->email);
	if(! $db->execute()){
		error_log('Could not set all tokens to valid=0 for user '. $data->email.' and in function change_password_by_token');        
	}
	$db->disconnect();

	// TODO send an email to the user indicating that password has been reset

	return $response;	
  }
  
  /**
   * change_password_by_password allows to change the current password to a new one by presenting
   * both old and new password
   *
   * @param $configuration		configuration object
   * @return $response of the form
   * {
   * 	ok		=>	TRUE|FALSE		indicating success or failure
   * 	message => "somestring"		indicating additional information especially for failures
   * }
   */
  function change_password_by_password($configuration){
	$data = json_decode(file_get_contents('php://input'));

	$response  = array(
		'ok' 		=> TRUE,
		'message' 	=> 'password changed'
	);
	
	# validate input
	$sanitizer = new Sanitizer();
	if(!isset($data->password_old) or !$sanitizer->isAsciiText($data->password_old)){
	    $response['ok'] 	 = FALSE;
		$response['message'] = "Your provided old password is not in a valid format";
		return $response;
	}	
	if(!isset($data->password_new) or !$sanitizer->isAsciiText($data->password_new)){
	    $response['ok'] 	 = FALSE;
		$response['message'] = "Your provided new password is not in a valid format";
		return $response;
	}
	
	# get user object by session id and password
	$user      = new User($configuration);
	$user_data = $user->getUser();
	
	# check that the old password was correct
	if(! $user->isPasswordCorrect($data->password_old)){
		$response['ok'] 	 = FALSE;
		$response['message'] = "Wrong password provided";
		return $response;
	}

	# generate new password salt/hash
	$new_salt          = rand(0, 65635);
	$new_password_hash = crypt($data->password_new, '$6$rounds=5000$'.$new_salt);
	
	# update password
	$db = new DBAccess($configuration);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
		error_log("Cannot connect to database");
		echo "Internal server error.";
		exit;
	}
	$query = "UPDATE user 
		SET password = NULL,
		password_salt = ?,
		password_hash = ?
	    WHERE id = ?";
	$db->prepare($query);
	$db->bind_param('isi', 
		$new_salt,  
		$new_password_hash,
		$user_data['id']
	);
	if(!$db->execute()){
		error_log('Change password did not work: ' . $query);
		$db->disconnect();
		HttpHeader::setResponseCode(500);
		echo "Internal error, password could not be changed. Please try again or contact us";
		exit;
	}

	$db->disconnect();
	return $response;	
  }
    
?>