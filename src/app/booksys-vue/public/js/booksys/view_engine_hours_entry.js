class BooksysViewEngineHoursEntry {

    // call this function to display the boat log entry dialog
    // - location              ID of the element where to place the dialog
    // - engineHoursEntryId    the ID of the boat log entry that should be displayed
    // - cb                    callback function upon modal closing (save/cancel)
    static displayEntry(location, engineHoursEntryId, cb){
        BooksysViewEngineHoursEntry.getEntry(engineHoursEntryId, function(entry){
            BooksysViewEngineHoursEntry.loadView(location, entry, cb);
        });
    }

    // get the boat log entry with the given id from the server and calls the
    // callback function cb
    static getEntry(id, cb){
        var params = new Object();
        params['id'] = id;

        $.ajax({
            type: "POST",
            url: "api/boat.php?action=get_engine_hours_entry",
            cache: false,
            data: JSON.stringify(params),
            success: function(data){
                var eh = $.parseJSON(data);
                console.log(eh);
                cb(eh);
            },
        });
    }

    static showApp(){
        $('#view_engine_hours_entry_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_engine_hours_entry_modal').modal('show');
    }

    static hideApp(){
        $('#view_engine_hours_entry_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays a specific entry at a given location
    // - location     DOM location
    // - entry        a fuel entry as retrieved by getEngineHoursEntry()
    // - cb           callback function that gets called upon changed fuel entry
    static loadView(location, entry, cb){
        console.log(entry);
        // simplify some of the information
        // full name
        entry.full_name = entry.user_first_name + " " + entry.user_last_name;
        
        // date to date string
        var date = new Date(entry.time*1000);
        entry.date_text = dateToString(date);

        // simplify course/default handling
        entry.is_course = false;
        if(entry.type_name == "course"){
            entry.is_course = true;
        }
        
        // setup errors array
        var errors = new Array();

        $('#'+location).load("res/view_engine_hours_entry.vue", function(){    
            var engineHoursEntryEditApp = new Vue({
                el: '#view_engine_hours_entry_modal',
                data: {
                    "entry": entry,
                    "errors": errors,
                },
                methods: {
                    save: function(event){
                        // validate the fuel entry
                        this.errors = BooksysViewEngineHoursEntry.validateEntry(this.entry);
                        if(this.errors.length > 0){
                            return;
                        }

                        // reconstruct the correct values
                        // if(this.entry.is_fuel_discount === false){
                        //     this.entry.cost_chf_brutto = null;
                        // }

                        // save fuel entry
                        BooksysViewEngineHoursEntry.saveEntry(this.entry, cb);
                    },
                    cancel: function(event){
                        BooksysViewEngineHoursEntry.hideApp();
                        cb();
                    },
                },
            });

            // Initialize Bootstrap Switch
			$("[name='isCourse']").bootstrapSwitch({
				onText: "Course",
				onColor: "info",
				offText: "Private",
				animate: true,
				size: 'normal',
				onSwitchChange: function(event, state){ 
                    entry.is_course = state;
                    if(entry.is_course == true){
                        entry.type = 1;
                        entry.type_name = "course";
                    }else{
                        entry.type = 0;
                        entry.type_name = "default";
                    }
                }
			});

            BooksysViewEngineHoursEntry.showApp();
        });
    }

    // removes the view from the DOM
    static destroyView(location){
        var child = document.getElementById(location).lastElementChild;
        while (child){
            document.getElementById(location).removeChild(child)
            child = document.getElementById(location).lastElementChild;
        }
    }

    // validate whether input is valid
    static validateEntry(entry){
        var errors = [];

        // no validation required

        return errors;
    }

    // save the updated fuel entry
    static saveEntry(entry, callbackFn){
        $.ajax({
            type: "POST",
            url: "api/boat.php?action=update_engine_hours_entry",
            cache: false,
            data: JSON.stringify(entry),
            success: function(data){
                BooksysViewEngineHoursEntry.hideApp(); 
                callbackFn();
            },
        });
    }
}