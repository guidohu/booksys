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
		// sanitize id
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
						 u.license, u.status, u.locked, u.comment, 
						 ur.id as user_role_id, ur.name as user_role_name
					FROM user u
					JOIN user_status us ON us.id = u.status
					JOIN user_role ur ON ur.id = us.user_role_id
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

	public function getUserBalance($id){
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
		
		// get heats
		$query = 'SELECT SUM(cost_chf) as cost
			FROM heat 
			WHERE user_id = ?;';
		$db->prepare($query);
		$db->bind_param('i', $id);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(! isset($res) or $res === FALSE){
			error_log('classes/User: Cannot retrieve user heats by id');
			return null;
		}
		$heats = $res[0]['cost'];

		// get payment refunds
		$query = 'SELECT SUM(amount_chf) as cost
			FROM expenditure 
			WHERE user_id = ?
			AND type_id = 4;';
		$db->prepare($query);
		$db->bind_param('i', $id);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(! isset($res) or $res === FALSE){
			error_log('classes/User: Cannot retrieve user session payment refunds by id');
			return null;
		}
		$refunds = $res[0]['cost'];

		// get payments
		$query = 'SELECT SUM(amount_chf) as payments
			FROM payment 
			WHERE user_id = ? 
			AND type_id = 4;';
		$db->prepare($query);
		$db->bind_param('i', $id);
		$db->execute();
		$res = $db->fetch_stmt_hash();
		if(! isset($res) or $res === FALSE){
			error_log('classes/User: Cannot retrieve user payments by id');
			return null;
		}
		$payments = $res[0]['payments'];

		$db->disconnect();

		// calculate the balance with a 2 digit precision
		$total_balance = $payments - $refunds - $heats;
		$total_balance = round($total_balance, 2);
		
		return $total_balance;
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