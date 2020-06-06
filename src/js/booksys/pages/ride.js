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
        // Add mobile style also in desktop mode
        BooksysBrowser.addMobileCSS();
        $("#body").load("res/ride_mobile.html", function(){
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

// Variable that gives access to the heatlist instance
let heatlist;

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
    $("#watch").hide();
    $("#btn_continue").hide();
    $("#status-panel").hide();

    // setup our beacon
    setupBeacon();

    // create a new stop watch
    watch = new BooksysViewStopWatch();
    heatlist = new BooksysViewHeatList(
        {
            heatDetailsCb: function(id){
                BooksysViewHeatEntry.displayHeatEntry("heatentry", id, function(){
                    BooksysViewHeatEntry.destroyView("heatentry");
                    refreshHeatlist();
                });
            },
        }
    );

    // check for existing session and switch to watch
    // if required
    if(watch.isActive()){
        let sessionId = watch.getSessionId();
        showWatch(sessionId);
    }else{
        updateBookings(date_start, date_end);
    }
}

// Requests the new data from the server and
// updates the view
function updateBookings(date_start, date_end){
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

            // Display daily info
            var date = moment(json.window_start, "X");
            $("#watch_date").html(date.tz(getTimeZone()).format("dddd DD.MM.YYYY"));
            $("#watch_time").html("-");
            // Calculate sunset time and display it
            var sunsetTime = moment(json.sunset, "X");
            $("#watch_sunset").html(sunsetTime.tz(getTimeZone()).format("HH:mm"));
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
    
    // Update stopwatch info
    var start = selectedSession.start;
    var end   = selectedSession.end;
    $("#watch_time").html(start.tz(getTimeZone()).format('HH:MM') + " - " + end.tz(getTimeZone()).format('HH:MM'));
    if(selectedSession.title){
        $("#watch_title").html(selectedSession.title);
    }else{
        $("#watch_title").html('-');
    }

    // In case this is a free slot, we like to provide some
    // presets for the creation of a new session
    var presets = null;
    var timezone = getTimeZone();
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

// send regular beacons
function setupBeacon(){
    // setup beacons
    let beacon = new BooksysBeacon(
        {
            timeout:    10000,
        },
        {
            success:    function(){
                var content = '<span class="glyhpicon glyphicon-ok"></span> connected';
                $('#status-panel').fadeOut(5000);
                $('#status-msg').html(content);
            },
            redirect:   function(url){
                window.location.href = url;
                return;
            },
            failure:    function(){
                var content = '<span class="glyphicon glyphicon-exclamation-sign"></span> bad connection';
                $('#status-msg').html(content);                    
                $('#status-panel').fadeIn(2000);
            }
        }
    );
}

function showSessionWatch(){
    showWatch(selectedSession.id);
}

function showWatch(sessionId){
    $("#select_session").hide();

    // define callback actions for the watch
    let callbacks = {
        back: function(){
            watch.destroyView();
            showSessionSelect();
        },
        success: function(){
            heatlist.display("heatlist", watch.getSessionId());
        }
    };

    // setup the watch and heatlist
    watch.display("stopwatch", sessionId, callbacks);
    heatlist.display("heatlist", sessionId);
    
    // display the watch
    $("#watch").show();       
}

// refresh the heat list
function refreshHeatlist(){
    heatlist.display("heatlist", watch.getSessionId());
}
    
// Shows the view to select a session
function showSessionSelect(){
    $("#watch").hide();
    $("#select_session").show();
    watch.destroyView();
    
    updateBookings(date_start, date_end);
}