<?php
    // automatically load all classes
    spl_autoload_register('boat_autoloader');
    function boat_autoloader($class){
        include '../classes/'.$class.'.php';
    }

    // Get configuration access
    $configuration = new Configuration();

    // Check if the user is already logged in and is of type admin
    $lc = new Login($configuration);
    if(!$lc->isLoggedInAndRole($configuration->admin_user_status_id, $configuration->login_page)){
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

    switch($_GET['action']){
        case 'get_engine_hours_log':
            $response = get_engine_hours_log($configuration);
            echo json_encode($response);
            exit;
        case 'get_engine_hours_latest':
            $response = get_engine_hours_latest($configuration);
            echo json_encode($response);
            exit;
        case 'get_engine_hours_entry':
            get_engine_hours_entry($configuration);
            exit;
        case 'update_engine_hours_entry':
            update_engine_hours_entry($configuration);
            exit;
        case 'update_engine_hours':
            $response = update_engine_hours($configuration);
            echo json_encode($response);
            exit;
        case 'update_fuel':
            $response = update_fuel($configuration);
            echo json_encode($response);
            exit;
        case 'get_fuel_entry':
            get_fuel_entry($configuration);
            exit;
        case 'update_fuel_entry':
            update_fuel_entry($configuration);
            exit;
        case 'get_fuel_log':
            $response = get_fuel_log($configuration);
            echo json_encode($response);
            exit;
        case 'get_maintenance_log':
            get_maintenance_log($configuration);
            exit;
        case 'update_maintenance_log':
            update_maintenance_log($configuration);
            exit;
    }

    HttpHeader::setResponseCode(400);
    return;

    /* Returns the log_book data */
    function get_engine_hours_log($configuration){
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('api/boat: Cannot connect to the database');
            return Status::errorStatus("Cannot connect to the database");
        }

        $query = 'SELECT blb.id as id, UNIX_TIMESTAMP(DATE_FORMAT(blb.timestamp, "%Y-%m-%dT%T+02:00")) as time, blb.before_hours as before_hours,
                        blb.after_hours as after_hours, blb.delta_hours as delta_hours,
                        blb.type as type, st.name as type_name, u.id as user_id,
                        u.first_name as user_first_name, u.last_name as user_last_name
                    FROM boat_engine_hours blb, user u, session_type st
                    WHERE blb.user_id = u.id
                    AND blb.type    = st.id
                    ORDER BY blb.timestamp DESC';

        // Limit the number of returned lines for mobile browsers
        if(MobileDevice::isMobileBrowser() or isset($_GET['mobile'])){
            $query .= "  LIMIT 0, 100";
        }

        $res = $db->fetch_data_hash($query);
        $db->disconnect();
        if(!isset($res)){
            return Status::errorStatus("Cannot retrieve engine hour logs from database.");
        }
        return Status::successDataResponse("Engine hour log retrieved.", $res);
    }

    /* Returns the fuel log data */
    function get_fuel_log($configuration){
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('api/boat: Cannot connect to the database');
            return Status::errorStatus("Cannot connect to the database");
        }

        $query = 'SELECT bf.id as id, UNIX_TIMESTAMP(DATE_FORMAT(bf.timestamp, "%Y-%m-%dT%TZ")) as timestamp, bf.engine_hours as engine_hours,
                        bf.liters as liters, bf.cost_chf as cost, bf.cost_chf_brutto as cost_brutto,
                        u.id as user_id,
                        u.first_name as user_first_name, u.last_name as user_last_name
                    FROM boat_fuel bf, user u
                    WHERE bf.user_id = u.id
                    ORDER BY bf.timestamp DESC';

        $res = $db->fetch_data_hash($query, -1);
        $db->disconnect();
        if(!isset($res)){
            HttpHeader::setResponseCode(500);
            return Status::errorStatus("Cannot get the fuel logs from the database.");
        }
        return Status::successDataResponse("success", $res);
    }

    /* Returns the maintenance log data */
    function get_maintenance_log($configuration){
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('api/boat: Cannot connect to the database');
            HttpHeader::setResponseCode(500);
            exit;
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
            HttpHeader::setResponseCode(500);
        }
        echo json_encode($res);
    }

    /* Returns the latest log_book entry, even if not complete yet */
    function get_engine_hours_latest($configuration){
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            error_log('api/boat: Cannot connect to the database');
            return Status::errorStatus("Cannot connect to the database");
        }

        $query = 'SELECT blb.id as id, blb.timestamp as time_stamp, blb.before_hours as before_hours,
                    blb.after_hours as after_hours, blb.delta_hours as delta_hours,
                    blb.type as type, st.name as type_name, u.id as user_id,
                    u.first_name as user_first_name, u.last_name as user_last_name
                    FROM boat_engine_hours blb
                    LEFT JOIN user u ON u.id = blb.user_id
                    LEFT JOIN session_type st ON st.id = blb.type
                    INNER JOIN (SELECT max(timestamp) as timestamp FROM boat_engine_hours) beh ON beh.timestamp = blb.timestamp';
        $res = $db->fetch_data_hash($query);
        $db->disconnect();
        if(!isset($res)){
            return Status::errorStatus("Cannot get the latest engine hours entry");
        }
        return Status::successDataResponse("success", $res);
    }

    function get_engine_hours_entry($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // input validation
        $sanitizer = new Sanitizer();
        if(!$post_data->id or !$sanitizer->isInt($post_data->id)){
            HttpHeader::setResponseCode(400);
            $res = array();
            $res["error"] = "Invalid engine hours log ID format";
            error_log('api/boat: Cannot connect to the database');
            echo json_encode($res);
            exit;
        }

        $db = new DBAccess($configuration);
        if(!$db->connect()){
            $res = array();
            $res["error"] = "Cannot connect to backend database";
            error_log('api/boat: Cannot connect to the database');
            echo json_encode($res);
            HttpHeader::setResponseCode(500);
            exit;
        }

        $query = 'SELECT blb.id as id, UNIX_TIMESTAMP(DATE_FORMAT(blb.timestamp, "%Y-%m-%dT%T+02:00")) as time, blb.before_hours as before_hours,
                blb.after_hours as after_hours, blb.delta_hours as delta_hours,
                blb.type as type, st.name as type_name, u.id as user_id,
                u.first_name as user_first_name, u.last_name as user_last_name
            FROM boat_engine_hours blb, user u, session_type st
            WHERE blb.user_id = u.id
            AND blb.type    = st.id
            AND blb.id = ';
        $query = $query . $post_data->id;

        $res = $db->fetch_data_hash($query);
        $db->disconnect();
        if(!isset($res[0])){
            HttpHeader::setResponseCode(500);
            return;
        }

        HttpHeader::setResponseCode(200);
        echo json_encode($res[0]);
        return;
    }

    // Note: currently we only allow to switch the type
    function update_engine_hours_entry($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(! isset($post_data->type) or !$sanitizer->isInt($post_data->type)){
            // currently we only allow 0 (private) or 1 (course)
            if($post_data->type != 0 and $post_data->type != 1){
                HttpHeader::setResponseCode(400);
                echo "No valid engine hours entry type selected";
                return;
            }
        }
        if(!$post_data->id or !$sanitizer->isInt($post_data->id)){
            HttpHeader::setResponseCode(400);
            echo "No valid engine hours entry ID selected";
            return;
        }

        // setup database access
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            HttpHeader::setResponseCode(500);
            exit;
        }

        // check that the entry does indeed exist
        $query = 'UPDATE boat_engine_hours 
            SET type = ?
            WHERE id = ?';
        $db->prepare($query);
        $db->bind_param('dd',
            $post_data->type,
            $post_data->id
        );
        if(!$db->execute()){
            $db->disconnect();
            HttpHeader::setResponseCode(500);
            return;
        }

        HttpHeader::setResponseCode(200);
        return;
    }

    function update_engine_hours($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
            return Status::errorStatus("No valid user selected, please select a user");
        }
        if(!$post_data->engine_hours_before or ! $sanitizer->isFloat($post_data->engine_hours_before)){
            return Status::errorStatus("No valid value for engine hours 'before' given.");
        }
        if(isset($post_data->engine_hours_after) and !$sanitizer->isFloat($post_data->engine_hours_after)){
            return Status::errorStatus("No valid value for engine hours 'after' given.");
        }
        if(isset($post_data->type) and !$sanitizer->isInt($post_data->type)){
            return Status::errorStatus("No valid type for the boat usage has been selected.");
        }
        // logical input validation
        if(isset($post_data->engine_hours_after)){
            if($post_data->engine_hours_after < $post_data->engine_hours_before){
                return Status::errorStatus("Engine hours afterwards cannot be smaller than before.");
            }
        }

        // check if there exists an empty log entry which needs an update
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            return Status::errorStatus("Cannot connect to database");
        }
        $query = 'SELECT count(*) as numberof FROM (SELECT blb.id as id, blb.before_hours as before_hours
                    FROM boat_engine_hours blb
                WHERE blb.after_hours IS NULL AND
                        blb.delta_hours IS NULL) as dummy';
        $res = $db->fetch_data_hash($query);
        if($res == FALSE){
            $db->disconnect();
            return Status::errorStatus("Cannot determine whether this is a new entry or one that can be updated.");
        }

        if($res[0]['numberof'] > 0 and $post_data->engine_hours_after){
            // update the existing entry
            $query = 'UPDATE boat_engine_hours blb
                        SET blb.after_hours = ?,
                            blb.delta_hours = ?
                    WHERE blb.after_hours IS NULL';
            $db->prepare($query);
            $db->bind_param('dd',
                            $post_data->engine_hours_after,
                            $post_data->engine_hours_after - $post_data->engine_hours_before);
            if(!$db->execute()){
                $db->disconnect();
                return Status::errorStatus("Cannot update existing entry. An error occurred.");
            }
        }elseif($res[0]['numberof'] > 0){
            $db->disconnect();
            return Status::errorStatus("To update the engine hour entry, please provide the engine hours after the trip.");
        }
        else{
            // check that type is valid
            $usage_type = 0;
            if(! isset($post_data->type)){
                return Status::errorStatus("No valid boat usage type selected.");
            }else{
                $usage_type = $post_data->type;
            }

            // create a new entry
            $date = new DateTime();
            $query = 'INSERT INTO boat_engine_hours
                    (timestamp, before_hours, type, user_id)
                    VALUES
                    (?,?,?,?)';
            $db->prepare($query);
            $db->bind_param('sdii',
                            $date->format('Y-m-d H:i:s'),
                            sprintf('%.2f',$post_data->engine_hours_before),
                            $usage_type,
                            $post_data->user_id);
            if(!$db->execute()){
                $db->disconnect();
                return Status::errorStatus("Internal Server Error, cannot add new engine hours entry.");
            }
        }
        $db->disconnect();
        return Status::successStatus("new engine hour entry added or existing one updated");
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

        // create a new entry
        $date = new DateTime();
        $query = "INSERT INTO boat_fuel
                (timestamp, engine_hours, liters, cost_chf, user_id)
                VALUES
                (?,?,?,?,?);";
        $db->prepare($query);
        $db->bind_param('sdddi',
            $date->format('Y-m-d H:i:s'),
            sprintf('%.2f', $post_data->engine_hours),
            sprintf('%.2f', $post_data->liters),
            sprintf('%.2f', $post_data->cost),
            $post_data->user_id
        );
        if(!$db->execute()){
            $db->disconnect();
            return Status::errorStatus("Internal Server Error, cannot add new fuel entry.");
        }
        return Status::successStatus("successfully added new fuel entry");
    }

    function get_fuel_entry($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(!$post_data->id or !$sanitizer->isInt($post_data->id)){
            HttpHeader::setResponseCode(400);
            echo "No valid fuel entry ID given";
            exit;
        }

        // setup database access
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            HttpHeader::setResponseCode(500);
            exit;
        }

        // get the entry
        $query = 'SELECT bf.id as id, UNIX_TIMESTAMP(DATE_FORMAT(bf.timestamp, "%Y-%m-%dT%TZ")) as timestamp, 
                bf.engine_hours as engine_hours,
                bf.liters as liters, bf.cost_chf as cost, 
                bf.cost_chf_brutto as cost_brutto,
                u.id as user_id,
                u.first_name as user_first_name, u.last_name as user_last_name
            FROM boat_fuel bf, user u
            WHERE bf.user_id = u.id 
                AND bf.id = ';
        $query = $query . $post_data->id;

        $res = $db->fetch_data_hash($query);
        $db->disconnect();
        if(!isset($res[0])){
            HttpHeader::setResponseCode(500);
            echo "No entry found";
            error_log("api/boat: cannot find a fuel entry for ID: $post_data->id");
            return;
        }

        $response = array();
        $response['fuel_entry'] = $res[0];
        $response['currency']   = $configuration->currency;

        echo json_encode($response);
        return;
    }

    function update_fuel_entry($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(!$post_data->id or !$sanitizer->isInt($post_data->id)){
            HttpHeader::setResponseCode(400);
            echo "No valid fuel entry ID given";
            return;
        }
        if(!$post_data->engine_hours or ! $sanitizer->isFloat($post_data->engine_hours)){
            HttpHeader::setResponseCode(400);
            echo "No valid value for engine hours given.";
            return;
        }
        if($post_data->liters and !$sanitizer->isFloat($post_data->liters)){
            HttpHeader::setResponseCode(400);
            echo "No valid value for the litres of fuel given.";
            return;
        }
        if($post_data->cost and !$sanitizer->isFloat($post_data->cost)){
            HttpHeader::setResponseCode(400);
            echo "No valid value for the cost (net) of fuel given.";
            return;
        }
        if($post_data->cost_brutto and !$sanitizer->isFloat($post_data->cost_brutto)){
            HttpHeader::setResponseCode(400);
            echo "No valid value for the cost (gros) of fuel given.";
            return;
        }

        // setup database access
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            HttpHeader::setResponseCode(500);
            exit;
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
            HttpHeader::setResponseCode(500);
            return;
        }

        HttpHeader::setResponseCode(200);
        return;
    }

    function update_maintenance_log($configuration){
        $post_data = json_decode(file_get_contents('php://input'));

        // general input validation
        $sanitizer = new Sanitizer();
        if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
            HttpHeader::setResponseCode(400);
            echo "No valid user selected, please select a user";
            return;
        }
        if(!$post_data->engine_hours or ! $sanitizer->isFloat($post_data->engine_hours)){
            HttpHeader::setResponseCode(400);
            echo "No valid value for engine hours given.";
            return;
        }
        if(!$post_data->description){
            HttpHeader::setResponseCode(400);
            echo "No valid value for the litres of fuel given.";
            return;
        }

        // create a new entry
        $date = new DateTime();
        $db = new DBAccess($configuration);
        if(!$db->connect()){
            HttpHeader::setResponseCode(500);
            exit;
        }
        $query = "INSERT INTO boat_maintenance
                (timestamp, engine_hours, description, user_id)
                VALUES
                (?,?,?,?);";
        $db->prepare($query);
        $db->bind_param('sdsi',
                        $date->format('Y-m-d H:i:s'),
                        sprintf('%.2f', $post_data->engine_hours),
                        $post_data->description,
                        $post_data->user_id);
        if(!$db->execute()){
            $db->disconnect();
            HttpHeader::setResponseCode(500);
            echo "Internal Server Error, cannot add new fuel entry";
            return;
        }
        HttpHeader::setResponseCode(200);
    }

?>
