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
    // handle returning from the user setup
    // let userSetupCallback = {
    //     success: function(){
    //         BooksysViewUserSignup.destroyView('sign_up_user');
    //         window.location.href = "/index.html";
    //     },
    //     cancel: function(){
    //         BooksysViewUserSignup.destroyView('sign_up_user');
    //         window.location.href = "/index.html";
    //     }
    // };
    BooksysViewSettings.display('settings', null);
}