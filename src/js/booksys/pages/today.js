// Keep start/end of displayed day available
var startDay = null;
var endDay   = null;

// Sessions that are displayed in the pie
var pieSessions = new Array();
// the ID of the pie that is currently selected
var selectedSessionIdx = -1;


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

        $("#body").load("res/today_mobile.html", function(){
            loadContent();
        });
    } else {
        $("#body").load("res/today.html", function(){
            loadContent();
        });
    }
});

// loads the content of the site
function loadContent(){
    $.stayInWebApp();

    // get date from GET parameter
    var dateStr = getUrlVars()['getdate'];
    if(dateStr == null){
        dateStr = moment().tz(getTimeZone()).format("YYYY-MM-DD");
    }

    // get start/end of day
    moment.tz.add("Europe/Berlin|CET CEST CEMT|-10 -20 -30|01010101010101210101210101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010|-2aFe0 11d0 1iO0 11A0 1o00 11A0 Qrc0 6i00 WM0 1fA0 1cM0 1cM0 1cM0 kL0 Nc0 m10 WM0 1ao0 1cp0 dX0 jz0 Dd0 1io0 17c0 1fA0 1a00 1ehA0 1a00 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1o00 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 WM0 1qM0 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 WM0 1qM0 WM0 1qM0 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 WM0 1qM0 11A0 1o00 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 WM0 1qM0 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 WM0 1qM0 11A0 1o00 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 WM0 1qM0 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 WM0 1qM0 WM0 1qM0 11A0 1o00 11A0 1o00|41e5");
    var day   = moment.tz(dateStr + "T00:00:00", getTimeZone());
    var start = day.tz(getTimeZone()).startOf('day').clone();
    var end   = day.tz(getTimeZone()).endOf('day').clone();
    startDay  = start.clone();
    endDay    = end.clone();

    updateBookings(start, end);
}

// returns the URL variables
function getUrlVars() {
    var vars = {};
    var parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
        vars[key] = value;
    });
    return vars;
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

// Go to the next day
function nextDay(){
    selectedSessionIdx = -1;
    BooksysPie.resetSelection();
    startDay = startDay.tz(getTimeZone()).add(1, 'day').startOf('day');
    endDay   = endDay.tz(getTimeZone()).add(1, 'day').endOf('day');
    updateBookings(startDay, endDay);
}

// Go to the previous day
function prevDay(){
    selectedSessionIdx = -1;
    BooksysPie.resetSelection();
    startDay = startDay.tz(getTimeZone()).subtract(1, 'day').startOf('day');
    endDay   = endDay.tz(getTimeZone()).subtract(1, 'day').endOf('day');
    updateBookings(startDay, endDay);
}

// Requests the new data from the server and
// updates the view
function updateBookings(start, end){
    // get the session-information from the database
    var data = {
        start:	start.tz(getTimeZone()).format("X"),
        end:    end.tz(getTimeZone()).format("X")
    };
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=get_booking_day",
        data: JSON.stringify(data),
        success: function(resp){
            var json = $.parseJSON(resp);

            // in case we are not logged in --> redirect to login page
            if(typeof json.redirect != "undefined"){
                window.location.replace = json.redirect;
                return;
            }
            
            // Some general info which has to be adjusted
            $("#titleDate").text(start.tz(getTimeZone()).format("dddd DD.MM.YYYY"));
            var sunrise = moment(json.sunrise, "X").tz(getTimeZone());
            var sunset  = moment(json.sunset, "X").tz(getTimeZone());
            $("#detail_sunrise").html(sunrise.format('HH:mm'));
            $("#detail_sunset").html(sunset.format('HH:mm'));
            
            // Draw the pie and get the pie content returned
            let properties = {
                containerHeight: 454,
                containerWidth:  700,
                circleX:         350,
                circleY:         199,
                circleRadius:    100,
                animation:       true,
                labels:          true,
            }
            if(BooksysBrowser.isMobile()){
                properties = {
                    containerHeight: 300,
                    containerWidth:  350,
                    circleX:         175,
                    circleY:         150,
                    circleRadius:    100,
                    animation:       false,
                    labels:          true,
                }
                BooksysViewSessionDetails.display("session_details", null, null, function(){});
            }
            
            pieSessions = BooksysPie.drawPie("pie", json, function(id){
                updateDetail(id);
            }, properties);

            // we might need to display a specific session
            if(selectedSessionIdx != -1){
                BooksysPie.selectSector(selectedSessionIdx);
                updateDetail(selectedSessionIdx);
            }else{
                // show the session details but without a selected session
                BooksysViewSessionDetails.display("session_details", null, null, function(){});
            }
        }
    });	
}

// Updates the details on the right side
function updateDetail(idx){
    selectedSessionIdx     = idx;
    let session            = pieSessions[idx];

    // display details mobile version
    let presets = null;
    let timezone = getTimeZone();
    if(session.id == null){
        presets             = new Object();
        presets.title       = null;
        presets.description = null;
        presets.start       = session.start.tz(getTimeZone());
        presets.end         = session.end.tz(getTimeZone());
        presets.maxRiders   = 10;
        presets.type        = 0; // 0: private
    }
    // Display the proper session details and regiser
    // required callback functions
    var callbacks = new Object();
    callbacks.createSession = function(){
        BooksysViewSessionEditor.display("user_dialog", session.id, presets, function(){
            BooksysViewSessionEditor.hideApp();
            BooksysViewSessionEditor.destroyView("user_dialog");
            updateBookings(startDay, endDay);
        });
    };
    callbacks.deleteSession = function(id){
        deleteSession(id);
    };
    callbacks.addRider = function(id){
        BooksysViewSessionRiderSelection.display("user_dialog", id, function(){
            BooksysViewSessionRiderSelection.destroyView("user_dialog");
            updateBookings(startDay, endDay, session.id);
        });
    };
    callbacks.removeRider = function(sessionId, userId){
        // TODO supply this function within BooksysViewSessionDetails to be used
        // directly instead of implementing it everywhere where this module is used
        deleteRiderSession(userId, sessionId);
    };
    BooksysViewSessionDetails.display("session_details", session.id, presets, callbacks);
}

// deletes a rider from a session
function deleteRiderSession(user_id, session_id){
    // delete the rider from the session
    var data = {
        session_id: session_id,
        user_id:    user_id
    }
    
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=delete_user",
        data: JSON.stringify(data),
        success: function(resp){
            if(resp != ""){
                alert(resp);
            }else{
                updateBookings(startDay, endDay, session_id);
            }
        },
        error: function(resp){
            alert("Ooops. There was an error.");
        }
    });
}

// Deletes a session
function deleteSession(){
    var req = {
        session_id:		pieSessions[selectedSessionIdx].id
    };

    // get the session-information from the database
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=delete_session",
        data: JSON.stringify(req),
        success: function(data){	
            updateBookings(startDay, endDay);
        },
        error: function(data){
            alert('Cannot delete session.');
        },
    });
}