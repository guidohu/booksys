<?php
  
  // automatically load all classes
  spl_autoload_register('user_autoloader');
  function user_autoloader($class){
    include '../classes/'.$class.'.php';
  }

  // Load configuration
  $configuration = new Configuration();
  
  // Check if the user is already logged in and is of type admin
  $lc = new Login($configuration);
  if(!$lc->isLoggedInRedirect($configuration->login_page)){
      HttpHeader::setResponseCode(403);
      exit;
  }
  
  // check if we have an action
  if(!isset($_GET['action'])){
    HttpHeader::setResponseCode(200);
    exit;
  }
  
  $_SESSION = Login::getSessionData($configuration);
  
  switch($_GET['action']){
    case 'get_admin_users':
        get_admin_users($configuration, $lc);
        exit;
    case 'get_session_users':
        get_session_users($configuration, $lc);
        exit;
    case 'get_all_users':
        $response = get_all_users($configuration, $lc);
        echo json_encode($response); 
        exit;
    case 'get_my_user':
        get_my_user($configuration);
        exit;
    case 'get_my_user_heats':
        get_my_user_heats($configuration);
        exit;
    case 'get_my_user_sessions':
        get_my_user_sessions($configuration);
        exit;
    case 'change_my_user_data':
        $response = change_my_user_data($configuration);
        echo json_encode($response);
        exit;
    case 'lock_user':
        $response = lock_user($configuration, $lc);
        echo json_encode($response);
        exit;
    case 'unlock_user':
        $response = unlock_user($configuration, $lc);
        echo json_encode($response);
        exit;
    case 'change_user_group_membership':
        $response = change_user_group_membership($configuration, $lc);
        echo json_encode($response);
        exit;
	case 'get_all_users_detailed':
        $response = get_all_users_detailed($configuration, $lc);
        echo json_encode($response);
        exit;
    case 'get_user_groups':
        $response = get_user_groups($configuration, $lc);
        echo json_encode($response);
        exit; 
    case 'get_user_roles':
        $response = get_user_roles($configuration, $lc);
        echo json_encode($response);
        exit;
    case 'save_user_group':
        $response = save_user_group($configuration, $lc);
        echo json_encode($response);
        exit;
    case 'delete_user_group':
        $response = delete_user_group($configuration, $lc);
        echo json_encode($response);
        exit;
    case 'delete_user':
        $response = delete_user($configuration, $lc);
        echo json_encode($response);
        exit;
  }
  
  HttpHeader::setResponseCode(400);
  return;

  function delete_user($configuration, $lc){
    // user needs to be admin
    if(!$lc->isAdmin()){
        return Status::errorStatus("permission denied");
    }

    $data = json_decode(file_get_contents('php://input'));

    // sanitize
    $sanitizer = new Sanitizer();
    if(! isset($data->id) or !$sanitizer->isInt($data->id)){
        return Status::errorStatus("No valid user id given.");
    }

    // do not delete yourself
    $current_user = $lc->getSessionData($configuration);
    if($data->id == $current_user['user_id']){
        return Status::errorStatus("Deleting your own user is not possible.");
    }

    // the user that should be deleted cannot be admin
    $userAPI = new User($configuration);
    $user    = $userAPI->getUserById($data->id);
    if($user['user_role_name'] == 'admin'){
        return Status::errorStatus("Administrator users cannot be deleted.");
    }

    // the user that should be deleted has to have a zero balance
    $balance = $userAPI->getUserBalance($user['id']);
    if($balance != 0){
        return Status::errorStatus('User can only be delete if they have a balance of 0.');
    }

    // delete user
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        error_log('api/user: Cannot connect to the database');
        return Status::errorStatus('Cannot connect to the database.');
    }
    $query = "UPDATE user
        SET deleted = 1,
            comment = NULL,
            locked  = 1,
            license = 0,
            email   = NULL,
            mobile  = NULL,
            plz     = NULL,
            city    = NULL,
            address = NULL,
            last_name = 'user',
            first_name = 'deleted',
            password_hash = NULL,
            password_salt = NULL,
            username = NULL
        WHERE id = ?;";
    $db->prepare($query);
    $db->bind_param('i', $user['id']);
    if(!$db->execute()){
        return Status::errorStatus("User could not be deleted.");
    }

    return Status::successStatus("success");
  }

  function delete_user_group($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("insufficient permissions");
    }

    // get all the values from the query
    $user_group = json_decode(file_get_contents('php://input'));

    // sanitize
    $sanitizer = new Sanitizer();
    if(! $user_group->user_group_id or !$sanitizer->isInt($user_group->user_group_id)){
        return Status::errorStatus("No valid user_group_id given");
    }

    // Delete the user group
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        error_log("delete_user_group: not able to connect to database");
        return Status::errorStatus("Cannot connect to the database");
    }

    // Check that no user is assigned to the user group
    $query = "SELECT count(*) as count FROM user WHERE status = ?";
    $db->prepare($query);
    $db->bind_param('i',
        $user_group->user_group_id
    );
    if(! $db->execute()){
        $db->disconnect();
        return Status::errorStatus("Cannot check whether users are still assigned to this user group.");
    }
    $res = $db->fetch_stmt_hash();
    if($res[0]['count'] > 0){
        return Status::errorStatus("Cannot delete group, still in use");
    }

    // Delete the price of this group
    $query = "DELETE FROM pricing WHERE user_status_id = ?";
    $db->prepare($query);
    $db->bind_param('i',
        $user_group->user_group_id
    );
    if(! $db->execute()){
        $db->disconnect();
        return Status::errorStatus("Cannot delete pricing entry related to the user group that should be deleted.");
    }

    // Delete the user group
    $query = "DELETE FROM user_status WHERE id = ?";
    $db->prepare($query);
    $db->bind_param('i',
        $user_group->user_group_id
    );
    if(! $db->execute()){
        $db->disconnect();
        return Status::errorStatus("Cannot delete user group");
    }

    $db->disconnect();
    return Status::successStatus("User group deleted");

  }

  function save_user_group($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("Insufficient permissions");
    }

    // get all the values from the query
    $user_group = json_decode(file_get_contents('php://input'));

    // two cases:
    // 1. new user group
    // 2. update to an existing user group
    if(! isset($user_group->user_group_id)){
        return create_user_group($configuration, $user_group);
    }else{
        return update_user_group($configuration, $user_group);
    }
  }

  function create_user_group($configuration, $data){

    // sanitize data
    $sanitizer = new Sanitizer();
    if(!$data->user_group_name or !$sanitizer->isAsciiText($data->user_group_name)){
        return Status::errorStatus("No valid user_group_name provided");
    }
    if(!$data->user_group_description){
        return Status::errorStatus("No valid user_group_description");
    }
    if(!$data->user_role_id or !$sanitizer->isUserRoleId($data->user_role_id)){
        return Status::errorStatus("No valid user_role_id selected");
    }
    if(!$data->price_min or !$sanitizer->isFloat($data->price_min)){
        return Status::errorStatus("No valid user_role_id selected");
    }
    if(! isset($data->price_description)){
        return Status::errorStatus("No valid price_description");
    }
    
    // Add the user group
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        error_log("create_user_group: not able to connect to database");
        return Status::errorStatus("Cannot connect to database");
    }

    // Create the user group
    $query = "INSERT INTO user_status 
              (name, description, user_role_id)
              VALUES
              (?,?,?)";
    $db->prepare($query);
    $db->bind_param('ssd',
        $data->user_group_name,
        $data->user_group_description,
        $data->user_role_id
    );
    if(! $db->execute()){
        $db->disconnect();
        return Status::errorStatus("Cannot add new group due an an unkown error.");
    }

    $id = $db->last_id();
    error_log("ID of last insert: " . $id);

    // Create the price for that user group
    $query = "INSERT INTO pricing
              (user_status_id, price_chf_min, comment)
              VALUES
              (?,?,?);";
    $db->prepare($query);
    $db->bind_param('ids',
        $id,
        $data->price_min,
        $data->price_description
    );
    if(! $db->execute()){
        $db->disconnect();
        return Status::errorStatus("Cannot add new user group due to an unknown error");
    }

    return Status::successStatus("User group created");
  }

  function update_user_group($configuration, $data){
    $user_group_price_update_status = update_user_group_price($configuration, $data);
    if($user_group_price_update_status['ok'] != TRUE){
        return $user_group_price_update_status;
    }
    
    // sanitize data
    $sanitizer = new Sanitizer();
    if(! $data->user_group_id or !$sanitizer->isInt($data->user_group_id)){
        return Status::errorStatus("No valid user_group_id given");
    }
    if(! $data->user_role_id or !$sanitizer->isUserRoleId($data->user_role_id)){
        return Status::errorStatus("No valid user_role_id given");
    }

    // update the user group
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        error_log("update_user_group_role: not able to connect to database");
        return Status::errorStatus("Cannot connect to the database");
    }

    $query = "UPDATE user_status
                SET name = ?,
                    description = ?,
                    user_role_id = ?
                WHERE id = ?";
    $db->prepare($query);
    $db->bind_param('ssdd',
        $data->user_group_name,
        $data->user_group_description,
        $data->user_role_id,
        $data->user_group_id
    );
    if($db->execute()){
        $db->disconnect();
        return Status::successStatus("user group updated");
    }
    $db->disconnect();
    return Status::errorStatus("could not update user group");
  }

  function update_user_group_price($configuration, $data){
    // sanitize data
    $sanitizer = new Sanitizer();
    if(! $data->price_id or !$sanitizer->isInt($data->price_id)){
        return Status::errorStatus("No valid price_id given");
    }
    if(! $data->price_min or !$sanitizer->isFloat($data->price_min)){
        return Status::errorStatus("No valid price_min given");
    }

    // update the price
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        error_log("update_user_group_price: not able to connect to database");
        return Status::errorStatus("Cannot connect to database");
    }

    $query = "UPDATE pricing 
              SET price_chf_min = ?, 
                  comment = ?
              WHERE id = ?;";
    $db->prepare($query);
    $db->bind_param('dsd', 
        $data->price_min, 
        $data->price_description,
        $data->price_id
    );
    if($db->execute()){
        $db->disconnect();
        return Status::successStatus("updated user group price");
    }
    $db->disconnect();
    return Status::errorStatus("could not update user group price");
  }

  function get_user_roles($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("Not sufficient permissions");
    }

    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("get_user_roles: not able to connect to database");
        return Status::errorStatus("Cannot connect to the database");
    }

    $ret = Array();

    // get user roles
    $query = "SELECT ur.id AS user_role_id,
                     ur.name AS user_role_name,
                     ur.description AS user_role_description
              FROM user_role ur";
    $db->prepare($query);
    $db->execute();
    $res = $db->fetch_stmt_hash();

    $i = 0;
    foreach ($res as $row){
        $i = $row['user_role_id'];
        $ret[$i] = $row;
    }

    return Status::successDataResponse("success", $ret);
  }

  // Returns all user types from the database
  function get_user_groups($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("API call not allowed");
    }

    $db = new DBAccess($configuration);
    if(!$db->connect()){
        error_log("get_user_types: not able to connect to database");
        return Status::errorStatus("Cannot connect to database");
    }

    $ret = Array();

    // get user types
    $query = "SELECT ur.id AS user_role_id, 
                     ur.name AS user_role_name, 
                     ur.description AS user_role_description, 
                     us.id AS user_group_id, 
                     us.name AS user_group_name, 
                     us.description AS user_group_description, 
                     p.id AS price_id, 
                     p.price_chf_min AS price_min,
                     p.comment AS price_description
              FROM user_role ur
              JOIN user_status us ON us.user_role_id = ur.id
              JOIN pricing p ON p.user_status_id = us.id";
    $db->prepare($query);
    $db->execute();
    $res = $db->fetch_stmt_hash();

    $i = 0;
    foreach ($res as $row){
        $i = $row['user_group_id'];
        $ret[$i]['user_group_id'] = $row['user_group_id'];
        $ret[$i]['user_group_name'] = $row['user_group_name'];
        $ret[$i]['user_group_description'] = $row['user_group_description'];
        $ret[$i]['user_role_id'] = $row['user_role_id'];
        $ret[$i]['user_role_name'] = $row['user_role_name'];
        $ret[$i]['user_role_description'] = $row['user_role_description'];
        $ret[$i]['price_id'] = $row['price_id'];
        $ret[$i]['price_min'] = $row['price_min'];
        $ret[$i]['price_description'] = $row['price_description'];
    }

    return Status::successDataResponse("success", $ret);
  }
  
  /* Returns all the details of all users (admin view) */
  function get_all_users_detailed($configuration, $lc){
	// only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        return Status::errorStatus("not authorized for this API call");
    }
	
	$db = new DBAccess($configuration);
	if(!$db->connect()){
        error_log("get_all_user_detailed: not able to connect to database");
        return Status::errorStatus("not able to connect to database");
	}
	
	$ret = Array();
    
    # get user data
    $query = "SELECT id, username, first_name, last_name,   
	                 address, city, plz, mobile, email,
					 license, status, locked, comment 
			  FROM user WHERE deleted = 0";
	$db->prepare($query);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	
	$i = 0;
	foreach ($res as $row){
	    $i = $row['id'];
		$ret[$i]['id'] = $row['id'];
        $ret[$i]['username']   = $row['username'];
        $ret[$i]['first_name'] = $row['first_name'];
        $ret[$i]['last_name']  = $row['last_name'];
        $ret[$i]['address']    = $row['address'];
        $ret[$i]['city']       = $row['city'];
        $ret[$i]['plz']        = $row['plz'];
        $ret[$i]['mobile']     = $row['mobile'];
        $ret[$i]['email']      = $row['email'];        
        $ret[$i]['license']    = $row['license'];
        $ret[$i]['locked']     = $row['locked'];
        $ret[$i]['status']     = $row['status'];
        $ret[$i]['comment']    = $row['comment'];
        $ret[$i]['total_heat_cost']   = 0;        
        $ret[$i]['total_heat_min']    = 0;
        $ret[$i]['total_payment']     = 0;
		$i++;
	}
    
    # get cost info from heats
    $query = "SELECT h.user_id as id, 
            sum(h.cost_chf) as cost, 
            sum(h.duration_s) as time
        FROM heat h 
        JOIN user u ON h.user_id = u.id 
        WHERE u.deleted = 0 
        GROUP BY h.user_id";
	$db->prepare($query);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	foreach ($res as $row){
        $id = $row['id'];
        $ret[$id]['total_heat_cost'] = floatval($row['cost']);
        $ret[$id]['total_heat_min'] = intval($row['time']/60);
	}
	    
    # get payment info
    $query = "SELECT p.user_id as id, 
            sum(p.amount_chf) as pay
        FROM payment p
        JOIN user u ON u.id = p.user_id
        WHERE p.type_id = 4 
        AND u.deleted = 0
        GROUP BY p.user_id";
	$db->prepare($query);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	foreach ($res as $row){
		$id = $row['id'];
		if(!isset($ret[$id])){
			error_log('lib_user: Payment from a user that does not exist: ' . $row['pay'] . ' ' . $configuration->currency . ' from user with ID: ' . $row['id']);
			continue;
		}
		$ret[$id]['total_payment'] = floatval($row['pay']);
	}
	
	# get payback info
	$query = "SELECT e.user_id as id,
            sum(e.amount_chf) as payback
        FROM expenditure e
        JOIN user u ON u.id = e.user_id
        WHERE e.type_id = 4
        AND u.deleted = 0 
        GROUP BY e.user_id";
	$db->prepare($query);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	foreach ($res as $row){
		$id = $row['id'];
		if(!isset($ret[$id])){
			error_log('lib_user: Packback to a user that does not exist: ' . $row['payback'] . ' ' . $configuration->currency . ' to user with ID: ' . $row['id']);
			continue;
        }
        error_log(print_r($row, TRUE));
        if(isset($ret[$id]['total_payment'])){
            $ret[$id]['total_payment'] = $ret[$id]['total_payment'] - floatval($row['payback']);
        }else{
            $ret[$id]['total_payment'] = floatval($row['payback']) * -1;
        }
    }
    
    $response = array();
    $response['users'] = $ret;
    $response['currency'] = $configuration->currency;
    
    return Status::successDataResponse("success", $response);
  }
  
  /* Returns all heats of a user */
  function get_my_user_heats($configuration){
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("get_my_user_heats: not able to connect to database");
        exit;
    }    
    $res = Array();
    
    // get the total minutes and costs
    $session = Login::getSessionData($configuration);
    $query = "SELECT sum(duration_s) as time, sum(cost_chf) as cost
              FROM heat 
              WHERE user_id = ?";
    $db->prepare($query);
    $db->bind_param('i', $session['user_id']);
    $db->execute();
    $sum = $db->fetch_stmt_hash();
    if(count($sum)>0){
        $res['heat_time_min'] = (int) ($sum[0]['time'] / 60);
        $res['heat_cost'] = round(floatval($sum[0]['cost']), 2);
    }

    // get the total minutes and costs (YTD)
    $query = "SELECT sum(duration_s) as time, sum(cost_chf) as cost
              FROM heat 
              WHERE year(heat.timestamp) = year(now())
                AND user_id = ?";
    $db->prepare($query);
    $db->bind_param('i', $session['user_id']);
    $db->execute();
    $sum = $db->fetch_stmt_hash();
    if(count($sum)>0){
        $res['heat_time_min_ytd'] = (int) ($sum[0]['time'] / 60);
        $res['heat_cost_ytd'] = round(floatval($sum[0]['cost']), 2);
    }
    
    // get the actual balance of the user
    $query = "SELECT sum(amount_chf) as total 
              FROM payment WHERE user_id = ?
              AND type_id = 4;";
    $db->prepare($query);
    $db->bind_param('i', $session['user_id']);
    $db->execute();
    $payment = $db->fetch_stmt_hash();
    if(count($payment)>0){
        $res['payment_total']   = round(floatval($payment[0]['total']), 2);
        $res['balance_current'] = $res['payment_total'] - $res['heat_cost'];
    }
    
    // deduct the payback from the user's balance_current
    $query = "SELECT sum(amount_chf) as payback
              FROM expenditure 
              WHERE user_id = ? AND type_id = 4;";
    $db->prepare($query);
    $db->bind_param('i', $session['user_id']);
    $db->execute();
    $payback = $db->fetch_stmt_hash();
    if(count($payback)>0){
        $res['payback_total']   = round(floatval($payback[0]['payback']), 2);
        $res['balance_current'] = $res['balance_current'] - $res['payback_total'];
    }
    
    // Get the user's last heats
    $query = "SELECT timestamp, duration_s, cost_chf
              FROM heat WHERE user_id = ? 
              ORDER BY timestamp DESC LIMIT 100;";
    $db->prepare($query);
    $db->bind_param('i', $session['user_id']);
    $db->execute();
    $heats = $db->fetch_stmt_hash();
    for($i=0; $i<count($heats); $i++){
        $res['heats'][$i]['date'] = date('d.m.Y', strtotime($heats[$i]['timestamp']));
        $res['heats'][$i]['cost'] = sprintf('%.2f', $heats[$i]['cost_chf']);
        $min = floor($heats[$i]['duration_s'] / 60);
        $sec = sprintf('%02d', round($heats[$i]['duration_s'] % 60));
        $res['heats'][$i]['duration'] =  $min . ':' . $sec;
    }

    $res['currency'] = $configuration->currency;
    $res['balance_current'] = round($res['balance_current'], 2);
    
    echo json_encode($res);    
  }
  
  # lock a user
  function lock_user($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("permission denied");
    }
    
    $post_data = json_decode(file_get_contents('php://input'));
    
    // sanitize input
    $sanitizer = new Sanitizer();
    if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
        return Status::errorStatus("No valid user_id provided, please select a user");
    }
    
    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        return Status::errorStatus("Cannot connect to the database");
    }
    
    $query = 'UPDATE user SET locked = 1 WHERE id = ?;';
    $db->prepare($query);
    $db->bind_param('d', $post_data->user_id);
    if(!$db->execute()){
        $db->disconnect();
        return Status::errorStatus("Attempt to lock user failed.");
    }
    $db->disconnect();
    return Status::successStatus("success");
  }
  
  # unlock a user
  function unlock_user($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("not sufficient permissions");
    }
    
    $post_data = json_decode(file_get_contents('php://input'));
    
    // sanitize input
    $sanitizer = new Sanitizer();
    if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
        return Status::errorStatus("No valid user_id selected, please select a user");
    }
    
    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        return Status::errorStatus("Cannot connect to database");
    }
    
    $query = 'UPDATE user SET locked = 0 WHERE id = ?;';
    $db->prepare($query);
    $db->bind_param('d', $post_data->user_id);
    if(!$db->execute()){
        $db->disconnect();
        return Status::errorStatus("Attempt to unlock user failed");
    }
    $db->disconnect();
    return Status::successStatus("success");
  }
  
  # change the status of a user
  # e.g. to 'member', 'admin', 'guest'
  function change_user_group_membership($configuration, $lc){
      // only admins are allowed to call this function
      if(!$lc->isAdmin()){
          return Status::errorStatus("Not sufficient permissions");
      }
      
      $post_data = json_decode(file_get_contents('php://input'));
    
      // sanitize input
      $sanitizer = new Sanitizer();
      if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
          return Status::errorStatus("No valid user selected, please select a user");
      }
      if(!$post_data->status_id or !$sanitizer->isInt($post_data->status_id)){
          return Status::errorStatus("No valid status selected, please select a status");
      }
      
      // connect to the database
      $db = new DBAccess($configuration);
      if(!$db->connect()){
        return Status::errorStatus("Cannot connect to database");
      }
    
      $query = 'UPDATE user SET status = ? WHERE id = ?;';
      $db->prepare($query);
      $db->bind_param('dd', 
          $post_data->status_id, 
          $post_data->user_id
      );
      if(!$db->execute()){
          $db->disconnect();
          return Status::errorStatus("Attempt to update a user's group has failed");
      }
      $db->disconnect();
      return Status::successStatus("success");
  }
  
  /* Returns all users */
  function get_all_users($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("No authorization");
    }
    
    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        return Status::errorStatus("Cannot connect to database");
    }
    
    $query = 'SELECT 
        id as id, 
        first_name as first_name, 
        last_name as last_name 
        FROM user 
        WHERE deleted = 0
        ORDER BY first_name, last_name;';
    $res = $db->fetch_data_hash($query);
    $db->disconnect();
    if(!$res){
        return Status::errorStatus("Cannot get users from database");
    }
    return Status::successDataResponse("Users retrieved", $res);
  }
  
  /* Returns the admin users*/
  function get_admin_users($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }
   
    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        exit;
    }
    
    $query = 'SELECT u.id AS id, u.first_name AS first_name, u.last_name AS last_name
        FROM user u
        JOIN user_status us ON u.status = us.id
        JOIN user_role ur ON us.user_role_id = ur.id
        WHERE ur.name = "admin" 
        AND u.deleted = 0
        ORDER BY u.first_name, u.last_name;';
    $res = $db->fetch_data_hash($query, 0);
    $db->disconnect();
    if(!$res){
        HttpHeader::setResponseCode(500);
    }
    echo json_encode($res);
  }

  function get_session_users($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }

    $data = json_decode(file_get_contents('php://input'));
    $sanitizer = new Sanitizer();
    if(!isset($data->session_id) or !$sanitizer->isInt($data->session_id)){
        error_log('api/user.php: No session_id provided: ' . $data->session_id);
        echo 'No valid session_id provided';
    }

    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        exit;
    }
    
    $query = 'SELECT u.id AS id, u.first_name AS first_name, u.last_name AS last_name
        FROM user u
        JOIN user_to_session us ON u.id = us.user_id
        WHERE us.session_id = ? 
        AND u.deleted = 0
        ORDER BY u.first_name, u.last_name;';
    $db->prepare($query);
    $db->bind_param('i', $data->session_id);
    $db->execute();
    $res = $db->fetch_stmt_hash();
    if(!isset($res)){
        HttpHeader::setResponseCode(500);
        exit;
    }
    $db->disconnect();
    echo json_encode($res);
  }
  
  function get_my_user($configuration){
  
    // sanitize cookie
    $sanitize = new Sanitizer();
    if(!$sanitize->isCookie($_COOKIE['SESSION'])){
        HttpHeader::setResponseCode(403);
        exit;
    }
  
    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        exit;
    }
    
    $query = 'SELECT u.id, u.username, u.first_name, u.last_name,
                     u.address, u.city, u.plz, u.mobile, u.email,
                     u.license, u.status, u.locked, u.comment, ur.id as user_role_id, ur.name as user_role_name
                FROM user u, browser_session bs, user_status us, user_role ur 
                WHERE bs.session_secret = ?
                  AND bs.user_id = u.id
                  AND u.status = us.id
                  AND us.user_role_id = ur.id;';
    $db->prepare($query);
    $db->bind_param('s', $_COOKIE['SESSION'] );
    $db->execute();
    $res = $db->fetch_stmt_hash();
    $db->disconnect();
    if(!$res){
        HttpHeader::setResponseCode(500);
    }
    echo json_encode($res[0]);
  }
  
  // Change a user's personal data
  function change_my_user_data($configuration){
    $data = json_decode(file_get_contents('php://input'));
    
    // input validation
    $sanitizer = new Sanitizer();
    if(! isset($data->first_name)){
        return Status::errorStatus("No first name provided.");
    }
    if(! isset($data->last_name)){
        return Status::errorStatus("No last name provided.");
    }
    if(! isset($data->address)){
        return Status::errorStatus("No address provided.");
    }
    if(!isset($data->mobile) or !$sanitizer->isMobileNumber($data->mobile)){
        return Status::errorStatus("The mobile number is not correct.");
    }
    if(! isset($data->plz) or !$sanitizer->isInt($data->plz)){
        return Status::errorStatus("The zip code is not of a valid format.");
    }
    if(! isset($data->city)){
        return Status::errorStatus("No city provided.");
    }
    if(! isset($data->email) or !$sanitizer->isEmail($data->email)){
        return Status::errorStatus("The email address is not correct.");
    }
    if(! isset($data->license) or !$sanitizer->isBoolean($data->license)){
        return Status::errorStatus("Invalid input for the license information.");
    }
    
    // get the user info
    $user = new User($configuration);
    $user_data = $user->getUser();
    
    // update the user in the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        return Status::errorStatus("Backend error: Cannot connect to the database");
    }
    $query = "UPDATE user SET
                      first_name = ?, 
                      last_name = ?, 
                      address = ?, 
                      city = ?, 
                      plz = ?, 
                      mobile = ?, 
                      email = ?, 
                      license = ?
                      WHERE id = ?;";
    $db->prepare($query);
    $db->bind_param('ssssissii',
                    $data->first_name,
                    $data->last_name,
                    $data->address,
                    $data->city,
                    $data->plz,
                    $data->mobile,
                    $data->email,
                    $data->license,
                    $user_data['id']);
    if(!$db->execute()){
        $db->disconnect();
        return Status::errorStatus("Backend error: Cannot update user data");
    }    
    $db->disconnect();
    
    return Status::successStatus("Successfully updated.");
  }
  
  function get_my_user_sessions($configuration){
    // the upcoming sessions
    $res['sessions']     = Array();
    // the past sessions
    $res['sessions_old'] = Array();
    
    // get the user's ID
    $session_data = Login::getSessionData($configuration);
    
    // database connection
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log('api/user.php: Cannot open database connection');
        return;
    }
    
    // select all sessions where the user is either creator, booked in session
    // invited or has to drive the course
    $query = 'SELECT DISTINCT a_s.id AS id, 
                    UNIX_TIMESTAMP(a_s.start_time) AS start_time, 
                    UNIX_TIMESTAMP(a_s.end_time) AS end_time, 
                    a_s.title AS title, 
                    a_s.type AS type,
                    a_s.status AS status,
                    a_s.start_time AS native_start_time
              FROM 
              (SELECT a.id, a.start_time, a.end_time, a.title, a.type, 3 as status
               FROM `session` a 
               WHERE a.creator_id = ? AND a.type = 1
               
               UNION ALL 
               SELECT  b.id, b.start_time, b.end_time, b.title, b.type, 3 as status
               FROM `session` b, `user_to_session` u2s  
               WHERE b.id = u2s.session_id AND u2s.user_id = ?
               
               UNION ALL
               SELECT c.id, c.start_time, c.end_time, c.title, c.type, inv.status
               FROM session c, invitation inv
               WHERE c.id = inv.session_id AND inv.user_id = ?
               AND (inv.status = 1 OR inv.status = 2)
              ) a_s ORDER BY a_s.start_time DESC LIMIT 100';
    $db->prepare($query);
    $db->bind_param('iii',
                    $session_data['user_id'],
                    $session_data['user_id'],
                    $session_data['user_id']);
    $db->execute();
    $db_res = $db->fetch_stmt_hash();
    if(!isset($db_res) or $db_res === FALSE){
        error_log('api/user.php: Cannot retrieve sessions for user: ' . $session_data['user_id']);
        HttpHeader::setResponseCode(500);
        echo "Cannot retrieve sessions";
        return;
    }
    
    $old_idx = 0;
    $new_idx = 0;
    for($i=0; $i<count($db_res); $i++){
        if($db_res[$i]['start_time'] < time()){
            $res['sessions_old'][$old_idx]['id']     = $db_res[$i]['id'];
            $res['sessions_old'][$old_idx]['title']  = utf8_decode($db_res[$i]['title']);
            $res['sessions_old'][$old_idx]['type']   = $db_res[$i]['type'];
            $res['sessions_old'][$old_idx]['status'] = $db_res[$i]['status'];
            $res['sessions_old'][$old_idx]['start']  = $db_res[$i]['start_time'];
            $res['sessions_old'][$old_idx]['end']    = $db_res[$i]['end_time'];
            $old_idx++;
        }else{
            $res['sessions'][$new_idx]['id']         = $db_res[$i]['id'];
            $res['sessions'][$new_idx]['title']      = utf8_decode($db_res[$i]['title']);
            $res['sessions'][$new_idx]['type']       = $db_res[$i]['type'];
            $res['sessions'][$new_idx]['status']     = $db_res[$i]['status'];
            $res['sessions'][$new_idx]['start']      = $db_res[$i]['start_time'];
            $res['sessions'][$new_idx]['end']        = $db_res[$i]['end_time'];
            $new_idx++;
        }
    }
    
    // get the confirmed riders for every session
    // get the riders for every upcoming session
    for($i = 0; $i< count($res['sessions']); $i++){
        $query = 'SELECT u.id AS id, u.first_name AS fn, u.last_name AS ln
                  FROM user_to_session u2s, user u
                  WHERE u2s.session_id = ? AND u2s.user_id = u.id';
        $db->prepare($query);

        $session_id = $res['sessions'][$i]['id'];
        $res['sessions'][$i]['riders'] = Array();       
        
        $db->bind_param('i', $session_id);
        $db->execute();
        $db_res = $db->fetch_stmt_hash();
        if(!isset($db_res) or $db_res === FALSE){
            error_log('api/user.php: Cannot retrieve riders for session: ' . $session_id);
        }else{
            for($j=0; $j<count($db_res); $j++){
                $res['sessions'][$i]['riders'][$j]['id']         = $db_res[$j]['id'];
                $res['sessions'][$i]['riders'][$j]['first_name'] = $db_res[$j]['fn'];
                $res['sessions'][$i]['riders'][$j]['last_name']  = $db_res[$j]['ln'];
            }
        }
    }
    
    // get the riders for every old session
    for($i = 0; $i< count($res['sessions_old']); $i++){
        $query = 'SELECT u.id AS id, u.first_name AS fn, u.last_name AS ln
                  FROM user_to_session u2s, user u
                  WHERE u2s.session_id = ? AND u2s.user_id = u.id';
        $db->prepare($query);

        $session_id = $res['sessions_old'][$i]['id'];
        $res['sessions_old'][$i]['riders'] = Array();

        $db->bind_param('i', $session_id);
        $db->execute();
        $db_res = $db->fetch_stmt_hash();
        if(!isset($db_res) or $db_res === FALSE){
            error_log('api/user.php: Cannot retrieve riders for session: ' . $session_id);
        }else{
            for($j=0; $j<count($db_res); $j++){
                $res['sessions_old'][$i]['riders'][$j]['id']         = $db_res[$j]['id'];
                $res['sessions_old'][$i]['riders'][$j]['first_name'] = $db_res[$j]['fn'];
                $res['sessions_old'][$i]['riders'][$j]['last_name']  = $db_res[$j]['ln'];
            }
        }
    }
    
    $json = json_encode($res, JSON_PARTIAL_OUTPUT_ON_ERROR, 0);
    if($json == FALSE){
        error_log('api/user.php: Cannot encode result to json');
        HttpHeader::setResponseCode(500);
        echo "Cannot retrieve sessions (parsing error)";
        return;
    }

    echo $json;
    
  }
?>