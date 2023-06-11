<?php
    /**
     * The API allows to check basic backend status 
     * and to perform actions on the backend such as 
     * a database
     * update.
     * 
     */

    // automatically load all classes that are needed
    spl_autoload_register('backend_autoloader');
    function backend_autoloader($class){
        if(file_exists(__DIR__.'/../../classes/'.$class.'.php')){
            include __DIR__.'/../../classes/'.$class.'.php';
        }
    }

    $statusResponse = array(
        "configFile"  => FALSE,
        "configDb"    => FALSE,
        "dbReachable" => FALSE,
        "adminExists" => FALSE
    );

    // check if we have an action
    if(!isset($_GET['action'])){
        echo json_encode(Status::errorStatus('no API action selected'));
        exit;
    }

    // test whether the config-file exists
    if (!file_exists('../../config/config.php') || !is_readable('../../config/config.php')) {
        error_log("No config file yet");
        echo json_encode(Status::successDataResponse("no config", $statusResponse));
        exit;
    }else{
        $statusResponse['configFile'] = TRUE;
    }

    // get configuration
    $configuration = new Configuration();
    $status = "";

    switch($_GET['action']){
        case 'get_status':
            $status = get_status($configuration);
            break;
        default:
            HttpHeader::setResponseCode(400);
            $status = Status::errorStatus('invalid action requested');
            break;
    }

    echo json_encode($status);
    return;

    /**
     * Returns the status of the backend. Which can be one of the following states:
     * - "all backend components are good"
     * 
     * Error cases are:
     * - "no config"
     * - "no database configured"
     * - "database connection is not working"
     * - "no user configured"
     * - "no mynautique configured"
     */
    function get_status($configuration){
        
        $statusResponse = array(
            "configFile"  => TRUE,
            "configDb"    => FALSE,
            "dbReachable" => FALSE,
            "adminExists" => FALSE,
            "myNautiqueConfigured" => FALSE,
        );

        // check if db is configured
        if(! _is_db_configured($configuration)){
            error_log('api/backend: No database configured');
            return Status::successDataResponse('no database configured', $statusResponse);
        }else{
            $statusResponse['configDb'] = TRUE;
        }
        
        // check if db access works
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('api/backend: Cannot connect to the database');
            return Status::successDataResponse('database connection is not working', $statusResponse);
        }else{
            $statusResponse['dbReachable'] = TRUE;
        }

        $res = $db->fetch_data_hash("SELECT count(*) AS users FROM user;");
        $db->disconnect();

        if($res == NULL || $res[0] == NULL || $res[0]['users'] == NULL || $res[0]['users'] < 1){
            error_log('api/backend: No user configured');
            return Status::successDataResponse('no user configured', $statusResponse);
        }else{
            $statusResponse['adminExists'] = TRUE;
        }

        // check if myNautique is configured
        if(! _is_my_nautique_configured($configuration)){
            error_log('api/backend: No myNautique configured');
            return Status::successDataResponse('no mynautique configured', $statusResponse);
        }else{
            $statusResponse['myNautiqueConfigured'] = TRUE;
        }

        return Status::successDataResponse("all backend components are good", $statusResponse);
    }

    // returns http status permission denied in case the user is not an admin
    function _is_admin_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn($configuration->admin_user_role_id)){
            HttpHeader::setResponseCode(403);
            exit;
        }
    }

    // returns whether the database is configured currently
    function _is_db_configured($configuration){
        // check if db_server is set
        if(!isset($configuration->db_server) or $configuration->db_server == ''){
            return FALSE;
        }
        // check if db_name is set
        if(!isset($configuration->db_name) or $configuration->db_name == ''){
            return FALSE;
        }
        // check if db_user is set
        if(!isset($configuration->db_user) or $configuration->db_user == ''){
            return FALSE;
        }

        return TRUE;
    }

    // returns whether myNautique is configured currently
    function _is_my_nautique_configured($configuration){
        // check if myNautique "enabled" is present
        if(!isset($configuration->mynautique_enabled)){
            return FALSE;
        }
        // check if user is set
        if(!isset($configuration->mynautique_user) or $configuration->mynautique_user == ''){
            return FALSE;
        }
        // check if password is set
        if(!isset($configuration->mynautique_password) or $configuration->mynautique_password == ''){
            return FALSE;
        }

        return TRUE;
    }

?>