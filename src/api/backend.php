<?php
    // automatically load all classes
    spl_autoload_register('backend_autoloader');
    function backend_autoloader($class){
        if(file_exists(__DIR__.'/../classes/'.$class.'.php')){
            include __DIR__.'/../classes/'.$class.'.php';
        }
    }

    // test whether the config-file exists
    if (!file_exists('../config/config.php') && is_readable('../config/config.php')) {
        error_log("No config file yet");
        $status = array();
        $status['ok'] = FALSE;
        $status['message'] = 'no config';
        $status['redirect'] = '/setup.html';
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

    switch($_GET['action']){
        case 'get_status':
            get_status($configuration);
            exit;
        case 'get_version':
            get_version($configuration);
            exit;
        case 'admin_check_database_update':
            admin_check_database_update($configuration);
            exit;
        case 'update_database':
            update_database($configuration);
            exit;
    }

    HttpHeader::setResponseCode(400);
    $status = array();
    $status['ok'] = FALSE;
    $status['message'] = 'invalid action requested';
    echo json_encode($status);
    return;

    // returns the status of the backend
    function get_status($configuration){
        $status = array();
        $status['ok'] = TRUE;
        $status['message'] = 'all backend components are good';

        // check if db is configured
        if(! _is_db_configured($configuration)){
            error_log('api/backend: No database configured');
            $status['ok'] = FALSE;
            $status['message'] = 'no database configured';
            echo json_encode($status);
            return;
        }
        
        // check if db access works
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('api/backend: Cannot connect to the database');
            $status['ok'] = FALSE;
            $status['message'] = 'database connection is not working';
        }

        $res = $db->fetch_data_hash("SELECT count(*) AS users FROM user;");
        error_log(print_r($res, TRUE));
        if($res == NULL || $res[0] == NULL || $res[0]['users'] == NULL || $res[0]['users'] < 1){
            error_log('api/backend: No user configured');
            $status['ok'] = FALSE;
            $status['message'] = 'no user configured';
        }

        $db->disconnect();
        echo json_encode($status);
        return;
    }

    // returns the version of the backend
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

        echo json_encode($response);
        return;
    }

    // check whether the database needs to be updated
    function admin_check_database_update($configuration){
        // check that the user is admin indeed
        is_admin_or_return($configuration);

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
        $db_versions          = explode(".", $db_version);
        $db_versions_required = explode(".", $db_version_required);
        $app_versions         = explode(".", $app_version);
        if(count($db_versions) != count($db_versions_required)){
            $response['updateAvailable'] = TRUE;
        }
        else{
            for($i = 0; $i < count($db_versions); $i++){
                if($db_versions[$i] < $db_versions_required[$i]){
                    $response['updateAvailable'] = TRUE;
                    break;
                }
            }
        }        
        
        echo json_encode($response);
        return;
    }

    function update_database($configuration){
        // only an administrator is supposed to do this
        is_admin_or_return($configuration);

        // get current schema version
        $db_schema_version = _get_database_version($configuration);
        $backend_version   = $configuration->version;

        // get the update files
        $update_files = _get_update_files($db_schema_version, $configuration->version);
        error_log(print_r($update_files, TRUE));

        $db = new DBAccess($configuration);
        if(!$db->connect()){
            $status['ok'] = FALSE;
            $status['msg'] = 'DB access not working with the provided settings. Please check for typos or make sure the database is reachable.';
            echo json_encode($status);
            return;
        }

        // apply changes for every update file
        foreach ($update_files as $file){
            // setup a new database
            if(!$db->apply_sql_file(__DIR__."/../config/db/updates/".$file)){
                $status['ok']  = FALSE;
                $status['msg'] = "DB update in file $file could not be applied to database, please check error logs";
                $db->disconnect();
                echo json_encode($status);
                return;
            }
        }

        // get the new database version
        $status['db_version'] = _get_database_version($configuration);

        // update was successful
        $db->disconnect();
        $status['ok'] = TRUE;
        $status['msg'] = 'database upgrade successful.';
        echo json_encode($status);
        return;
    }

    // returns http status permission denied in case the user is not an admin
    function is_admin_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn($configuration->admin_user_status_id)){
            HttpHeader::setResponseCode(403);
            exit;
        }
    }

    function _get_database_version($configuration){
        // connect to database
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('Cannot connect to database');
            return NULL;
        }

        // check whether db version tracking table exists
        // if not: an update is most probably required
        $query = "SELECT * 
            FROM information_schema.tables
            WHERE table_schema = ? 
            AND table_name = 'configuration'
            LIMIT 1;";
        $db->prepare($query);
        $db->bind_param('s',
            $configuration->db_name
        );
        if(!$db->execute()){
            $db->disconnect();
            return NULL;
        }
        $result = $db->fetch_stmt_hash();
        if(count($result) == 0){
            return 0;
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

    // returns an array the sql update files
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
        preg_match('/^(\d)+\.(\d)+\.*(\d+)*$/', $a, $version_a);
        preg_match('/^(\d)+\.(\d)+\.*(\d+)*$/', $b, $version_b);

        // go through every version segment
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