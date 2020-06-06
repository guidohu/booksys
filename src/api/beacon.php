<?php
    /**
     * The beacon API allows to know whether a client is connected or not and to
	 * update the session of the user
	 *  
     */

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
		Status::errorStatus("you are not logged in");
	} else {
		Status::successStatus("you are logged in");
	}
	
	echo json_encode($status);
	exit;

?>