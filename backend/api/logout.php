<?php
	// Logout the user
	include __DIR__.'/../classes/Configuration.php';
	include __DIR__.'/../classes/Login.php';
	$configuration = new Configuration();
	if(Login::logout($configuration)){
        echo json_encode([
            "ok" => TRUE,
            "message" => 'logged out'
        ]);
    }else{
        echo json_encode([
            "ok" => FALSE,
            "message" => 'error while logging out'
        ]);
    }
    exit;
?>