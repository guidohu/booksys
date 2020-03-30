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
        get_all_users($configuration, $lc);
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
        change_my_user_data($configuration);
        exit;
    case 'lock_user':
        lock_user($configuration, $lc);
        exit;
    case 'unlock_user':
        unlock_user($configuration, $lc);
        exit;
    case 'change_user_group_membership':
        change_user_group_membership($configuration, $lc);
        exit;
	case 'get_all_user_detailed':
		get_all_user_detailed($configuration, $lc);
        exit;
    case 'get_user_groups':
        get_user_groups($configuration, $lc);
        exit; 
    case 'get_user_roles':
        get_user_roles($configuration, $lc);
        exit;
    case 'save_user_group':
        save_user_group($configuration, $lc);
        exit;
    case 'delete_user_group':
        delete_user_group($configuration, $lc);
        exit;
  }
  
  HttpHeader::setResponseCode(400);
  return;

  function delete_user_group($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }

    // get all the values from the query
    $user_group = json_decode(file_get_contents('php://input'));

    // sanitize
    $sanitizer = new Sanitizer();
    if(! $user_group->user_group_id or !$sanitizer->isInt($user_group->user_group_id)){
        HttpHeader::setResponseCode(400);
        echo "No valid user_group_id given";
        return FALSE;
    }

    // Delete the user group
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("delete_user_group: not able to connect to database");
        return FALSE;
    }

    // Check that no user is assigned to the user group
    $query = "SELECT count(*) as count FROM user WHERE status = ?";
    $db->prepare($query);
    $db->bind_param('i',
        $user_group->user_group_id
    );
    if(! $db->execute()){
        $db->disconnect();
        HttpHeader::setResponseCode(500);
        return FALSE;
    }
    $res = $db->fetch_stmt_hash();
    error_log(var_dump($res));
    if($res[0]['count'] > 0){
        HttpHeader::setResponseCode(400);
        echo "Cannot delete group, still in use";
        return FALSE;
    }

    // Delete the price of this group
    $query = "DELETE FROM pricing WHERE user_status_id = ?";
    $db->prepare($query);
    $db->bind_param('i',
        $user_group->user_group_id
    );
    if(! $db->execute()){
        $db->disconnect();
        HttpHeader::setResponseCode(500);
        return FALSE;
    }

    // Delete the user group
    $query = "DELETE FROM user_status WHERE id = ?";
    $db->prepare($query);
    $db->bind_param('i',
        $user_group->user_group_id
    );
    if(! $db->execute()){
        $db->disconnect();
        HttpHeader::setResponseCode(500);
        return FALSE;
    }

    $db->disconnect();
    HttpHeader::setResponseCode(200);
    return true;

  }

  function save_user_group($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }

    // get all the values from the query
    $user_group = json_decode(file_get_contents('php://input'));

    // two cases:
    // 1. new user group
    // 2. update to an existing user group
    if(! isset($user_group->user_group_id)){
        if(! create_user_group($configuration, $user_group)){
            return FALSE;
        }
    }else{
        if(! update_user_group($configuration, $user_group)){
            return FALSE;
        }
    }

    HttpHeader::setResponseCode(200);
    return TRUE;
  }

  function create_user_group($configuration, $data){

    // sanitize data
    $sanitizer = new Sanitizer();
    if(!$data->user_group_name or !$sanitizer->isAsciiText($data->user_group_name)){
        HttpHeader::setResponseCode(400);
        echo "No valid user_group_name selected";
        return;
    }
    if(!$data->user_group_description){
        HttpHeader::setResponseCode(400);
        echo "No valid user_group_description";
        return;
    }
    if(!$data->user_role_id or !$sanitizer->isUserRoleId($data->user_role_id)){
        HttpHeader::setResponseCode(400);
        echo "No valid user_role_id selected";
        return;
    }
    if(!$data->price_min or !$sanitizer->isFloat($data->price_min)){
        HttpHeader::setResponseCode(400);
        echo "No valid user_role_id selected";
        return;
    }
    if(! isset($data->price_description)){
        HttpHeader::setResponseCode(400);
        echo "No valid price_description";
        return;
    }
    
    // Add the user group
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("create_user_group: not able to connect to database");
        return FALSE;
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
        HttpHeader::setResponseCode(500);
        return FALSE;
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
        HttpHeader::setResponseCode(500);
        return FALSE;
    }

    return TRUE;
  }

  function update_user_group($configuration, $data){
    if(! update_user_group_price($configuration, $data)){
        return FALSE;
    }
    
    // sanitize data
    $sanitizer = new Sanitizer();
    if(! $data->user_group_id or !$sanitizer->isInt($data->user_group_id)){
        HttpHeader::setResponseCode(400);
        echo "No valid user_group_id given";
        return FALSE;
    }
    if(! $data->user_role_id or !$sanitizer->isUserRoleId($data->user_role_id)){
        HttpHeader::setResponseCode(400);
        echo "No valid user_role_id given";
        return FALSE;
    }

    // update the user group
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("update_user_group_role: not able to connect to database");
        return FALSE;
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
        return TRUE;
    }
    $db->disconnect();
    HttpHeader::setResponseCode(500);
    return FALSE;
  }

  function update_user_group_price($configuration, $data){
    // sanitize data
    $sanitizer = new Sanitizer();
    if(! $data->price_id or !$sanitizer->isInt($data->price_id)){
        HttpHeader::setResponseCode(400);
        echo "No valid price_id given";
        return FALSE;
    }
    if(! $data->price_min or !$sanitizer->isFloat($data->price_min)){
        HttpHeader::setResponseCode(400);
        echo "No valid price_min given";
        return FALSE;
    }

    // update the price
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("update_user_group_price: not able to connect to database");
        return FALSE;
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
        return TRUE;
    }
    $db->disconnect();
    HttpHeader::setResponseCode(500);
    return FALSE;
  }

  function get_user_roles($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }

    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("get_user_roles: not able to connect to database");
        exit;
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

    echo json_encode($ret);
  }

  // Returns all user types from the database
  function get_user_groups($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }

    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        error_log("get_user_types: not able to connect to database");
        exit;
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

    $response = array();
    $response["user_groups"] = $ret;
    $response["currency"]    = $configuration->currency;

    echo json_encode($response);
  }
  
  /* Returns all the details of all users (admin view) */
  function get_all_user_detailed($configuration, $lc){
	// only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }
	
	$db = new DBAccess($configuration);
	if(!$db->connect()){
		HttpHeader::setResponseCode(500);
        error_log("get_all_user_detailed: not able to connect to database");
        exit;
	}
	
	$ret = Array();
    
    # get user data
    $query = "SELECT id, username, first_name, last_name,   
	                 address, city, plz, mobile, email,
					 license, status, locked, comment 
			  FROM user";
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
    $query = "SELECT user_id as id, 
                sum(cost_chf) as cost, 
                sum(duration_s) as time
              FROM heat GROUP BY user_id";
	$db->prepare($query);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	foreach ($res as $row){
		$id = $row['id'];
		$ret[$id]['total_heat_cost'] = floatval($row['cost']);
		$ret[$id]['total_heat_min'] = intval($row['time']/60);
	}
	    
    # get payment info
    $query = "SELECT user_id as id, 
               sum(amount_chf) as pay
              FROM payment WHERE type_id = 4 GROUP BY user_id";
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
	$query = "SELECT user_id as id,
	          sum(amount_chf) as payback
			FROM expenditure WHERE type_id = 4 GROUP BY user_id";
	$db->prepare($query);
	$db->execute();
	$res = $db->fetch_stmt_hash();
	foreach ($res as $row){
		$id = $row['id'];
		if(!isset($ret[$id])){
			error_log('lib_user: Packback to a user that does not exist: ' . $row['payback'] . ' ' . $configuration->currency . ' to user with ID: ' . $row['id']);
			continue;
		}
		$ret[$id]['payment'] = $ret[$id]['payment'] - floatval($row['payback']);
    }
    
    $response = array();
    $response['users'] = $ret;
    $response['currency'] = $configuration->currency;
    
    echo json_encode($response);
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
        $res['heat_cost'] = (int) $sum[0]['cost'];
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
        $res['heat_cost_ytd'] = (int) $sum[0]['cost'];
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
        $res['payment_total'] = $payment[0]['total'];
        $res['balance_current'] = $payment[0]['total'] - $res['heat_cost'];
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
        $res['payback_total']   = $payback[0]['payback'];
        $res['balance_current'] = $res['balance_current'] - $payback[0]['payback'];
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
    
    echo json_encode($res);    
  }
  
  # lock a user
  function lock_user($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }
    
    $post_data = json_decode(file_get_contents('php://input'));
    
    // sanitize input
    $sanitizer = new Sanitizer();
    if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
        HttpHeader::setResponseCode(400);
        echo "No valid user selected, please select a user";
        return;
    }
    
    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        exit;
    }
    
    $query = 'UPDATE user SET locked = 1 WHERE id = ?;';
    $db->prepare($query);
    $db->bind_param('d', $post_data->user_id);
    if(!$db->execute()){
        $db->disconnect();
        HttpHeader::setResponseCode(500);
        return;
    }
    $db->disconnect();
    HttpHeader::setResponseCode(200);
  }
  
  # unlock a user
  function unlock_user($configuration, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        HttpHeader::setResponseCode(403);
        exit;
    }
    
    $post_data = json_decode(file_get_contents('php://input'));
    
    // sanitize input
    $sanitizer = new Sanitizer();
    if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
        HttpHeader::setResponseCode(400);
        echo "No valid user selected, please select a user";
        return;
    }
    
    // connect to the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        exit;
    }
    
    $query = 'UPDATE user SET locked = 0 WHERE id = ?;';
    $db->prepare($query);
    $db->bind_param('d', $post_data->user_id);
    if(!$db->execute()){
        $db->disconnect();
        HttpHeader::setResponseCode(500);
        return;
    }
    $db->disconnect();
    HttpHeader::setResponseCode(200);
  }
  
  # change the status of a user
  # e.g. to 'member', 'admin', 'guest'
  function change_user_group_membership($configuration, $lc){
      // only admins are allowed to call this function
      if(!$lc->isAdmin()){
          HttpHeader::setResponseCode(403);
          exit;
      }
      
      $post_data = json_decode(file_get_contents('php://input'));
    
      // sanitize input
      $sanitizer = new Sanitizer();
      if(!$post_data->user_id or !$sanitizer->isInt($post_data->user_id)){
          HttpHeader::setResponseCode(400);
          echo "No valid user selected, please select a user";
          return;
      }
      if(!$post_data->status_id or !$sanitizer->isInt($post_data->status_id)){
          HttpHeader::setResponseCode(400);
          echo "No valid status selected, please select a status";
          return;
      }
      
      // connect to the database
      $db = new DBAccess($configuration);
      if(!$db->connect()){
          HttpHeader::setResponseCode(500);
          exit;
      }
    
      $query = 'UPDATE user SET status = ? WHERE id = ?;';
      $db->prepare($query);
      $db->bind_param('dd', 
          $post_data->status_id, 
          $post_data->user_id
      );
      if(!$db->execute()){
          $db->disconnect();
          HttpHeader::setResponseCode(500);
          return;
      }
      $db->disconnect();
      HttpHeader::setResponseCode(200);
  }
  
  /* Returns all users */
  function get_all_users($configuration, $lc){
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
    
    $query = 'SELECT 
        id as id, 
        first_name as first_name, 
        last_name as last_name 
        FROM user 
        ORDER BY first_name, last_name;';
    $res = $db->fetch_data_hash($query);
    $db->disconnect();
    if(!$res){
        HttpHeader::setResponseCode(500);
    }
    echo json_encode($res);
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
        HttpHeader::setResponseCode(400);
        echo "No first name provided.";
        return;
    }
    if(! isset($data->last_name)){
        HttpHeader::setResponseCode(400);
        echo "No last name provided";
        return;
    }
    if(! isset($data->address)){
        HttpHeader::setResponseCode(400);
        echo "No address given";
        return;
    }
    if(!isset($data->mobile) or !$sanitizer->isMobileNumber($data->mobile)){
        HttpHeader::setResponseCode(400);
        echo "The mobile number is not correct.";
        return;
    }
    if(! isset($data->plz) or !$sanitizer->isInt($data->plz)){
        HttpHeader::setResponseCode(400);
        echo "Your Zip Code is wrong";
        return;
    }
    if(! isset($data->city)){
        HttpHeader::setResponseCode(400);
        echo "No city provided";
        return;
    }
    if(! isset($data->email) or !$sanitizer->isEmail($data->email)){
        HttpHeader::setResponseCode(400);
        echo "Please check your Email again.";
        return;
    }
    if(! isset($data->license) or !$sanitizer->isBoolean($data->license)){
        HttpHeader::setResponseCode(400);
        echo "The License input is wrong.";
        return;
    }
    
    // get the user info
    $user = new User($configuration);
    $user_data = $user->getUser();
    
    // update the user in the database
    $db = new DBAccess($configuration);
    if(!$db->connect()){
        HttpHeader::setResponseCode(500);
        return;
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
        HttpHeader::setResponseCode(500);
        echo "Internal Server Error, cannot update user data";
        return;
    }    
    $db->disconnect();
    
    error_log("Changed user information for user_id: " . $user_data['id']);
    HttpHeader::setResponseCode(200);
    return;
  }
  
  function get_my_user_sessions($configuration){
    $res['sessions']     = Array();
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
    $query = 'SELECT a_s.id AS id, 
                    UNIX_TIMESTAMP(a_s.start_time) AS start_time, 
                    UNIX_TIMESTAMP(a_s.end_time) AS end_time, 
                    a_s.title AS title, 
                    a_s.type AS type,
                    a_s.status AS status
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