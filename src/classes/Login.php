<?php
spl_autoload_register('Login::autoloader');

abstract class LoginStatus {
	const Locked          = 0;
	const WrongPassword   = 1;
	const CorrectPassword = 2;
	const UnknownError    = 99;
}

// Browser Session
class SessionInfo {

	public $logged_in;
	public $user_id;
	public $user_status;
	public $locked;
	public $username;
	public $first_name;
	public $last_name;
	public $last_activity;
	public $remember;
	public $valid_thru;
	public $session_secret;
	public $user_role_id;
	public $user_role_name;

	public function __construct($dbh){
		$this->dbh = $dbh;
	}

	public function serialize(){
		$string_representation = array();
		$string_representation['logged_in']      = $this->logged_in;
		$string_representation['user_id']        = $this->user_id;
		$string_representation['user_status']    = $this->user_status;
		$string_representation['locked']         = $this->locked;
		$string_representation['username']       = $this->username;
		$string_representation['first_name']     = $this->first_name;
		$string_representation['last_name']      = $this->last_name;
		$string_representation['last_activity']  = $this->last_activity;
		$string_representation['remember']       = $this->remember;
		$string_representation['valid_thru']     = $this->valid_thru;
		$string_representation['session_secret'] = $this->session_secret;
		$string_representation['user_role_id']   = $this->user_role_id;
		$string_representation['user_role_name'] = $this->user_role_name;

		return json_encode($string_representation);
	}

	public function insert_session(){
		$query = 'INSERT INTO browser_session
						  (session_secret, valid_thru, last_activity, remember,
						   user_id, username, first_name, last_name, user_status, user_role_id, session_data)
						  VALUES
						  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);';
		$this->dbh->prepare($query);
		$this->dbh->bind_param('sssssssssis',
			$this->session_secret,
			date('Y-m-d H:i:s', $this->valid_thru),
			date('Y-m-d H:i:s', $this->last_activity),
			(int) $this->remember,
			$this->user_id,
			$this->username,
			$this->first_name,
			$this->last_name,
			$this->user_status,
			$this->user_role_id,
			$this->serialize()
		);

		if($this->dbh->execute() === TRUE){
			return TRUE;
		}else{
			throw new Exception("Cannot add new session to database");
		}
	}
}

class Login{

	private $role = 0;

	public static function autoloader($class){
		include $class . '.php';
	}

	// Pass a Configuration object
	public function __construct($configuration){
		$this->config = $configuration;
	}

	// Returns the session data or returns null
	// In case of an error it directly sets the Http Response Header
	static public function getSessionData($configuration, $db=NULL){
        // check whether cookie is set
		if(! isset($_COOKIE['SESSION'])){
			HttpHeader::setLocation($configuration->login_page.'?val=1');
            return null;
		}

	    // DB Connection
		$disconnect = FALSE;
		if(!isset($db)){
			$db = new DBAccess($configuration);
			if(!$db->connect()){
				HttpHeader::setResponseCode(500);
				exit;
			}
			$disconnect = TRUE;
		}

		// Check whether there is an active session
		$query = 'SELECT session_data FROM browser_session
		          WHERE session_secret = ?';
		$db->prepare($query);
		$db->bind_param('s', $_COOKIE['SESSION']);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(isset($res) and $res){
			$session_data = (array) json_decode($res[0]['session_data']);
			if($disconnect){
				$db->disconnect();
			}
			return $session_data;
		}else{
			error_log("classes/Login: No session data found");
			if($disconnect){
				$db->disconnect();
			}
			return null;
		}
	}

	// configuration	A Configuration object
	static public function updateSessionData($data, $configuration, $db=NULL){
		// DB Connection
		$disconnect_again = FALSE;
		if(!isset($db)){
			$db = new DBAccess($configuration);
			if(!$db->connect()){
				HttpHeader::setResponseCode(500);
				exit;
			}
			$disconnect_again = TRUE;
		}

		$query = 'UPDATE browser_session SET
						  session_data = ?,
						  valid_thru = ?,
						  last_activity = ?,
		                  WHERE session_secret = ?;';
		$db->prepare($query);
		$db->bind_param('ssss', json_encode($data),
		                        date('Y-m-d H:i:s', time()+$configuration->browser_session_timeout_max),
								date('Y-m-d H:i:s', time()),
								$_COOKIE['SESSION']);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if($res){
			$_SESSION = $data;
			if($disconnect_again === TRUE) $db->disconnect();
			return $_SESSION;
		}else{
			if($disconnect_again === TRUE) $db->disconnect();
			return null;
		}
	}

