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

    $response = null;

    switch($_GET['action']){
        case 'get_db_config':
            $response = get_db_config($configuration);
            break;
        case 'get_recaptcha_key':
            $response = get_recaptcha_key($configuration);
            break;
        case 'get_logo_file':
            $response = get_logo_file($configuration);
            break;
        case 'get_configuration':
            $response = get_configuration($configuration);
            break;
        case 'set_configuration':
            $response = set_configuration($configuration);
            break;
        case 'setup_db_config':
            $response = setup_db_config($configuration);
            break;
        case 'setup_mynautique_config':
            $response = setup_mynautique_config($configuration);
            break;
        case 'make_user_admin':
            $response = make_user_admin($configuration);
            break;
        default:
            HttpHeader::setResponseCode(400);
            $response = Status::errorStatus("Action not supported");
            break;
    }

    echo json_encode($response);
    return;

    // returns the current database configuration
    function get_db_config($configuration){
        $db_config = array();
        $db_config['is_configured'] = false;
        $db_config['db_server']     = null;
        $db_config['db_name']       = null;
        $db_config['db_user']       = null;
        $db_config['db_password']   = null; // we do not reveal the password via API

        if(!_is_db_configured($configuration)){
            return Status::successDataResponse("not configured", $db_config);
        }else{
            // to reveal already configured info
            // the user needs to be admin
            _is_admin_or_return($configuration);

            // return db configuration
            $db_config['is_configured'] = TRUE;
            $db_config['db_server']     = $config['db_server'];
            $db_config['db_name']       = $config['db_name'];
            $db_config['db_user']       = $config['db_user'];
            $db_config['db_password']   = null; // we do not reveal the password via API

            return Status::successDataResponse("success", $db_config);
        }
    }

    function get_recaptcha_key($configuration){
        $response = array();
        $response['key'] = $configuration->recaptcha_publickey;
        return Status::successDataResponse("success", $response);
    }

    function get_logo_file($configuration){
        $response = array();
        $response['uri'] = $configuration->logo_file;
        return Status::successDataResponse("success", $response);
    }

    function get_configuration($configuration){
        // either there is no config yet
        // or the user is admin user
        if(!_is_db_configured($configuration)){
            return Status::errorStatus("no db configuration found");
        }

        // only logged in users can retrieve configuration
        _is_logged_in_or_return($configuration);

        // configuration parameters that are public
        // to all logged in users
        $response = [
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
            "recaptcha_publickey"       => $configuration->recaptcha_publickey,
            "logo_file"                 => $configuration->logo_file,
        ];

        // configuration parameters that are only available to admins
        if(_is_admin($configuration)){
            $response["smtp_sender"]          = $configuration->smtp_sender;
            $response["smtp_server"]          = $configuration->smtp_server;
            $response["smtp_username"]        = $configuration->smtp_username;
            $response["smtp_password"]        = "hidden"; // the proper value is hidden, $configuration->smtp_password
            $response["recaptcha_privatekey"] = $configuration->recaptcha_privatekey;
            $response["engine_hour_format"]   = $configuration->engine_hour_format;
            $response["fuel_payment_type"]    = $configuration->fuel_payment_type;
            $response["mynautique_enabled"]   = $configuration->mynautique_enabled;
            $response["mynautique_user"]      = $configuration->mynautique_user;
            $response["mynautique_password"]  = "hidden"; // the proper value is hidden, $configuration->mynautique_password
            $response["mynautique_boat_id"]   = $configuration->mynautique_boat_id;
            $response["mynautique_fuel_capacity"] = $configuration->mynautique_fuel_capacity;
        }
        
        return Status::successDataResponse('success', $response);
    }

    function set_configuration($configuration){
        $data = json_decode(file_get_contents('php://input'));

        // default value replacement
        if(isset($data->smtp_sender) && $data->smtp_sender == ""){
            $data->smtp_sender = NULL;
        }
        if(isset($data->smtp_server) && $data->smtp_server == ""){
            $data->smtp_server = NULL;
        }
        if(isset($data->recaptcha_privatekey) && $data->recaptcha_privatekey == ""){
            $data->recaptcha_privatekey = NULL;
        }
        if(isset($data->recaptcha_publickey) && $data->recaptcha_publickey == ""){
            $data->recaptcha_publickey = NULL;
        }

        // sanitize input
        $sanitizer = new Sanitizer();
        if(!isset($data->location_time_zone) || !$sanitizer->isTimeZone($data->location_time_zone)){
            return Status::errorStatus("The provided timezone is missing or not in a valid format (valid e.g., Europe/Berlin)");
        }
        if(!isset($data->location_longitude) || !$sanitizer->isLongitude($data->location_longitude)){
            return Status::errorStatus("The provided longitude is missing or not in a valid format.");
        }
        if(!isset($data->location_latitude) || !$sanitizer->isLatitude($data->location_latitude)){
            return Status::errorStatus("The provided latitude is missing or not in a valid format.");
        }
        if(isset($data->location_map) && !$sanitizer->isGoogleMapsURL($data->location_map)){
            return Status::errorStatus("The provided google maps embedded URL is not in a valid format.");
        }        
        if(!isset($data->currency) || !$sanitizer->isCurrency($data->currency)){
            return Status::errorStatus("The provided currency is not in a valid format.");
        }
        if(isset($data->payment_account_iban) && !$sanitizer->isIBAN($data->payment_account_iban)){
            return Status::errorStatus("The provided IBAN is not in a valid format.");
        }
        if(isset($data->payment_account_bic) && !$sanitizer->isBIC($data->payment_account_bic)){
            return Status::errorStatus("The provided BIC/SWIFT code is not in a valid format.");
        }
        if(isset($data->smtp_sender) && !$sanitizer->isEmail($data->smtp_sender)){
            return Status::errorStatus("The provided sender email address is not in a valid format.");
        }
        if(isset($data->smtp_server) && !$sanitizer->isServerAddress($data->smtp_server)){
            return Status::errorStatus("The provided SMTP server address is not in a valid format. (e.g., smtp.gmail.com:587");
        }
        if(isset($data->recaptcha_privatekey) && !$sanitizer->isRecaptcha($data->recaptcha_privatekey)){
            return Status::errorStatus("The provided recaptcha private key is not in a valid format.");
        }
        if(isset($data->recaptcha_publickey) && !$sanitizer->isRecaptcha($data->recaptcha_publickey)){
            return Status::errorStatus("The provided recaptcha public key is not in a valid format.");
        }
        if(isset($data->engine_hour_format) && !$sanitizer->isEngineHourFormat($data->engine_hour_format)){
            return Status::errorStatus("The provided engine hour format is not valid.");
        }
        if(isset($data->fuel_payment_type) && !$sanitizer->isFuelPaymentType($data->fuel_payment_type)){
            return Status::errorStatus("The provided fuel payment type value is not valid.");
        }
        // no sanitization for location address
        //                     payment account owner
        //                     payment account comment
        //                     smtp_username
        //                     smtp_password
        // parameters will be stored as text and will be escaped where used

        // update database
        $keys = array(
            'currency',
            'engine_hour_format',
            'fuel_payment_type',
            'location_address',
            'location_latitude',
            'location_longitude',
            'location_map',
            'location_time_zone',
            'logo_file',
            'payment_account_bic',
            'payment_account_comment',
            'payment_account_iban',
            'payment_account_owner',
            'recaptcha_privatekey',
            'recaptcha_publickey',
            'smtp_password',
            'smtp_sender',
            'smtp_server',
            'smtp_username',
            'mynautique_boat_id',
            'mynautique_fuel_capacity',
        );
        foreach ($keys as $key) {
            // if the key is not defined, we do not store anything
            if(!isset($data->$key)){
                // do not delete password, it is only submitted
                // upon change
                if($key != 'smtp_password'){
                    $configuration->set_configuration_property($key, NULL);
                    continue;
                }else{
                    continue;
                }
            }

            if(! $configuration->set_configuration_property($key, $data->$key)){
                error_log("Cannot update property $key with value " . $data->$key);
                return Status::errorStatus("Cannot save property $key with value " . $data->$key);
            }
        }

        // change non database configuration
        error_log("set mynautique config");
        $status = _set_mynautique_config($configuration, $data);
        if(Status::isError($status)){
            return $status;
        }

        return Status::successStatus("successfully updated");
    }

    function setup_db_config($configuration){
        // either there is no config yet
        // or the user is admin user
        if(_is_db_configured($configuration)){
            _is_admin_or_return($configuration);
        }

        // get and validate user input
        $data = json_decode(file_get_contents('php://input'));
        
        $host = parse_url($data->db_server, PHP_URL_HOST);
        if(!isset($host)){
            return Status::errorStatus('no proper hostname specified');
        }
        $port = parse_url($data->db_server, PHP_URL_PORT);
        if(!isset($port)){
            return Status::errorStatus('no port specified');
        }
        if(!preg_match('/^[a-zA-Z0-9-_\.\:\]\[]+:[0-9]+$/', $data->db_server)){
            return Status::errorStatus('only [ip|hostname]:[port] format supported for database host');
        }
        if(!isset($data->db_name) || !preg_match('/^[a-zA-Z0-9-_]+$/', $data->db_name)){
            return Status::errorStatus('invalid database name');
        }
        if(!isset($data->db_user) || !preg_match('/^[a-zA-Z0-9-_]+$/', $data->db_user)){
            return Status::errorStatus('invalid database user');
        }
        if(!isset($data->db_password) || preg_match('/\s+/', $data->db_password)){
            return Status::errorStatus('invalid database password (containing spaces)');
        }

        // verify that connection works
        $new_config = new Configuration();
        $new_config->db_server   = $data->db_server;
        $new_config->db_name     = $data->db_name;
        $new_config->db_user     = $data->db_user;
        $new_config->db_password = $data->db_password;
        $db = new DBAccess($new_config);
        if(!$db->connect()){
            return Status::errorStatus('DB access not working with the provided settings. Please check for typos, make sure that the database is reachable and that your username and password are correct.');
        }

        // check whether tables already exist
        $query = 'SHOW TABLES LIKE "configuration"';
        $res = $db->fetch_data_hash($query, 0);
        if(!$res){
            // if the table does not exist, we setup the database
            // setup a new database
            if(!file_exists("../config/db/schema.sql") || !$db->apply_sql_file("../config/db/schema.sql")){
                return Status::errorStatus("DB scheme could not be applied to database, please check server error logs");
            }
        }

        // write configuration file (only after schema is applied)// build the configuration for the configuration file
        $configuration->config_file_variables['db_server'] = $data->db_server;
        $configuration->config_file_variables['db_name'] = $data->db_name;
        $configuration->config_file_variables['db_user'] = $data->db_user;
        $configuration->config_file_variables['db_password'] = $data->db_password;
        if(_write_config_file($configuration->config_file_variables)){
            return Status::successStatus("configuration applied and database setup");
        }
        return Status::errorStatus("db could not be configured");
    }

    function setup_mynautique_config($configuration){
        // basic db setup needs to be in place and
        // we need to be admin to configure myNautique
        if(_is_mynautique_configured($configuration)){
            _is_admin_or_return($configuration);
        }

        // get and validate user input
        $data = json_decode(file_get_contents('php://input'));
        return _set_mynautique_config($configuration, $data);
    }

    function _set_mynautique_config($configuration, $values){
        if(!isset($values->mynautique_enabled) || !($values->mynautique_enabled == FALSE || $values->mynautique_enabled == TRUE)){
            return Status::errorStatus('mynautique_enabled must be TRUE/FALSE');
        }
        if(!isset($values->mynautique_user) || !preg_match('/^[a-zA-Z0-9-_@\.]*$/', $values->mynautique_user )){
            return Status::errorStatus('mynautique_user must be a valid username');
        }

        // handle new account and password change
        $set_password = TRUE;
        if(_is_mynautique_configured($configuration) && !isset($values->mynautique_password)){
            // if mynautique is already configured and we do not change the password
            $set_password = FALSE;
        }
        else if(!isset($values->mynautique_password) || !preg_match('/^[a-zA-Z0-9-_@]*$/', $values->mynautique_password )){
            return Status::errorStatus('mynautique_password must be a valid password');
        }

        // update configuration file
        $configuration->config_file_variables['mynautique_enabled'] = $values->mynautique_enabled == TRUE ? TRUE : FALSE;
        $configuration->config_file_variables['mynautique_user'] = $values->mynautique_user;
        if($set_password == TRUE){
            $configuration->config_file_variables['mynautique_password'] = $values->mynautique_password;
        }else{
            $configuration->config_file_variables['mynautique_password'] = $configuration->mynautique_password;
        }

        if(_write_config_file($configuration->config_file_variables)){
            return Status::successStatus("mynautique configured");
        }
        return Status::errorStatus("myNautique could not be configured, an error occurred");
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

        // if there is already one admin user
        // we do not allow this operation
        if(_is_admin_user_configured($configuration)){
            return Status::errorStatus("one admin user already exists, this operation is not permitted");
        }

        // input validation
        $sanitizer = new Sanitizer();
        if(! isset($data->user_id) or !$sanitizer->isInt($data->user_id)){
			return Status::errorStatus("No valid user_id specified.");
        }

        // unlock user and make admin
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            return Status::errorStatus("Cannot connect to database.");
        }
        $query = 'UPDATE user 
            SET status = ?, 
            locked = ? 
            WHERE id = ?;';
        $db->prepare($query);
        $db->bind_param('iii',
            $configuration->default_admin_user_status_id,
            0,
            $data->user_id
        );
        if(!$db->execute()){
            $db->disconnect();
            return Status::errorStatus("Cannot make user admin");
        }

        return Status::successStatus("user promoted to admin user");
    }

    function _is_admin_user_configured($configuration){
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            throw new Exception('api/configuration/_is_admin_user_configured: cannot connect to database');
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
    
    function _is_db_configured($configuration){
        // get the configuration if it exists
        if(is_null($configuration)){
            return FALSE;
        }

        return $configuration->is_db_configured();
    }

    function _is_mynautique_configured($configuration){
        if(is_null($configuration)){
            return FALSE;
        }
        return $configuration->is_mynautique_configured();
    }

    // returns http status permission denied in case the user is not an admin
    function _is_admin_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn($configuration->admin_user_role_id)){
            echo json_encode(Status::errorStatus("not sufficient permissions"));
            exit;
        }
    }

    function _is_admin($configuration){
        $lc = new Login($configuration);
        return $lc->isAdmin();
    }

    function _is_logged_in_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn()){
            echo json_encode(Status::errorStatus("not logged in"));
            exit;
        }
    }

    function _write_config_file($values){
        // set some defaults in cause $values are not defined
        $db_server = $values['db_server'] ?: "";
        $db_name   = $values['db_name']   ?: "";
        $db_user   = $values['db_user']   ?: "";
        $db_password = $values['db_password'] ?: "";
        $mynautique_enabled = (isset($values['mynautique_enabled']) && $values['mynautique_enabled'] == TRUE) ? "TRUE" : "FALSE";
        $mynautique_user = isset($values['mynautique_user']) && $mynautique_enabled == "TRUE" ? $values['mynautique_user'] : "";
        $mynautique_password = isset($values['mynautique_password']) && $mynautique_enabled == "TRUE" ? $values['mynautique_password'] : "";

        // write configuration file (only after schema is applied)// build the configuration for the configuration file
        $config_string =  '<?php'."\n";
        $config_string .= '// Database Configuration'."\n";
        $config_string .= '//==================================================='."\n";
        $config_string .= '// database server'."\n";
        $config_string .= '$config[\'db_server\']   = "'. $db_server .'";'."\n";
        $config_string .= '// database'."\n";
        $config_string .= '$config[\'db_name\']     = "'. $db_name .'";'."\n";
        $config_string .= '// database user'."\n";
        $config_string .= '$config[\'db_user\']     = "'. $db_user .'";'."\n";
        $config_string .= '// database user password'."\n";
        $config_string .= '$config[\'db_password\'] = "'. $db_password .'";'."\n";
        $config_string .= "\n";

        if ( isset($values['mynautique_enabled'])) {
            $config_string .= '// MyNautique Configuration'."\n";
            $config_string .= '//==================================================='."\n";
            $config_string .= '// myNautique enabled'."\n";
            $config_string .= '$config[\'mynautique_enabled\'] = '.$mynautique_enabled .';'."\n";
            $config_string .= '// myNautique user'."\n";
            $config_string .= '$config[\'mynautique_user\']     = "'. $mynautique_user .'";'."\n";
            $config_string .= '// myNautique password'."\n";
            $config_string .= '$config[\'mynautique_password\']     = "'. $mynautique_password .'";'."\n";
        }

        $config_string .=  '?>'."\n";

        $bytes_written = file_put_contents ("../config/config.php", $config_string);

        if($bytes_written == FALSE) {
            return $bytes_written;
        }
        return TRUE;
    }



?>