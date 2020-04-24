// Check if mobile browser first,
// then continue with checks
$(function() { 
    if(BooksysBrowser.isMobile()){
        // Make it behave like an app
        BooksysBrowser.setViewportMobile();
        BooksysBrowser.setManifest();
        BooksysBrowser.setMetaMobile();
        // Add mobile style dynamically
        BooksysBrowser.addMobileCSS();

        // $("#body").load("res/index_mobile.html", function(){
        //     backendStatus();
        // });
        loadContent();
    } else {
        // $("#body").load("res/index.html", function(){
        //     backendStatus();
        // });
        loadContent();
    }
});

// loads the content
function loadContent(){
    // handle returning from the user setup
    let userSetupCallback = {
        success: function(){
            BooksysViewAdminSetup.destroyView('body');
            window.location.href = "/index.html";
        },
        cancel: function(){
            BooksysViewAdminSetup.destroyView('body');
            window.location.href = "/index.html";
        }
    };

    // handle returning from the database setup
    let dbSetupCallback = {
        success: function(){
            BooksysViewDatabaseSetup.destroyView('body');
            BooksysViewAdminSetup.displayEntry('body', userSetupCallback);
        },
        cancel: function(){
            BooksysViewDatabaseSetup.destroyView('body');
            window.location.href = "/index.html";
        }
    }

    BooksysViewDatabaseSetup.displayEntry('body', dbSetupCallback);
}