	// Login with username and password
	// - username	either the username or the email
	// - password   the password hash
	// - remember   true for the maximum session timeout, false for default session timeout
	// - config     the configuration object
	//
	// exits with Code 500 in case of an error or returns
	//
	// Returns:
	// -1: wrong username/password
	// -2: user is locked
	//  1: in case the login was successful and the session is stored in the database
	public function login($username, $password, $remember, $configuration){
		$password_salt = 0;

		// DB connection
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			HttpHeader::setResponseCode(500);
			error_log('classes/Login.php: Cannot connect to database');
			exit;
		}

		// get salt of this user (or create one if the user has none)
		$query = 'SELECT id AS user_id, 
			password_salt AS password_salt
			FROM user
			WHERE (username = ? OR email = ?)';
		$db->prepare($query);
		$db->bind_param('ss',
			$username,
			$username
		);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(! isset($res)){
			HttpHeader::setResponseCode(500);
			$db->disconnect();
			error_log('classes/Login.php: Cannot find user, SQL error for ' . $username . ' and password ***');
			exit;
		}else{
			if(sizeof($res) != 1){
				// this user is not found
				$db->disconnect();
				error_log('classes/Login.php: Wrong username/password (user not found): ' . $username . '/' . $password);
				return -1;
			}else{
				// check whether this user has a salt already
				$password_salt = $res[0]['password_salt'];
				if(! isset($password_salt)){
					$password_salt = self::generate_password_salt($username, $db);
				}
			}
		}

		// check if the user can login with salted password
		//  -> yes: use new flow
		//  -> no : use old flow
		if(self::has_salted_password($username, $db)){
			try {
				$login_status = self::verify_password_hash($username, $password, $password_salt, $db);
				if($login_status == LoginStatus::WrongPassword){
					return -1;
				}elseif($login_status == LoginStatus::Locked){
					return -2;
				}
			} catch (Exception $e){
				HttpHeader::setResponseCode(500);
				$db->disconnect();
				error_log('classes/Login.php: Cannot login, error: ' . $e->getMessage());
			}
		}else{
			try {
				$login_status = self::verify_password($username, $password, $password_salt, $db);
				if($login_status == LoginStatus::WrongPassword){
					error_log("Wrong password for user $username");
					return -1;
				}elseif($login_status == LoginStatus::Locked){
					error_log("User locked for user $username");
					return -2;
				}
			} catch (Exception $e){
				HttpHeader::setResponseCode(500);
				$db->disconnect();
				error_log('classes/Login.php: Cannot login, error: ' . $e->getMessage());
			}
		}

		// update session information
		try {
			$user_info = self::get_user_info($username, $db);
			$session_info = new SessionInfo($db);
			$session_info->logged_in      = TRUE;
			$session_info->user_id        = $user_info['id'];
			$session_info->user_status    = $user_info['status'];
			$session_info->locked         = $user_info['locked'];
			$session_info->username       = $user_info['username'];
			$session_info->first_name     = $user_info['first_name'];
			$session_info->last_name      = $user_info['last_name'];
			$session_info->user_role_id   = $user_info['user_role_id'];
			$session_info->user_role_name = $user_info['user_role_name'];
			$session_info->last_activity  = time();
			if($remember){
				$session_info->remember   = TRUE;
				$session_info->valid_thru = time() + $configuration->browser_session_timeout_max;
			}else{
				$session_info->remember   = FALSE;
				$session_info->valid_thru = time() + $configuration->browser_session_timeout_max;
			}
			$session_info->session_secret = self::get_session_secret();

			// insert new session information to database
			$session_info->insert_session();

			// set the session cookie in the browser
			// secure flag cannot be set in non HTTPS environments (such as testing)
			// domain is unknown generally
			setcookie(
				"SESSION", 
				$session_info->session_secret, 
				[   'expires' => time()+(3600*24*7*56), 
					'path'    => '/', 
					'domain'  => '', 
					'secure'  => FALSE,
					'httponly' => TRUE,
					'samesite'=> "strict"
				]
			);
		} catch (Exception $e){
			HttpHeader::setResponseCode(500);
			$db->disconnect();
			error_log('classes/Login.php: Cannot update session information for ' . $username . ' and password ***: '. $e->getMessage());
			setcookie('SESSION', '', -1, '/');
			exit;
		}

