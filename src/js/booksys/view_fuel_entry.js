class BooksysViewFuelEntry {

    // call this function to display the fuel entry dialog
    // - location    ID of the element where to place the dialog
    // - fuelEntryId the ID of the fuel entry that should be displayed
    // - cb          callback function upon modal closing (save/cancel)
    static displayFuelEntry(location, fuelEntryId, cb){
        BooksysViewFuelEntry.getFuelEntry(fuelEntryId, function(entry){
            BooksysViewFuelEntry.loadView(location, entry, cb);
        });
    }

    // get the fuel entry with the given id from the server and calls the
    // callback function cb
    static getFuelEntry(id, cb){
        var params = new Object();
        params['id'] = id;

        $.ajax({
            type: "POST",
            url: "api/boat.php?action=get_fuel_entry",
            cache: false,
            data: JSON.stringify(params),
            success: function(data){
                var entry = $.parseJSON(data);
                cb(entry);
            },
        });
    }

    static showApp(){
        $('#view_fuel_entry_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_fuel_entry_modal').modal('show');
    }

    static hideApp(){
        $('#view_fuel_entry_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays a specific entry at a given location
    // - location     DOM location
    // - data         a fuel entry as retrieved by getFuelEntry()
    // - cb           callback function that gets called upon changed fuel entry
    static loadView(location, data, cb){
        let entry    = data.fuel_entry;
        let currency = data.currency;

        // simplify some of the information
        // full name
        entry.full_name = entry.user_first_name + " " + entry.user_last_name;
        
        // date to date string
        var date = new Date(entry.timestamp*1000);
        entry.date_text = dateToString(date);
        
        // setup errors array
        var errors = new Array();
        
        // handle discounted/non-discounted entries
        if(typeof entry.cost_brutto === 'undefined' || entry.cost_brutto == null){
            entry.is_fuel_discount = false;
        } else {
            entry.is_fuel_discount = true;
        }

        // TODO append it to the body / remove it upon cancel
        //var targetLocation = document.getElementById(location);

        $('#'+location).load("res/view_fuel_entry.vue", function(){    
            var fuelEntryEditApp = new Vue({
                el: '#view_fuel_entry_modal',
                data: {
                    "entry":    entry,
                    "currency": currency,
                    "errors":   errors,
                },
                methods: {
                    save: function(event){
                        // validate the fuel entry
                        this.errors = BooksysViewFuelEntry.validateFuelEntry(this.entry);
                        if(this.errors.length > 0){
                            return;
                        }

                        // reconstruct the correct values
                        if(this.entry.is_fuel_discount === false){
                            this.entry.cost_brutto = null;
                        }

                        // save fuel entry
                        BooksysViewFuelEntry.saveFuelEntry(this.entry, cb);
                    },
                    cancel: function(event){
                        BooksysViewFuelEntry.hideApp();
                        cb();
                    },
                },
            });

            // Initialize Bootstrap Switch
			$("[name='isFuelDiscount']").bootstrapSwitch({
				onText: "Yes",
				onColor: "info",
				offText: "No",
				animate: true,
				size: 'normal',
				onSwitchChange: function(event, state){ 
                    entry.is_fuel_discount = state;

                    if(entry.is_fuel_discount === true && entry.cost_brutto == null){
                        // discount gets activated, we might need to move cost NET to cost GROS
                        entry.cost_brutto = entry.cost;
                        entry.cost = null;
                    } else if(entry.is_fuel_discount === false){
                        // discount gets deactivated, we might need to move cost GROS to cost NET
                        if(entry.cost == null && entry.cost_brutto != null) {
                            entry.cost = entry.cost_brutto;
                            entry.cost_brutto = null;
                        }
                    }
                }
			});

            BooksysViewFuelEntry.showApp();
        });
    }

    static destroyView(location){
        //$('#'+location).removeChild();
        var child = document.getElementById(location).lastElementChild;
        while (child){
            document.getElementById(location).removeChild(child)
            child = document.getElementById(location).lastElementChild;
        }
    }

    // validate whether input is valid
    static validateFuelEntry(entry){
        var errors = [];

        if(entry.is_fuel_discount === true && entry.cost_brutto == null){
            errors.push("Cost (Gros) needs to be specified if discount is delected");
        }
        if(entry.is_fuel_discount){
            if(! entry.cost_brutto || ! entry.cost_brutto.toString().match(/^\d+(\.\d{1,2})?$/)){
                errors.push("Cost (Gros) needs to be a number (e.g. 19.23)");
            }
        }
        if(! entry.cost || ! entry.cost.toString().match(/^\d+(\.\d{1,2})?$/)){
            errors.push("Cost (Net) needs to be a number (e.g. 19.23)");
        }
        if(! entry.liters || ! entry.liters.toString().match(/^\d+(\.\d{1,2})?$/)){
            errors.push("Fuel needs to be a number (e.g. 19.23)");
        }
        if(! entry.engine_hours || ! entry.engine_hours.toString().match(/^\d+(\.\d{1,2})?$/)){
            errors.push("Engine needs to be a number (e.g. 1912.2)");
        }

        return errors;
    }

    // save the updated fuel entry
    static saveFuelEntry(entry, callbackFn){
        $.ajax({
            type: "POST",
            url: "api/boat.php?action=update_fuel_entry",
            cache: false,
            data: JSON.stringify(entry),
            success: function(data){
                BooksysViewFuelEntry.hideApp(); 
                callbackFn();
            },
        });
    }
}