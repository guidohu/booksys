// Check if mobile browser first,
// then continue with checks
$(function() {
    if(BooksysBrowser.isMobile()){
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
