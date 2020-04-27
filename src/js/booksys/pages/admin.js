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

        $("#body").load("res/admin_mobile.html");
    } else {
        $("#body").load("res/admin.html");
    }
});

// Stay in the WebApp
$(function() {
    $.stayInWebApp();
});