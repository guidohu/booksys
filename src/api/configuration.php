<?php
    // automatically load all classes
    spl_autoload_register('configuration_autoloader');
    function configuration_autoloader($class){
        if(file_exists(__DIR__.'/../classes/'.$class.'.php')){
            include __DIR__.'/../classes/'.$class.'.php';
        }
    }

    // test whether the config-file exists
    $configuration = new Configuration();

    // check if we have an action
    if(!isset($_GET['action'])){
        HttpHeader::setResponseCode(200);
        exit;
    }

    switch($_GET['action']){
        case 'get_db_config':
            get_db_config($configuration);
            exit;
        case 'get_recaptcha_key':
            get_recaptcha_key($configuration);
            exit;
        case 'get_customization_parameters':
            get_customization_parameters($configuration);
            exit;
        case 'get_configuration':
            $response = get_configuration($configuration);
            echo json_encode($response);
            exit;
        case 'set_configuration':
            $response = set_configuration($configuration);
            echo json_encode($response);
            exit;
        case 'setup_db_config':
            setup_db_config($configuration);
            exit;
        case 'is_admin_user_configured':
            is_admin_user_configured($configuration);
            exit;
        case 'make_user_admin':
            make_user_admin($configuration);
            exit;
    }

    HttpHeader::setResponseCode(400);
    $status = array();
    $status['ok'] = FALSE;
    $status['message'] = 'invalid action requested';
    echo json_encode($status);
    return;

    // returns the current database configuration
    function get_db_config($configuration){
        if(!_is_db_configured($configuration)){
            $db_config = array();
            $db_config['is_configured'] = FALSE;
            echo json_encode($db_config);
            return;
        }else{
            // to reveal already configured info
            // the user needs to be admin
            _is_admin_or_return($configuration);

            // return db configuration
            $db_config = array();
            $db_config['is_configured'] = TRUE;
            $db_config['db_server']     = $config['db_server'];
            $db_config['db_name']       = $config['db_name'];
            $db_config['db_user']       = $config['db_user'];
            $db_config['db_password']   = null; // we do not reveal the password via API

            echo json_encode($db_config);
            return;
        }
    }

    function get_recaptcha_key($configuration){
        $response = array();
        $response['key'] = $configuration->recaptcha_publickey;
        echo json_encode($response);
        return;
    }

    function get_customization_parameters($configuration){
        // only available to users that are logged in
        _is_logged_in_or_return($configuration);

        // return customization parameters
        $response = array();
        $response['currency']                = $configuration->currency;
        $response['location_map_iframe']     = $configuration->location_map;
        $response['location_address']        = $configuration->location_address;
        $response['payment_account_owner']   = $configuration->payment_account_owner;
        $response['payment_account_iban']    = $configuration->payment_account_iban;
        $response['payment_account_bic']     = $configuration->payment_account_bic;
        $response['payment_account_comment'] = $configuration->payment_account_comment;
        echo json_encode($response);
        return;
    }

    function get_configuration($configuration){
        $response = [
            "ok" => TRUE,
        ];
        // either there is no config yet
        // or the user is admin user
        if(!_is_db_configured($configuration)){
            $response["ok"] = FALSE;
            $response["msg"] = "no db configuration found";
            return $response;
        }
        
        // only admins can retrieve configuration data
        // returns with a HTTP 403
        _is_admin_or_return($configuration);

        $response = [
            "ok"                        => TRUE,
            "location_time_zone"        => $configuration->location_time_zone,
            "location_longitude"        => $configuration->location_longitude,
            "location_latitude"         => $configuration->location_latitude,
            "location_map"              => $configuration->location_map,
            "location_address"          => $configuration->location_address,
            "currency"                  => $configuration->currency,
            "payment_account_owner"     => $configuration->payment_account_owner,
            "payment_account_iban"      => $configuration->payment_account_iban,
            "payment_account_bic"       => $configuration->payment_account_bic,
            "payment_account_comment"   => $configuration->payment_account_comment,
            "smtp_sender"               => $configuration->smtp_sender,
            "smtp_server"               => $configuration->smtp_server,
            "smtp_username"             => $configuration->smtp_username,
            "smtp_password"             => "hidden", // $configuration->smtp_password,
            "recaptcha_publickey"       => $configuration->recaptcha_publickey,
            "recaptcha_privatekey"      => $configuration->recaptcha_privatekey
        ];

        return $response;
    }

    function set_configuration($configuration){
        $response = [
            "ok" => TRUE,
            "msg" => "configuration updated",
            "errors" => [],
        ];

        $data = json_decode(file_get_contents('php://input'));

        // sanitize input
        $sanitizer = new Sanitizer();
        if(!isset($data->location_time_zone) || !$sanitizer->isTimeZone($data->location_time_zone)){
            return _get_status(FALSE, "The provided timezone is missing or not in a valid format (valid e.g., Europe/Berlin)");
        }
        if(!isset($data->location_longitude) || !$sanitizer->isLongitude($data->location_longitude)){
            return _get_status(FALSE, "The provided longitude is missing or not in a valid format.");
        }
        if(!isset($data->location_latitude) || !$sanitizer->isLatitude($data->location_latitude)){
            return _get_status(FALSE, "The provided latitude is missing or not in a valid format.");
        }
        if(isset($data->location_map) && !$sanitizer->isGoogleMapsURL($data->location_map)){
            return _get_status(FALSE, "The provided google maps embedded URL is not in a valid format.");
        }        
        if(!isset($data->currency) || !$sanitizer->isCurrency($data->currency)){
            return _get_status(FALSE, "The provided currency is not in a valid format.");
        }
        if(isset($data->payment_account_iban) && !$sanitizer->isIBAN($data->payment_account_iban)){
            return _get_status(FALSE, "The provided IBAN is not in a valid format.");
        }
        if(isset($data->payment_account_bic) && !$sanitizer->isBIC($data->payment_account_bic)){
            return _get_status(FALSE, "The provided BIC/SWIFT code is not in a valid format.");
        }
        if(isset($data->smtp_sender) && !$sanitizer->isEmail($data->smtp_sender)){
            return _get_status(FALSE, "The provided sender email address is not in a valid format.");
        }
        if(isset($data->smtp_server) && !$sanitizer->isServerAddress($data->smtp_server)){
            return _get_status(FALSE, "The provided SMTP server address is not in a valid format. (e.g., smtp.gmail.com:587");
        }
        if(isset($data->recaptcha_privatekey) && !$sanitizer->isAlphaNum($data->recaptcha_privatekey)){
            return _get_status(FALSE, "The provided recaptcha private key is not in a valid format.");
        }
        if(isset($data->recaptcha_publickey) && !$sanitizer->isAlphaNum($data->recaptcha_publickey)){
            return _get_status(FALSE, "The provided recaptcha public key is not in a valid format.");
        }
        // no sanitization for location address
        //                     payment account owner
        //                     payment account comment
        //                     smtp_username
        //                     smtp_password
        // parameters will be stored as text and will be escaped where used

        // update database
        $keys = array(
            'location_time_zone',
            'location_longitude',
            'location_latitude',
            'location_map',
            'location_address',
            'currency',
            'payment_account_owner',
            'payment_account_iban',
            'payment_account_bic',
            'payment_account_comment',
            'smtp_sender',
            'smtp_server',
            'smtp_username',
            'smtp_password',
            'recaptcha_privatekey',
            'recaptcha_publickey'
        );
        foreach ($keys as $key) {
            // if the key is not defined, we do not store anything
            if(!isset($data->$key)){
                continue;
            }

            if(! $configuration->set_configuration_property($key, $data->$key)){
                error_log("Cannot update property $key with value " . $data->$key);
                $response["ok"] = FALSE;
                $response["msg"] = "Errors with some of the properties";
                array_push($response["errors"], "Cannot update property $key with value " . $data->$key);
            }
        }

        return $response;
    }

    function setup_db_config($configuration){
        // either there is no config yet
        // or the user is admin user
        if(_is_db_configured($configuration)){
            _is_admin_or_return($configuration);
        }

        // get and validate user input
        $data = json_decode(file_get_contents('php://input'));
        $status = array();
        $status['OK'] = TRUE;
        $status['msg'] = "configuration applied and database setup";
        
        $host = parse_url($data->db_server, PHP_URL_HOST);
        if(!isset($host)){
            $status['OK'] = FALSE;
            $status['msg'] = 'no proper hostname specified';
            echo json_encode($status);
            return;
        }
        $port = parse_url($data->db_server, PHP_URL_PORT);
        if(!isset($port)){
            $status['OK'] = FALSE;
            $status['msg'] = 'no port specified';
            echo json_encode($status);
            return;
        }
        if(!preg_match('/^[a-zA-Z0-9-_]+:[0-9]+$/', $data->db_server)){
            $status['OK'] = FALSE;
            $status['msg'] = 'only [ip|hostname]:[port] format supported for database host';
            echo json_encode($status);
            return;
        }
        if(!isset($data->db_name) || !preg_match('/^[a-zA-Z0-9-_]+$/', $data->db_name)){
            $status['OK'] = FALSE;
            $status['msg'] = 'invalid database name';
            echo json_encode($status);
            return;
        }
        if(!isset($data->db_user) || !preg_match('/^[a-zA-Z0-9-_]+$/', $data->db_user)){
            $status['OK'] = FALSE;
            $status['msg'] = 'invalid database user';
            echo json_encode($status);
            return;
        }
        if(!isset($data->db_password) || preg_match('/\s+/', $data->db_password)){
            $status['OK'] = FALSE;
            $status['msg'] = 'invalid database password (containing spaces)';
            echo json_encode($status);
            return;
        }

        // verify that connection works
        $new_config = new Configuration();
        $new_config->db_server   = $data->db_server;
        $new_config->db_name     = $data->db_name;
        $new_config->db_user     = $data->db_user;
        $new_config->db_password = $data->db_password;
        $db = new DBAccess($new_config);
        if(!$db->connect()){
            $status['OK'] = FALSE;
            $status['msg'] = 'DB access not working with the provided settings. Please check for typos, make sure that the database is reachable and that your username and password are correct.';
            echo json_encode($status);
            return;
        }

        // setup a new database
        if(!file_exists("../config/db/schema.sql") || !$db->apply_sql_file("../config/db/schema.sql")){
            $status['OK']  = FALSE;
            $status['msg'] = "DB scheme could not be applied to database, please check error logs";
            echo json_encode($status);
            return;
        }

        // write configuration file (only after schema is applied)// build the configuration for the configuration file
        $config_string =  '<?php'."\n";
        $config_string .= '// Database Configuration'."\n";
        $config_string .= '//==================================================='."\n";
        $config_string .= '// database server'."\n";
        $config_string .= '$config[\'db_server\']   = "'. $data->db_server .'";'."\n";
        $config_string .= '// database'."\n";
        $config_string .= '$config[\'db_name\']     = "'. $data->db_name .'";'."\n";
        $config_string .= '// database user'."\n";
        $config_string .= '$config[\'db_user\']     = "'. $data->db_user .'";'."\n";
        $config_string .= '// database user password'."\n";
        $config_string .= '$config[\'db_password\'] = "'. $data->db_password .'";'."\n";
        $config_string .=  '?>'."\n";
        file_put_contents ("../config/config.php", $config_string);

        // return status
        echo json_encode($status);
        return;
    }

    // Tests whether there is at least one admin present in the database
    // (if an admin exists the user setup has been done)
    function is_admin_user_configured($configuration){
        $res = array();
        try{
            $res['OK'] = TRUE;
            $res['admin_exists'] = _is_admin_user_configured($configuration);
        }
        catch(Exception $e){
            $res['OK'] = FALSE;
            error_log($e);
            $res['msg'] = "Error querying the database to check whether an admin already exists";
        }
        
        // return result
        echo json_encode($res);
        return;
    }

    /**
     * make_user_admin allows a user to become admin during the
     * initial setup of the application. As soon as an admin user
     * exits this function will not perform any action
     *
     * @param $configuration    configurtaion object
     * @return response of the form
     * {
     *      ok      =>  TRUE|FALSE
     *      message =>  "description message"
     * }
     */
    function make_user_admin($configuration){
        $data = json_decode(file_get_contents('php://input'));

        $status = array();
        $status['ok'] = TRUE;

        // if there is already one admin user
        // we do not allow this operation
        if(_is_admin_user_configured($configuration)){
            _print_status(FALSE, "one admin user already exists, this operation is not permitted");
            return;
        }

        // input validation
        $sanitizer = new Sanitizer();
        if(! isset($data->user_id) or !$sanitizer->isInt($data->user_id)){
			_print_status(FALSE, "No valid user_id specified.");
			return;
        }

        // unlock user and make admin
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            _print_status(FALSE, "Cannot connect to database.");
			return;
        }
        $query = 'UPDATE user 
            SET status = ?, 
            locked = ? 
            WHERE id = ?;';
        $db->prepare($query);
        $db->bind_param('iii',
            $configuration->admin_user_status_id,
            0,
            $data->user_id
        );
        if(!$db->execute()){
            $db->disconnect();
            _print_status(FALSE, "Cannot make user admin");
            return;
        }

        echo json_encode($status);
        return;
    }

    function _is_admin_user_configured($configuration){
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            throw new Exception('api/configuration/is_admin_user_configured: cannot connect to database');
        }

        $query = 'SELECT u.id 
            FROM user u 
            JOIN user_status us ON u.status = us.id 
            JOIN user_role ur ON us.user_role_id = ur.id 
            WHERE ur.name LIKE \'admin\';';
        $db->prepare($query);
        $db->execute();
        $res = $db->fetch_stmt_hash();
        $db->disconnect();
        if(isset($res) && count($res) > 0){
            return TRUE;
        }
        return FALSE;
    }

    // helper function to print the status
	function _print_status($ok, $msg){
		$status = array();
		$status['ok'] = $ok;
		$status['msg'] = $msg;
		echo json_encode($status);
		return;
    }
    
    function _get_status($ok, $msg){
		$status = array();
		$status['ok'] = $ok;
		$status['msg'] = $msg;
        return $status;
    }
    
    function _is_db_configured($configuration){
        // get the configuration if it exists
        if(is_null($configuration)){
            return FALSE;
        }

        $db_config = $configuration->get_db_config();

        // check if db_server is set
        if(!isset($db_config['db_server']) or $db_config['db_server'] == ''){
            return FALSE;
        }
        // check if db_name is set
        if(!isset($db_config['db_name']) or $db_config['db_name'] == ''){
            return FALSE;
        }
        // check if db_user is set
        if(!isset($db_config['db_user']) or $db_config['db_user'] == ''){
            return FALSE;
        }

        return TRUE;
    }

    // returns http status permission denied in case the user is not an admin
    function _is_admin_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn($configuration->admin_user_status_id)){
            HttpHeader::setResponseCode(403);
            exit;
        }
    }

    function _is_logged_in_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn()){
            HttpHeader::setResponseCode(403);
            exit;
        }
    }



?>