// Initialize reCAPTCHA
var recaptchaOnloadCallback = function(){
    recaptcha.render(document.getElementById("captcha"));
};
let recaptcha = new BooksysRecaptcha("recaptchaOnloadCallback");

// Check if mobile browser first,
// then continue with the page logic
$(function() { 
    if(BooksysBrowser.isMobile()){
        // Make it behave like an app
        BooksysBrowser.setViewportMobile();
        BooksysBrowser.setManifest();
        BooksysBrowser.setMetaMobile();
        // Add mobile style dynamically
        BooksysBrowser.addMobileCSS();
        $.stayInWebApp();

        loadContent();
    } else {
        // Do not add mobile style in desktop mode
        loadContent();
    }
});

function loadContent(){
    let email = "";
    let pwResetEmail = new BooksysViewPasswordResetEmail();
    let pwResetEmailCallback = {
        cancel: function(){
            pwResetEmail.destroyView();
            window.location.href = "/index.html";
        },
        success: function(emailProvided){
            pwResetEmail.destroyView();
            email = emailProvided;
            showResetPassword(email);
        }
    };
    pwResetEmail.display("password_reset", pwResetEmailCallback, recaptcha);
}

function showResetPassword(email){
    let pwResetToken = new BooksysViewPasswordResetToken(email);
    let pwResetPasswordCallback = {
        cancel: function(){
            pwResetToken.destroyView();
            window.location.href = "/index.html";
        },
        success: function(){
            pwResetToken.destroyView();
            window.location.href = "/index.html";
        }
    }

    // display the reset password view
    pwResetToken.display("password_reset", pwResetPasswordCallback, recaptcha);
}

