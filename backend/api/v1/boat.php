<?php
    // automatically load all classes
    spl_autoload_register('boat_autoloader');
    function boat_autoloader($class){
        include '../../classes/'.$class.'.php';
    }

    // Get configuration access
    $configuration = new Configuration();

    // Check if the user is already logged in and is of type admin
    $lc = new Login($configuration);
    if(!$lc->isAdmin($configuration->admin_user_role_id)){
        $response = array();
        $response["redirect"] = $configuration->login_page;
        echo json_encode($response);
        exit;
    }

    // check if we have an action
    if(!isset($_GET['action'])){
        HttpHeader::setResponseCode(200);
        exit;
    }

    $response = null;

    switch($_GET['action']){
        case 'get_maintenance_log':
            $response = get_maintenance_log($configuration);
            break;
        case 'update_maintenance_log':
            $response = update_maintenance_log($configuration);
            break;
        default:
            $response = Status::errorStatus("Action not supported");
            break;
    }

    echo json_encode($response);
    return;

    /* Returns the maintenance log data */
    function get_maintenance_log($configuration){
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            return Status::errorStatus("Cannot connect to the database");
        }

        $query = 'SELECT bm.id as id, UNIX_TIMESTAMP(DATE_FORMAT(bm.timestamp, "%Y-%m-%dT%TZ")) as timestamp, bm.engine_hours,
                    bm.description as description, u.id as user_id,
                    u.first_name as user_first_name, u.last_name as user_last_name
                    FROM   boat_maintenance bm, user u
                    WHERE bm.user_id = u.id
                    ORDER BY bm.timestamp DESC;';
        $res = $db->fetch_data_hash($query, -1);
        $db->disconnect();
        if(!isset($res)){
            return Status::errorStatus("Cannot retrieve maintenance log entries from the database");
        }
        return Status::successDataResponse("success", $res);
    }

    function update_maintenance_log($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
            return Status::errorStatus("No valid user selected, please select a user");
        }
        if(!$post_data->engine_hours or ! $sanitizer->isFloat($post_data->engine_hours)){
            return Status::errorStatus("No valid value for engine hours given.");
        }
        if(!$post_data->description){
            return Status::errorStatus("No valid value for the litres of fuel given.");
        }

        // create a new entry
        $date = new DateTime();
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            return Status::errorStatus("Cannot connect to database.");
        }
        $query = "INSERT INTO boat_maintenance
                (timestamp, engine_hours, description, user_id)
                VALUES
                (?,?,?,?);";
        $db->prepare($query);
        $db->bind_param('sdsi',
                        $date->format('Y-m-d H:i:s'),
                        sprintf('%.5f', $post_data->engine_hours),
                        $post_data->description,
                        $post_data->user_id);
        if(!$db->execute()){
            $db->disconnect();
            return Status::errorStatus("Cannot add the new maintenance entry to the database due to an error.");
        }

        return Status::successStatus("maintenance entry added successfully");
    }

?>
