// This view is used to edit application settings

class BooksysViewSettings {

    // call this function to display the settings dialog
    // - location              ID of the element where to place the dialog
    // - cb                    callback function upon modal closing (success/cancel)
    static display(location, cb){
        BooksysViewSettings.getData(function(settings){
            BooksysViewSettings.loadView(location, settings, cb);
        });
    }

    static getData(callback){
        $.ajax({
            type: "GET",
            url: "api/configuration.php?action=get_configuration",
            cache: false,
            success: function(data){
                console.log(data);
                let response = JSON.parse(data);
                if(response.ok){
                    callback(response);
                }else{
                    alert(response.msg);
                }
            },
            error: function(xhr, status, error){
                console.log(xhr);
                console.log(status);
                console.log(error);
                var errorMsg = $.parseJSON(xhr.responseText);
                alert(errorMsg.error);
            }
        });
    }

    static showApp(){
        $('#view_settings_form_modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: true,
        });
        $('#view_settings_form_modal').modal('show');
    }

    static hideApp(){
        $('#view_settings_form_modal').modal('hide');
		$('.modal-backdrop').remove();
    }

    // displays the dialog at a given location
    // - location     DOM location (id)
    // - settings     the settings to display
    // - cb           callback function that gets called upon user has been added
    static loadView(location, settings, cb){
        // setup errors array
        let errors = new Array();

        // display the vue-html content
        $('#'+location).load("res/view_settings.vue", function(){    
            var settingsApp = new Vue({
                el: '#view_settings_modal',
                data: {
                    "settings":     settings,
                    "errors":       errors,
                },
                methods: {
                    save: function(event){
                        // validate the db configuration
                        this.errors = BooksysViewSettings.validateEntry(this.settings);
                        if(this.errors.length > 0){
                            return;
                        }

                        // setup db configuration and database
                        let saveCallback = {
                            success: function(){
                                cb.success();
                            },
                            failure: function(sErrors){
                                settingsApp.errors = sErrors;
                            }
                        };
                        BooksysViewSettings.saveEntry(this.settings, saveCallback);
                    },
                    cancel: function(event){
                        BooksysViewSettings.hideApp();
                        cb.cancel();
                    },
                    done: function(event){
                        cb.success();
                    },
                },
            });

            BooksysViewSettings.showApp();
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
    static validateEntry(settings){
        let errors = [];

        // lazy validation (proper validation done on backend)

        // timezone
        let timezoneRegex = /^[A-Za-z]+\/[A-Za-z]+$/;
        if(settings.location_time_zone.match(timezoneRegex) == null){
            errors.push("Please define a time zone with the timezone database name (e.g., Europe/Berlin)");
        }
        // longitude/latitude
        let longitudeRegex = /^[+-]*\d+\.\d+$/;
        if(settings.location_longitude.match(longitudeRegex) == null){
            errors.push("Please define a proper longitude (e.g., 8.542939)");
        }
        if(settings.location_latitude.match(longitudeRegex) == null){
            errors.push("Please define a proper latitude (e.g., 47.367658)");
        }
        // google maps embedded map
        let mapIframeRegex = /https:\/\/www\.google\.com\/maps\/embed\?pb=/;
        // Note: an embedded link could look like '<iframe src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d2700.749935537599!2d8.534332816397074!3d47.3973117791714!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x47900a6544fd6e55%3A0xa1fcf8d958534aec!2sGuggachstrasse%2038-48%2C%208057%20Z%C3%BCrich!5e0!3m2!1sen!2sch!4v1588080267795!5m2!1sen!2sch" width="800" height="600" frameborder="0" style="border:0;" allowfullscreen="" aria-hidden="false" tabindex="0"></iframe>'
        if(settings.location_map != null && settings.location_map.match(mapIframeRegex) == null){
            errors.push("Please select a proper embedded iframe format from maps.google.com");
        }
        // we do not enforce any format for the address

        // currency
        let currencyRegex = /^[A-Z]{1,3}$/;
        if(settings.currency.match(currencyRegex) == null){
            errors.push("Please select a proper currency (e.g., CHF)");
        }
        // we do not enforce any format for the account owner
        // iban
        let ibanRegex = /^([A-Z]{2}[ \-]?[0-9]{2})(?=(?:[ \-]?[A-Z0-9]){9,30}$)((?:[ \-]?[A-Z0-9]{3,5}){2,7})([ \-]?[A-Z0-9]{1,3})?$/;
        if(settings.payment_account_iban != null && settings.payment_account_iban.match(ibanRegex) == null){
            errors.push("Please provide a proper IBAN number");
        }
        // bic
        let bicRegex = /^[a-zA-Z]{4}[a-zA-Z]{2}[a-zA-Z0-9]{2}[XXX0-9]{0,3}/
        if(settings.payment_account_bic != null && settings.payment_account_bic.match(bicRegex) == null){
            errors.push("Please provide a proper BIC/SWIFT Code");
        }
        // we do not enforce any format for the comment

        // sender address
        let emailRegex = /^(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])$/;
        if(settings.smtp_sender != null && settings.smtp_sender.match(emailRegex) == null){
            errors.push("Please provide a proper email address for the sender of emails");
        }
        // smtp server
        let serverRegex = /^[\w0-9]+(\.[\w0-9]+)*(:[0-9]+)*$/;
        if(settings.smtp_server != null && settings.smtp_server.match(serverRegex) == null){
            errors.push("Please provide a proper SMTP server address for where emails should be sent to");
        }
        // we do not enforce any format for the username
        // we do not enforce any format for the password

        // recaptcha private key
        let recaptchaRegex = /^[a-zA-Z0-9]+$/;
        if(settings.recaptcha_privatekey != null && settings.recaptcha_privatekey.match(recaptchaRegex) == null){
            errors.push("Please provide a proper recaptcha private key");
        }
        if(settings.recaptcha_publickey != null && settings.recaptcha_publickey.match(recaptchaRegex) == null){
            errors.push("Please provide a proper recaptcha public key");
        }

        return errors;
    }

    // save the database configuration
    static saveEntry(settings, callbackFn){
        // get the URL from the iframe (if given as entire iframe)
        console.log(settings);
        let mapUrl = [];
        if(settings.location_map != null && settings.location_map != ''){
            let mapUrlRegex = /https:\/\/www\.google\.com\/maps\/embed\?pb=[^"\s]+/;
            mapUrl = [...settings.location_map.match(mapUrlRegex)];
            console.log(mapUrl);
            if(mapUrl.length != 1){
                callbackFn.failure(["The Google Maps IFrame URL could not be recognized"]);
                return;
            }
        }
        

        // build the data object
        let request = {
            "location_time_zone"        : settings.location_time_zone,
            "location_longitude"        : settings.location_longitude,
            "location_latitude"         : settings.location_latitude,
            "location_map"              : mapUrl[0],
            "location_address"          : settings.location_address,
            "currency"                  : settings.currency,
            "payment_account_owner"     : settings.payment_account_owner,
            "payment_account_iban"      : settings.payment_account_iban,
            "payment_account_bic"       : settings.payment_account_bic,
            "payment_account_comment"   : settings.payment_account_comment,
            "smtp_sender"               : settings.smtp_sender,
            "smtp_server"               : settings.smtp_server,
            "smtp_username"             : settings.smtp_username,
            "smtp_password"             : settings.smtp_password,
            "recaptcha_privatekey"      : settings.recaptcha_privatekey,
            "recaptcha_publickey"       : settings.recaptcha_publickey,
        };

        $.ajax({
            type: "POST",
            url: "/api/configuration.php?action=set_configuration",
            cache: false,
            data: JSON.stringify(request),
            success: function(resp){
                let r = JSON.parse(resp);
                if(r.ok){
                    callbackFn.success();
                }else{
                    if(r.errors.length > 0){
                        callbackFn.failure(r.errors);
                    }else{
                        callbackFn.failure(r.msg);
                    }
                }
            },
            error: function(xhr, textStatus, error){
                console.error(xhr);
                console.error(textStatus);
                console.error(error);
                callbackFn.failure([error]);
            }
        });
    }
}