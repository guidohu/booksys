class BooksysViewSessionEditor {

    // call this function to display the session details
    // - location       ID of the element where to place it
    // - id             the ID of the session to display
    // - sessionPresets presets to display (e.g. in case of a new session)
    // - cb             callback function
    static display(location, id, sessionPresets, cb){
        // load session and display
        BooksysViewSessionEditor.getData(id, function(session){
            BooksysViewSessionEditor.loadView(location, session, sessionPresets, cb);
        });
    }

    // loads the data from the server in case an id is specified
    static getData(id, cb){
        if(id == null){
            // the selected slot is available
            // use presets
            cb(null);
            return;
        }

        $.ajax({
            type: "GET",
            url: "api/booking.php?action=get_session&id=" + id,
            cache: false,
            contentType: "application/json; charset=utf-8",
            success: function(data){
                var details = $.parseJSON(data);
                var timezone = getTimeZone();
                details.start_time = moment(details.start_time).tz(timezone);
                details.end_time   = moment(details.end_time).tz(timezone);
                cb(details);
            },
            error: function(request, status, error){
                console.log(request);
                console.log(status);
                console.log(error);
                alert("cannot get session information from server");
            }
        });

        return;
    }

    // get timezone settings
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

    // displays the App
    static showApp(){
        $('#view_session_editor_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_session_editor_modal').modal('show');
    }

    // hides the App
    static hideApp(){
        $('#view_session_editor_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // destroys what have been loaded
    static destroyView(location){
        var child = document.getElementById(location).lastElementChild;
        while (child){
            document.getElementById(location).removeChild(child)
            child = document.getElementById(location).lastElementChild;
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
        if(session == null){
            // this is a new session, use presets if any
            session             = new Object();
            session.id          = null;
            session.type        = 0;
            if(sessionPresets != null){
                session.title       = sessionPresets.title;
                session.description = sessionPresets.description;
                session.start       = sessionPresets.start;
                session.end         = sessionPresets.end;
                session.maxRiders   = sessionPresets.maxRiders;
                session.type        = sessionPresets.type;
            }
        }

        if(session.type == 1){
            session.isCourse = true;
        } else {
            session.isCourse = false;
        }

        $('#'+location).load("res/view_session_editor.vue", function(){
            // load the vue App
            var App = new Vue({
                el: "#view_session_editor_modal",
                data: {
                    "session": session,
                    "errors" : [],
                },
                methods: {
                    save: function(event){
                        let callbacks = {
                            success: cb,
                            failure: function(errs){
                                App.errors = errs
                            }
                        };
                        BooksysViewSessionEditor.save(this.session, callbacks);
                    },
                    cancel: function(){
                        BooksysViewSessionEditor.hideApp();
                    }
                }
            });

            // Initialize datepickers
            $('#startPicker').datetimepicker({
                format: 'D. MMMM YYYY HH:mm',
                //format: 'HH:mm',
                sideBySide: false,
                //showClose: true,
                date: session.start,
                timeZone: BooksysViewSessionEditor.getTimeZone(),
                defaultDate: moment().tz(getTimeZone()),
            });
            $('#startPicker').datetimepicker().on('dp.change', function (event) {
                session.start = event.date;
            });
            $('#endPicker').datetimepicker({
                format: 'D. MMMM YYYY HH:mm',
                //format: 'HH:mm',
                sideBySide: false,
                //showClose: true,
                date: session.end,
                timeZone: BooksysViewSessionEditor.getTimeZone(),
                defaultDate: moment().add(2, 'hours').tz(getTimeZone()),
            });
            $('#endPicker').datetimepicker().on('dp.change', function (event) {
                session.end = event.date;
            });


            // Initialize Bootstrap Switch
			$("#isCourse").bootstrapSwitch({
				onText: "Course",
				onColor: "info",
				offText: "Private",
				animate: true,
				size: 'normal',
				onSwitchChange: function(event, state){ 
                    session.isCourse = state;
                    if(session.isCourse == true){
                        session.type = 1;
                    } else {
                        session.type = 0;
                    }
                }
			});

            BooksysViewSessionEditor.showApp();
        });
    }

    // save
    static save(session, cb){
        // in case we do not have an ID, create a new session
        if(session.id == null){
            var req = new Object();
            req.title      = session.title;
            req.comment    = session.description;
            req.max_riders = session.maxRiders;
            req.start      = session.start.format('X');
            req.end        = session.end.format('X');
            req.type       = session.type;
            $.ajax({
                type: "POST",
                url:  "api/booking.php?action=add_session",
                data: JSON.stringify(req),
                success: function(response){
                    var json = $.parseJSON(response);
                    cb.success(json.session_id);
                },
                error: function(xhr,status,error){
                    console.log(xhr);
                    console.log(status);
                    console.log(error);
                    var errorMsg = $.parseJSON(xhr.responseText);
                    cb.failure([errorMsg.error]);
                },
            });
        }
        // TODO add case to also handle updates of a given
        // session
    }
}