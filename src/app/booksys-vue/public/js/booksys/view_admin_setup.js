// This view is used for an initial admin user setup

class BooksysViewAdminSetup {

    // call this function to display the setup admin user dialog
    // - location              ID of the element where to place the dialog
    // - cb                    callback function upon modal closing (success/cancel)
    static displayEntry(location, cb){
        BooksysViewAdminSetup.checkSetup(function(config){
            BooksysViewAdminSetup.loadView(location, config, cb);
        });
    }

    // get the current database configuration (if any) and calls the
    // callback function cb
    static checkSetup(cb){
        $.ajax({
            type: "GET",
            url: "api/configuration.php?action=is_admin_user_configured",
            cache: false,
            success: function(data){
                let response = $.parseJSON(data);
                console.log(response);
                cb(response);
            },
            error: function(data){
                console.log(data);
                cb(null);
            }
        });
    }

    static showApp(){
        $('#view_admin_setup_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_admin_setup_modal').modal('show');
    }

    static hideApp(){
        $('#view_admin_setup_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays the dialog at a given location
    // - location     DOM location (id)
    // - bResponse    the response from the backend
    // - cb           callback function that gets called upon changed fuel entry
    static loadView(location, bResponse, cb){
        
        // setup errors array
        let errors = new Array();

        // user
        let user = new Array();
        let isCreation = false;

        // error and response handling
        if(bResponse == null){
            errors = [ "Unknown error on backend" ];
        }
        else {
            console.log(bResponse);
            if(bResponse.OK == true && bResponse.admin_exists == false){
                isCreation = true;
            }
            else if(bResponse.OK == false){
                errors = [ bResponse.msg ];
            }
        }

        // display the vue-html content
        $('#'+location).load("res/view_admin_setup.vue", function(){    
            var databaseSetupApp = new Vue({
                el: '#view_admin_setup_modal',
                data: {
                    "user":         user,
                    "errors":       errors,
                    "isCreation":   isCreation,
                },
                methods: {
                    save: function(event){
                        // validate the db configuration
                        this.errors = BooksysViewAdminSetup.validateEntry(this.user);
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
                        BooksysViewAdminSetup.saveEntry(this.user, saveCallback);
                    },
                    cancel: function(event){
                        BooksysViewAdminSetup.hideApp();
                        cb.cancel();
                    },
                },
            });

            BooksysViewAdminSetup.showApp();
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
    static validateEntry(user){
        let errors = [];

        // lazy validation (proper validation done on backend)
        if(! user.email){
            errors.push("Please select a proper email as your username");
        }
        if(user.password != user.passwordConfirm){
            errors.push("The two passwords do not match, please check for typos");
        }
        if(! user.password){
            errors.push("Please select a password");
        }
        if(! user.firstName){
            errors.push("Please select a first name");
        }
        if(! user.lastName){
            errors.push("Please select a last name");
        }
        if(! user.address){
            errors.push("Please enter your address");
        }
        if(! user.zipCode){
            errors.push("Please enter your zip code");
        }
        if(! user.city){
            errors.push("Please select a place / city");
        }
        if(! user.phone){
            errors.push("Please select a phone number");
        }

        // password strength
        let pwUpperRegex = /[A-Z]+/;
        let pwLowerRegex = /[a-z]+/;
        let pwDigitRegex = /[0-9]+/;
        if(user.password.match(pwUpperRegex) == null){
            errors.push("The password needs to contain at least one upper case letter (A-Z)");
        }
        if(user.password.match(pwLowerRegex) == null){
            errors.push("The password needs to contain at least one lower case letter (a-z)");
        }
        if(user.password.match(pwDigitRegex) == null){
            errors.push("The password needs to contain at least one digit (0-9)");
        }
        if(user.password.length < 10){
            errors.push("The password needs to be at least 10 characters long");
        }

        return errors;
    }

    // save the database configuration
    static saveEntry(user, callbackFn){
        // calculate password hash
        let password = '';
        try {
			var hashInput  = user.password;
			var hashObj    = new jsSHA(hashInput, "TEXT");
			password       = hashObj.getHash("SHA-256","HEX");	
		} catch(e) {
			callbackFn.failure([ "Password could not be transformed to a Sha256 hash: " + e ]);
			return;
		};

        // construct user data
        let userData = {
            username:   user.email,
            password:   password,
            first_name: user.firstName,
            last_name:  user.lastName,
            address:    user.address,
            mobile:     user.phone,
            plz:        user.zipCode,
            city:       user.city,
            email:      user.email,
            license:    true
        };

        // sign-up user and make it admin
        $.ajax({
            type: "POST",
            url: "api/sign_up.php?action=sign_up",
            cache: false,
            data: JSON.stringify(userData),
            success: function(resp){
                let data = $.parseJSON(resp);
                console.log(data);
                if(data.ok == true){
                    console.log("User created");
                    BooksysViewAdminSetup.makeUserAdmin(data.user_id, callbackFn);
                }else{
                    console.log("Admin user not created");
                    callbackFn.failure([ data.msg ]);
                }    
            },
        });
    }

    // unlocks the newly created user and makes it admin
    static makeUserAdmin(userId, callbackFn){
        // make the user admin
        let data = {
            user_id:    userId,
        };

        $.ajax({
            type: "POST",
            url: "api/configuration.php?action=make_user_admin",
            cache: false,
            data: JSON.stringify(data),
            success: function(resp){
                let data = $.parseJSON(resp)
                if(data.ok == true){
                    console.log("User changed to admin");
                    BooksysViewAdminSetup.hideApp(); 
                    callbackFn.success();
                }else{
                    console.log("User not changed to admin");
                    callbackFn.failure([ data.msg ]);
                }    
            },
        });
    }
}