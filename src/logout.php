<?php
	// Logout the user and redirect to login page
	include __DIR__.'/classes/Configuration.php';
	include __DIR__.'/classes/Login.php';
	$configuration = new Configuration();
	Login::logout($configuration);

	header("Location: /index.html?valid=1");
    exit;
?>

