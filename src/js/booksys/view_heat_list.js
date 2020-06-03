class BooksysViewHeatList {

    constructor(parameters){
        this.heatDetailsCb = parameters.heatDetailsCb
    }

    // call this function to display the session details
    // - location       ID of the element where to place it
    // - sessionId      the ID of the session to display
    display(location, sessionId){
        // load session/user data and display
        this.location  = location;
        this.sessionId = sessionId;

        let that = this;
        this.getData(sessionId, function(heats){
            that.loadView(location, heats);
        });
    }

    // loads the data from the server
    getData(sessionId, cb){

        if(sessionId == null){
            let query = {
                from:   moment().tz(getTimeZone()).startOf('day').format('X'),
                to:     moment().tz(getTimeZone()).endOf('day').format('X'),
            };
            // get the last few heats of the admins
            $.ajax({
                type:   "POST",
                url:    "api/heat.php?action=get_heats",
                data:   JSON.stringify(query),
                cache:  false,
                success: function(resp){
                    let heats = $.parseJSON(resp);
                    if(cb != null){
                        cb(heats);
                    }
                    // TODO error handling
                },
                error: function(xhr, status, error){
                    console.log(xhr);
                    console.log(status);
                    console.log(error);
                    alert(error);
                }
            });
        }else{
            // get the heats of the given session
            let query = {
                "session_id": sessionId
            };
            $.ajax({
                type:   "POST",
                url:    "api/heat.php?action=get_session_heats",
                cache:  false,
                data:   JSON.stringify(query),
                success: function(resp){
                    let heats = $.parseJSON(resp);
                    if(cb != null){
                        cb(heats);
                    }
                    // TODO error handling
                },
                error: function(xhr, status, error){
                    console.log(xhr);
                    console.log(status);
                    console.log(error);
                    alert(error);
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

    // loads the vue template and creates a vue instance
    // - location        where to display the view
    // - data            the heats and currency data
    loadView(location, data){
        let that         = this;
        that.heats       = data.heats;
        that.currency    = data.currency;
        that.hasComments = false;

        // process the heats for display
        for (let idx in that.heats){
            // pretty format the time
            let minutes = Math.floor(that.heats[idx].duration_s / 60);
            if(minutes < 10){
                minutes = "0" + minutes;
            }
            let seconds = that.heats[idx].duration_s % 60;
            if(seconds < 10){
                seconds = "0" + seconds;
            }
            that.heats[idx].duration_string = minutes + ":" + seconds;

            // check whether we have a comment in one of the heats
            if(that.heats[idx].comment != null && that.heats[idx].comment != ''){
                that.hasComments = true;
            }
        }

        $('#'+location).load("res/view_heat_list.vue", function(){
            var App = new Vue({
                el: "#view_heat_list",
                data: {
                    heats:       that.heats,
                    hasComments: that.hasComments,
                    currency:    that.currency,
                },
                methods: {
                    showDetails: function(heatIndex){
                        // callback to show heat details
                        that.heatDetailsCb(that.heats[heatIndex].heat_id);
                    }
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
    }
}