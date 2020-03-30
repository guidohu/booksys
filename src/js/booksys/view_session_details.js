class BooksysViewSessionDetails {

    // call this function to display the session details
    // - location       ID of the element where to place it
    // - id             the ID of the session to display
    // - sessionPresets presets to display (e.g. in case of a new session)
    // - cb             callback functions
    //                  cb.create --> called when a new session gets created
    static display(location, id, sessionPresets, cb){
        // load session and display
        BooksysViewSessionDetails.getData(id, sessionPresets, function(details){
            BooksysViewSessionDetails.loadView(location, details, sessionPresets, cb);
        });  
    }

    // loads the data from the server
    static getData(id, sessionPresets, cb){
        var details = new Object();

        if(id == null && sessionPresets != null){
            // the selected slot is available
            // use presets
            details.id = id;
            cb(details);
            return;
        } else if(id == null){
            // no session selected
            cb(null);
            return;
        }

        $.ajax({
            type: "GET",
            url: "api/booking.php?action=get_session&id=" + id,
            cache: false,
            contentType: "application/json; charset=utf-8",
            success: function(data){
                details = $.parseJSON(data);
                var timezone = getTimeZone();
                details.start_time = moment(details.start_time, "X").tz(timezone);
                details.end_time   = moment(details.end_time, "X").tz(timezone);
                cb(details);
            },
            error: function(request, status, error){
                console.log(request);
                console.log(status);
                console.log(error);
                alert("An error occurred and the session cannot be retrieved.");
            }
        });

        return;
    }

    // gets the timezone
    static getTimeZone(){
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

    // loads the vue template and creates a vue instance
    // - location        where to display the view
    // - session         the session object itself (needs a session.id of -1 as minimum)
    // - sessionPresets  used in case a session should be created (in case session.id is -1)
    //                   sessionPresets allows the following key-values:
    //                   title, 
    // - cb              used in case a new session is created / or an existing session got altered
    static loadView(location, session, sessionPresets, cb){
                
        if((typeof session == 'undefined' || session == null) && (typeof sessionPresets == 'undefined' || sessionPresets == null)){
            // no session selected
            // - nothing to do currently
        }
        else if(session.id == null){
            // this is not a session yet, but display presets might have been provided
            // that are uses within a "create" session window
            session.title       = "no session yet";
            session.date_string = sessionPresets.start.format('dddd') + ' ' + sessionPresets.start.format('DD.MM.YYYY');
            session.start       = sessionPresets.start;
            session.end         = sessionPresets.end;
            session.time_string = session.start.format('HH:mm') + " - " + session.end.format('HH:mm');
            session.riders      = new Array();
            session.isCreate    = true;
            session.exists      = false;
        }
        else {
            // create the corresponding output format
            session.date_string = session.start_time.format('dddd') + ' ' + session.start_time.format('L');
            session.start       = session.start_time;
            session.end         = session.end_time;
            session.time_string = session.start.format('HH:mm') + " - " + session.end.format('HH:mm');
            session.isCreate    = false;
            session.exists      = true;
        }

        $('#'+location).load("res/view_session_details.vue", function(){
            var App = new Vue({
                el: "#view_session_details",
                data: {
                    "session": session
                },
                methods: {
                    "createSession": function(){
                        cb.createSession()
                    },
                    "deleteSession": function(){
                        cb.deleteSession(session.id);
                    },
                    "addRider": function(){
                        cb.addRider(session.id);
                    },
                }
            });
        });
    }
}