// Check if mobile browser first,
// then continue with checks
$(function() {
    if(BooksysBrowser.isMobile()){
        // Make it behave like an app
        BooksysBrowser.setViewportMobile();
        BooksysBrowser.setManifest();
        BooksysBrowser.setMetaMobile();
        // Add mobile style dynamically
        BooksysBrowser.addMobileCSS();

        $("#body").load("res/index_mobile.html", function(){
            $.stayInWebApp();
            setupEventListeners();
            backendStatus();
        });
    } else {
        $("#body").load("res/index.html", function(){
            backendStatus();
            setupEventListeners();
        });
    }
});

function setupEventListeners(){
    // register custom submit action
    document.getElementById("login_form").addEventListener("submit", function(event){
        event.preventDefault();
        login();
    });
}

// Check backend status and abort or continue
function backendStatus(){
    $("#status_panel").hide();

    BooksysBackend.getStatus({
        green: function(text, redirect){
            console.log("backend status ok");
            isLoggedIn();
            setFocusToUsername();
        },
        red: function(text, redirect){
            switch(text){
                case "no database configured":
                    console.log("Server reports: " + text); 
                    window.location.href = "/setup.html";
                    break;
                case "no user configured":
                    console.log("Server reports: " + text); 
                    window.location.href = "/setup.html";
                    break;
                case "no config":
                    console.log("Server reports: " + text);
                    window.location.href = "/setup.html";
                    break;
                default:
                    $('#login_form').hide();
                    $("#status_panel").show();
                    console.error("backend status not ok");
                    console.error(text);
                    var div = document.createElement('div');
                    div.className = 'alert alert-danger';
                    div.innerHTML = "The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible.";
                    document.getElementById("status_body").appendChild(div);
            }
        }
    });
}

// Check if we are already logged in
function isLoggedIn() {
    let login = new BooksysLoginCheck();
    login.check(function(){
        window.location.href = "dashboard.html";
    });
}

// Sets the focus to the username input
function setFocusToUsername() {
    $('#username').focus();
}

// Calculate the SHA-256 hash of the password on user-side
// to never transmit a user's password in clear text
function calcHash(password) {
    try {
        var hashObj = new jsSHA(password, "TEXT");
        var hashOutput = hashObj.getHash("SHA-256","HEX");
        return hashOutput;
    } catch(e) {
        alert("Password could not be sent encrypted.");
        return false;
    }
}

// Perform the actual login
function login(){
    // protect logins through recaptcha
    let userInput     = document.getElementById("username").value;
    let passwordInput = document.getElementById("password").value;
    let hashPwd       = calcHash(passwordInput);
    if(!hashPwd){
        alert("Cannot calculate hash");
        return;
    }

    let data = new Object();
    data['username']  = userInput;
    data['password']  = hashPwd;
    data['remember']  = 1;

    // check that a username and password was specified
    if(userInput == "" || passwordInput == ""){
        setStatus("no username or password given");
        return;
    }

    $.ajax({
        type: "POST",
        url: 'api/login.php?action=login',
        data: JSON.stringify(data),
        success: function(resp){
            console.log(resp);
            var json = $.parseJSON(resp);
            if(json.login_status_code==-1){
                setStatus("incorrect username or password");
            } else if(json.login_status_code==-2){
                setStatus("please be patient until we activate your account");
            } else if(json.login_status_code==1){
                // successful login
                storeMyUser(function(){
                    window.location.href = "dashboard.html";
                });
            } else{
                setStatus("incorrect username or password");
            }
        },
        error: function(data, text, errorCode){
            // something went wrong
            setStatus("unknown error"); 
        }

    });
}

function storeMyUser(callback){
    $.ajax({
        type: "GET",
        url: 'api/user.php?action=get_my_user',
        success: function(resp){
            console.log(resp);
            var json = $.parseJSON(resp);
            
            // store the role in local storage
            storeUserInStorage(json);
            callback();
        },
        error: function(data, text, errorCode){
            // something went wrong
            // var div = document.createElement('div');
            // div.className = 'label label-danger';
            // div.innerHTML = "unknown error";
            // document.getElementById("login_label").appendChild(div);  
        }

    });
}

// add user information in local storage
function storeUserInStorage(userInfo){
    if (typeof(Storage) !== "undefined") {
        localStorage.setItem('userRole', userInfo.user_role_id);
        localStorage.setItem('userRoleName', userInfo.user_role_name);
    } else {
        // Sorry! No Web Storage support..
    }
}

// remove user information in local storage
function removeUserInStorage(){
    if (typeof(Storage) !== "undefined") {
        localStorage.removeItem('userRole');
        localStorage.removeItem('userRoleName');
    }
}

function setStatus(text){
    // remove a previous status
    let status = document.getElementById("login_label"); 
    let child = status.lastElementChild;  
    while (child) { 
        status.removeChild(child); 
        child = status.lastElementChild; 
    }

    var div = document.createElement('label');
    div.className = 'label label-danger';
    div.innerHTML = text;
    document.getElementById("login_label").appendChild(div);
    return;
}