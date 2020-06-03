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
        if(file_exists(__DIR__.'/../classes/'.$class.'.php')){
            include __DIR__.'/../classes/'.$class.'.php';
        }
    }

    // test whether the config-file exists
    if (!file_exists('../config/config.php') && is_readable('../config/config.php')) {
        error_log("No config file yet");
        $status = array(
            'ok' => FALSE,
            'message' => 'no config',
            'redirect' => '/setup.html',
        );
        echo json_encode($status);
        exit;
    }

    // check if we have an action
    if(!isset($_GET['action'])){
        HttpHeader::setResponseCode(200);
        exit;
    }

    // get configuration
    $configuration = new Configuration();
    $status = "";

    switch($_GET['action']){
        case 'get_status':
            $status = get_status($configuration);
            break;
        case 'get_version':
            $status = get_version($configuration);
            break;
        case 'admin_check_database_update':
            $status = admin_check_database_update($configuration);
            break;
        case 'update_database':
            $status = update_database($configuration);
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
     *
     * @param $configuration    configuration object
     * @return status which contains the following structure:
     * {
     *      ok      => TRUE|FALSE,
     *      message => "text description"
     * }
     */
    function get_status($configuration){
        $status = array();
        $status['ok'] = TRUE;
        $status['message'] = 'all backend components are good';

        // check if db is configured
        if(! _is_db_configured($configuration)){
            error_log('api/backend: No database configured');
            return Status::errorStatus('no database configured');
        }
        
        // check if db access works
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('api/backend: Cannot connect to the database');
            return Status::errorStatus('database connection is not working');
        }

        $res = $db->fetch_data_hash("SELECT count(*) AS users FROM user;");
        if($res == NULL || $res[0] == NULL || $res[0]['users'] == NULL || $res[0]['users'] < 1){
            error_log('api/backend: No user configured');
            $status = Status::errorStatus('no user configured');
        }

        $db->disconnect();
        return $status;
    }

    /**
     * get_version returns the version of the different components of the backend
     *
     * @param $configuration    configuration object
     * @return $response containing:
     * {
     *      app_version         => "1.0.0"
     *      db_version          => "1.0.0"
     *      db_version_required => "1.0.0"
     * }
     */
    function get_version($configuration){
        $response = array();

        // get app version
        if(isset($configuration->version)){
            $response["app_version"] = $configuration->version;
        } else{
            $response["app_version"] = null;
        }

        // get database schema version
        $response["db_version_required"] = $configuration->required_db_version;
        $response["db_version"] = _get_database_version($configuration);

        return $response;
    }

    /**
     * Checks whether the database needs to be updated or not
     *
     * @param $configuration    configuration object
     * @return $response of the form:
     * {
     *      ok              => TRUE|FALSE
     *      message         => "description"
     *      updateAvailable => TRUE|FALSE
     * }
     */
    function admin_check_database_update($configuration){
        // check that the user is admin indeed
        _is_admin_or_return($configuration);

        // prepare response
        $response = array();
        $response['updateAvailable'] = FALSE;

        // get versions
        $app_version         = $configuration->version;
        $db_version_required = $configuration->required_db_version;
        $db_version          = _get_database_version($configuration);
        if(!isset($db_version)){
            $response['ok'] = FALSE;
            $response['message'] = "Cannot get db version";
            $response['updateAvailable'] = FALSE;
        }

        // compare versions
        if( _version_cmp($db_version, $db_version_required) == -1 ){
            $response['updateAvailable'] = TRUE;
        }
        
        return $response;
    }

    /**
     * Updates the database, admin priviledges are required to call this function.
     *
     * @param $configuration        the configuration object
     * @return response containing:
     * {
     *      ok          => TRUE|FALSE
     *      msg         => "message describing the failure or success"
     *      db_version  => "1.0.0"                  // database version after upgrade
     * }
     */
    function update_database($configuration){
        // only an administrator is supposed to do this
        _is_admin_or_return($configuration);

        // get current schema version
        $db_schema_version = _get_database_version($configuration);

        // get the update files
        $update_files = _get_update_files($db_schema_version, $configuration->required_db_version);
        
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            return Status::errorStatus('DB access not working with the provided settings. Please check for typos or make sure the database is reachable.');
        }

        // apply changes for every update file
        foreach ($update_files as $file){
            // setup a new database
            if(!$db->apply_sql_file(__DIR__."/../config/db/updates/".$file)){
                $db->disconnect();
                return Status::errorStatus("DB update in file $file could not be applied to database, please check error logs");
            }
        }

        // get the new database version
        $status['db_version'] = _get_database_version($configuration);

        // update was successful
        $db->disconnect();
        return Status::successStatus('database upgrade successful.');
    }

    // returns http status permission denied in case the user is not an admin
    function _is_admin_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn($configuration->admin_user_status_id)){
            HttpHeader::setResponseCode(403);
            exit;
        }
    }

    // returns the database version
    function _get_database_version($configuration){
        // connect to database
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('Cannot connect to database');
            return NULL;
        }

        // check the current version of the database schema
        // Note: if we reach this part of the code, the configuration table exists
        $query = 'SELECT value 
            FROM configuration
            WHERE property = "schema.version";';
        $db_result = $db->fetch_data_hash($query);
        if(count($db_result) != 1){
            error_log('no value found for schema.version');
            $db->disconnect();
            return NULL;
        }

        return $db_result[0]["value"];
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

    // returns all update files in sorted order that are required to get from 
    // current to desired version
    function _get_update_files($current_schema_version, $desired_schema_version){
        // get all schema files
        $files = _get_all_schema_update_files();
        // sort them by version
        usort($files, "_schema_version_cmp");

        $update_files = array();
        foreach ($files as $file){
            // get version from file
            $version = array();
            preg_match('/.*_v(.*)\.sql$/', $file, $version);

            $file_version = $version[1];

            // file_version is smaller or equal than current version
            if(_version_cmp($file_version, $current_schema_version) == 0 
                || _version_cmp($file_version, $current_schema_version) == -1)
            {
                continue;
            }

            // file_version is bigger than desired version
            if(_version_cmp($file_version, $desired_schema_version) == 1){
                continue;
            }

            // add file to the list uf update files
            array_push($update_files, $file);
        }

        return $update_files;
    }

    // returns an array containing the sql update files
    function _get_all_schema_update_files(){
        $files = scandir (__DIR__.'/../config/db/updates/');
        $sql_files = array();
        foreach($files as $f){
            if(preg_match('/\.sql$/', $f)){
                array_push($sql_files, $f);
            }
        }
        return $sql_files;
    }

    // sorts the sql update files based on their version
    function _schema_version_cmp($a, $b){
        // extract version from $a and $b
        $version_a = array();
        $version_b = array();
        preg_match('/.*_v(.*)\.sql$/', $a, $version_a);
        preg_match('/.*_v(.*)\.sql$/', $b, $version_b);

        return _version_cmp($version_a[1], $version_b[1]);
    }

    // sorts version numbers
    // returns:
    //  0: if version $a == $b
    //  1: if version $a > $b
    // -1: if version $a < $b
    function _version_cmp($a, $b){
        $version_a = array();
        $version_b = array();
        // in index 1 there will be the entire preg_match
        preg_match('/^(\d+)\.(\d+)\.?(\d+)?$/', $a, $version_a);
        preg_match('/^(\d+)\.(\d+)\.?(\d+)?$/', $b, $version_b);

        // go through every version segment
        // start with 1 as 0 is reserveved for the entire match in the regex
        for ($i=1; $i<max(count($version_a), count($version_b)); $i++){
            if(isset($version_a[$i]) && isset($version_b[$i])){
                if ($version_a[$i] < $version_b[$i]){
                    return -1;
                } else if ($version_a[$i] > $version_b[$i]){
                    return 1;
                }
            }
            if(!isset($version_a[$i]) && isset($version_b[$i])){
                return -1;
            }
            if(isset($version_a[$i]) && !isset($version_b[$i])){
                return 1;
            }
        }

        return 0;
    }



?>