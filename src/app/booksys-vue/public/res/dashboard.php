<?php
  include '../classes/Configuration.php';
  include '../classes/Login.php';

  $configuration = new Configuration();
  $login = new Login($configuration);
  $session_data = Login::getSessionData($configuration);

?>

<div class="display">
	<header class="welcome">
		Welcome, <?php echo '<a href=\'account.html\'>' . $session_data['first_name'] . " " . $session_data['last_name'] . "</a>"; ?>
	</header>
	<div class="vcenter_placeholder vcenter-dashboard">
		<div class="vcenter_content">
			<?php
			    // depending on the user-status we show different dashboards
				if($session_data['user_role_id'] == $configuration->admin_user_status_id){
					include('dashboard_admin.html');
				}else if($session_data['user_role_id'] == $configuration->member_user_status_id){
					include('dashboard_member.html');
				}else if($session_data['user_role_id'] == $configuration->guest_user_status_id){
					include('dashboard_guest.html');
				}
			?>
		</div>
	</div>
</div>

<footer>
	Copyright 2013-2020 by Guido Hungerbuehler / <span id="version"></span>
</footer>