class BooksysRecaptcha {

    // Google Recaptcha v2 helper library that only loads and renders
    // it in case the keys are actually present in the configuration of the App

    // Typical usage:
    // var recaptchaOnloadCallback = function(){
    //     recaptcha.render(document.getElementById("captcha"));
    // };
    // let recaptcha = new BooksysRecaptcha("recaptchaOnloadCallback");
    
    constructor(onloadCallbackName){
        this.key = null;
        this.token = null;
        this.getRecaptchaKey(onloadCallbackName);
    }

    // try to get a recaptcha key from the server
    // and load the javascript library
    getRecaptchaKey(onloadCallbackName){
        let that = this;
        $.ajax({
            type:   "GET",
            url:    "api/configuration.php?action=get_recaptcha_key",
            cache:  false,
            success: function(resp){
                let r = $.parseJSON(resp);
                that.key = r.key;
                if(that.key != null){
                    that.getRecaptchaLibrary(onloadCallbackName);
                }
            },
            error: function(xhr, status, error){
                console.log(xhr);
                console.log(status);
                console.log(error);
                var errorMsg = $.parseJSON(xhr.responseText);
                that.key = null;
            }
        });
    }

    // downloads the recaptcha library
    getRecaptchaLibrary(onloadCallbackName){
        $.ajax({
            url: "https://www.google.com/recaptcha/api.js?onload="+onloadCallbackName+"&render=explicit",
            dataType: 'script',
            cache: false,
            success: function(){},
            async: true
        });
    }

    // renders the recaptcha in the given domElement
    render(domElement){
        let that = this;
        grecaptcha.render(
            domElement, 
            {
                'sitekey' : that.key,
                'theme' : 'light',
                'callback': function(response){that.storeToken(response)}
            }
        );
    }

    // stores the token
    storeToken(response){
        this.token = response;
    }

    // returns the token
    getToken(){
        return this.token;
    }
}