		$db->disconnect();
		return 1;
	}

	// Logout the current user
	static public function logout($configuration){
		// check that we have a session cookie
		if(!isset($_COOKIE['SESSION'])){
			return FALSE;
		}

		// DB Connection
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			error_log("Login/logout: cannot connect to database");
			HttpHeader::setResponseCode(500);
			exit;
		}
		$query = 'UPDATE browser_session SET
						  valid_thru = ?
		                  WHERE session_secret = ?;';
		$db->prepare($query);
		$db->bind_param('ss', date("Y-m-d H:i:s", time()), $_COOKIE['SESSION']);
		if($db->execute()){
			setcookie('SESSION', '', -1, '/');
			$db->disconnect();
			return TRUE;
		}else{
			$db->disconnect();
			return FALSE;
		}
	}

	// Check if the user has a valid session
	// Return TRUE if the user is already logged in, else return FALSE
    // Use for roles: 0 guest, 1 member, 2 admin, 3 course
	public function isLoggedIn($role=NULL){
		// use sanitizer
		$filter = new Sanitizer();

		// get user cookie
		if(!isset($_COOKIE['SESSION']) or !$filter->isCookie($_COOKIE['SESSION'])){
			return FALSE;
		}

		// check if the user has still a valid session
		$db = new DBAccess($this->config);
		if(!$db->connect()){
			error_log("classes/Login: Cannot connect to database");
			return FALSE;
		}

		// query for the user's session in the database
		$query = 'SELECT user_id as uid, username as uname, first_name as fn,
		                 last_name as ln, user_status as status, user_role_id as user_role_id, remember as remember
				  FROM browser_session
				  WHERE session_secret = ?
				    AND valid_thru > ?;';
		$db->prepare($query);
		$db->bind_param('ss', $_COOKIE['SESSION'], date("Y-m-d H:i:s", time()));
		$db->execute();
		$result = $db->fetch_stmt_hash();
		if(!isset($result)){
			error_log("classes/Login: Cannot execute: " . $query);
		    $db->disconnect();
			return FALSE;
		}

		// only one entry should be found
		if(count($result) > 1){
			error_log('More than one active session. Destroy all for :' . $_COOKIE['SESSION']);
			$query = 'UPDATE browser_session
			          SET valid_thru = ?
                      WHERE session_secret = ?';
			$db->prepare($query);
			$db->bind_param('ss', date("Y-m-d H:i:s", time()-1), $_COOKIE['SESSION']);
			$db->execute();
			$db->disconnect();
			setcookie('SESSION', null, -1, '/');
			return FALSE;
		}
		if(count($result) == 0){
			error_log('classes/Login: No session was found in the database');
			$db->disconnect();
			setcookie('SESSION', null, -1, '/');
			return FALSE;
		}

		// If the role does not match, the user is logged out because he
		// tampered with some sites he should not have access to
		if(isset($role) and $role != $result[0]['user_role_id']){
			error_log('Access to resources which need more permissions. Destroy all for :' . $_COOKIE['SESSION']);
			$query = 'UPDATE browser_session
			          SET valid_thru = ?
                      WHERE session_secret = ?';
			$db->prepare($query);
			$db->bind_param('ss',
				date("Y-m-d H:i:s", time()-1),
				$_COOKIE['SESSION']
			);
			$db->execute();
			$db->disconnect();
			setcookie('SESSION', null, -1, '/');
			return FALSE;
		}

		// The user is logged in:
		// prolong user-session in database by the default time
		$query = 'UPDATE browser_session
			      SET valid_thru = ?,
				      last_activity = ?
				  WHERE session_secret = ?;';
		$db->prepare($query);
		$db->bind_param('sss',
			date("Y-m-d H:i:s", time()+ $this->config->browser_session_timeout_max),
			date("Y-m-d H:i:s", time()),
			$_COOKIE['SESSION']
		);
		$db->execute();
		$db->disconnect();
		return TRUE;
	}

	public function isLoggedInAndRole($role=0){
		if($this->isLoggedIn($role)){
			return TRUE;
		}else{
			return FALSE;
		}
	}

	public function isLoggedInRedirect($redirect_url=NULL){
		if($this->isLoggedIn()){
			return TRUE;
		}else{
			if($redirect_url){
				header("Location: ".$redirect_url);
			}
			return FALSE;
		}
	}

	public function isLoggedInAndRoleRedirect($role, $redirect_url){
		if($this->isLoggedIn($role)){
			return TRUE;
		}else{
			header("Location:".$redirect_url);
			return FALSE;
		}
	}

	public function isAdmin(){
		if(isset($this->role) and $this->role == $this->config->admin_user_status_id){
			return TRUE;
		}elseif(!$this->role){
			return $this->isLoggedInAndRole($this->config->admin_user_status_id);
		}else{
			return FALSE;
		}
	}

	private function generate_password_salt($username, $db){
		$password_salt = rand(0, 65635);

		$query = 'UPDATE user SET password_salt = ? WHERE username = ? OR email = ?';
		$db->prepare($query);
		$db->bind_param('sss',
			$password_salt,
			$username,
			$username
		);
		$db->execute();

		return $password_salt;
	}

	private function has_salted_password($username, $db){
		$query = 'SELECT id, password, password_hash FROM user WHERE username = ? OR email = ?';
		$db->prepare($query);
		$db->bind_param('ss',
			$username,
			$username
		);
		$db->execute();
		$result = $db->fetch_stmt_hash();
		if(!isset($result)){
			error_log("Cannot query database");
			return FALSE;
		}else{
			if(isset($result[0]['password_hash']) and !isset($result[0]['password'])){
				return TRUE;
			}else{
				return FALSE;
			}
		}
	}

	private function verify_password_hash($username, $password, $password_salt, $db){
		// generate password hash
		$password_hash = self::get_password_hash($password, $password_salt);

		$query = 'SELECT id, username, first_name, last_name, email, status, locked
		          FROM user
				  WHERE (username = ? OR email = ?)
					AND password_hash = ? AND password_hash IS NOT NULL';
		$db->prepare($query);
		$db->bind_param('sss',
			$username,
			$username,
			$password_hash
		);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(! isset($res)){
			throw new Exception('classes/Login.php: Cannot login, SQL error for ' . $username . ' and password ***');
		}
		if(sizeof($res) > 1){
			throw new Exception('classes/Login.php: Username not unique for user ' . $username);
		}
		if(sizeof($res) == 0){
			return LoginStatus::WrongPassword;
		}
		if(sizeof($res) == 1){
			if($res[0]['locked'] == 1){
				return LoginStatus::Locked;
			}else{
				return LoginStatus::CorrectPassword;
			}
		}

		throw new Exception('classes/Login.php: Unknown Error while checking login');
	}


	private function verify_password($username, $password, $password_salt, $db){
		$query = 'SELECT id, username, first_name, last_name, email, status, locked
		          FROM user
				  WHERE (username = ? OR email = ?)
					AND password = ? AND password IS NOT NULL';
		$db->prepare($query);
		$db->bind_param('sss',
			$username,
			$username,
			$password
		);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(! isset($res)){
			throw new Exception('classes/Login.php: Cannot login, SQL error for ' . $username . ' and password ***');
		}
		if(sizeof($res) > 1){
			throw new Exception('classes/Login.php: Username not unique for user ' . $username);
		}
		if(sizeof($res) == 0){
			return LoginStatus::WrongPassword;
		}
		if(sizeof($res) == 1){
			// set the password_hash which is salted
			self::set_password_hash($username, $password, $password_salt, $db);
			if($res[0]['locked'] == 1){
				return LoginStatus::Locked;
			}else{
				return LoginStatus::CorrectPassword;
			}
		}

		throw new Exception('classes/Login.php: Unknown Error while checking login');
	}

	private function set_password_hash($username, $password, $password_salt, $db){
		$password_hash = self::get_password_hash($password, $password_salt);

		$query = 'UPDATE user SET password_hash = ? WHERE username = ? OR email = ?';
		$db->prepare($query);
		$db->bind_param('sss',
			$password_hash,
			$username,
			$username
		);
		if($db->execute()){
			// delete the old password
			$query = 'UPDATE user SET password = NULL WHERE username = ? OR email = ?';
			$db->prepare($query);
			$db->bind_param('ss',
				$username,
				$username
			);
			$db->execute();
		}else{
			error_log("Cannot set password_hash for user $username");
		}
	}

	private function get_password_hash($password, $salt){
		$password_hash = crypt($password, '$6$rounds=5000$'.$salt);
		return $password_hash;
	}

	private function get_user_info($username, $db){
		$query = 'SELECT u.id, u.username, u.first_name, u.last_name, u.email, u.status, u.locked, ur.id as user_role_id, ur.name as user_role_name
			      FROM user u
				  JOIN user_status us ON us.id = u.status
				  JOIN user_role ur ON ur.id = us.user_role_id
			      WHERE (username = ? OR email = ?)';
		$db->prepare($query);
		$db->bind_param('ss',
			$username,
			$username
		);
		if($db->execute()){
			$result = $db->fetch_stmt_hash();
			if(!isset($result) or sizeof($result) != 1){
				throw new Exception("User information not found or not unique");
			}
			$user = $result[0];
			return $user;
		}else{
			throw new Exception("Cannot query database for user information");
		}
	}

	private function get_session_secret(){
		$random = random_bytes(2048);
		return hash('sha256', $random);
	}
}
?>
