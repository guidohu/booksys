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
        $.stayInWebApp();

        $("#body").load("res/dashboard_mobile.php", function(){
            finalizeUI();
        });
    } else {
        $("#body").load("res/dashboard.php", function(){
            finalizeUI();
        });
    }
});

function finalizeUI() {
    isLoggedIn();
    updateBookings();
    checkDatabaseUpdate();
    getVersion();
}

// Check if we are already logged in
function isLoggedIn() {
    let login = new BooksysLoginCheck();
    login.check(null, function(){
        window.location.href = "index.html";
    });
}

function checkDatabaseUpdate(){
    // check whether a database update is required
    if(!isAdminUser()){
        return;
    }

    $.ajax({
        type: "GET",
        url: "api/backend.php?action=admin_check_database_update",
        cache: false,
        success: function(resp){
            var json = $.parseJSON(resp);

            if(json.updateAvailable == true){
                BooksysViewDatabaseUpdate.display("user_dialog", function(){
                    BooksysViewDatabaseUpdate.destroyView("user_dialog");
                });
            }
            console.log(json);
        },
        error: function(xhr, status, errorCode){
            console.log(xhr);
            console.log(status);
            console.log(errorCode);
            return;
        }
    });
}

// get the versions from the backend
function getVersion(){
    $.ajax({
        type: "GET",
        url: "api/backend.php?action=get_version",
        cache: false,
        success: function(resp){
            var json = $.parseJSON(resp);
            console.log(json);
            if(document.getElementById("version") != null){
                document.getElementById("version").innerHTML = "v"+json.app_version;
            }
        },
        error: function(xhr, status, errorCode){
            console.log(xhr);
            console.log(status);
            console.log(errorCode);
            return;
        }
    });
}

// Requests the new data from the server and
// updates the view
function updateBookings(){
    var date_start = moment().tz(getTimeZone()).startOf('day');
    var date_end   = moment().tz(getTimeZone()).endOf('day');

    var data = Object();
    data['start'] = date_start.format("X");
    data['end']   = date_end.format("X");

    // get the session-information from the database
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=get_booking_day",
        data: JSON.stringify(data),
        cache: false,
        success: function(resp){
            var json = $.parseJSON(resp);

            // in case we are not logged in --> redirect to login page
            if(typeof json.redirect != "undefined"){
                window.location.replace(json.redirect);
                return;
            }

            var properties = {
                containerWidth: 80,
                containerHeight: 85,
                animate: false,
            };

            BooksysPie.drawPie("today", json, null, properties);
        }
    });
}

// gets the user's timezone or returns the default
function getTimeZone(){
    if (typeof(Storage) !== "undefined") {
        var timezone = localStorage.getItem("timezone");
        if(timezone != null){
            return timezone;
        }
        else{
            // TODO let user set his timezone
            return "Europe/Berlin";
        }
    } else {
        // Sorry! No Web Storage support..
        return "Europe/Berlin";
    }
}

// returns true in case the user is an admin user
function isAdminUser(){
    if (typeof(Storage) !== "undefined") {
        // id: 3 -> admin
        // id: 2 -> member
        // id: 1 -> guest
        var roleId = localStorage.getItem('userRole');
        if(roleId == 3){
            return true;
        }
    }
    return false;
}