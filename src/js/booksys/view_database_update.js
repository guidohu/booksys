class BooksysViewDatabaseUpdate {

    // call this function to display the database update dialog
    // - location    ID of the element where to place the dialog
    // - cb          callback function upon modal closing
    static display(location, cb){
        BooksysViewDatabaseUpdate.getData(function(info){
            BooksysViewDatabaseUpdate.loadView(location, info, cb);
        });
    }

    // get the fuel entry with the given id from the server and calls the
    // callback function cb
    static getData(cb){
        // get the version info from the backend
        $.ajax({
            type: "GET",
            url: "api/backend.php?action=get_version",
            cache: false,
            success: function(resp){
                var versionInfo = $.parseJSON(resp);
                cb(versionInfo);
            },
        });
    }

    // displays the app (internal function)
    static showApp(){
        $('#view_database_update_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_database_update_modal').modal('show');
    }

    // hide the app again
    static hideApp(){
        $('#view_database_update_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays a specific entry at a given location
    // - location     DOM location
    // - entry        a fuel entry as retrieved by getFuelEntry()
    // - cb           callback function that gets called upon changed fuel entry
    static loadView(location, info, cb){
        info.isUpdated = true;
        if(info.app_version != info.db_version){
            info.isUpdated = false;
        }

        info.ok = true;  // currently we are in a healthy state
        info.isUpdating = false; // currently not updating yet

        $('#'+location).load("res/view_database_update.vue", function(){    
            var databaseUpateApp = new Vue({
                el: '#view_database_update_modal',
                data: {
                    "info": info,
                    "entry": {},
                    "errors": [],
                },
                methods: {
                    update: function(event){
                        BooksysViewDatabaseUpdate.update(info);
                    },
                    cancel: function(event){
                        BooksysViewDatabaseUpdate.hideApp();
                        cb();
                    },
                },
            });

            BooksysViewDatabaseUpdate.showApp();
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

    static update(info){
        info.isUpdating = true;

        $.ajax({
            type: "GET",
            url: "api/backend.php?action=update_database",
            cache: false,
            success: function(resp){
                var json = $.parseJSON(resp);
                console.log(json);

                // get server info
                info.queries = json.queries;
                info.ok = json.ok;
                info.message = json.message;
                if(json.ok == true){
                    info.isUpdated = true;
                    info.db_version = json.db_version;
                }
                info.isUpdating = false;

                console.log(info);
            },
            error: function(xhr, status, info){
                console.log(xhr);
                console.log(status);
                console.log(info);
            }
        });
    }
}