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

        $("#body").load("res/info_mobile.html", function(){
            // Stay in the WebApp
            $(function() {
                $.stayInWebApp();
            });
            loadContent();
        });
    } else {
        $("#body").load("res/info.html", function(){
            loadContent();
        });
    }
});

function loadContent(){
    displayLocationMap();
}

function displayLocationMap(){
    $.ajax({
        type: "GET",
        url: 'api/configuration.php?action=get_customization_parameters',
        success: function(data){
            var json = JSON.parse(data);

            // add the address
            if(json.location_address != null){
                $('#location-address').html(json.location_address);
            }else{
                $('#location-address').html("< no address provided for display >");
            }

            // create the map
            if(json.location_map_iframe != null){
                $('#location-map').html(json.location_map_iframe);
                let iframe = document.getElementById("location-map").firstElementChild;
                if(BooksysBrowser.isMobile()){
                    iframe.setAttribute("width", "340");
                    iframe.setAttribute("height", "350");
                }else{
                    iframe.setAttribute("width", "600");
                    iframe.setAttribute("height", "300");
                }
            }else{
                $('#location-map').html("< no map available >");
            }
        },
        error: function(data, text, errorCode){
            // something went wrong
            console.error("Cannot call get_customization_parameters");
            console.error(text);
            console.log(data);
            console.log(errorCode);
            $('#location-address').html("< Address cannot be displayed >");
            $('#location-map').html("< Location cannot be displayed on a map ></br>You might need to login first.");
        }
    });
}
