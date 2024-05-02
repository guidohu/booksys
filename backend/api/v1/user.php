<?php
  
  // automatically load all classes
  spl_autoload_register('user_autoloader');
  function user_autoloader($class){
    include '../../classes/'.$class.'.php';
  }

  // Load configuration
  $configuration = new Configuration();
  
  // Check if the user is already logged in
  $lc = new Login($configuration);
  if(!$lc->isLoggedIn()){
      echo json_encode(Status::errorStatus("login required"));
      return;
  }
  
  // check if we have an action
  if(!isset($_GET['action'])){
    HttpHeader::setResponseCode(200);
    exit;
  }
  
  $_SESSION = Login::getSessionData($configuration);
  
  $response = null;

  switch($_GET['action']){
    case 'save_user_group':
        $response = save_user_group($configuration, $lc);
        break;
    case 'delete_user_group':
        $response = delete_user_group($configuration, $lc);
        break;
    case 'delete_user':
        $response = delete_user($configuration, $lc);
        break;
    default:
        HttpHeader::setResponseCode(400);
        $response = Status::errorStatus("Action not supported");
        break;
  }
  
  echo json_encode($response);
  return;

  function delete_user($configuration, $lc){
    // user needs to be admin
    if(!$lc->isAdmin()){
        return Status::errorStatus("insufficient permissions");
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
        return Status::errorStatus("insufficient permissions");
    }

    // get all the values from the query
    $user_group = json_decode(file_get_contents('php://input'));

    // two cases:
    // 1. new user group
    // 2. update to an existing user group
    if(! isset($user_group->user_group_id)){
        return create_user_group($configuration, $user_group, $lc);
    }else{
        return update_user_group($configuration, $user_group, $lc);
    }
  }

  function create_user_group($configuration, $data, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("insufficient permissions");
    }

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
    if(!isset($data->price_min) or !$sanitizer->isFloat($data->price_min)){
        return Status::errorStatus("No valid price value provided: '$data->price_min'");
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

  function update_user_group($configuration, $data, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("insufficient permissions");
    }

    $user_group_price_update_status = update_user_group_price($configuration, $data, $lc);
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

  function update_user_group_price($configuration, $data, $lc){
    // only admins are allowed to call this function
    if(!$lc->isAdmin()){
        return Status::errorStatus("insufficient permissions");
    }

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
  
//   function get_my_user($configuration){
  
//     // sanitize cookie
//     $sanitize = new Sanitizer();
//     if(!$sanitize->isCookie($_COOKIE['SESSION'])){
//         return Status::errorStatus("invalid session cookie format");
//     }
  
//     // connect to the database
//     $db = new DBAccess($configuration);
//     if(!$db->connect()){
//         return Status::errorStatus("Cannot connect to database");
//     }
    
//     $query = 'SELECT u.id, u.username, u.first_name, u.last_name,
//                      u.address, u.city, u.plz, u.mobile, u.email,
//                      u.license, u.status, u.locked, u.comment, ur.id as user_role_id, ur.name as user_role_name
//                 FROM user u, browser_session bs, user_status us, user_role ur 
//                 WHERE bs.session_secret = ?
//                   AND bs.user_id = u.id
//                   AND u.status = us.id
//                   AND us.user_role_id = ur.id;';
//     $db->prepare($query);
//     $db->bind_param('s', $_COOKIE['SESSION'] );
//     $db->execute();
//     $res = $db->fetch_stmt_hash();
//     $db->disconnect();
//     if(!$res){
//         return Status::errorStatus("Cannot get user information");
//     }

//     return Status::successDataResponse("success", $res[0]);
//   }
?>