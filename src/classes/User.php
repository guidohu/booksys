<?php  

class User{

	private $role = 0;
	private $user;

    // Automatically load classes from the classes folder
	function __autoload($class_name){
		include $class_name.'.php';
	}
	
	public function __construct($conf){
		$this->config = $conf;
	}
	
	// Returns a hash with the user's profile
	public function getUser(){
		// sanitize cookie
		$sanitize = new Sanitizer();
		if(!$sanitize->isCookie($_COOKIE['SESSION'])){
			return null;
		}
		
		// connect to the database
		$db = new DBAccess($this->config);
		if(!$db->connect()){
			error_log('classes/User: Cannot connect to the database');
			return null;
		}
		
		$query = 'SELECT u.id, u.username, u.first_name, u.last_name,
				         u.address, u.city, u.plz, u.mobile, u.email,
						 u.license, u.status, u.locked, u.comment
				  FROM user u, browser_session bs 
				  WHERE bs.session_secret = ?
				  AND bs.user_id = u.id;';
		$db->prepare($query);
		$db->bind_param('s', $_COOKIE['SESSION']);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		$db->disconnect();
		if(! isset($res) or $res === FALSE){
			error_log('classes/User: Cannot retrieve user by cookie');
			return null;
		}
		$this->user = $res[0];
		
		return $res[0];
	}
	
	public function getUserById($id){
		// sanitize cookie
		$sanitize = new Sanitizer();
		if(!$sanitize->isInt($id)){
			error_log('classes/User: No valid id provided');
			return null;
		}
		
		// connect to the database
		$db = new DBAccess($this->config);
		if(!$db->connect()){
			error_log('classes/User: Cannot connect to the database');
			return null;
		}
		
		$query = 'SELECT u.id, u.username, u.first_name, u.last_name,
						 u.address, u.city, u.plz, u.mobile, u.email,
						 u.license, u.status, u.locked, u.comment
					FROM user u
					WHERE u.id = ?;';
		$db->prepare($query);
		$db->bind_param('i', $id);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		$db->disconnect();
		if(! isset($res) or $res === FALSE){
			error_log('classes/User: Cannot retrieve user by id');
			return null;
		}
		$this->user = $res[0];
		
		return $res[0];
	}
	
	// Returns true if the password is correct
	public function isPasswordCorrect($password){
	
		// connect to the database
		$db = new DBAccess($this->config);
		if(!$db->connect()){
			return FALSE;
		}
		
		$query = 'SELECT u.id, u.password, u.password_salt, u.password_hash
				  FROM user u, browser_session bs 
				  WHERE u.id = ?';
		$db->prepare($query);
		$db->bind_param('i', $this->user['id']);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		$db->disconnect();
		if(! isset($res) or $res===FALSE){
			return FALSE;
		}

		// get the password from the DB
		$password_db = '';
		$password_provided = '';
		if(!is_null($res[$this->user['id']]['password_hash'])){
			$password_db = $res[$this->user['id']]['password_hash'];
			$password_provided = $this->get_password_hash(
				$password,
				$res[$this->user['id']]['password_salt']
			);
		}else{
			$password_db       = $res[$this->user['id']]['password'];
			$password_provided = $password;
		}

		if($password_db == $password_provided){
			return TRUE;
		}else{
			return FALSE;
		}
	}	

	private function get_password_hash($password, $salt){
		$password_hash = crypt($password, '$6$rounds=5000$'.$salt);
		return $password_hash;
	}
}
?>