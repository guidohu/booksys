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

        $("#body").load("res/account_mobile.html", function(){
            loadContent();
        });
    } else {
        $("#body").load("res/account.html", function(){
            loadContent();
        });
    }
});

function loadContent(){
    isLoggedIn();
    getAccount();
    validation();
    getPaymentInfo();
}

var account;

// Set validator defaults
$.validator.setDefaults({
    errorElement: "validate_msg",
});

// Check if we are already logged in
function isLoggedIn() {
    $.ajax({
        type: "GET",
        url: 'api/login.php?action=isLoggedIn',
        success: function(data){
            var json = $.parseJSON(data);
            if(!json.loggedIn){
                // we are not logged in, redirect to login page
                window.location.href = "/index.html";
            }
        },
        error: function(data, text, errorCode){
            // something went wrong
            console.error("Cannot call isLoggedIn(), error returned");
        }
    });
}

// Setup validation functions
function validation(){
    $('#change_user_data').validate({
        submitHandler: function(form){
            submitEditUserData();
        },
        rules: {
            user_first_name: {
                required: true,
            },
            user_last_name: {
                required: true,
            },
            user_address: {
                required: true,
            },
            user_mobile: {
                required: true,
                mobileNumber: true,
            },
            user_plz: {
                required: true,
                number: true,
            },
            user_city: {
                required: true,
            },
            user_email: {
                required: true,
                email: true,
            },
        },
        messages: {
            user_first_name: {
                required: 'no first name given',
            },
            user_last_name: {
                required: 'no last name given',
            },
            user_address: {
                required: 'no address given',
            },
            user_mobile: {
                required: 'no mobile number given',
                number: 'this is no valid phone number',
            },
            user_plz: {
                required: 'no zip code given',
                number: 'this is no zip code',
            },
            user_city: {
                required: 'no city given',
            },
            user_email: {
                required: 'no email given',
                email: 'this is no valid email',
            }
        },		
    });
    
    $('#change_user_password').validate({
        submitHandler: function(form){
            submitEditPassword();
        },
        rules: {				
            user_password_old: {
                required: true,
            },
            user_password_new: {
                required: true,
            },
            user_password_new_2: {
                equalTo: '#user_password_new',
                required: true,
            },
        },
        messages: {
            user_password_old: {
                required: 'please provide your current password',
            },
            user_password_new: {
                required: 'no new password provided',
            },
            user_password_new_2: {
                equalsTo: 'password not identical',
            },
        },
    });
}

function getPaymentInfo(){
    $.ajax({
        type: "GET",
        url: 'api/configuration.php?action=get_customization_parameters',
        success: function(data){
            var json = JSON.parse(data);

            // add the owner
            if(json.payment_account_owner != null){
                $('#payment_account_owner').html(escapeHTML(json.payment_account_owner));
            }else{
                $('#payment_account_owner').html("not available");
            }
            // add the owner
            if(json.payment_account_iban != null){
                $('#payment_account_iban').html(escapeHTML(json.payment_account_iban));
            }else{
                $('#payment_account_iban').html("not available");
            }
            // add the owner
            if(json.payment_account_bic != null){
                $('#payment_account_bic').html(escapeHTML(json.payment_account_bic));
            }else{
                $('#payment_account_bic').html("not available");
            }
            // add the owner
            if(json.payment_account_comment != null){
                $('#payment_account_comment').html(escapeHTML(json.payment_account_comment));
            }else{
                $('#payment_account_comment').html("not available");
            }
        },
        error: function(data, text, errorCode){
            // something went wrong
            console.error("Cannot call get_customization_parameters");
            console.error(text);
            console.log(data);
            console.log(errorCode);
            $('#payment_account_owner').html("not available");
            $('#payment_account_iban').html("not available");
            $('#payment_account_bic').html("not available");
            $('#payment_account_comment').html("not available");
        }
    });
}

// Returns the account information
function getAccount(){
    // get the account information
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_footer').hide();
    $('#status_header').html('<h4><span class="glyphicon glyphicon-refresh"></span> Loading</h4>');
    $('#status_body').html('Please wait while user data is retrieved...');
    $('#status_modal').modal('show');
    
    $.ajax({
        type: "GET",
        url: "api/user.php?action=get_my_user",
        success: function(data){
            var json = $.parseJSON(data);
            updateAccount(json);
            $('#status_modal').modal('hide');
        }
    });
    
    $.ajax({
        type: "GET",
        url: "api/user.php?action=get_my_user_heats",
        success: function(data){
            var data = $.parseJSON(data);
            updateStatistics(data);
            $('#status_modal').modal('hide');
        }
    });
}

function updateAccount(info){
    $("#first_name").html(info.first_name);
    $("#last_name").html(info.last_name);
    $("#address").html(info.address);
    $("#plz").html(info.plz);
    $("#city").html(info.city);
    $("#email").html(info.email);
    $("#mobile").html(info.mobile);
    if(info.license == 1){
        $("#license").html("Driving License: Yes");
    }else{
        $("#license").html("Driving License: No");
    }
    account = info;
}

