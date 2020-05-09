<?php

/* Handles database access to the booking system database */
class DBAccess{

	// configuration:	Configuration object
	public function __construct($configuration){
		$this->config = $configuration;
	}

	/* Connect to the database */
    public function connect(){
		try{
			// redefine error handler to not print PHP warnings to the console
			set_error_handler(function($errno, $errstr, $errfile, $errline, $errcontext) {
				throw new ErrorException($errstr, 0, $errno, $errfile, $errline);
			});

			$this->dbh = new mysqli(
				$this->config->db_server,
				$this->config->db_user,
				$this->config->db_password,
				$this->config->db_name
			);
			if($this->dbh->connect_errno){
				error_log("Could not connect to database: " . $this->dbh->connect_errno
						. " - " . $this->dbh->connect_error);
				restore_error_handler();
				return FALSE;
			}
			if($this->dbh->ping()){
				$this->dbh->set_charset("utf8");
				return TRUE;
			}
		} catch(Exception $e){
			error_log("Exception thrown while connecting with mysqli: " . $e);
		}

		restore_error_handler();
		return FALSE;
	}

	/* Disconnect the database */
	public function disconnect(){
		return $this->dbh->close();
	}

	/* Executes an *.sql file */
	public function apply_sql_file($file){
		$commands = file_get_contents($file);
		if(!$this->dbh->multi_query($commands)){
			error_log("First Query failed: " . $this->dbh->error);
			return FALSE;
		}

		error_log($this->dbh->more_results());

		do {
			if(!$this->dbh->next_result()){
				error_log("Specific query failed: " . $this->dbh->error);
				return FALSE;
			}
			if($result = $this->dbh->store_result()){
				// not required currently
				error_log("Query applied");
				$result->free();
			}else{
				error_log("Ran query without result set");
			}
		} while($this->dbh->more_results());

		if($this->dbh->error){
			error_log("Specific query failed: " . $this->dbh->error);
			return FALSE;
		}

		return TRUE;
	}

	/* Send a query to the database.
       Returns the resource, TRUE or FALSE
    */
	public function query($query){
		$res = $this->dbh->query($query);
		if(!isset($res) or $res === FALSE or !$res){
			error_log("Cannot execute: '$query'\nThe Following Error Appeared: " . $this->dbh->error);
			return FALSE;
		}
		return $res;
	}

	/* Returns the query as hash with the given column as id.
	   - col 	number of the row

	   On Error: returns FALSE
	   Otherwise at least an empty array
	   The index of the array is either the specified column or an integer (if -1 option)
	*/
	public function fetch_data_hash($query, $col=-1){
		$res = $this->query($query);
		if(!isset($res) or !$res){
			error_log("DEBUG: res was not defined");
			return FALSE;
		}

		$hash = Array();
		$i=0;
		while( $arr = $res->fetch_array(MYSQLI_BOTH)){
			foreach ($arr as &$value) {
				$value = $value;
			}
			if($col == -1){
				$hash[$i . ''] = $arr;
			}else{
				$hash[$arr[$col] . ''] = $arr;
			}
			$i++;
		}
		return $hash;
	}

	/* Returns true in case the query returns 1 or more results
	   - query		the SQL query to execute

	   On Error            : returns false
	   In case of no result: returns false
	*/
	public function exists($query){
		$res = $this->query($query);
		if(!isset($res) or !$res){
			error_log('DBAccess: Error to execute ' . $query);
			return FALSE;
		}

		if($res->num_rows > 0){
			return TRUE;
		}
		return FALSE;
	}

	/* Prepare a DB Query
	   - query		the SQL query

	   On Error		: returns false
	   Otherwise    : returns the statement
	*/
	public function prepare($query){
		// if we have already a stmt, close it
		if(isset($this->stmt)){
			$this->stmt->close();
		}

		// start with a new statement
		$stmt = $this->dbh->stmt_init();
		$stmt = $this->dbh->prepare($query);
		if(!$stmt){
			error_log('DBAccess: Cannot prepare: ' . $this->dbh->error . ' on query ' . $query);
			return FALSE;
		}
		if(!isset($stmt)){
			error_log('DBAccess: Cannot prepare: ' . $this->dbh->error . ' on query ' . $query);
			return FALSE;
		}
		$this->stmt = $stmt;
		return $this->stmt;
	}


	public function bind_param(){
		if(!isset($this->stmt)){
			error_log('DBAccess: bind_param called without defined stmt');
			return FALSE;
		}
		call_user_func_array(array(&$this->stmt, 'bind_param'), $this->refValues(func_get_args()));
	}

	function refValues($arr){
		if (strnatcmp(phpversion(),'5.3') >= 0) //Reference is required for PHP 5.3+
		{
			$refs = array();
			foreach($arr as $key => $value)
				$refs[$key] = &$arr[$key];
			return $refs;
		}
		return $arr;
	}


	public function execute(){
		if(!isset($this->stmt)){
			error_log('DBAccess: Nothing to execute, no statement defined');
			return FALSE;
		}
		if(!$this->stmt->execute()){
			error_log('DBAccess: Cannot execute prepared statement: ' . $this->stmt->error);
			return FALSE;
		}
		return TRUE;
	}

	public function last_id(){
		if(!isset($this->dbh)){
			error_log("DBAccess: No db handler");
			return FALSE;
		}
		return $this->dbh->insert_id;
	}

	public function stmt_close(){
		$this->stmt->close();
	}

	/* Returns a hash_array of of the result. They array key is either numeric or with
	   a given $key it is the specified row
	   - key		the row to use for the array-key (DEFAULT: NULL -> numeric)
	*/
	public function fetch_stmt_hash($key=NULL){
		$array = array();
		$result = $this->stmt;

		if($result instanceof mysqli_stmt)
		{
			$result->store_result();

			$variables = array();
			$data = array();
			$meta = $result->result_metadata();

			while($field = $meta->fetch_field())
				$variables[] = &$data[$field->name]; // pass by reference

			call_user_func_array(array($result, 'bind_result'), $variables);

			$i=0;
			while($result->fetch()){
				if($key){
					$name = $data[$key];
				}else{
					$name = $i;
				}
				$array[$name] = array();
				foreach($data as $k=>$v){
					$array[$name][$k] = $v;
				}

				$i++;
			}
			$result->close();
			unset($this->stmt);
		}
		elseif($result instanceof mysqli_result)
		{
			while($row = $result->fetch_assoc())
				$array[] = $row;
		}

		return $array;
	}
}
?>
