<?php

class Email{
	
	public static function sendMail($to, $subject, $message, $configuration){
		// if smtp configuration is not existing
		if(!isset($configuration->smtp_sender)){
			error_log("Cannot send email, no SMTP sender configured");
			return FALSE;
		}
		if(!isset($configuration->smtp_server)){
			error_log("Cannot send email, no SMTP server configured");
			return FALSE;
		}
		if(!isset($configuration->smtp_username)){
			error_log("Cannot send email, no SMTP username configured");
			return FALSE;
		}
		if(!isset($configuration->smtp_password)){
			error_log("Cannot send email, no SMTP password configured");
			return FALSE;
		}

		// send email
		if(isset($to)){
			$command = "/usr/bin/sendemail";
			$command .= " -f "  . escapeshellarg($configuration->smtp_sender);
			$command .= " -t "  . escapeshellarg($to);
			$command .= " -u "  . escapeshellarg($subject);
			$command .= " -m "  . escapeshellarg($message);
			$command .= " -s "  . escapeshellarg($configuration->smtp_server);
			$command .= " -xu " . escapeshellarg($configuration->smtp_username);
			$command .= " -xp " . escapeshellarg($configuration->smtp_password);
			$command .= " -o tls=yes";

			$result = "";
			$return_code = 0;

			// execute command to send an email
			exec($command, $result, $return_code);

			if($return_code != 0){
				error_log("Cannot send email");
				error_log(print_r($result, TRUE));
				return FALSE;
			}

			return TRUE;
		}
		
		return FALSE;
	}
}
?>