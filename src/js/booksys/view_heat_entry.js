class BooksysViewHeatEntry {

    // call this function to display a heat entry dialog
    // - location    ID of the element where to place the dialog
    // - heatId      the ID of the heat entry that should be displayed
    // - cb          callback function upon modal closing (save/cancel)
    static displayHeatEntry(location, heatId, cb){
        BooksysViewHeatEntry.getHeatEntry(heatId, function(entry){
            BooksysViewHeatEntry.loadView(location, entry, cb);
        });
    }

    // get the heat entry with the given id from the server and calls the
    // callback function cb
    static getHeatEntry(id, cb){
        var params = new Object();
        params['heat_id'] = id;

        $.ajax({
            type: "POST",
            url: "api/heat.php?action=get_heat",
            cache: false,
            data: JSON.stringify(params),
            success: function(data){
                var entry = $.parseJSON(data);
                cb(entry);
            },
        });
    }

    // displays the app
    static showApp(){
        $('#view_heat_entry_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_heat_entry_modal').modal('show');
    }

    // hides the app
    static hideApp(){
        $('#view_heat_entry_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays a specific entry at a given location
    // - location     DOM location
    // - data         a heat entry as retrieved by getHeatEntry()
    // - cb           callback function that gets called upon changed heat entry
    static loadView(location, data, cb){
        console.log(data);
        let heat     = data.heat;
        let currency = data.currency;

        // simplify some of the information
        // full name
        heat.full_name = heat.first_name + " " + heat.last_name;
        
        // date to date string
        let date = moment(heat.timestamp, "X").tz(BooksysViewHeatEntry.getTimeZone()).format("DD.MM.YYYY HH:mm");
        heat.date_text = date;

        // price to price string
        heat.price_per_min_text = sprintf("%.2f", heat.price_per_min);

        // duration to duration string
        var duration_s     = heat.duration_s % 60;
        var duration_min   = parseInt(heat.duration_s / 60);
        heat.duration_text = sprintf("%d:%02d", duration_min, duration_s);

        // cost to cost string
        heat.cost_text     = sprintf("%.2f", heat.cost);
        
        // setup errors array
        var errors = new Array();

        $('#'+location).load("res/view_heat_entry.vue", function(){    
            var heatEntryEditApp = new Vue({
                el: '#view_heat_entry_modal',
                data: {
                    "heat":     heat,
                    "currency": currency,
                    "errors":   errors,
                    "showSave": false,
                },
                methods: {
                    save: function(event){
                        // save heat entry
                        let callback = {
                            success: cb,
                            failure: function(errors){
                                this.errors = errors;
                            },
                        };
                        BooksysViewHeatEntry.saveHeatEntry(this.heat, callback);
                    },
                    cancel: function(event){
                        BooksysViewHeatEntry.hideApp();
                        cb();
                    },
                    durationChange: function(event){
                        let duration_text = this.heat.duration_text;

                        // check that it is a valid format
                        let re = /^\d+:\d{2}$/;
                        if(re.test(duration_text)){
                            let duration_arr     = duration_text.split(":");
                            let duration_s       = 60*parseInt(duration_arr[0]) + parseInt(duration_arr[1]);
                            this.heat.duration_s = duration_s;
                            
                            // round to .05
                            let cost             = this.heat.price_per_min * duration_s / 60;
                            cost                 = Math.ceil(cost / 0.05) * 0.05;
                            this.heat.cost       = cost
                            this.heat.cost_text  = sprintf("%.2f", cost);

                            // activate save button
                            this.showSave        = true;
                        }else{
                            this.showSave        = false;
                            console.log("waiting for duration format to be correct: " + duration_text);
                        }
                    },
                },
            });

            BooksysViewHeatEntry.showApp();
        });
    }

    static destroyView(location){
        var child = document.getElementById(location).lastElementChild;
        while (child){
            document.getElementById(location).removeChild(child)
            child = document.getElementById(location).lastElementChild;
        }
    }

    // save the updated fuel entry
    static saveHeatEntry(heat, callbackFn){
        let request = {
            heat_id:    heat.heat_id,
            user_id:    heat.user_id,
            duration_s: heat.duration_s
        };

        $.ajax({
            type: "POST",
            url: "api/heat.php?action=update_heat",
            cache: false,
            data: JSON.stringify(request),
            success: function(data){
                let response = JSON.parse(data);
                if(response.ok){
                    BooksysViewHeatEntry.hideApp(); 
                    callbackFn.success();
                }else{
                    callbackFn.failure([ response.msg ]);
                } 
            },
            error: function(request, status, error){
                console.error(request);
                console.error(status);
                console.error(error);
                callbackFn.failure([ error ]);
            }
        });
    }

    /* Creates a dateString as we use it everywhere */
    static dateToString(date){
        var string = '';
        // Date
        string += date.getFullYear();
        string += '.';
        if(date.getMonth()< 9){
            string+= '0';
        }
        string += date.getMonth()+1;
        string += '.';
        if(date.getDate() < 10){
            string += 0;
        }
        string += date.getDate();
        
        // Time
        string += ' ';
        if(date.getHours() < 10){
            string += 0;
        }
        string += date.getHours();
        string += ':';
        if(date.getMinutes() < 10){
            string += 0;
        }
        string += date.getMinutes();
        return string;		
    }

    // gets the user's timezone or returns the default
    static getTimeZone(){
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
}