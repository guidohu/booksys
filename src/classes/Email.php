<?php

class Email{
	
	public static function sendMail($to, $subject, $message){
		// $from
		// $smtp_server
		// $username
		// $password
		if(isset($to)){
			return mail($to, $subject, $message, 'From: Wake and Surf <noreply@wakeandsurf.ch>', '-r noreply@wakeandsurf.ch');
		}
		// sendemail -f ghungerbuehler@gmail.com -t ghungerbuehler@gmail.com -u "Test" -m "no text" -s smtp.gmail.com:587 -xu ghungerbuehler@gmail.com -xp 'wfhcueybbpelugue'
		
		return FALSE;
	}
}
?>