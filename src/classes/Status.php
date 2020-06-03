<?php

spl_autoload_register('Status::autoloader');

/* Sanitizes user input */
class Status{
    public static function autoloader($class){
		include $class . '.php';
    }

    public function __construct(){
		// nothing to do
	}
    
    public static function errorStatus($error){
        return Status::status(FALSE, $error);
    }

    public static function successStatus($msg){
		return Status::status(TRUE, $msg);
	}

	public static function status($flag, $msg){
		$status = array();
		$status['ok'] = $flag;
		$status['msg'] = $msg;
		return $status;
	}
}
?>