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
        case 'update_fuel':
            $response = update_fuel($configuration);
            break;
        case 'update_fuel_entry':
            $response = update_fuel_entry($configuration);
            break;
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

    function update_fuel($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
            return Status::errorStatus("No valid user selected, please select a user");
        }
        if(!$post_data->engine_hours or ! $sanitizer->isFloat($post_data->engine_hours)){
            return Status::errorStatus("No valid value for engine hours given.");
        }
        if($post_data->liters and !$sanitizer->isFloat($post_data->liters)){
            return Status::errorStatus("No valid value for the litres of fuel given.");
        }
        if($post_data->cost and !$sanitizer->isFloat($post_data->cost)){
            return Status::errorStatus("No valid value for the cost of fuel given.");
        }

        // get the maximum engine hour log
        $query = 'SELECT max(blb.after_hours) as max_hours FROM boat_engine_hours blb;';
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            return Status::errorStatus("Cannot retrieve maximum engine hours.");
        }
        $res = $db->fetch_data_hash($query);
        $max_hours = 0;
        if($res[0]['max_hours']){
            $max_hours = $res[0]['max_hours'];
        }

        // handle the fuel payment type
        // note: either we pay instantly -> contributes_to_balance = 1
        //       or we pay when we get a bill -> contributes_to_balance = 0
        $contributes_to_balance = 1;
        if($configuration->fuel_payment_type == "billed"){
            $contributes_to_balance = 0;
        }

        // create a new entry
        $date = new DateTime();
        $query = "INSERT INTO boat_fuel
                (timestamp, engine_hours, liters, cost_chf, user_id, contributes_to_balance)
                VALUES
                (?,?,?,?,?,?);";
        $db->prepare($query);
        $db->bind_param('sdddii',
            $date->format('Y-m-d H:i:s'),
            sprintf('%.5f', $post_data->engine_hours),
            sprintf('%.3f', $post_data->liters),
            sprintf('%.3f', $post_data->cost),
            $post_data->user_id,
            $contributes_to_balance
        );
        if(!$db->execute()){
            $db->disconnect();
            return Status::errorStatus("Internal Server Error, cannot add new fuel entry.");
        }
        return Status::successStatus("successfully added new fuel entry");
    }

    function update_fuel_entry($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(!isset($post_data->id) or !$sanitizer->isInt($post_data->id)){
            return Status::errorStatus("No valid fuel entry ID given");
        }
        if(!isset($post_data->engine_hours) or ! $sanitizer->isFloat($post_data->engine_hours)){
            return Status::errorStatus("No valid value for engine hours given.");
        }
        if(!isset($post_data->liters) or !$sanitizer->isFloat($post_data->liters)){
            return Status::errorStatus("No valid value for the litres of fuel given.");
        }
        if(!isset($post_data->cost) or !$sanitizer->isFloat($post_data->cost)){
            return Status::errorStatus("No valid value for the cost (net) of fuel given.");
        }
        if(isset($post_data->cost_brutto) and !$sanitizer->isFloat($post_data->cost_brutto)){
            return Status::errorStatus("No valid value for the cost (gros) of fuel given.");
        }

        // setup database access
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            return Status::errorStatus("Cannot connect to the database.");
        }

        // check that the entry does indeed exist
        $query = 'UPDATE boat_fuel 
            SET engine_hours = ?,
                liters = ?,
                cost_chf = ?,
                cost_chf_brutto = ?
            WHERE id = ?';
        $db->prepare($query);
        $db->bind_param('ddddd',
            $post_data->engine_hours,
            $post_data->liters,
            $post_data->cost,
            $post_data->cost_brutto,
            $post_data->id
        );
        if(!$db->execute()){
            $db->disconnect();
            return Status::errorStatus("Cannot update fuel entry due to unknown reasons. Please try again.");
        }

        return Status::successStatus("successfully updated");
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
