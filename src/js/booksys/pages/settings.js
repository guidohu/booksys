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
    let settingsCallback = {
        success: function(){
            BooksysViewSettings.destroyView('settings');
            window.location.href = "/admin.html";
        },
        cancel: function(){
            BooksysViewSettings.destroyView('settings');
            window.location.href = "/admin.html";
        }
    };
    BooksysViewSettings.display('settings', settingsCallback);
}