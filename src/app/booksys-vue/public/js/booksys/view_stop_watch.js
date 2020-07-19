class BooksysViewStopWatch {

    constructor(){
        this.users = new Array();
        this.heats = new Map();
        this.watch = {
            display:     "00:00",
            running:     false,
            isPaused:    false,
            user:        null,
            heatQueue:   0,
            comment:     null,
        };

        // load from local storage in case we 
        // are currently stopping a user
        if(typeof(Storage) !== "undefined"){
            // resume existing state if any
            if(localStorage.getItem("watch.running")){
                this.watch.running    = true;
                this.watch.startTime  = moment(localStorage.getItem("watch.startTime"), "X");
                this.watch.timeOffset = Number(localStorage.getItem("watch.timeOffset"));
                this.watch.user       = JSON.parse(localStorage.getItem('watch.user'));
                this.watch.sessionId  = localStorage.getItem("watch.sessionId");
                this.watch.comment    = localStorage.getItem("watch.comment");
                if(this.watch.sessionId == "null"){
                    this.watch.sessionId = null;
                }
                if(localStorage.getItem("watch.isPaused")){
                    this.watch.isPaused  = true;
                    this.watch.pauseTime = moment(localStorage.getItem("watch.pauseTime"), "X");
                }else{
                    this.watch.isPaused  = false;
                    this.timeout = setTimeout(this.updateDisplay.bind(this), 1000);
                }
                this.watch.display = this.getDisplayContent();
            }
            // resume sending of heats that have not been sent to the backend yet
            if(localStorage.getItem("watch.heats")){
                this.heats = new Map(Object.entries(JSON.parse(localStorage.getItem("watch.heats"))));
                this.heatsTimeout = null;
                if(this.heats.size > 0){
                    this.heatQueue = this.heats.size;
                    this.addHeats();
                }
            }
        }else{
            alert("Watch will not work properly with this browser's setting. Local Storage is not available.");
        }
    }

    // call this function to display the session details
    // - location       ID of the element where to place it
    // - sessionId      the ID of the session to display
    // - cb             callback functions
    display(location, sessionId, cb){
        // load session/user data and display
        this.location  = location;
        this.watch.sessionId = sessionId;
        localStorage.setItem("watch.sessionId", sessionId);

        let that = this;
        this.getData(sessionId, function(users){
            that.loadView(location, users, cb);
        });
    }

    isActive(){
        return this.watch.running;
    }

    isWatchPaused(){
        if(this.watch.isPaused == null || this.watch.isPaused == false){
            return false;
        }
        return true;
    }

    getUser(){
        return this.watch.user;
    }

    getSessionId(){
        return this.watch.sessionId;
    }

    // loads the data from the server
    getData(sessionId, cb){
        var details = new Object();

        if(sessionId == null){
            console.log("Get private users");
            // get the private riders
            $.ajax({
                type:   "GET",
                url:    "api/user.php?action=get_admin_users",
                cache:  false,
                success: function(resp){
                    let users = $.parseJSON(resp);
                    let privateUsers = Object.values(users);
                    if(cb != null){
                        cb(privateUsers);
                    }
                },
                error: function(xhr, status, error){
                    console.log(xhr);
                    console.log(status);
                    console.log(error);
                    var errorMsg = $.parseJSON(xhr.responseText);
                    alert(errorMsg.error);
                }
            });
        }else{
            // get the riders of this session
            // get the private riders
            let query = {
                "session_id": sessionId
            };
            $.ajax({
                type:   "POST",
                url:    "api/user.php?action=get_session_users",
                cache:  false,
                data:   JSON.stringify(query),
                success: function(resp){
                    let users = $.parseJSON(resp);
                    let sessionUsers = Object.values(users);
                    if(cb != null){
                        cb(sessionUsers);
                    }
                },
                error: function(xhr, status, error){
                    console.log(xhr);
                    console.log(status);
                    console.log(error);
                    var errorMsg = $.parseJSON(xhr.responseText);
                    alert(errorMsg.error);
                }
            });
        }
    }

    // gets the timezone
    getTimeZone(){
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

    // start the timer
    start(){
        // continue from where we started
        if(this.watch.running == false){
            // check that a user is selected
            if(this.watch.user.id == -1){
                alert("Please select a user");
            }
            document.getElementById("watchRiderSelect").disabled = true;
            this.watch.startTime = moment();
            this.watch.timeOffset = 0;
            
            this.timeout = setTimeout(this.updateDisplay.bind(this), 1000);
            this.watch.running = true;

            localStorage.setItem('watch.startTime', this.watch.startTime.format("X"));
            localStorage.setItem('watch.timeOffset', 0);
            localStorage.setItem('watch.running', true);
            localStorage.setItem('watch.user', JSON.stringify(this.watch.user));
        }else{
            // TODO
            alert("TODO");
        }
    }

    // pause taking time
    pause(){
        clearTimeout(this.timeout);
        this.watch.pauseTime = moment();
        this.watch.isPaused  = true;
        localStorage.setItem('watch.pauseTime', this.watch.pauseTime.format("X"));
        localStorage.setItem('watch.isPaused', this.watch.isPaused);
    }

    // resume taking time
    resume(){
        this.watch.timeOffset += Number(moment().format("X")) - Number(this.watch.pauseTime.format("X"));
        this.watch.isPaused    = false;
        localStorage.setItem('watch.timeOffset', this.watch.timeOffset);
        localStorage.removeItem('watch.pauseTime');
        localStorage.removeItem('watch.isPaused');
        this.timeout = setTimeout(this.updateDisplay.bind(this), 1000);    
    }

    // finish taking time and add the current heat to the database
    // cb.success -> called if heats have been added
    // cb.error   -> called if not all heats have been added
    finish(cb){
        // add heat to database
        let heat = {
            user_id:    this.watch.user.id,
            session_id: this.watch.sessionId,
            duration_s: this.getDuration(),
            comment:    this.watch.comment,
        };
        this.heats.set(moment().tz(getTimeZone()).format("X"), heat);
        localStorage.setItem('watch.heats', JSON.stringify(Object.fromEntries(this.heats)));
        this.addHeats(cb);
        this.cancel();
    }

    // cancel taking time and do not add heat to database
    cancel(){
        clearTimeout(this.timeout);
        this.watch.running = false;
        this.watch.isPaused = false;
        this.watch.startTime = null;
        this.watch.display = "00:00";
        
        this.watch.timeOffset = 0;
        this.watch.pauseTime  = null;
        this.selected = this.watch.user.id;
        
        localStorage.removeItem("watch.startTime");
        localStorage.removeItem("watch.pauseTime");
        localStorage.removeItem("watch.isPaused");
        localStorage.removeItem("watch.timeOffset");
        localStorage.removeItem("watch.running");
        localStorage.removeItem("watch.user");
    }  

    // loads the vue template and creates a vue instance
    // - location        where to display the view
    // - users           the users of this session
    // - cb              used in case a new session is created / or an existing session got altered
    loadView(location, users, cb){
        let that = this;
        that.users        = users;
        //that.watch.user   = that.users[0];
        that.selected     = -1;
        $('#'+location).load("res/view_stop_watch.vue", function(){
            var App = new Vue({
                el: "#view_stop_watch",
                data: {
                    "selected": that.selected,
                    "users":    that.users,
                    "watch":    that.watch,
                    "isPaused": that.isPaused,
                },
                created () {
                    // display the watch user as selected
                    if(this.watch.user != null){
                        this.selected = this.watch.user.id;
                    }
                },
                methods: {
                    "onChange": function(event){
                        console.log("onChange called");
                        // get user with the given ID (the list is small, thus traverse array)
                        for(let user of users){
                            if(user.id == event.target.value){
                                this.watch.user = user;
                            }
                        }
                        // reset comment upon user-change
                        this.watch.comment = null;
                        localStorage.removeItem("watch.comment");
                    },
                    "start": function(){
                        that.start();
                    },
                    "pause": function(){
                        that.pause();
                    },
                    "resume": function(){
                        that.resume();
                    },
                    "finish": function(){
                        that.finish(cb);
                    },
                    "addComment": function(){
                        console.log("addComment Clicked");
                        $('#addUserModal').modal('show');
                    },
                    "saveComment": function(){
                        if(this.watch.comment != null && this.watch.comment != ''){
                            localStorage.setItem("watch.comment", this.watch.comment);
                        }else{
                            localStorage.removeItem("watch.comment");
                        }
                        $('#addUserModal').modal('hide');
                    },
                    "removeComment": function(){
                        this.watch.comment = null;
                        localStorage.removeItem("watch.comment");
                    },
                    "back": function(){
                        cb.back();
                    }
                }
            });

            // add modal event listener
            $('#addUserModal').modal({
                backdrop: 'static',
                keyboard: false,
                show: false,
            });
            // focus input form upon opening modal
            $('#addUserModal').on('shown.bs.modal', function () {
                $('#watchCommentInput').focus()
            });
            // prevent enter from sending a GET request
            $('#watch_comment_form').on('keyup keypress', function(e) {
                var keyCode = e.keyCode || e.which;
                if (keyCode === 13) { 
                    e.preventDefault();
                    App.saveComment();
                    return false;
                }
            });
        });
    }

    destroyView(){
        let location = this.location;
        let child = document.getElementById(location).lastElementChild;
        while (child){
            document.getElementById(location).removeChild(child)
            child = document.getElementById(location).lastElementChild;
        }

        // tidy up
        if(localStorage.getItem("watch.heats") == '{}'){
            localStorage.removeItem("watch.heats");
            localStorage.removeItem("watch.sessionId");
        }
        localStorage.removeItem("watch.comment");
        this.watch.comment = null;
        this.watch.user    = null;
    }

    updateDisplay(){
        this.watch.display = this.getDisplayContent();
        this.timeout = setTimeout(this.updateDisplay.bind(this), 1000);
    }

    getDisplayContent(){
        let timediff = this.getDuration();
        let minutes_passed = Math.floor(timediff / 60);
        if(minutes_passed < 10){
            minutes_passed = "0" + minutes_passed;
        }
        let seconds_passed = timediff % 60;
        if(seconds_passed < 10){
            seconds_passed = "0" + seconds_passed;
        }
        let displayContent = minutes_passed+":"+seconds_passed;
        return displayContent;
    }

    getDuration(){
        if(this.watch.isPaused == true){
            this.watch.currentTime = this.watch.pauseTime;
        }else{
            this.watch.currentTime = moment();
        }
        let timediff = this.watch.currentTime.format("X") - this.watch.startTime.format("X") - this.watch.timeOffset;
        return timediff;
    }

    addHeats(cb){
        let that = this;
        if(this.heatsTimeout != null){
            clearTimeout(this.heatsTimeout);
        }

        $.ajax({
            type: "POST",
            url: "api/heat.php?action=add_heats",
            timeout: 10000,                   // 10 seconds timeout
            data: JSON.stringify(Object.fromEntries(this.heats)),
            success: function(data){
                let response = JSON.parse(data);
                for(let key in response){
                    if(response[key].ok == true){
                        if(cb != null && cb.success != null){
                            cb.success(that.heats[key]);
                        }
                        that.heats.delete(key);
                        localStorage.setItem("watch.heats", JSON.stringify(Object.fromEntries(that.heats)));
                    }else{
                        if(cb != null && cb.error != null){
                            cb.error(response[key].msg, that.heats[key]);
                        }
                    }
                };
                that.watch.heatQueue = that.heats.size;
            },
            error: function(data, err_txt, err_thrown){
                console.log("Error while submitting heats");
                console.log(data);
                console.log(err_txt);
                console.log(err_thrown);
                if(cb != null && cb.error != null){
                    cb.error("Cannot add heat to db, will retry automatically.");

                }
                that.watch.heatQueue = that.heats.size;
                that.heatsTimeout = setTimeout(that.addHeats.bind(that), 5000);
            }
        });
    }
}