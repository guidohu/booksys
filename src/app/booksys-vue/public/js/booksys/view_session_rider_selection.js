class BooksysViewSessionRiderSelection {

    // call this function to display the rider selection window
    // - location       ID of the element where to place it
    // - id             the ID of the session
    // - cb             callback upon exit function (called when done)                 
    static display(location, id, cb){
        // load session and display
        BooksysViewSessionRiderSelection.getData(function(users){
            BooksysViewSessionRiderSelection.loadView(location, id, users, cb);
        });  
    }

    // loads the users from the server
    static getData(cb){
        var users = [];

        $.ajax({
            type: "GET",
            url: "api/user.php?action=get_all_users",
            cache: false,
            contentType: "application/json; charset=utf-8",
            success: function(data){
                users = $.parseJSON(data);
                cb(users);
            },
            error: function(request, status, error){
                console.log(request);
                console.log(status);
                console.log(error);
                alert("Cannot get users from server");
            }
        });

        return;
    }

    // displays the App
    static showApp(){
        $('#view_session_rider_selection_modal').on('shown.bs.modal', function () {
            $('#search').focus();
        });

        $('#view_session_rider_selection_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_session_rider_selection_modal').modal('show');
    }

    // hides the App
    static hideApp(){
        $('#view_session_rider_selection_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // removes the content that has been added
    static destroyView(location){
        var child = document.getElementById(location).lastElementChild;
        while (child){
            document.getElementById(location).removeChild(child)
            child = document.getElementById(location).lastElementChild;
        }
    }

    // loads the vue template and creates a vue instance
    // - location        where to display the view
    // - id              the ID of the session
    // - users           array of users to select from
    // - cb              called upon cancel or adding of a user to a session
    static loadView(location, id, users, cb){
        var searchString = "";
        var filteredUsers = users;

        $('#'+location).load("res/view_session_rider_selection.vue", function(){
            // prevent the form to do its default submit on
            // enter
            $(window).keydown(function(event){
                if(event.keyCode == 13) {
                    event.preventDefault();
                    return false;
                }
            });
            
            var App = new Vue({
                el: "#view_session_rider_selection_modal",
                data: {
                    "users": users,
                    "filteredUsers": filteredUsers,
                },
                methods: {
                    "cancel": function(){
                        BooksysViewSessionRiderSelection.hideApp();
                        if(cb != null){
                            cb();
                        }
                    },
                    "add": function(){
                        BooksysViewSessionRiderSelection.addUsers(id, cb);
                    },
                    "filter": function(element, v1, v2){
                        if(element.key == "Enter"){
                            BooksysViewSessionRiderSelection.addUsers(id, cb);
                        }else{
                            searchString = element.srcElement.value;
                            this.filteredUsers = BooksysViewSessionRiderSelection.filterUsers(searchString, users);
                        }
                    }
                }
            });

            BooksysViewSessionRiderSelection.showApp();
        });       
    }

    // adds a user to a session
    // - sessionId      ID of the session to add users to
    // - cb             callback function to call upon successfully adding a user to a session
    static addUsers(sessionId, cb){
        var riders = new Array();
		$("#rider_users option:selected").each( function(){
				riders.push($(this).val());
			}
		)
		
		if(riders.length == 0){
			alert("Please select at least one rider.");
			return;
		}

        // build request data according to API
        var req = {
            user_ids: riders,
            session_id: sessionId
        };	
	
		// send the request to the server
		$.ajax({
			type: "POST",
			url: "api/booking.php?action=add_users",
			data: JSON.stringify(req),
			success: function(data){
                BooksysViewSessionRiderSelection.hideApp();
                if(cb != null){
                    cb();
                }
			},
			error: function(data){
				alert("Ooops. There was an error: " + data.responseText);
			}
		});
    }

    // filters for users that match the search string
    static filterUsers(searchString, users){
        searchString = searchString.toLowerCase();
        var filteredUsers = [];

        // filter for users that match the search
        for(var i=0; i<users.length; i++){
            var u = users[i];
            if(searchString == ""){
                u.visible = true;
                filteredUsers.push(u);
                continue;
            }

            var username = u.first_name + " " + u.last_name;
            username = username.toLowerCase();
            // console.log("'" + username + "' match against '" + searchString + "'")
            if(username.includes(searchString)){
                u.visible = true;
                filteredUsers.push(u);
            }else{
                u.visible = false;
            }
        }

        return filteredUsers;
    }
}