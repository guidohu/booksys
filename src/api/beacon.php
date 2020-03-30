<?php	
	// automatically load all classes
	spl_autoload_register('beacon_autoloader');
	function beacon_autoloader($class){
		include '../classes/'.$class.'.php';
	}

	// get configuration access
	$configuration = new Configuration();
	$status = array();
	
	// only users which are logged in are allowed to see booking information
	$lc = new Login($configuration);
	if(!$lc->isLoggedIn()){
		$status['ok'] = FALSE;
		$status['redirect'] = '/index.html';
		$status['msg'] = "you are not logged in";
		echo json_encode($status);
		exit;
	}
	
	$status['ok'] = TRUE;
	$status['msg'] = "you are logged in";
	echo json_encode($status);
	exit;

?>