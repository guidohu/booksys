<?php  

class Configuration{
	// TODO: remove this in case we switch to JSON or a better
	// configuration format
	public $config_file_variables;

	// database config
	public $db_server;
	public $db_name;
	public $db_user;
	public $db_password;

	// user status
	public $guest_user_status_id;
	public $member_user_status_id;
	public $admin_user_status_id;

	// properties
	public $browser_session_timeout_max;
	public $location_time_zone;
	public $location_longitude;
	public $location_latitude;
	public $business_day_start;
	public $business_day_end;
	public $session_cancel_graceperiod;
	public $recaptcha_publickey;
	public $recaptcha_privatekey;

	// customization
	public $currency;
	public $location_map;
	public $location_address;
	public $payment_account_owner;
	public $payment_account_iban;
	public $payment_account_bic;
	public $payment_account_comment;


	// TODO: remove dependencies on these
	public $login_page;
	public $version;
	public $required_db_version;

    // Automatically load classes from the classes folder
	function __autoload($class_name){
		if(file_exists($class_name.'.php')){
			include $class_name.'.php';
		}
	}
	
	public function __construct(){
		// read the config file
		if(file_exists(dirname(__DIR__).'/config/config.php')){
			include dirname(__DIR__).'/config/config.php';
			$this->config_file_variables = $config;
		}else{
			//throw new Exception("config file does not exist");
			error_log("No config file found, return empty config");
			return NULL;
		}

		// get database configuration parameters
		// from the configuration file
		$this->get_configuration_file($this->config_file_variables);

		// get all the user status that
		// are required for the code to run
		// TODO: that dependency should be removed
		$this->get_user_status($this->config_file_variables);

		// get all the configuration settings from the
		// database
		$this->get_configuration_properties($this->config_file_variables);

		// set hard defaults
		$this->login_page          = '/index.html';
		$this->version             = "1.8.2";
		$this->required_db_version = "1.8";
	}

	// TODO: deprecate this function
	public function get_old_config(){
		return $this->config_file_variables;
	}

	// TODO: deprecate this function
	public function get_db_config(){
		$db_conf = array();
		$db_conf['db_server']   = $this->db_server;
		$db_conf['db_name']     = $this->db_name;
		$db_conf['db_user']     = $this->db_user;
		$db_conf['db_password'] = $this->db_password;

		return $db_conf;
	}

	// Gets the content from the configuration file into
	// the public variables of this class
	private function get_configuration_file($config){
		$this->db_server   = $config['db_server'];
		$this->db_name     = $config['db_name'];
		$this->db_user     = $config['db_user'];
		$this->db_password = $config['db_password'];
	}
	
	// Returns a hash with the user's profile
	private function get_user_status($config){
		// connect to the database
		$db = new DBAccess($this);
		if(!$db->connect()){
			error_log('classes/Configuration: Cannot connect to the database');
			return null;
		}
		
		// get different user_status from database (the default ones)
		$query = 'SELECT id, name FROM user_status;';
		$db->prepare($query);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		$db->disconnect();
		if(! isset($res) or $res === FALSE){
			error_log('classes/Configuration: Cannot retrieve user status');
			return $config;
		}

		// go through result
		foreach($res as $user_status){
			switch($user_status['name']){
				case 'Guest':
					$this->guest_user_status_id = $user_status['id'];
					break;
				case 'Member':
					$this->member_user_status_id = $user_status['id'];
					break;
				case 'Admin':
					$this->admin_user_status_id = $user_status['id'];
					break;
			}
		}
		
		return;
	}

	private function get_configuration_properties($config){
		// connect to the database
		$db = new DBAccess($this);
		if(!$db->connect()){
			error_log('classes/Configuration: Cannot connect to the database');
			return null;
		}
		
		// get different user_status from database (the default ones)
		$query = 'SELECT property as name, value FROM configuration;';
		$db->prepare($query);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		$db->disconnect();
		if(! isset($res) or $res === FALSE){
			error_log('classes/Configuration: Cannot retrieve user status');
			return $config;
		}

		// go through result
		foreach($res as $property){
			switch($property['name']){
				case 'browser.session.timeout.max':
					$this->browser_session_timeout_max = (int) $property['value'];
					break;
				case 'location.timezone':
					$this->location_time_zone = $property['value'];
					break;
				case 'location.longitude':
					$this->location_longitude = $property['value'];
					break;
				case 'location.latitude':
					$this->location_latitude = $property['value'];
					break;
				case 'business.day.start':
					$this->business_day_start = $property['value'];
					break;
				case 'business.day.end':
					$this->business_day_end = $property['value'];
					break;
				case 'session.cancel.graceperiod':
					$this->session_cancel_graceperiod = $property['value'];
					break;
				case 'recaptcha.publickey':
					$this->recaptcha_publickey = $property['value'];
					break;
				case 'recaptcha.privatekey':
					$this->recaptcha_privatekey = $property['value'];
					break;
				case 'currency':
					$this->currency = $property['value'];
					break;
				case 'location.map':
					$this->location_map = $property['value'];
					break;
				case 'location.address':
					$this->location_address = $property['value'];
					break;
				case 'payment.account.owner':
					$this->payment_account_owner = $property['value'];
					break;
				case 'payment.account.iban':
					$this->payment_account_iban = $property['value'];
					break;
				case 'payment.account.bic':
					$this->payment_account_bic = $property['value'];
					break;
				case 'payment.account.comment':
					$this->payment_account_comment = $property['value'];
					break;
			}
		}

		// set defaults for values that are not specified yet
		if(is_null($this->browser_session_timeout_max)){
			$this->browser_session_timeout_max = 604800;
		}
		if(is_null($this->location_time_zone)){
			$this->location_time_zone = "Europe/Berlin";
		}
		if(is_null($this->location_longitude)){
			$this->location_longitude = 8.542939;
		}
		if(is_null($this->location_latitude)){
			$this->location_latitude = 47.367658;
		}
		if(is_null($this->business_day_start)){
			$this->business_day_start = "08:00:00";
		}
		if(is_null($this->business_day_end)){
			$this->business_day_end = "21:00:00";
		}
		if(is_null($this->session_cancel_graceperiod)){
			$this->session_cancel_graceperiod = 86400;
		}
		if(is_null($this->currency)){
			$this->currency = "CHF";
		}
	}

// ("browser.session.timeout.default", "10800"),
// DONE ("browser.session.timeout.max",     "604800"),
// DONE ("location.longitude",              "8.542939"),
// DONE ("location.latitude",               "47.367658"),
// ("location.gmt_offset",             "1"),
// DONE ("location.timezone",               "Europe/Berlin"),
// DONE ("business.day.start",              "08:00:00"),
// DONE ("business.day.end",                "21:00:00"),
// ("business.day.startatsunrise",     "false"),
// ("business.day.endatsunset",        "false"),
// DONE ("session.cancel.graceperiod",      "86400");
	

}
?>