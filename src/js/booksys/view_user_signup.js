// This view is used for to create a user

class BooksysViewUserSignup {

    // call this function to display the user creation dialog
    // - location              ID of the element where to place the dialog
    // - cb                    callback function upon modal closing (success/cancel)
    static display(location, cb, captcha){
        BooksysViewUserSignup.loadView(location, cb, captcha);
    }

    static showApp(){
        $('#view_user_signup_form_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_user_signup_form_modal').modal('show');
    }

    static hideApp(){
        $('#view_user_signup_form_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays the success
    static showSignUpSuccessful(){
        BooksysViewUserSignup.hideApp();
        $('#view_user_signup_created_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_user_signup_created_modal').modal('show');
    }

    // displays the dialog at a given location
    // - location     DOM location (id)
    // - cb           callback function that gets called upon user has been added
    static loadView(location, cb, captcha){
        // setup errors array
        let errors = new Array();

        // user
        let user = new Array();
        user.driving_license = false;

        // display the vue-html content
        $('#'+location).load("res/view_user_signup.vue", function(){    
            var databaseSetupApp = new Vue({
                el: '#view_user_signup_modal',
                data: {
                    "user":         user,
                    "errors":       errors,
                },
                methods: {
                    save: function(event){
                        // validate the db configuration
                        this.errors = BooksysViewUserSignup.validateEntry(this.user);
                        if(this.errors.length > 0){
                            return;
                        }

                        // setup db configuration and database
                        let saveCallback = {
                            success: function(){
                                BooksysViewUserSignup.showSignUpSuccessful();
                            },
                            failure: function(sErrors){
                                databaseSetupApp.errors = sErrors;
                            }
                        };
                        BooksysViewUserSignup.saveEntry(this.user, saveCallback, captcha);
                    },
                    cancel: function(event){
                        BooksysViewUserSignup.hideApp();
                        cb.cancel();
                    },
                    done: function(event){
                        cb.success();
                    },
                },
            });

            BooksysViewUserSignup.showApp();
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
        if(! user.accepted_no_school){
            errors.push("Please confirm that you know that this is not a wakeboard school.");
        }
        if(! user.accepted_terms_and_conditions){
            errors.push("Please accept our conditions. If you are not familiar with them, please contact one of the team members.");
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
    static saveEntry(user, callbackFn, captcha){
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
        let drivingLicense  = false;
        if(user.driving_license){
            drivingLicense = true;
        }
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
            license:    drivingLicense,
        };
        if(captcha != null){
            userData.recaptcha_token = captcha.getToken();
        }

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
                    callbackFn.success();
                }else{
                    console.log("User not created");
                    callbackFn.failure([ data.msg ]);
                }    
            },
            error: function(resp){
                console.log(resp);
                callbackFn.failure("Cannot create user due to unknown error.");
            }
        });
    }
}