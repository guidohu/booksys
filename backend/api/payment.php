<?php	
	// automatically load all classes
	spl_autoload_register('payment_autoloader');
	function payment_autoloader($class){
		include '../classes/'.$class.'.php';
	}

	// get configuration access
	$configuration = new Configuration();
	
	// Check if the user is already logged in and is of type admin
	$lc = new Login($configuration);
	if(!$lc->isAdmin()){
		echo json_encode(Status::errorStatus("Insufficient permissions"));
		exit;
	}
	
	// check if we have an action
	if(!isset($_GET['action'])){
		echo json_encode(Status::errorStatus("API action not known"));
		exit;
	}
	
	$response = '';
	switch($_GET['action']){
		case 'get_transactions':
			$response = get_transactions($configuration);
			break;
		case 'get_statistics':
			$response = get_statistics($configuration);
			break;
		case 'get_payment_types':
			$response = get_payment_types($configuration);
			break;
		case 'get_expenditure_types':
			$response = get_expenditure_types($configuration);
			break;
		case 'add_expenditure':
		    $response = add_expenditure($configuration);
			break;
		case 'add_payment':
			$response = add_payment($configuration);
			break;
		case 'get_years':
			$response = get_years($configuration);
			break;
		case 'delete_transaction':
			$response = delete_transaction($configuration);
			break;
		default:
			$respone = Status::errorStatus("action unknown for this API call");
			break;
	}

	echo json_encode($response);
	return;
	
	// FUNCTIONS
	//---------------------------------------------------	
	
	/* Returns the transactions */
	function get_transactions($configuration){
		$post_data = json_decode(file_get_contents('php://input'));
		
		$sanitizer = new Sanitizer();
		
		// general input validation and date selection
		$year = 'any';
		$year_constraint = " 1 ";
		if(isset($post_data->year) and $sanitizer->isInt($post_data->year)){
			$year = intval($post_data->year);
			$year_constraint = " year(acc.timestamp) = $year ";
		}
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to the database");
		}
		
		$query = "SELECT acc.tbl as tbl, acc.id as id, acc.user_id as user_id, u.first_name as fn, u.last_name as ln, 
                         acc.type_id as type_id, et.name as type_name, acc.timestamp as timestamp, 
                         acc.amount as amount, acc.comment as comment, '$configuration->currency' as currency
			      FROM 
				    (SELECT 0 as tbl, e.id, e.user_id, e.type_id, 
                            e.timestamp, e.amount_chf * -1 as amount, e.comment 
                     FROM expenditure e
                     UNION ALL
                     SELECT 1 as tbl, p.id, p.user_id, p.type_id as type_id, 
                            p.timestamp, p.amount_chf as amount, p.comment as comment 
                     FROM payment p
                     UNION ALL
                     SELECT 2 as tbl, b.id, b.user_id, 0 as type_id,
                            b.timestamp, b.cost_chf * -1 as amount, 
							CONCAT(b.liters, 'L fuel') as comment
                     FROM boat_fuel b WHERE contributes_to_balance = 1
                    ) as acc 
				  LEFT JOIN user u ON u.id = acc.user_id
				  LEFT JOIN expenditure_type et ON et.id = acc.type_id 
                  WHERE " . $year_constraint . "
                  ORDER BY acc.timestamp DESC";			  
		$res = $db->fetch_data_hash($query);
		if(!isset($res)){
		    error_log("api/payment: Cannot retrieve payment information");
			$db->disconnect();
			return Status::errorStatus("Cannot retrieve payment information from the database");
		}
		$db->disconnect();

		return Status::successDataResponse("success", $res);
	}
	
	/* Delete a payment entry from the database
       Parameters:
        - table_id    the id of the table (0,1,2)
		- row_id      the id of the row to delete
	*/
	function delete_transaction($configuration){
		$post_data = json_decode(file_get_contents('php://input'));
		
		$sanitizer = new Sanitizer();
		
		// general input validation
		if(! isset($post_data->table_id) or !$sanitizer->isInt($post_data->table_id)){
			error_log("api/payment: request without table_id");
			return Status::errorStatus("No table_id provided");
		}
		if(! isset($post_data->row_id) or !$sanitizer->isInt($post_data->row_id)){
			error_log("api/payment: request without row_id");
			return Status::errorStatus("No row_id provided");
		}
		
		// get translate the table_id to the actual
		// database table that is affected
		$row_id = intval($post_data->row_id);
		if($post_data->table_id == 0){
			$table_name = "expenditure";
		}elseif($post_data->table_id == 1){
			$table_name = "payment";
		}elseif($post_data->table_id){
			$table_name = "boat_fuel";
		}else{
			error_log("api/payment: invalid table_id");
			return Status::errorStatus("The provided table_id is not a valid one.");
		}
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database.");
		}
		
		// build the query
		$query = "DELETE FROM " . $table_name . " WHERE id = ?";
		$db->prepare($query);
		$db->bind_param('i', $row_id);
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/payment: Cannot delete entry');
			return Status::errorStatus("Cannot delete transaction.");
		}
		$db->disconnect();
		
		return Status::successStatus("deleted transaction successfully");	
	}
	
	/* Returns some statistics about balance and minute usage */
	function get_statistics($configuration){
		$post_data = json_decode(file_get_contents('php://input'));
		
		$sanitizer = new Sanitizer();
		
		// general input validation
		$year = 'any';
		if(isset($post_data->year) and $sanitizer->isInt($post_data->year)){
			$year = intval($post_data->year);
		}
		
		$statistic = Array();
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database");
		}
		
		// get total of all payments that we received
		$query = "SELECT coalesce(sum(amount_chf),0) as sum FROM payment";
		$res = $db->fetch_data_hash($query);
		$statistic['total_payment'] = $res[0]['sum'];
		
		// get total of all payments in the given year
		if($year != 'any'){
			$query = "SELECT coalesce(sum(amount_chf),0) as sum FROM payment WHERE year(timestamp) = $year";
			$res = $db->fetch_data_hash($query);
			$statistic['total_payment_selected_year'] = $res[0]['sum'];
		}else{
			$statistic['total_payment_selected_year'] = $statistic['total_payment'];
		}
		$res = $db->fetch_data_hash($query);
		$statistic['total_payment_selected_year'] = $res[0]['sum'];
		
		// get total of all expenses
		$query = "SELECT sum(t.tot) as sum 
		          FROM
				  ( SELECT coalesce(sum(e.amount_chf),0) as tot FROM expenditure e
				    UNION ALL
					SELECT coalesce(sum(b.cost_chf), 0) as tot FROM boat_fuel b WHERE b.contributes_to_balance = 1
				  ) as t";
		$res = $db->fetch_data_hash($query);
		$statistic['total_expenditure'] = $res[0]['sum'];
		
		// get total of all expenditures in the given year
		if($year !='any'){
			$query = "SELECT sum(t.tot) as sum 
		          FROM
				  ( SELECT coalesce(sum(e.amount_chf),0) as tot FROM expenditure e WHERE year(e.timestamp) = $year
				    UNION ALL
					SELECT coalesce(sum(b.cost_chf), 0) as tot FROM boat_fuel b WHERE year(b.timestamp) = $year AND b.contributes_to_balance = 1
				  ) as t";
			$res = $db->fetch_data_hash($query);
			$statistic['total_expenditure_selected_year'] = $res[0]['sum'];
		}else{
			$statistic['total_expenditure_selected_year'] = $statistic['total_expenditure'];
		}
		
		// get total of all expenditures without refunds
		$query = "SELECT sum(t.tot) as sum 
		          FROM
				  ( SELECT coalesce(sum(e.amount_chf),0) as tot FROM expenditure e WHERE e.type_id <> 8
				    UNION ALL
					SELECT coalesce(sum(b.cost_chf), 0) as tot FROM boat_fuel b WHERE b.contributes_to_balance = 1
				  ) as t";
		$res = $db->fetch_data_hash($query);
		$statistic['total_expenditure_no_refunds'] = $res[0]['sum'];
		
		// get total of all expenditures in the given year
		if($year !='any'){
			$query = "SELECT sum(t.tot) as sum 
		          FROM
				  ( SELECT coalesce(sum(e.amount_chf),0) as tot FROM expenditure e WHERE e.type_id <> 8 AND year(e.timestamp) = $year
				    UNION ALL
					SELECT coalesce(sum(b.cost_chf), 0) as tot FROM boat_fuel b WHERE year(b.timestamp) = $year AND b.contributes_to_balance = 1
				  ) as t";
			$res = $db->fetch_data_hash($query);
			$statistic['total_expenditure_selected_year_no_refunds'] = $res[0]['sum'];
		}else{
			$statistic['total_expenditure_selected_year_no_refunds'] = $statistic['total_expenditure_no_refunds'];
		}
		
		// get total used heat cost
		$query = "SELECT sum(cost_chf) as sum FROM heat";
		$res = $db->fetch_data_hash($query);
		$statistic['total_used'] = $res[0]['sum'];
		
		// get total used heat cost per year
		if($year != 'any'){
			$query = "SELECT coalesce(sum(cost_chf),0) as sum FROM heat WHERE year(timestamp) = $year";
			$res = $db->fetch_data_hash($query);
			$statistic['total_used_selected_year'] = $res[0]['sum'];
		}else{
			$statistic['total_used_selected_year'] = $statistic['total_used'];
		}
		
		// get total session payments
		$query = "SELECT coalesce(sum(amount_chf),0) as sum FROM payment WHERE type_id = 4";
		$res = $db->fetch_data_hash($query);
		$statistic['total_session_payment'] = $res[0]['sum'];
		
		// get total session payments this year
		if($year != 'any'){
			$query = "SELECT coalesce(sum(amount_chf),0) as sum FROM payment WHERE type_id = 4 AND year(timestamp) = $year";
			$res = $db->fetch_data_hash($query);
			$statistic['total_session_payment_selected_year'] = $res[0]['sum'];
		}else{
			$statistic['total_session_payment_selected_year'] = $statistic['total_session_payment'];
		}
		
		// get total open cost (payment for sessions - heat costs)
		$statistic['total_open'] = $statistic['total_session_payment'] - $statistic['total_used'];
		
		// get current balance
		$statistic['current_balance'] = $statistic['total_payment'] - $statistic['total_expenditure'];
		
		// get current profit
		$statistic['current_session_profit'] = $statistic['total_used'] - $statistic['total_expenditure_no_refunds'];
		$statistic['current_session_profit_selected_year'] = $statistic['total_used_selected_year'] - $statistic['total_expenditure_selected_year_no_refunds'];
		
		// get admin minutes
		$query = "SELECT coalesce(sum(duration_s), 0) as a_sum FROM heat h, user u WHERE h.user_id = u.id AND u.status = 2";
		$res = $db->fetch_data_hash($query);
		$statistic['admin_minutes'] = intval($res[0]['a_sum']/60);
		
		// get admin minutes selected year
		if($year != 'any'){
			$query = "SELECT coalesce(sum(duration_s), 0) as a_sum FROM heat h, user u WHERE h.user_id = u.id AND u.status = 2 AND year(h.timestamp) = $year";
			$res = $db->fetch_data_hash($query);
			$statistic['admin_minutes_selected_year'] = intval($res[0]['a_sum']/60);
		}else{
			$statistic['admin_minutes_selected_year'] = $statistic['admin_minutes'];
		}
		
		// get member minutes
		$query = "SELECT coalesce(sum(duration_s), 0) as m_sum FROM heat h, user u WHERE h.user_id = u.id AND u.status = 1";
		$res = $db->fetch_data_hash($query);
		$statistic['member_minutes'] = intval($res[0]['m_sum']/60);
		
		// get member minutes selected year
		if($year != 'any'){
			$query = "SELECT coalesce(sum(duration_s), 0) as a_sum FROM heat h, user u WHERE h.user_id = u.id AND u.status = 1 AND year(h.timestamp) = $year";
			$res = $db->fetch_data_hash($query);
			$statistic['member_minutes_selected_year'] = intval($res[0]['a_sum']/60);
		}else{
			$statistic['member_minutes_selected_year'] = $statistic['member_minutes'];
		}
		
		// get guest minutes
		$query = "SELECT coalesce(sum(duration_s), 0) as g_sum FROM heat h, user u WHERE h.user_id = u.id AND u.status = 0";
		$res = $db->fetch_data_hash($query);
		$statistic['guest_minutes'] = intval($res[0]['g_sum']/60);
		
		// get guest minutes selected year
		if($year != 'any'){
			$query = "SELECT coalesce(sum(duration_s), 0) as a_sum FROM heat h, user u WHERE h.user_id = u.id AND u.status = 0 AND year(h.timestamp) = $year";
			$res = $db->fetch_data_hash($query);
			$statistic['guest_minutes_selected_year'] = intval($res[0]['a_sum']/60);
		}else{
			$statistic['guest_minutes_selected_year'] = $statistic['guest_minutes'];
		}
		
		$db->disconnect();

		$statistic['currency'] = $configuration->currency;
		
		return Status::successDataResponse("success", $statistic);		
	}
	
	/* Returns the different expenditure types and their name */
	function get_expenditure_types($configuration){
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database");
		}
		
		// get types for expenditures
		$query = "SELECT id, name FROM expenditure_type ORDER BY name";
		$res = $db->fetch_data_hash($query);
		
		$db->disconnect();
		
		return Status::successDataResponse("success", $res);	
	}
	
	/* Returns the different payment types and their name */
	function get_payment_types($configuration){		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database");
		}
		
		// get types for payments
		$query = "SELECT id, name FROM expenditure_type 
		            WHERE id NOT IN (0, 1, 2) 
					ORDER BY name";
		$res = $db->fetch_data_hash($query);
		
		$db->disconnect();
		
		return Status::successDataResponse("success", $res);
	}
	
	/* Returns the years as an array for which we have payment data */
	function get_years($configuration){
		$years = Array();
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database");
		}
		
		// get years for payments
		$query = "SELECT year(timestamp) AS year FROM payment
				  UNION
				  SELECT year(timestamp) AS year FROM expenditure
				  UNION
				  SELECT year(timestamp) AS year FROM boat_fuel
				  UNION
				  SELECT 'any' AS year
				  ORDER BY year ASC";
		$res = $db->fetch_data_hash($query);
		
		foreach ($res as $row){
			array_push($years, $row['year']);
		}
		
		$db->disconnect();
		return Status::successDataResponse('success', $years);		
	}
	
	function add_payment($configuration){
		$post_data = json_decode(file_get_contents('php://input'));
		
		$sanitizer = new Sanitizer();
		
		// general input validation
		if(!isset($post_data->user_id) or !$sanitizer->isInt($post_data->user_id)){
			return Status::errorStatus("No valid user selected, please select a user");
		}
		if(!isset($post_data->type_id) or !$sanitizer->isInt($post_data->type_id)){
			return Status::errorStatus("No valid income type");
		}
		if(!isset($post_data->date) or !$sanitizer->isDate($post_data->date)){
			return Status::errorStatus("No valid date");
		}
		if(!isset($post_data->amount) or !$sanitizer->isFloat($post_data->amount)){
			return Status::errorStatus("No valid amount");
		}
		if($post_data->type_id != 4 and $post_data->type_id != 6
		   and (!isset($post_data->comment) or $post_data->comment == '')){
			return Status::errorStatus("No comment provided");
		}else if((!isset($post_data->comment) or $post_data->comment == '') and $post_data->type_id == 4){
			$post_data->comment = 'Session Payment';
		}else if((!isset($post_data->comment) or $post_data->comment == '') and $post_data->type_id == 6){
			$post_data->comment = 'Membership Fee';
		}
		
		# add a new payment to the database
		$query = "INSERT INTO payment (user_id, type_id, timestamp, amount_chf, comment)
		          VALUES ( ?,?,?,?,? );";
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to the database");
		}
		
		$db->prepare($query);
		$db->bind_param('iisds',
		                $post_data->user_id,
						$post_data->type_id,
						$post_data->date,
						sprintf('%.2f', $post_data->amount),
						$post_data->comment);
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/payment: Cannot add payment');
			return Status::errorStatus("Income transaction cannot be added, due to an unknown error.");
		}
		$db->disconnect();
		
		# in case it was a user who paid for sessions
		# we write an email to the user
		if($post_data->type_id == 4){
			$user = new User($configuration);
			$user_data = $user->getUserById($post_data->user_id);
			$msg = 'Dear ' . $user_data['first_name'] . ' ' . $user_data['last_name'] . "\n\n"
			     . 'Thank you for your payment. We received ' . $post_data->amount . ' ' . $configuration->currency . ' and' . "\n"
				 . 'added the amount to your account balance.'."\n\n"
				 . 'See you soon on the boat' . "\n"
				 . 'wakeandsurf Crew';
			Email::sendMail($user_data['email'], 'Payment received', $msg, $configuration);
		}
		
		return Status::successStatus("Income added");
	}
	
	function add_expenditure($configuration){
		$post_data = json_decode(file_get_contents('php://input'));
		
		$sanitizer = new Sanitizer();
		
		// general input validation
		if(!isset($post_data->user_id) or !$sanitizer->isInt($post_data->user_id)){
			return Status::errorStatus("No valid user selected, please select a user");
		}
		if(!isset($post_data->type_id) or !$sanitizer->isInt($post_data->type_id)){
			return Status::errorStatus("No valid expenditure type");
		}
		if(!isset($post_data->date) or !$sanitizer->isDate($post_data->date)){
			return Status::errorStatus("No valid date");
		}
		if(!isset($post_data->amount) or !$sanitizer->isFloat($post_data->amount)){
			return Status::errorStatus("No valid amount");
		}
		if(!isset($post_data->comment) or $post_data->comment == ''){
			return Status::errorStatus("No comment provided");
		}
		if($post_data->type_id == 0){
			return Status::errorStatus("Use the boat API for fuel expenses");
		}
		
		# add a new expenditure to the database
		$query = "INSERT INTO expenditure (user_id, type_id, timestamp, amount_chf, comment)
		          VALUES (?,?,?,?,?);";
		
		$db = new DBAccess($configuration);
		if(!$db->connect()){
			return Status::errorStatus("Cannot connect to database");
		}
		
		$db->prepare($query);
		$db->bind_param('iisds',
		                $post_data->user_id,
						$post_data->type_id,
						$post_data->date,
						sprintf('%.2f', $post_data->amount),
						$post_data->comment);
		if(!$db->execute()){
			$db->disconnect();
			error_log('api/payment: Cannot add expenditure');
			return Status::errorStatus("Internal Server Error, cannot add the expenditure.");
		}
		
		$db->disconnect();
		return Status::successStatus("Expense added");
	}

?>