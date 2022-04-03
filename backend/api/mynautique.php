<?php
    // automatically load all classes
    spl_autoload_register('mynautique_autoloader');
    function mynautique_autoloader($class){
        if(file_exists(__DIR__.'/../classes/'.$class.'.php')){
            include __DIR__.'/../classes/'.$class.'.php';
        }
    }

    // test whether the config-file exists
    $configuration = new Configuration();
    $api_info = [
        'MYNAUTIQUE_API_URL' => "https://mynautique.azurewebsites.net/api/v2",
        'MYNAUTIQUE_AUTH_URL' => "https://identitytoolkit.googleapis.com/v1",
        'MYNAUTIQUE_AUTH_API_KEY' => "AIzaSyAb8S3Owvo8k-gI8eK_DEztFOn0FcZFRxw",
    ];
    

    // check if we have an action
    if(!isset($_GET['action'])){
        HttpHeader::setResponseCode(200);
        exit;
    }

    switch($_GET['action']){
        case 'get_boat_info':
            $response = get_boat_info($configuration, $api_info);
            echo json_encode($response);
            exit;
        case 'get_boats_info':
            $response = get_boats_info($configuration, $api_info);
            echo json_encode($response);
            exit;
    }

    HttpHeader::setResponseCode(400);
    $status = array();
    $status['ok'] = FALSE;
    $status['message'] = 'invalid action requested';
    echo json_encode($status);
    return;

    function get_boat_info($configuration, $api_info){
         _is_admin_or_return($configuration);

        // get boat id
        $sanitizer = new Sanitizer();
        $post_data = json_decode(file_get_contents('php://input'));
        if(!isset($post_data->boat_id) || !$sanitizer->isInt($post_data->boat_id)){
            return Status::errorStatus("The boat_id is not provided in the right format.");
        }

        $boat_id = intval($post_data->boat_id);
        $boat_info = _get_boats_info($configuration, $api_info);
        $boat_info['boat']['info'] = $boat_info['boats'][$boat_id]['info'];
        $boat_info['boat']['telemetry'] = $boat_info['boats'][$boat_id]['telemetry'];
        $boat_info['boat']['metainfo']['fuel_capacity'] = $configuration->mynautique_fuel_capacity;
        unset($boat_info['boats']);

        return Status::successDataResponse('success', $boat_info);
    }

    function get_boats_info($configuration, $api_info){
        _is_admin_or_return($configuration);

        return Status::successDataResponse('success', _get_boats_info($configuration, $api_info));
    }

    // returns the current database configuration
    function _get_boats_info($configuration, $api_info){

        $response_data = [];
        $tokenId = "";
        $post_data = json_decode(file_get_contents('php://input'));
        $sanitizer = new Sanitizer();
        if(!isset($post_data->token) || !isset($post_data->token_expiry) || $post_data->token_expiry < time() - 10){
            // perform login
            $token = _login($configuration, $api_info);
            $tokenId = $token['idToken'];
            $cipher_token = _encrypt($configuration, $token['idToken']);
            $response_data['token'] = $cipher_token;
            $response_data['token_expiry'] = time() + $token['expiresIn'];
        }else{
            // decrypt token
            $tokenId = _decrypt($configuration, $post_data->token);
            $response_data['token'] = $post_data->token;
            $response_data['token_expiry'] = $post_data->token_expiry;
        }

        // get fleet
        $fleet = _get_fleet($api_info, $tokenId);
        // get fuel for each boat
        for($i=0; $i < count($fleet); $i++) {
            $boat = $fleet[$i];
            $boat_id = $boat->hin;
            $boat_telemetry = _get_boat_telemetry($api_info, $tokenId, $boat_id);
            $response_data['boats'][$boat_id]['telemetry'] = $boat_telemetry;
            $response_data['boats'][$boat_id]['info'] = $boat;
        }

        return $response_data;
    }

    function _login($configuration, $api_info) {
        $call = curl_init($api_info['MYNAUTIQUE_AUTH_URL']."/accounts:signInWithPassword?key=".$api_info['MYNAUTIQUE_AUTH_API_KEY']);

        $payload = [
            "email" => $configuration->mynautique_user,
            "password" => $configuration->mynautique_password,
            "returnSecureToken" => true
        ];
        curl_setopt( $call, CURLOPT_POSTFIELDS, json_encode($payload) );
        curl_setopt( $call, CURLOPT_HTTPHEADER, array('Content-Type: application/json'));
        curl_setopt( $call, CURLOPT_RETURNTRANSFER, true );
        $result_json = curl_exec($call);
        curl_close($call);

        $result = json_decode($result_json);
        $token = [
            "idToken" => $result->idToken,
            "expiresIn" => $result->expiresIn,
        ];
        return $token;
    }

    function _get_fleet($api_info, $token) {
        $call = curl_init();
        curl_setopt($call, CURLOPT_URL, $api_info['MYNAUTIQUE_API_URL']."/fleet/get-fleet");
        curl_setopt( $call, CURLOPT_HTTPHEADER, array(
            'Content-Type: application/json',
            'token: ' . $token
        ));
        curl_setopt($call, CURLOPT_RETURNTRANSFER, 1);
        $result_json = curl_exec($call);
        curl_close($call);
        $fleet = json_decode($result_json);
        return $fleet;
    }

    function _get_boat_telemetry($api_info, $token, $boat_id) {
        $call = curl_init();
        curl_setopt($call, CURLOPT_URL, $api_info['MYNAUTIQUE_API_URL']."/boat/get-boat-telemetry/".$boat_id);
        curl_setopt( $call, CURLOPT_HTTPHEADER, array(
            'Content-Type: application/json',
            'token: ' . $token
        ));
        curl_setopt($call, CURLOPT_RETURNTRANSFER, 1);
        $result_json = curl_exec($call);
        curl_close($call);
        $boat_telemetry = json_decode($result_json);
        return $boat_telemetry;
    }

    function _encrypt($configuration, $plaintext) {
        $cipher = "AES-128-CTR";
        $iv_length = openssl_cipher_iv_length($cipher);
        $options = 0;
        $iv = '1234567891011121';
        $key = $configuration->mynautique_password;
        $ciphertext = openssl_encrypt(
            $plaintext, 
            $cipher,
            $key, 
            $options, 
            $iv
        );
        return $ciphertext;
    }

    function _decrypt($configuration, $ciphertext) {
        $cipher = "AES-128-CTR";
        $iv_length = openssl_cipher_iv_length($cipher);
        $options = 0;
        $iv = '1234567891011121';
        $key = $configuration->mynautique_password;
        $plaintext=openssl_decrypt (
            $ciphertext, 
            $cipher, 
            $key, 
            $options, 
            $iv
        );
        return $plaintext;
    }

    // returns http status permission denied in case the user is not an admin
    function _is_admin_or_return($configuration){
        $lc = new Login($configuration);
        if(!$lc->isLoggedIn($configuration->admin_user_role_id)){
            HttpHeader::setResponseCode(403);
            exit;
        }
    }

?>