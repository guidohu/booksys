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

	public function isCurrency($text){
		return $this->match($text, '/^[A-Z]{1,3}$/');
	}
	
	public function isFloat($float){
		return $this->match($float, '/^-?\d+((\.|,)\d*)?$/');
	}

	public function isAlphaNum($text){
		return $this->match($text, '/^[a-zA-Z0-9]+$/');
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

	public function isIBAN($iban){
		return $this->match($iban, '/^([A-Z]{2}[ \-]?[0-9]{2})(?=(?:[ \-]?[A-Z0-9]){9,30}$)((?:[ \-]?[A-Z0-9]{3,5}){2,7})([ \-]?[A-Z0-9]{1,3})?$/');
	}

	public function isBIC($bic){
		return $this->match($bic, '/^[a-zA-Z]{4}[a-zA-Z]{2}[a-zA-Z0-9]{2}[XXX0-9]{0,3}/');
	}
	
	public function isDate($date){
		return $this->match($date, '/\d{4}-\d{2}-\d{2}/'); 
	}
	
	public function isTimeOfDay($time){
		return $this->match($time, '/\d{1,2}[\.:]\d{1,2}([\.:]\d{1,2})?/');
	}

	public function isServerAddress($address){
		return $this->match($address, '/^[\w0-9]+(\.[\w0-9]+){0,}(:[0-9]+)?$/');
	}

	public function isGoogleMapsURL($url){
		return $this->match($url, '/^https:\/\/www\.google\.com\/maps\/embed\?pb=[^"\s]+$/');
	}

	public function isLongitude($longitude){
		if($this->isFloat($longitude) && $longitude <= 180 && $longitude >= -180){
			return TRUE;
		}
		return FALSE;
	}

	public function isLatitude($latitude){
		if($this->isFloat($latitude) && $latitude <= 90 && $latitude >= -90){
			return TRUE;
		}
		return FALSE;
	}

	public function isTimeZone($timezoneString){
		// check with regex first
		if(! $this->match($timezoneString, '/^[A-Za-z]+\/[A-Za-z]+$/')){
			error_log("Does not match regex");
			return FALSE;
		}
		$response = TRUE;
		try {
			// try to build a DateTimeZone object
			$timezone = new DateTimeZone($timezoneString);
			if($timezone == FALSE){
				$response = FALSE;
			}
		} catch(Exception $e){
			$response = FALSE;
		}
		return $response;
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

	public function isEngineHourFormat($format){
		switch ($format) {
			case "hh:mm":
				return TRUE;
			case "hh.h":
				return TRUE;
			default:
				return FALSE;
		}
	}
	
	function match($var, $regex){
		if(isset($var) and preg_match($regex, $var)){
			return TRUE;
		}
		return FALSE;
	}
	
	
	

}
?>