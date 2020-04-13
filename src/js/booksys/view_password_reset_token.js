class BooksysViewPasswordResetToken {

    constructor(email){
        this.email = email;
        console.log(this.email);
        this.token = null;
        this.password = "";
        this.password2 = "";
    }

    // call this function to display the new password dialog
    // - location       ID of the element where to place it
    // - cb             the callback functions (cancel/success)
    display(location, cb){
        // load session/user data and display
        this.location  = location;
        this.cb        = cb;
        this.errors    = [];

        this.loadView();
    }

    // loads the vue template and creates a vue instance
    loadView(){
        let that = this;
        let reset = {
            email:      that.email,
            token:      that.token,
            password:   that.password,
            password2:  that.password2,
            errors:     that.errors,
        }

        $('#'+this.location).load("res/view_password_reset_token.vue", function(){
            var App = new Vue({
                el: "#view_password_reset_token",
                data: {
                    reset:       reset,
                },
                methods: {
                    "back": function(){
                        that.cb.back();
                    },
                    "save": function(){
                        this.reset.errors = that.validate(reset);

                        if(this.reset.errors.length == 0){
                            that.email = this.reset.email;
                            that.token = this.reset.token;
                            that.password = this.reset.password;
                            
                            let callbacks = {
                                success: that.cb.success,
                                failure: function(errors){
                                    App.reset.errors = errors;
                                }
                            }
                            that.changePassword(callbacks);
                        }
                    },
                }
            });

            that.showApp();
        });
    }

    validate(data){
        let errors = [];

        // token
        let tokenRegex = /^[0-9]+$/;
        if(data.token == null){
            errors.push("Please provide a valid token");
        }else if(data.token.match(tokenRegex) == null){
            errors.push("Please provide a valid token");
        }

        // password strength
        let pwUpperRegex = /[A-Z]+/;
        let pwLowerRegex = /[a-z]+/;
        let pwDigitRegex = /[0-9]+/;
        if(data.password.match(pwUpperRegex) == null){
            errors.push("The password needs to contain at least one upper case letter (A-Z)");
        }
        if(data.password.match(pwLowerRegex) == null){
            errors.push("The password needs to contain at least one lower case letter (a-z)");
        }
        if(data.password.match(pwDigitRegex) == null){
            errors.push("The password needs to contain at least one digit (0-9)");
        }
        if(data.password.length < 10){
            errors.push("The password needs to be at least 10 characters long");
        }

        // password
        if(data.password == ""){
            errors.push("Please provide a password");
        }
        if(data.password != data.password2){
            errors.push("Password and its confirmation are not identical");
        }

        return errors;
    }

    // requests a password reset token from the backend
    changePassword(callbacks){
        // send the ajax call to the api
        let request = {
            email:  this.email,
            token:  this.token,
        };

        // calculate password hash
        request.password = this.calculateHash(this.password);
        if(request.password == null){
            callbacks.failure(["Cannot generate hash of password"]);
            return;
        }

        let that = this;

        $.ajax({
            type: "POST",
            url: "api/password.php?action=change_password_by_token",
            data: JSON.stringify(request),
            success: function(data, textStatus, jqXHR){
                let json = $.parseJSON(data);
                if(json.ok){
                    callbacks.success(that.email);
                }else{
                    console.error(json);
                    callbacks.failure([ json.message ]);
                }
            },
            error: function(data, text, errorCode){
                console.log(data);
                console.log(text);
                console.log(errorCode);
                callbacks.failure([text]);
            },			
        });
    }

    calculateHash(password){
        try {
            let hashInput = password;
            let hashObj = new jsSHA(hashInput, "TEXT");
            let hashOutput = hashObj.getHash("SHA-256","HEX");
            return hashOutput;
        } catch(e) {
            console.error("Password could not be sent encrypted: " + e);
            return null;
        }
    }

    showApp(){
        $('#view_password_reset_token_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_password_reset_token_modal').modal('show');
    }

    hideApp(){
        $('#view_password_reset_token_modal').modal('hide');
		$('.modal-backdrop').remove();
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