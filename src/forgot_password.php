<!DOCTYPE html>
<html lang="en">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>Booking System for Boat Communities</title>

<link rel="stylesheet" type="text/css" href="bootstrap/css/bootstrap.min.css" />
<link rel="stylesheet" type="text/css" href="bootstrap/css/bootstrap-theme-bc.css" />
<link rel="stylesheet" type="text/css" href="style.css" />
<link rel="icon" type="image/png" sizes="32x32" href="/img/apple-touch-icon-32x32.png" />
<script src="js/jquery-latest.min.js"></script>
<script src="bootstrap/js/bootstrap.min.js"></script>
<script src="js/shared_functions.js"></script>
<script src="js/jquery.validate.min.js"></script>
<script src="js/jquery.validate.custom.rules.js"></script>
<script src="js/sha.js"></script>
<script src="js/jquery.stayInWebApp.min.js"></script>

</head>

<body>

<script>
    // Stay in the WebApp
    $(function() {
        $.stayInWebApp();
    });

    // the user id
	$user_id = 0;
	
	// Set validator defaults
	$.validator.setDefaults({
		errorElement: "validate_msg",
		invalidHandler: function(event, validator){
			var errors = validator.numberOfInvalids();
			if(errors){		
				$('#status_header').html('<h4>Error while signing up</h4>');
				$('#status_body').html('Please check your input again');
				$('#status_body').show();
				$('#status_footer').show();
				$('#error').show();
			}
		},
	});

    $(document).ready(function() {
	    $('#password_form').hide();
		
		$('#email_form').validate({
			//debug: true,
			submitHandler: function(form){
				getTokenMail();
			}, 
			rules: {
				email: {
					required:true,
					minlength:1,
					maxlength:100,
					email: true,
				},
			},
			messages: {
				email: {
					required: "please provide your email",
					email: "please enter a valid email"
				}
			},
		});
		
		$('#password_form').validate({
			//debug: true,
			submitHandler: function(form){
				changePassword();
			},
			rules: {
				token: {
					required:true,
					number: true,
				},
				password: {
					required:true,
				},
				password_conf: {
					equalTo:'#password',
				}
			},
			messages: {
				token: {
					required: 'please provide the token',
					number:   'token has to be a number',
				},
				password: {
					required: 'no password provided',
				},
				password_conf: {
					equalTo: 'passwords are not identical'
				}
			}
		});
	});	
	
	function getTokenMail(){
		$('#status_header').html('<h4>Requesting Token by Mail</h4>');
	    $('#status_body').hide();
		$('#status_footer').hide();
		$('#error').hide();
		$('#success').hide();
		$('#status_modal').modal({
			backdrop: 'static',
			keyboard: false,
			show:     true,
		});
		
		// send the ajax call to the api
		var data = new Object();
		data['email'] = $('#email').val();
		$.ajax({
			type: "POST",
			url: "api/password.php?action=token_request",
			data: JSON.stringify(data),
			success: function(data, textStatus, jqXHR){
				$('#status_header').html("<h4>Successfully sent Token Email</h4>");
				$('#status_header').show();
				$('#status_footer').show();
				$('#status_body').html('Please check your Email Inbox for the Password reset Code.');
				$('#status_body').show();
				$('#error').hide();
				$('#success').show();
				json = $.parseJSON(data);
				$user_id = json.user_id;
			},
			error: function(data, text, errorCode){
				$('#status_header').html('<h4>Error while sending Token Email</h4>');
				$('#status_body').html(data.responseText);
				$('#status_body').show();
				$('#status_footer').show();
				$('#error').show();
			},			
		});
	}
	
	function changePassword(){		
		$('#status_header').html('<h4>Resetting your Password</h4>');
	    $('#status_body').hide();
		$('#status_footer').hide();
		$('#error').hide();
		$('#success').hide();
		$('#status_modal').modal({
			backdrop: 'static',
			keyboard: false,
			show:     true,
		});
		
		// send the ajax call to the api
		var data         = new Object();
		data['token']    = $('#token').val();
		data['password'] = calcHash();
		data['user_id']  = $user_id;
		
		$.ajax({
			type: "POST",
			url: "api/password.php?action=change_password_by_token",
			data: JSON.stringify(data),
			success: function(data, textStatus, jqXHR){
				$('#status_header').html("<h4>Successfully changed password</h4>");
				$('#status_header').show();
				$('#status_footer').show();
				$('#status_body').html('Please check your Email Inbox for the Password reset Code.');
				$('#status_body').show();
				$('#error').hide();
				$('#success').show();
				window.location.href = 'index.html?valid=4';
			},
			error: function(data, text, errorCode){
				$('#status_header').html('<h4>Error while changing password</h4>');
				$('#status_body').html(data.responseText);
				$('#status_body').show();
				$('#status_footer').show();
				$('#error').show();
			},			
		});
	}
	
	function showPasswordForm(){
		$('#email_form').hide();
		$('#password_form').show();
		hideStatus();
	}
	
	// Hides the status window
	function hideStatus(){
		$("#status_modal").modal("hide");
	}
	
	// Calculates the hash of the password.
  	function calcHash() {
		try {
			var hashInput = document.getElementById("password").value;
			var hashObj = new jsSHA(hashInput, "TEXT");
			var hashOutput = hashObj.getHash("SHA-256","HEX");
			return hashOutput;
		} catch(e) {
			alert("Password could not be sent encrypted: " + e);
			return;
		}
	}
	
