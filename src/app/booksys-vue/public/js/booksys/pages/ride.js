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

        $("#body").load("res/ride_mobile.html", function(){
            loadContent();
        });
    } else {
        $("#body").load("res/ride.html", function(){
            loadContent();
        });
    }
});

// Stay in the WebApp
$(function() {
    $.stayInWebApp();
});

// get date parameters (to query TODAY)
let date_start = moment().tz(getTimeZone()).startOf('day');
let date_end   = moment().tz(getTimeZone()).endOf('day');

// store the session array that is displayed in the pie
// and keep track of the currently selected index
let pieSessions = new Array();
let selectedSessionIdx = -1;

// Variable that gives access to the watch instance
let watch;

// gets the user's timezone or returns the default
function getTimeZone(){
    if (typeof(Storage) !== "undefined") {
        let timezone = localStorage.getItem("timezone");
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

// Load dynamic content
function loadContent(){
    $("#btn_continue").hide();
    $("#status-panel").hide();

    // create a new stop watch
    watch = new BooksysViewStopWatch();

    // get the specific session ID from the URL
    let searchParams   = new URLSearchParams(window.location.search);
    let uriSessionId   = searchParams.get('sessionId');
    let selectedSector = searchParams.get('selectedSessionId');

    // switch to stop watch for the following reasons:
    // - the user is already taking the time currently
    // - the link is to a specific session
    // otherwise display the selection of sessions for today
    if(watch.isActive()){
        let sessionId = watch.getSessionId();
        window.location.href = '/watch.html?sessionId='+sessionId;
    }else if(uriSessionId != null){
        window.location.href = '/watch.html?sessionId='+sessionId;
    }else{
        updateBookings(date_start, date_end, selectedSector);
    }
}

// Requests the new data from the server and
// updates the view
function updateBookings(date_start, date_end, selectedSectorSessionId){
    // get the session-information from the database
    var data = {
        start: date_start.format('X'),
        end:   date_end.format('X')
    };
    
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=get_booking_day",
        data: JSON.stringify(data),
        success: function(resp){
            var json = $.parseJSON(resp);

            // in case we are not logged in --> redirect to login page
            if(typeof json.redirect != "undefined"){
                window.location.replace(json.redirect);
                return;
            }

            pieSessions = BooksysPie.drawPie("pie", json, function(id){
                updateDetail(id);
            }, {
                containerHeight:  300,
                containerWidth:   300,
                circleX:          150,
                circleY:          150,
                circleRadius:     100,
                animate:          true,
                labels:           false
            });
            $("#status-msg-select-session").html('Loading data finished.');
            $("#status-panel-select-session").fadeOut();
            
            if(selectedSessionIdx != -1){
                BooksysPie.selectSector(selectedSessionIdx);
                updateDetail(selectedSessionIdx);
            }else if(selectedSectorSessionId != null){
                // get the index of the selected session ID
                for(let i=0; i< pieSessions.length; i++){
                    if(pieSessions[i].id == selectedSectorSessionId){
                        BooksysPie.selectSector(i);
                        updateDetail(i);
                    }
                }
            }else{
                // show the session details but without a selected session
                BooksysViewSessionDetails.display("session_details", null, null, function(){});
            }
        }
    });
}

// Updates session information
function updateDetail(idx){
    selectedSessionIdx = idx;
    selectedSession    = pieSessions[idx];
    
    // Display the continue button if riders are available
    if(typeof selectedSession.riders != "undefined" && selectedSession.riders.length > 0){
        $("#btn_continue").fadeIn();
    }else{
        $("#btn_continue").fadeOut();
    }

    // In case this is a free slot, we like to provide some
    // presets for the creation of a new session
    var presets = null;
    if(selectedSession.id == null){
        presets             = new Object();
        presets.title       = null;
        presets.description = null;
        presets.start       = selectedSession.start.tz(getTimeZone());
        presets.end         = selectedSession.end.tz(getTimeZone());
        presets.maxRiders   = 10;
        presets.type        = 0; // 0: private
    }

    // Display the proper session details and regiser
    // required callback functions
    var callbacks = new Object();
    callbacks.createSession = function(){
        BooksysViewSessionEditor.display("user_dialog", selectedSession.id, presets, function(){
            BooksysViewSessionEditor.hideApp();
            BooksysViewSessionEditor.destroyView("user_dialog");
            updateBookings(date_start, date_end);
        });
    };
    callbacks.deleteSession = function(id){
        deleteSession(id);
    }
    callbacks.addRider = function(id){
        BooksysViewSessionRiderSelection.display("user_dialog", id, function(){
            BooksysViewSessionRiderSelection.destroyView("user_dialog");
            updateBookings(date_start, date_end);
        });
    }
    callbacks.removeRider = function(sessionId, userId){
        removeRider(userId, sessionId);
    };
    BooksysViewSessionDetails.display("session_details", selectedSession.id, presets, callbacks);
}

// delete a session
function deleteSession(id){
    var req = {
        session_id: id
    };
    $.ajax({
        type: "POST",
        url:  "api/booking.php?action=delete_session",
        data: JSON.stringify(req),
        success: function(data){
            updateBookings(date_start, date_end);
        },
        error: function(xhr,status,error){
            console.log(xhr);
            console.log(status);
            console.log(error);
            var errorMsg = $.parseJSON(xhr.responseText);
            alert(errorMsg.error);
        },
    });
}

// delete user from session
function removeRider(userId, sessionId){
    // delete the rider from the session
    var data = {
        session_id: sessionId,
        user_id:    userId
    }
    
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=delete_user",
        data: JSON.stringify(data),
        success: function(resp){
            if(resp != ""){
                alert(resp);
            }else{
                updateBookings(date_start, date_end);
            }
        },
        error: function(resp){
            alert("Ooops. There was an error.");
        }
    });
}

function continueToWatch(){
    window.location.href = '/watch.html?sessionId='+selectedSession.id;
}