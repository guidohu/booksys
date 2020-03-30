<?php

class Email{
	
	public static function sendMail($to, $subject, $message){
		if(isset($to)){
			return mail($to, $subject, $message, 'From: Wake and Surf <noreply@wakeandsurf.ch>', '-r noreply@wakeandsurf.ch');
		}
		
		return FALSE;
	}
}
?>