</script>

<div class="display">
    <div class="page-header">
		<h3>Reset Password</h3>
	</div>
	<form id="email_form">
		<div class="vcenter_placeholder vcenter_with_title">
			<div class="vcenter_content">
				<div class="form-horizontal">
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4">
							Please enter your Email address to reset your password.
						</div>
					</div>
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4">
							<label class="label-sm" for="email">Email</label>
							<input type="email" class="form-control input-sm" name="email" id="email"/>
						</div>
					</div>
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4 text-right">
							<a href="index.html" class="btn btn-default" id="cancel_btn">
								<span class="glyphicon glyphicon-remove"></span>Cancel
							</a>
							<button type="submit" class="btn btn-default" id="next_btn">
								<span class="glyphicon glyphicon-arrow-right"></span>Next
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</form> 
	
	<form id="password_form">
		<div class="vcenter_placeholder vcenter_with_title">
			<div class="vcenter_content">
				<div class="form-horizontal">
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4">
							Please enter your Token and the new password.
						</div>
					</div>
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4">
							<label class="label-sm" for="email">Token</label>
							<input type="text" class="form-control input-sm" name="token" id="token"/>
						</div>
					</div>
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4">
							<label class="label-sm" for="email">Password</label>
							<input type="password" class="form-control input-sm" name="password" id="password"/>
						</div>
					</div>
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4">
							<label class="label-sm" for="email">Confirm Password</label>
							<input type="password" class="form-control input-sm" name="password_conf" id="password_conf"/>
						</div>
					</div>
					<div class="form-group">
						<div class="col-sm-4 col-sm-offset-4 col-xs-4 col-xs-offset-4 text-right">
							<a href="index.html" class="btn btn-default" id="cancel_btn">
								<span class="glyphicon glyphicon-remove"></span>Cancel
							</a>
							<button type="submit" class="btn btn-default" id="next_btn">
								<span class="glyphicon glyphicon-arrow-right"></span>Reset Password
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</form> 
	
	<!-- Status information -->
	<div class="modal fade" id="status_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content" id="status_modal_content">
				<div class="modal-header" id="status_header">
					<h4>
						<span class="glyphicon glyphicon-upload"></span>Sending Data...
					</h4>
				</div>
				<div class="modal-body" id="status_body">
					<span class="glyphicon glyphicon-upload"></span>Sending Data...
				</div>
				<div class="modal-footer" id="status_footer">
					<button type="button" class="btn btn-default" id="success" onclick='showPasswordForm()'>OK</button>
					<button type="button" class="btn btn-default" id="error" onclick='hideStatus()'>OK</button>
				</div>
			</div>
		</div>
	</div>
</div>


</body>
</html>
