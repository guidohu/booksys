<?php
    include '../classes/Configuration.php';
    include '../classes/Login.php';

    $configuration = new Configuration();
    $login = new Login($configuration);
    $session_data = Login::getSessionData($configuration);
    
    // depending on the user-status we show different dashboards
    if($session_data['user_role_id'] == $configuration->admin_user_status_id){
        include('dashboard_admin_mobile.html');
    }else if($session_data['user_role_id'] == $configuration->member_user_status_id){
        include('dashboard_member_mobile.html');
    }else if($session_data['user_role_id'] == $configuration->guest_user_status_id){
        include('dashboard_guest_mobile.html');
    }
?>
