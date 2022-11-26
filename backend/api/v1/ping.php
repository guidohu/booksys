<?php
	// automatically load all classes
	spl_autoload_register('resources_autoloader');
	function resources_autoloader($class){
		include '../../classes/'.$class.'.php';
    }

    // check if we have an action
	if(!isset($_GET['action'])){
		echo json_encode(Status::errorStatus("API action not known"));
		exit;
    }

    $response = '';
    switch($_GET['action']){
        case 'ping':
            $response = Status::successStatus("ok");
            break;
        default:
            $response = Status::errorStatus("action unknown for this API call");
            break;
    }

    echo json_encode($response);
    return;

?>
    
