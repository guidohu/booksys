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

	// default user status / the ones that are shipped with the installation by default
	// Note: they have to match the IDs in the default
	//       schema definition
	public $default_guest_user_status_id = 1;
	public $default_member_user_status_id = 2;
	public $default_admin_user_status_id = 3;

	// user roles
	public $guest_user_role_id;
	public $member_user_role_id;
	public $admin_user_role_id;

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
	public $smtp_sender;
	public $smtp_server;
	public $smtp_username;
	public $smtp_password;

	// customization
	public $currency;
	public $location_map;
	public $location_address;
	public $payment_account_owner;
	public $payment_account_iban;
	public $payment_account_bic;
	public $payment_account_comment;
	public $logo_file;
	public $engine_hour_format;
	public $fuel_payment_type;


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

		// get all the user roles from the database that
		// are required for the code to run
		$this->get_user_roles($this->config_file_variables);

		// get all the configuration settings from the
		// database
		$this->get_configuration_properties($this->config_file_variables);

		// set hard defaults
		$this->login_page          = '/index.html';
		$this->version             = "2.3.3";
		$this->required_db_version = "1.16";
	}

	public function is_db_configured(){
		// check if configuration is set
        if(!isset($this->db_server) or $this->db_server == ''){
            return FALSE;
        }
        // check if db_name is set
        if(!isset($this->db_name) or $this->db_name == ''){
            return FALSE;
        }
        // check if db_user is set
        if(!isset($this->db_user) or $this->db_user == ''){
            return FALSE;
        }

        return TRUE;
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
	private function get_user_roles($config){
		// connect to the database
		$db = new DBAccess($this);
		if(!$db->connect()){
			error_log('classes/Configuration: Cannot connect to the database');
			return null;
		}
		
		// get different user_status from database (the default ones)
		$query = 'SELECT id, name FROM user_role;';
		$db->prepare($query);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		$db->disconnect();
		if(! isset($res) or $res === FALSE){
			error_log('classes/Configuration: Cannot retrieve user roles');
			return $config;
		}

		// go through result and handle all supported user roles
		foreach($res as $user_role){
			switch($user_role['name']){
				case 'guest':
					$this->guest_user_role_id = $user_role['id'];
					break;
				case 'member':
					$this->member_user_role_id = $user_role['id'];
					break;
				case 'admin':
					$this->admin_user_role_id = $user_role['id'];
					break;
				default:
				    error_log("Unsupported user role: " . $user_role['name'] . " make sure this one is handled properly in the code");
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
				case 'browser.session.timeout.max': // TODO
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
				case 'business.day.start': // TODO
					$this->business_day_start = $property['value'];
					break;
				case 'business.day.end': // TODO
					$this->business_day_end = $property['value'];
					break;
				case 'session.cancel.graceperiod': // TODO
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
				case 'smtp.sender':
					$this->smtp_sender = $property['value'];
					break;
				case 'smtp.server': 
					$this->smtp_server = $property['value'];
					break;
				case 'smtp.username':
					$this->smtp_username = $property['value'];
					break;
				case 'smtp.password':
					$this->smtp_password = $property['value'];
					break;
				case 'logo.file':
					$this->logo_file = $property['value'];
					break;
				case 'engine.hour.format':
					$this->engine_hour_format = $property['value'];
					break;
				case 'fuel.payment.type':
					$this->fuel_payment_type = $property['value'];
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

	public function set_configuration_property($key, $value){
		$db = new DBAccess($this);
        if(!$db->connect()){
			error_log("Cannot connect to database");
            return FALSE;
        }
		$db->prepare('REPLACE INTO configuration (property, value) VALUES (?, ?);');
		
		$db_key = '';
		switch ($key) {
			case 'engine_hour_format':
				$db_key = "engine.hour.format";
				break;
			case 'fuel_payment_type':
				$db_key = "fuel.payment.type";
				break;
			case 'logo_file':
				$db_key = "logo.file";
				break;
			case 'location_time_zone':
				$db_key = "location.timezone";
				break;
			case 'location_latitude':
				$db_key = "location.latitude";
				break;
			case 'location_longitude':
				$db_key = "location.longitude";
				break;
			case 'location_map':
				$db_key = "location.map";
				break;
			case 'location_address':
				$db_key = "location.address";
				break;
			case 'currency':
				$db_key = "currency";
				break;
			case 'payment_account_owner':
				$db_key = "payment.account.owner";
				break;
			case 'payment_account_iban':
				$db_key = "payment.account.iban";
				break;
			case 'payment_account_bic':
				$db_key = "payment.account.bic";
				break;
			case 'payment_account_comment':
				$db_key = "payment.account.comment";
				break;
			case 'smtp_sender':
				$db_key = "smtp.sender";
				break;
			case 'smtp_server':
				$db_key = "smtp.server";
				break;
			case 'smtp_username':
				$db_key = "smtp.username";
				break;
			case 'smtp_password':
				$db_key = "smtp.password";
				break;
			case 'recaptcha_privatekey':
				$db_key = "recaptcha.privatekey";
				break;
			case 'recaptcha_publickey':
				$db_key = "recaptcha.publickey";
				break;
			default:
				return FALSE;
		}

		$db->bind_param('ss',
			$db_key,
			$value
		);
		if(! $db->execute()){
			return FALSE;
		}
		return TRUE;
	}
}
?>