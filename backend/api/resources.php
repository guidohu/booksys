<?php
	// automatically load all classes
	spl_autoload_register('resources_autoloader');
	function resources_autoloader($class){
		include '../classes/'.$class.'.php';
    }
    
    // get configuration access
    $configuration = new Configuration();
    
 	// Check if the user is already logged in and is of type admin
    $lc = new Login($configuration);
    if(!$lc->isLoggedInAndRole($configuration->admin_user_status_id)){
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
        case 'upload_logo':
            $response = upload_logo($configuration);
            break;
        default:
            $response = Status::errorStatus("action unknown for this API call");
            break;
    }

    echo json_encode($response);
    return;

    //--------------------------
    // API functions
    //--------------------------
    function upload_logo($configuration){
        $target_dir = "/var/www/html/uploads/logos/";
        error_log(print_r($_FILES, true));
        $target_file_name = basename($_FILES['logo']['name']);
        $target_file = $target_dir . $target_file_name;

        // max 500kB
        if($_FILES['logo']['size'] > 500000){
            return Status::errorStatus('provided file is too large (max 500kB allowed)');
        }

        $check = getimagesize($_FILES["logo"]["tmp_name"]);
        if($check !== false){
            // check for allowed mime-types
            if($check['mime'] != "image/png"
                && $check['mime'] != "image/jpeg"
                && $check['mime'] != "image/jpg"
                && $check['mime'] != "image/gif"
            ){
                return Status::errorStatus('image type not supported (use: jpeg, jpg, gif, png)');
            }

            // move file to uploads
            if(move_uploaded_file($_FILES["logo"]['tmp_name'], $target_file)){
                return Status::successDataResponse('logo uploaded', [ 'uri' => '/uploads/logos/'.$target_file_name]);
            }else{
                return Status::errorStatus('Cannot upload image, unknown error while storing the file');
            }
        } else {
            return Status::errorStatus('provided file is not an image');
        }  
    }

    function get_logos($configuration){

    }

?>
    
