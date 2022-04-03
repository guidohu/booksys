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

	public static function isError($status){
		return ($status['ok'] == FALSE);
	}
    
    public static function errorStatus($error){
        return Status::status(FALSE, $error);
    }

    public static function successStatus($msg){
		return Status::status(TRUE, $msg);
	}

	public static function successDataResponse($msg, $data){
		return Status::status(TRUE, $msg, $data);
	}

	public static function redirectStatus($msg, $redirect){
		$status = array();
		$status['ok'] = TRUE;
		$status['type'] = 'redirect';
		$status['msg'] = $msg;
		$status['target'] = $redirect;

		return $status;
	}

	public static function status($flag, $msg, $data=NULL){
		$status = array();
		$status['ok'] = $flag;
		$status['msg'] = $msg;

		if(isset($data)){
			$status['data'] = $data;
		}

		return $status;
	}
}
?>