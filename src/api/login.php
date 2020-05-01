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
		login($configuration);
		exit;
	case 'isLoggedIn':
		isLoggedIn($configuration);
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
		$res['loggedIn'] = TRUE;
	}

	echo json_encode($res);
	return;
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
		HttpHeader::setResponseCode(400);
		error_log('No valid username provided: "' . $post_data->username . '"');
		echo 'No valid username given';
		return;
	}
	if(!isset($post_data->password) or !$sanitizer->isAsciiText($post_data->password)){
		HttpHeader::setResponseCode(400);
		error_log('No valid password provided for user: "' . $post_data->username . '"');
		echo 'No valid password given';
		return;
	} 
	
	// input is valid and we can try to login
	$login  = new Login($configuration);
	$result = $login->login($post_data->username, $post_data->password, $configuration);
	$response = array();
	switch ($result){
		case -1:
			// wrong username/password
			$response['login_successful']  = false;
			$response['login_status_code'] = -1;
			echo json_encode($response);
			return;
		case -2:
			// user is locked
			$response['login_successful']  = false;
			$response['login_status_code'] = -2;
			echo json_encode($response);
			return;
		case 1:
			// successful login
			$response['login_successful']  = true;
			$response['login_status_code'] = 1;
			echo json_encode($response);
			return;
	}
  }
  
?>