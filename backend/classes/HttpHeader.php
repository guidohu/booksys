<?php

class HttpHeader{
	
	function __autoload($class_name){
		set_include_path(get_include_path(). PATH_SEPARATOR . '../classes');
		include_once $class_name.'.php';
	}
	
	// Sets the http response code
	static public function setResponseCode($code){
		$response = 'HTTP/1.1 ' . $code . ' ';
		
		switch($code){
			case 200:
				$response .= 'OK';
				break;
			case 302:
				$response .= "Found";
				break;
			case 400:
				$response .= 'Bad Request';
				break;
			case 401:
				$response .= 'Unauthorized';
				break;
			case 403:
				$response .= 'Forbidden';
				break;
			case 500:
				$response .= 'Internal Server Error';
				break;
			default:
				error_log('HttpHeader setResponseCode for '.$code.' not supported');
				return;
		}
		header($response);
		return;
	}
	
	// Sets the location of the HTTP header
	static public function setLocation($location){	
		header("Location: ".$location);
		return;
	}
}

?>