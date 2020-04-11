// Initialize reCAPTCHA
var recaptchaOnloadCallback = function(){
    recaptcha.render(document.getElementById("captcha"));
};
let recaptcha = new BooksysRecaptcha("recaptchaOnloadCallback");

// Check if mobile browser first,
// then continue with the page logic
$(function() { 
    if(BooksysBrowser.isMobile()){
        // Make it behave like an app
        BooksysBrowser.setViewportMobile();
        BooksysBrowser.setManifest();
        BooksysBrowser.setMetaMobile();
        // Add mobile style dynamically
        BooksysBrowser.addMobileCSS();
        $.stayInWebApp();

        loadContent();
    } else {
        // Do not add mobile style in desktop mode
        loadContent();
    }
});

function loadContent(){
    let pwResetEmail = new BooksysViewPasswordResetEmail();
    let pwResetEmailCallback = {
        cancel: function(){
            pwResetEmail.destroyView();
            window.location.href = "/index.html";
        },
    };
    pwResetEmail.display("password_reset", pwResetEmailCallback, recaptcha);
}


// Stay in the WebApp
$(function() {
    $.stayInWebApp();
});

// the user id
$user_id = 0;

// Set validator defaults
$.validator.setDefaults({
    errorElement: "validate_msg",
    invalidHandler: function(event, validator){
        var errors = validator.numberOfInvalids();
        if(errors){		
            $('#status_header').html('<h4>Error while signing up</h4>');
            $('#status_body').html('Please check your input again');
            $('#status_body').show();
            $('#status_footer').show();
            $('#error').show();
        }
    },
});

$(document).ready(function() {
    $('#password_form').hide();
    
    $('#password_form').validate({
        //debug: true,
        submitHandler: function(form){
            changePassword();
        },
        rules: {
            token: {
                required:true,
                number: true,
            },
            password: {
                required:true,
            },
            password_conf: {
                equalTo:'#password',
            }
        },
        messages: {
            token: {
                required: 'please provide the token',
                number:   'token has to be a number',
            },
            password: {
                required: 'no password provided',
            },
            password_conf: {
                equalTo: 'passwords are not identical'
            }
        }
    });
});	

function changePassword(){		
    $('#status_header').html('<h4>Resetting your Password</h4>');
    $('#status_body').hide();
    $('#status_footer').hide();
    $('#error').hide();
    $('#success').hide();
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    
    // send the ajax call to the api
    var data         = new Object();
    data['token']    = $('#token').val();
    data['password'] = calcHash();
    data['user_id']  = $user_id;
    
    $.ajax({
        type: "POST",
        url: "api/password.php?action=change_password_by_token",
        data: JSON.stringify(data),
        success: function(data, textStatus, jqXHR){
            $('#status_header').html("<h4>Successfully changed password</h4>");
            $('#status_header').show();
            $('#status_footer').show();
            $('#status_body').html('Please check your Email Inbox for the Password reset Code.');
            $('#status_body').show();
            $('#error').hide();
            $('#success').show();
            window.location.href = 'index.html?valid=4';
        },
        error: function(data, text, errorCode){
            $('#status_header').html('<h4>Error while changing password</h4>');
            $('#status_body').html(data.responseText);
            $('#status_body').show();
            $('#status_footer').show();
            $('#error').show();
        },			
    });
}

function showPasswordForm(){
    $('#email_form').hide();
    $('#password_form').show();
    hideStatus();
}

// Hides the status window
function hideStatus(){
    $("#status_modal").modal("hide");
}

// Calculates the hash of the password.
function calcHash() {
    try {
        var hashInput = document.getElementById("password").value;
        var hashObj = new jsSHA(hashInput, "TEXT");
        var hashOutput = hashObj.getHash("SHA-256","HEX");
        return hashOutput;
    } catch(e) {
        alert("Password could not be sent encrypted: " + e);
        return;
    }
}