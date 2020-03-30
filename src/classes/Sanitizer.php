<?php

spl_autoload_register('Sanitizer::autoloader');

/* Sanitizes user input */
class Sanitizer{

	public static function autoloader($class){
		include $class . '.php';
	}

	public function __construct($conf=NULL){
		// nothing to do
	}

	public function isCookie($cookie){
		return $this->match($cookie, '/^[a-f0-9]+$/');
	}
	
	public function isInt($int){
		return $this->match($int, '/^\d+$/');
	}
	
	public function isAsciiText($text){
		return $this->match($text, '/^[a-zA-Z0-9\s]+$/');	
	}
	
	public function isFloat($float){
		return $this->match($float, '/^\d+((\.|,)\d*)?$/');
	}
	
	public function isAlphaNumEmail($text){
		return $this->match($text, '/^[0-9a-zA-Z@_\.-]+$/');
	}
	
	public function isBoolean($bool){
		return $this->match($bool, '/^(true)|(false)|1|0$/');
	}
	
	public function isMobileNumber($number){
		return $this->match($number, '/^\+?[0-9]+$/');
	}
	
	public function isEmail($email){
	    return (filter_var($email, FILTER_VALIDATE_EMAIL));
	}
	
	public function isDate($date){
		return $this->match($date, '/\d{4}-\d{2}-\d{2}/'); 
	}
	
	public function isTimeOfDay($time){
		return $this->match($time, '/\d{1,2}[\.:]\d{1,2}([\.:]\d{1,2})?/');
	}

	public function isUserRoleId($id){
		if(! $this->match($id, '/^\d{1}$/')){
			return FALSE;
		}
		if($id < 1 or $id > 3){
			return FALSE;
		}
		return TRUE;
	}
	
	function match($var, $regex){
		if(isset($var) and preg_match($regex, $var)){
			return TRUE;
		}
		return FALSE;
	}
	
	
	

}
?>