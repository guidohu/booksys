<?php
  
  // automatically load all classes
  spl_autoload_register('login_autoloader');
  function login_autoloader($class){
	include '../classes/'.$class.'.php';
  }
  
  // check if we have an action
  if(!isset($_GET['action'])){
	HttpHeader::setResponseCode(200);
	exit;
  }

  // configuration
  $configuration = new Configuration();
  
  switch($_GET['action']){
	case 'login':
		$response = login($configuration);
		echo json_encode($response);
		exit;
	case 'isLoggedIn':
		$response = isLoggedIn($configuration);
		echo json_encode($response);
		exit;
	default:
		exit;
  }
  
  // The functions
  //------------------------------------------------------
  
  // Returns either loggedIn = true or loggedIn = false
  function isLoggedIn($configuration){
	// Prepare default result
	$res = array();
	$res['loggedIn'] = false;

	$login = new Login($configuration);
	if($login->isLoggedIn()){
		return Status::successDataResponse("logged in", [ "loggedIn" => TRUE ]);
	}

	return Status::successDataResponse("not logged in", [ "loggedIn" => FALSE ]);
  }
  
  // performs a user login
  function login($configuration){
	$post_data = json_decode(file_get_contents('php://input'));
	
	$sanitizer = new Sanitizer();
	
	// general input validation
	if(!isset($post_data->username) or 
	   (!$sanitizer->isEmail($post_data->username) 
	    and !$sanitizer->isAsciiText($post_data->username)
	   )
	  ){
		error_log('No valid username provided: "' . $post_data->username . '"');
		return Status::errorStatus("Please use a valid username");
	}
	if(!isset($post_data->password) or !$sanitizer->isAsciiText($post_data->password)){
		error_log('No valid password provided for user: "' . $post_data->username . '"');
		return Status::errorStatus("Please use a valid password");
	} 
	
	// input is valid and we can try to login
	$login  = new Login($configuration);
	$result = $login->login($post_data->username, $post_data->password, $configuration);
	$response = array();
	switch ($result){
		case -1:
			// wrong username/password
			return Status::errorStatus("invalid username/password");
		case -2:
			// user is locked
			return Status::errorStatus("user account is locked");
		case 1:
			// successful login
			return Status::successStatus("login successful");
	}
  }
  
?>