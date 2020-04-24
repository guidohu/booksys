$.stayInWebApp();
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

        $("#body").load("res/log_mobile.html", function(){
            loadContent();
        });
    } else {
        $("#body").load("res/log.html", function(){
            loadContent();
        });
    }
});

function loadContent(){
    getLogs();
}

// Returns the account information
function getLogs(){
    // get the account information
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_body').html('Retrieving data...');
    $('#status_modal').modal('show');
    $.ajax({
        type: "GET",
        cache: false,
        url: "api/log.php?action=get_logs",
        success: function(data){
            var json = $.parseJSON(data);
            updateLog(json);
            $('#status_modal').modal('hide');
        }
    });
}

function updateLog(log){
    var content = '';
    
    if(log.length == 0){
        content += 'No logs to display';
    }
    
    for(var i=0; i<log.length; i++){
        if(i%2==0){
            content += "<div class=\"panel-row-odd row\" style=\"padding-top:10px;\">";
        }else{
            content += "<div class=\"panel-row-even row\" style=\"padding-top:10px;\">";
        }
        
        content += "<div class='row'>";
        content += "<div class=\"col-sm-12 col-xs-12 text-left\">";
        content += "<div class='col-sm-3 col-xs-3'>";
        content +=     log[i].time;
        content += "</div>";
        content += "<div class='col-sm-9 col-xs-9'>";
        content +=     encodeHTMLEntities(log[i].log);
        content += "</div>";
        content += "</div>";
        content += "</div>";
        
        content += "</div>";
    }
    
    $("#logs").html(content);
}

// encode HTML code into plain text
function encodeHTMLEntities(text) {
      var textArea = document.createElement('textarea');
      textArea.innerText = text;
      return textArea.innerHTML;
}