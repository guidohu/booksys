class BooksysViewDatabaseSetup {

    // call this function to display the setup database dialog
    // - location              ID of the element where to place the dialog
    // - cb                    callback function upon modal closing (success/cancel)
    static displayEntry(location, cb){
        BooksysViewDatabaseSetup.getEntry(function(config){
            BooksysViewDatabaseSetup.loadView(location, config, cb);
        });
    }

    // get the current database configuration (if any) and calls the
    // callback function cb
    static getEntry(cb){
        $.ajax({
            type: "GET",
            url: "api/configuration.php?action=get_db_config",
            cache: false,
            success: function(data){
                let dbConfig = $.parseJSON(data);
                console.log(dbConfig);
                cb(dbConfig);
            },
            error: function(data){
                console.log("database setup already done? - return empty db config");
                cb(null);
            }
        });
    }

    static showApp(){
        $('#view_database_setup_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_database_setup_modal').modal('show');
    }

    static hideApp(){
        $('#view_database_setup_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays the dialog at a given location
    // - location     DOM location (id)
    // - dbConfig     the current database configuration as retrieved from the backend
    // - cb           callback function that gets called upon changed fuel entry
    static loadView(location, dbConfig, cb){
        console.log(dbConfig);

        if(dbConfig == null){
            cb.success();
            return;
        }
        
        // setup errors array
        let errors = new Array();

        $('#'+location).load("res/view_database_setup.vue", function(){    
            var databaseSetupApp = new Vue({
                el: '#view_database_setup_modal',
                data: {
                    "dbConfig":     dbConfig,
                    "errors":       errors,
                },
                methods: {
                    save: function(event){
                        // validate the db configuration
                        this.errors = BooksysViewDatabaseSetup.validateEntry(this.dbConfig);
                        if(this.errors.length > 0){
                            return;
                        }

                        // setup db configuration and database
                        let saveCallback = {
                            success: cb.success,
                            failure: function(sErrors){
                                databaseSetupApp.errors = sErrors;
                            }
                        };
                        BooksysViewDatabaseSetup.saveEntry(this.dbConfig, saveCallback);
                    },
                    cancel: function(event){
                        BooksysViewDatabaseSetup.hideApp();
                        cb.cancel();
                    },
                },
            });

            BooksysViewDatabaseSetup.showApp();
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
    static validateEntry(config){
        let errors = [];

        // lazy validation (proper validation done on backend)
        if(! config.host || ! config.host.toString().match(/^[a-zA-Z0-9-\.\[\]_:]+$/)){
            errors.push("Please select a proper hostname or IP for the database host");
        }
        if(! config.name || ! config.name.toString().match(/^[a-zA-Z0-9-_]+$/)){
            errors.push("Please select a proper database name");
        }
        if(! config.user || ! config.user.toString().match(/^[a-zA-Z0-9-_]+$/)){
            errors.push("Please select a proper database user");
        }
        if(! config.password){
            errors.push("Please select a password for the database user");
        }

        return errors;
    }

    // save the database configuration
    static saveEntry(dbConfig, callbackFn){
        let config = {
            db_server:   dbConfig.host,
            db_name:     dbConfig.name,
            db_user:     dbConfig.user,
            db_password: dbConfig.password,
        };

        $.ajax({
            type: "POST",
            url: "api/configuration.php?action=setup_db_config",
            cache: false,
            data: JSON.stringify(config),
            success: function(resp){
                let data = $.parseJSON(resp)
                if(data.OK == true){
                    console.log("DB configuration applied");
                    BooksysViewDatabaseSetup.hideApp(); 
                    callbackFn.success();
                }else{
                    console.log("DB configuration not applied");
                    callbackFn.failure([ data.msg ]);
                }
                
            },
        });
    }
}