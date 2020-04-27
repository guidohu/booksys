// Keep start/end of displayed day available
var startDay = null;
var endDay   = null;

// Sessions that are displayed in the pie
var pieSessions = [];
// the ID of the pie that is currently selected
var selectedSessionIdx = null;


// Check if mobile browser first,
// then continue with checks
$(function() { 
    if(BooksysBrowser.isMobile()){
        // Make it behave like an app
        // BooksysBrowser.setViewportMobile();
        // BooksysBrowser.setManifest();
        // BooksysBrowser.setMetaMobile();
        // Add mobile style dynamically
        // BooksysBrowser.addMobileCSS(); (not required here)

        $("#body").load("res/today.html", function(){
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

    $("#detail_session_create").hide();
    $("#detail_rider_book").hide();
    $("#create_session_modal").modal("hide");
    updateBookings(start, end);
    $('#accordion').collapse();
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
    startDay = startDay.tz(getTimeZone()).add(1, 'day').startOf('day');
    endDay   = endDay.tz(getTimeZone()).add(1, 'day').endOf('day');
    $("#detail_session_create").hide();
    updateBookings(startDay, endDay);
}

// Go to the previous day
function prevDay(){
    startDay = startDay.tz(getTimeZone()).subtract(1, 'day').startOf('day');
    endDay   = endDay.tz(getTimeZone()).subtract(1, 'day').endOf('day');
    $("#detail_session_create").hide();
    updateBookings(startDay, endDay);
}

// Requests the new data from the server and
// updates the view
// - idx	the position of the session in the pieSession array that should be updated and highlighted
function updateBookings(start, end, sessionId){
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
                location.href = json.redirect;
                return;
            }
            
            // Some general info which has to be adjusted
            $("#titleDate").text(start.tz(getTimeZone()).format("dddd DD.MM.YYYY"));
            var sunrise = moment(json.sunrise, "X").tz(getTimeZone());
            var sunset  = moment(json.sunset, "X").tz(getTimeZone());
            $("#detail_sunrise").html(sunrise.format('HH:mm'));
            $("#detail_sunset").html(sunset.format('HH:mm'));
            $("#detail_date").html("please select a session");
            $("#detail_description").html("-");
            $("#detail_rider_count").html(0);
            $("#detail_rider_max_space").html(0);
            $("#detail_rider_list").html("");
            $("#menu_booking").hide();
            $("#menu_session").hide();
            
            // Draw the pie and get the pie content returned
            var properties = {
                containerHeight: 454,
                containerWidth:  700,
                circleX:         350,
                circleY:         199,
                circleRadius:    100,
                animation:       true,
                labels:          true,
            };
            pieSessions = BooksysPie.drawPie("pie", json, updateDetail, properties);

            // we might need to display a specific session
            if(sessionId != null){
                selectedSessionIdx = getPieIdBySessionId(sessionId);
                BooksysPie.selectSector(selectedSessionIdx);
                updateDetail(selectedSessionIdx);
            }else{
                selectedSessionIdx = null;
            }
        }
    });	
}

// Updates the details on the right side
function updateDetail(idx){
    selectedSessionIdx     = idx;
    var session            = pieSessions[idx];

    // display session information
    if(session.title){
        $("#detail_title").text(session.title.substring(0,24));
    }else{
        $("#detail_title").text('-');
    }
    if(session.start){
        $("#detail_date").text(session.start.tz(getTimeZone()).format('HH:mm') + " - " + session.end.tz(getTimeZone()).format('HH:mm'));
    }else{
        $("#detail_date").text("-");
    }
    if(session.comment){
        $("#detail_description").text(session.comment.substring(0,24));
    }else{
        $("#detail_description").text("-");
    }
    
    // TODO (if this page gets displayed to others than admins) check if user is admin and if we have to display the 'create' button
    $("#menu_session").show();
    if(session.id == null){
        // no session yet
        $("#detail_session_create").show();
        $("#menu_session").hide();
        $("#detail_rider_book").hide();
        $("#menu_booking").hide();
    }else if(session.free != null && session.free == 0){
        // session is full
        $("#detail_session_create").hide();
        $("#detail_session").show();
        $("#menu_session").show();
        $("#detail_rider_book").hide();
        $("#menu_booking").show();
    }else{
        // show the create session button
        $("#detail_session_create").hide();
        $("#detail_rider_book").show();
        $("#detail_session").show();
        $("#menu_session").show();
        $("#menu_booking").show();
    }
    
    // update badges for rider and rider information
    var riders = "";
    if(session.riders != null){
        $("#detail_rider_count").text(session.riders.length);
        var max_space = parseInt(session.riders.length) + parseInt(session.free);
        $("#detail_rider_max_space").text(max_space);

        for(var i=0; i<session.riders.length; i++){
        var id = "delete_" + i;
        riders += "<div class=\"col-sm-10 col-xs-10\">"
                    + session.riders[i].name
                    + "</div>"
                    + "<div class=\"col-sm-2 col-xs-2\">"
                    + "<button class=\"btn btn-default btn-xs text-center\" onclick=\"deleteRiderSession(" + session.riders[i].id + ", " + session.id + ")\">"
                    + "<span class=\"glyphicon glyphicon-remove\"></span>"
                    + "</button>"
                    + "</div>"
                    + "<div class='clear'></div>";
        }
        $("#detail_rider_list").html(riders);
    }	
}

// Just books the actual user to the selected session
function db_simpleBook(){
    // send a booking request to the server
    var data = new Object();
    data['session_id'] = pieSessions[selectedSessionIdx].id;
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=add_self",
        data: JSON.stringify(data),
        success: function(resp){
            updateBookings(startDay, endDay, pieSessions[selectedSessionIdx].id);
        },
        error: function(resp){
            alert("Ooops. There was an error.");
        }
    });
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
            updateBookings(startDay, endDay, session_id);
        },
        error: function(resp){
            alert("Ooops. There was an error.");
        }
    });
}

// Shows the create session dialog
function showCreateSession(){
    let presets = {
        title:			null,
        description: 	null,
        maxRiders:      10,
        type:           0, 		// 0: private
    };
    // if a session is pre-selected -> use its start/end time
    if(selectedSessionIdx != null){
        session = pieSessions[selectedSessionIdx];
        presets.start = session.start.tz(getTimeZone());
        presets.end   = session.end.tz(getTimeZone());
    }
    else{
        presets.start = startDay;
        presets.end   = endDay;
    }

    BooksysViewSessionEditor.display("user_dialog", null, presets, function(newSessionId){
        BooksysViewSessionEditor.hideApp();
        BooksysViewSessionEditor.destroyView("user_dialog");
        updateBookings(startDay, endDay, newSessionId);
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

// Shows the add rider dialog where users can be added to the session
function showAddRider(){
    var sessionId = pieSessions[selectedSessionIdx].id;
    BooksysViewSessionRiderSelection.display("user_dialog", sessionId, function(){
        BooksysViewSessionRiderSelection.destroyView("user_dialog");
        updateBookings(startDay, endDay, sessionId);
    });
}

function showEditSession(){
    alert("TODO");	
}

// Given a session_id it returns the position in pieSessions of that session or null if not found
function getPieIdBySessionId(session_id){
    // get the position in pieSession that is affected
    for(var i = 0; i<pieSessions.length; i++){
        if(pieSessions[i].id != null && pieSessions[i].id == session_id){
            return i;
        }
    }
    return null;	
}