function updateStatistics(info){
    $("#statistic_time").html(
        "Riding Time: " + info.heat_time_min_ytd + " min</br>(Total: " + info.heat_time_min + " min)"
        + "</br>Riding Cost: " + info.heat_cost_ytd + " " + info.currency + "</br>(Total: " + info.heat_cost + " " + info.currency + ")"
    );
                            
    $("#balance_current").html("Your Balance: " + info.balance_current + " " + info.currency);
    
    var content = '';
    
    for(var i=0; i<info.heats.length; i++){
        if(i%2==0){
            content += '<tr class="bc_bg_even bc_bg_hover">';
        }else{
            content += '<tr class="bc_bg_odd bc_bg_hover">';
        }
        content += '<td>' + info.heats[i].date + '</td>';
        content += '<td>' + info.heats[i].cost + ' ' + info.currency + '</td>';
        content += '<td>' + info.heats[i].duration + ' min</td>';
    }
    content += '</tr>';
    
    $('#heat_list').html(content);
}

// Hides the status window
function hideStatus(){
    $('#status_modal').modal('hide');
}

// Shows a status message with the given title and message
function showStatus(title, message){
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    
    var header = '<h4>'
                + title
                + '</h4>';		
    $('#status_header').html(header);
    $('#status_body').html(message);
    $("#status_footer").show();
    
    $('#status_modal').modal('show');
}

// Display the payment information
function showPaymentInfo(){
    // show the dialog
    $('#payment_info_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $("#payment_info_modal").modal('show');
}

// Hide the payment information again
function hidePaymentInfo(){
    $("#payment_info_modal").modal('hide');
}

function showHeatInfo(){
    // show the dialog
    $('#heat_info_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $("#heat_info_modal").modal('show');
}

function hideHeatInfo(){
    $("#heat_info_modal").modal('hide');
}

function showEditUserData(){
    // fill form fields
    $('#user_first_name').val(account.first_name);
    $('#user_last_name').val(account.last_name);
    $('#user_address').val(account.address);
    $('#user_plz').val(account.plz);
    $('#user_city').val(account.city);
    $('#user_email').val(account.email);
    $('#user_mobile').val(account.mobile);
    $('#change_user_data').valid();
    
    if(account.license == 1){
        $('#user_license').val('1');
    }else{
        $('#user_license').val('0');
    }

    // show the user data dialog
    $('#user_data_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#user_data_modal').modal('show');
}

function hideEditUserData(){
    $('#user_data_modal').modal('hide');
}

function submitEditUserData(){
    // collect all the info
    var data = new Object();
    data['first_name'] = $('#user_first_name').val();
    data['last_name'] = $('#user_last_name').val();
    data['address'] = $('#user_address').val();
    data['plz'] = $('#user_plz').val();
    data['city'] = $('#user_city').val();
    data['email'] = $('#user_email').val();
    data['mobile'] = $('#user_mobile').val();
    data['license'] = $('#user_license').val();
    
    $('#user_edit_save').html('Sending...');
    $.ajax({
        type: "POST",
        url:  "api/user.php?action=change_my_user_data",
        data: JSON.stringify(data),
        success: function(data, textStatus, jqXHR){
            $('#user_edit_save').html('OK');
            hideEditUserData();
            getAccount();				
        },
        error: function(data, text, errorCode){
            $('#user_edit_save').html('OK');
            hideEditUserData();
            showStatus('Error', data.responseText);
        },
    });		
}

function showEditPassword(){
    // show the password dialog
    $('#user_password_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#user_password_modal').modal('show');
}

function hideEditPassword(){
    $('#user_password_old').val('');
    $('#user_password_new').val('');
    $('#user_password_new_2').val('');
    $('#user_password_modal').modal('hide');
}

function submitEditPassword(){	
    // calculate hashes
    var data = new Object();
    data['password_old'] = calcHash('user_password_old');
    data['password_new'] = calcHash('user_password_new');
    
    
    $('#user_password_save').html('Sending...');
    $.ajax({
        type: "POST",
        url: "api/password.php?action=change_password_by_password",
        data: JSON.stringify(data),
        success: function(data, textStatus, jqXHR){
            let response = JSON.parse(data);
            if(response.ok){
                $('#user_password_save').html('OK');
                hideEditPassword();
                showStatus('Password Change Successful', 'Your password has been changed.');
            }else{
                if(response.message == "You are not logged in."){
                    console.log('user not logged in, redirecting to login page');
                    window.location.href = "/index.html";
                }else{
                    $('#user_password_save').html('OK');
                    hideEditPassword();
                    showStatus('Error', response.message);
                }
            }
        },
        error: function(data, text, errorCode){
            $('#user_password_save').html('OK');
            hideEditPassword();
            showStatus('Error', data.responseText);
        },
    });
    
}

// Replaces the  password string with it's SHA-256 hash.
function calcHash(id) {
    try {
        var hashInput = document.getElementById(id).value;
        var hashObj = new jsSHA(hashInput, "TEXT");
        var hashOutput = hashObj.getHash("SHA-256","HEX");
        return hashOutput;
    } catch(e) {
        alert("Password could not be sent encrypted: " + e);
        return;
    }
}

function escapeHTML(text) {
    return text.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}