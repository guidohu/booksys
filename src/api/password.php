<?php  
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
  
  switch($_GET['action']){
	case 'token_request':
		token_request($configuration);
		exit;
	case 'change_password_by_token':
		change_password_by_token($configuration);
		exit;
	case 'change_password_by_password':
	    // only a user which is logged in is allowed to change a 
		// password
		if(! $lc->isLoggedIn()){
		    HttpHeader::setResponseCode('401');
			echo "You have been logged out. Please login again.";
		    exit;
	    }
		change_password_by_password($configuration);
		exit;
  }
  
  HttpHeader::setResponseCode(400);
  return;
  
  function token_request($configuration){
	$data = json_decode(file_get_contents('php://input'));
  
    # input validation
	$sanitizer = new Sanitizer();
	if(!isset($data->email) or !$sanitizer->isEmail($data->email)){
	    HttpHeader::setResponseCode(400);
		echo "Your Email is not a valid email address";
		return;
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
		HttpHeader::setResponseCode(400);
		echo "No user with the given email address.";
		$db->disconnect();
		exit;
	}elseif(count($res)>1){
		HttpHeader::setResponseCode(400);
		echo "The provided email is not unique. Please contact an administrator.";
		$db->disconnect();
		exit;
	}
	
	$user_id    = $res[0]['id'];
	$username   = $res[0]['username'];
	$first_name = $res[0]['first_name'];
	$last_name  = $res[0]['last_name'];
	
	# generate random token
	$token = rand(100000, 999999);
	$token_hash = hash('sha256', $token);
	
	# store token in db
	$query = sprintf('INSERT INTO password_reset (user_id, token, valid)
	                  VALUES ( "%d", "%s", 1);',
					  $user_id,
					  $token_hash);
	$res = $db->query($query);
	if($res == FALSE){
		echo "Internal server error";
		$db->disconnect();
		HttpHeader::setResponseCode(500);
		exit;
	}
	$db->disconnect();
	
	# send token by mail
	$mail = new Email();
	$message = <<< _END
Dear $first_name $last_name

Please use this token to reset your password:

 Username: $username
 Token   : $token
 
See you on the lake soon
 
_END;
 
	if(!$mail->sendMail($data->email, 'Password Reset Token', $message)){
		echo "Password reset token could not be sent, please verify your<br>input or contact us.";
		http_reponse_code(500);
	}
	
	$ret = Array();
	$ret['user_id'] = $user_id;
	echo json_encode($ret);
	return;
  }
  
  /* Change the password of a user */
  function change_password_by_token($configuration){
	$data = json_decode(file_get_contents('php://input'));
	
	# validate input
	$sanitizer = new Sanitizer();
	if(!isset($data->token) or !$sanitizer->isCookie($data->token)){
	    HttpHeader::setResponseCode(400);
		echo "Wrong token format";
		return;
	}	
	if(!isset($data->user_id) or !$sanitizer->isInt($data->user_id)){
	    HttpHeader::setResponseCode(400);
		echo "Illegal user_id request";
		return;
	}
	if(!isset($data->password) or !$sanitizer->isCookie($data->password)){
	    HttpHeader::setResponseCode(400);
		echo "Wrong password format";
		return;
	}
	
	# calculate hash of token
	$token_hash = hash('sha256', $data->token);
	
	# check if we have a user_id with this token
	$db = new DBAccess($configuration);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
		echo "Internal server error.";
		exit;
	}	
	$query = 'SELECT id FROM password_reset
	          WHERE user_id = ?
				AND token   = ?
				AND valid   = 1';
	$db->prepare($query);
	$db->bind_param('is', $data->user_id, $token_hash);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	if(!isset($res) or count($res)<1){
		echo "The token code is not correct.";
		HttpHeader::setResponseCode(400);
		$db->disconnect();
		exit;
	}
	
	# set all tokens of this user to invalid
	$query = 'UPDATE password_reset SET valid = 0 WHERE user_id = ?';
	$db->prepare($query);
	$db->bind_param('i', $data->user_id);
	if(! $db->execute()){
		error_log('Could not set all tokens to valid=0 for user '. $data->user_id.' and in function change_password_by_token');        
	}
	
	# change password
	$query = 'UPDATE user SET password = ? WHERE id = ?';
	$db->prepare($query);
	$db->bind_param('si', $data->password, $data->user_id);
	if(!$db->execute()){
		error_log('Change password did not work: ' . $query);
		echo "Internal error, password could not be changed.<br>Please try again or contact us";
		$db->disconnect();
		HttpHeader::setResponseCode(500);
	}
	
	HttpHeader::setResponseCode(200);
	$db->disconnect();
	return;	
  }
  
  /* Change the password of a user */
  function change_password_by_password($configuration){
	$data = json_decode(file_get_contents('php://input'));
	
	# validate input
	$sanitizer = new Sanitizer();
	if(!isset($data->password_old) or !$sanitizer->isAsciiText($data->password_old)){
	    HttpHeader::setResponseCode(400);
		echo "Illegal password provided";
		return;
	}	
	if(!isset($data->password_new) or !$sanitizer->isAsciiText($data->password_new)){
	    HttpHeader::setResponseCode(400);
		echo "Illegal password provided";
		return;
	}
	
	# get user object by session id and password
	$user = new User($configuration);
	$user_data = $user->getUser();
	
	# check that the old password was correct
	if(! $user->isPasswordCorrect($data->password_old)){
		HttpHeader::setResponseCode(400);
		echo "Wrong password provided";
		return;
	}

	# generate new password salt/hash
	$new_salt          = rand(0, 65635);
	$new_password_hash = crypt($data->password_new, '$6$rounds=5000$'.$new_salt);
	
	# update password
	$db = new DBAccess($configuration);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
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
		$user_data['id']);
	if(!$db->execute()){
		error_log('Change password did not work: ' . $query);
		$db->disconnect();
		echo "Internal error, password could not be changed.<br>Please try again or contact us";
		HttpHeader::setResponseCode(500);
	}
	
	HttpHeader::setResponseCode(200);
	$db->disconnect();
	return;	
  }
    
?>