class BooksysViewPasswordResetEmail {

    constructor(){
        
    }

    // call this function to display the email address dialog
    // - location       ID of the element where to place it
    // - cb             the callback functions (cancel/success)
    // - captcha        reference to a reCaptcha object
    display(location, cb, captcha){
        // load session/user data and display
        this.location  = location;
        this.cb        = cb;
        this.captcha   = captcha;
        this.email     = "";
        this.errors    = [];

        this.loadView();
    }

    // loads the vue template and creates a vue instance
    loadView(){
        let that = this;

        $('#'+this.location).load("res/view_password_reset_email.vue", function(){
            var App = new Vue({
                el: "#view_password_reset_email",
                data: {
                    email:       that.email,
                    errors:      that.errors,
                },
                methods: {
                    "cancel": function(){
                        that.cb.cancel();
                    },
                    "save": function(){
                        that.email  = this.email;
                        this.errors = that.validate();

                        if(this.errors.length == 0){
                            let callbacks = {
                                success: that.cb.success,
                                failure: function(errors){
                                    App.errors = errors;
                                }
                            }
                            that.requestToken(callbacks);
                        }
                    },
                }
            });

            that.showApp();
        });
    }

    validate(){
        let errors = [];

        // lazy isEmail check
        let mailRegex = /(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])/;
        if(this.email.match(mailRegex) == null){
            errors.push("Please provide a valid Email address");
        }

        return errors;
    }

    requestToken(callbacks){
        // send the ajax call to the api
        let request = {
            email:  this.email
        }

        // get the recaptcha_token
        if(this.captcha != null){
            request.recaptcha_token = this.captcha.getToken();
        }

        let that = this;

        $.ajax({
            type: "POST",
            url: "api/password.php?action=token_request",
            data: JSON.stringify(request),
            success: function(data, textStatus, jqXHR){
                let json = $.parseJSON(data);
                if(json.ok){
                    callbacks.success(that.email);
                }else{
                    console.log(json);
                    callbacks.failure([ json.message ]);
                }
            },
            error: function(data, text, errorCode){
                console.log(data);
                console.log(text);
                console.log(errorCode);
                that.errors = [text];
            },			
        });
    }

    showApp(){
        $('#view_password_reset_email_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_password_reset_email_modal').modal('show');
    }

    hideApp(){
        $('#view_password_reset_email_modal').modal('hide');
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