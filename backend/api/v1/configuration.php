<?php
    // automatically load all classes
    spl_autoload_register('configuration_autoloader');
    function configuration_autoloader($class){
        if(file_exists(__DIR__.'/../../classes/'.$class.'.php')){
            include __DIR__.'/../../classes/'.$class.'.php';
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
        case 'setup_mynautique_config':
            $response = setup_mynautique_config($configuration);
            break;
        default:
            HttpHeader::setResponseCode(400);
            $response = Status::errorStatus("Action not supported");
            break;
    }

    echo json_encode($response);
    return;

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
        if($values->mynautique_enabled == TRUE){
            if(!isset($values->mynautique_user) || !preg_match('/^[a-zA-Z0-9-_@\.]*$/', $values->mynautique_user )){
                return Status::errorStatus('mynautique_user must be a valid username');
            }
        }

        // handle new account and password change
        $set_password = TRUE;
        if(_is_mynautique_configured($configuration) && !isset($values->mynautique_password)){
            // if mynautique is already configured and we do not change the password
            $set_password = FALSE;
        }
        else if($values->mynautique_enabled == TRUE){
            if(!isset($values->mynautique_password) || !preg_match('/^[a-zA-Z0-9-_@]*$/', $values->mynautique_password )){
                return Status::errorStatus('mynautique_password must be a valid password');
            }
        }

        // update configuration file
        $configuration->config_file_variables['mynautique_enabled'] = $values->mynautique_enabled == TRUE ? TRUE : FALSE;
        $configuration->config_file_variables['mynautique_user'] = $values->mynautique_enabled == TRUE ? $values->mynautique_user : "";
        if($set_password == TRUE){
            $configuration->config_file_variables['mynautique_password'] = $values->mynautique_enabled == TRUE ? $values->mynautique_password : "";
        }else{
            $configuration->config_file_variables['mynautique_password'] = $values->mynautique_enabled == TRUE ? $configuration->mynautique_password : "";
        }

        if(_write_config_file($configuration->config_file_variables)){
            return Status::successStatus("mynautique configured");
        }
        return Status::errorStatus("myNautique could not be configured, an error occurred");
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

        $bytes_written = file_put_contents ("../../config/config.php", $config_string);

        if($bytes_written == FALSE) {
            return $bytes_written;
        }
        return TRUE;
    }



?>