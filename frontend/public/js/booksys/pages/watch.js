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

        $("#body").load("res/watch_mobile.html", function(){
            loadContent();
        });
    } else {
        $("#body").load("res/watch.html", function(){
            loadContent();
        });
    }
});

// Stay in the WebApp
$(function() {
    $.stayInWebApp();
});

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
    $('#status-panel').hide()

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

    // get the specific session ID from the URL
    let searchParams = new URLSearchParams(window.location.search);
    let uriSessionId = searchParams.get('sessionId');

    // switch to stop watch for the following reasons:
    // - the user is already taking the time currently
    // - the link is to a specific session
    // otherwise display the selection of sessions for today
    if(watch.isActive()){
        let sessionId = watch.getSessionId();
        updateWatchDetails(uriSessionId);
        showWatch(sessionId);
    }else if(uriSessionId != null){
        updateWatchDetails(uriSessionId);
        showWatch(uriSessionId);
    }else{
        alert("No session selected");
    }
}

function updateWatchDetails(sessionId){
    $.ajax({
        type: "GET",
        url: "api/booking.php?action=get_session&id=" + sessionId,
        cache: false,
        contentType: "application/json; charset=utf-8",
        success: function(data){
            let details = $.parseJSON(data);

            let date  = moment(details.start_time, "X");
            $("#watch_date").html(date.tz(getTimeZone()).format("dddd DD.MM.YYYY"));
            let start = moment(details.start_time, "X");
            let end   = moment(details.end_time, "X");
            $("#watch_time").html(start.tz(getTimeZone()).format('HH:mm') + " - " + end.tz(getTimeZone()).format('HH:mm'));
            let sunsetTime = moment(details.sunset, "X");
            $("#watch_sunset").html(sunsetTime.tz(getTimeZone()).format("HH:mm"));

            if(details.title != null){
                $("#watch_title").html(details.title);
            }else{
                $("#watch_title").html("-");
            }
        },
        error: function(request, status, error){
            console.log(request);
            console.log(status);
            console.log(error);
            alert("An error occurred and the session cannot be retrieved.");
        }
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

function showWatch(sessionId){
    // define callback actions for the watch
    let callbacks = {
        back: function(){
            watch.destroyView();
            window.location.href = '/ride.html?selectedSessionId='+sessionId;
        },
        success: function(){
            heatlist.display("heatlist", watch.getSessionId());
        }
    };

    // setup the watch and heatlist
    watch.display("stopwatch", sessionId, callbacks);
    heatlist.display("heatlist", sessionId);       
}

// refresh the heat list
function refreshHeatlist(){
    heatlist.display("heatlist", watch.getSessionId());
}