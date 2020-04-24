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

        $("#body").load("res/my_cal_mobile.html", function(){
            loadContent();
        });
    } else {
        // Add mobile style also in desktop mode
        $("#body").load("res/my_cal.html", function(){
            loadContent();
        });
    }
});

// Stay in the WebApp
$(function() {
    $.stayInWebApp();
});

function loadContent(){
    getSchedule();
}

// Returns the schedule for this user
function getSchedule(){
    // get schedule information for this user
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_body').html('Retrieving data...');
    $('#status_modal').modal('show');
    $.ajax({
        type: "POST",
        url: "api/user.php?action=get_my_user_sessions",
        cache: false,
        success: function(data){
            //console.log(data);
            let schedule = $.parseJSON(data);
            updateScheduleView(schedule);
            $('#status_modal').modal('hide');
        }
    });
}

function updateScheduleView(info){

    // add all the past sessions to the schedule
    var content = '';
    
    if(info.sessions_old.length == 0){
        content += "No sessions in your personal history";
    }
    
    for(var i=0; i<info.sessions_old.length; i++){
        if(i%2==0){
            content += "<div class=\"panel-row-odd row\" style=\"padding-top:10px;\">";
        }else{
            content += "<div class=\"panel-row-even row\" style=\"padding-top:10px;\">";
        }

        let start = moment(info.sessions_old[i].start, "X").tz(getTimeZone());
        let end   = moment(info.sessions_old[i].end, "X").tz(getTimeZone());
        
        content += "<div class=\"col-sm-12 col-xs-12  text-left\">";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-2 col-xs-6\">";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += start.format('ddd') + ", " + start.format('DD.MM.YYYY');
        content += "</div>";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += start.format('HH:mm') + " - " + end.format('HH:mm');
        content += "</div>";
        content += "</div>";
        content += "</div>";
        content += "<div class=\"col-sm-8 col-xs-6\">";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += '<strong>' + info.sessions_old[i].title + '</strong>';
        content += "</div>";
        content += "</div>";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += '<div class="col-sm-2 hidden-xs">';
        content += "Riders: ";
        content += '</div>';            
        content += '<div class="col-sm-10 col-xs-12">';
        content += '<small>';
        if(info.sessions_old[i].riders.length < 1){            
            content += "none";
        }else{
            for(var r=0; r<info.sessions_old[i].riders.length; r++){
                content += info.sessions_old[i].riders[r].first_name;
                content += ' ';
                content += info.sessions_old[i].riders[r].last_name;
                if(r+1 < info.sessions_old[i].riders.length){
                    content += ', ';
                }
            }
        }
        content += '</small>';
        content += '</div>';
        content += "</div>";
        content += "</div>";
        content += "</div>";
        content += "<div class=\"col-sm-2 col-xs-6 hidden-xs text-right\">";
        
        // It is past, we cannot change sessions anymore
        // thus we disable it for the non xs view
        content += "<a class=\"btn btn-default disabled\" href=\"#\"><span class=\"glyphicon glyphicon-ok\"></span></a>";
        content += "<a class=\"btn btn-default disabled\" href=\"#\"><span class=\"glyphicon glyphicon-remove\"></span></a>";
        content += "</div>";
        content += "</div>";
        content += "</div>";
        content += "</div>";
    }
    
    $('#past_sessions').html(content);
    
    // add all the upcoming sessions to the schedule
    content  = '';
    
    if(info.sessions.length == 0){
        content += "No upcoming sessions or invitations";
    }
    
    for(var i=0; i<info.sessions.length; i++){
        if(i%2==0){
            content += "<div class=\"panel-row-odd row\" style=\"padding-top:10px;\">";
        }else{
            content += "<div class=\"panel-row-even row\" style=\"padding-top:10px;\">";
        }

        let start = moment(info.sessions[i].start, "X").tz(getTimeZone());
        let end   = moment(info.sessions[i].end, "X").tz(getTimeZone());
        
        content += "<div class=\"col-sm-12 col-xs-12 text-left\">";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-2 col-xs-6\">";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += start.format('ddd') + ", " + start.format('DD.MM.YYYY');
        content += "</div>";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += start.format('HH:mm') + " - " + end.format('HH:mm');
        content += "</div>";
        content += "</div>";
        content += "</div>";
        content += "<div class=\"col-sm-8 col-xs-6 hidden-xs\">";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += '<strong>' + info.sessions[i].title + '</strong>';
        content += "</div>";
        content += "</div>";
        content += "<div class=\"row\">";
        content += "<div class=\"col-sm-12 col-xs-12\">";
        content += '<div class="col-sm-2">';
        content += "Riders: ";
        content += '</div>';            
        content += '<div class="col-sm-10 col-xs-10">';
        content += '<small>';
        if(info.sessions[i].riders.length < 1){            
            content += "none";
        }else{
            for(var r=0; r<info.sessions[i].riders.length; r++){
                content += info.sessions[i].riders[r].first_name;
                content += ' ';
                content += info.sessions[i].riders[r].last_name;
                if(r+1 < info.sessions[i].riders.length){
                    content += ', ';
                }
            }
        }
        content += '</small>';
        content += '</div>';
        content += "</div>"
        content += "</div>"
        content += "</div>"
        content += "<div class=\"col-sm-2 col-xs-6 text-right\">"
        
        // If it is a course -> we are driver and can't cancel, unfortunately
        if(info.sessions[i].type == 1){
            content += "<button class=\"btn btn-default disabled\"><span class=\"glyphicon glyphicon-ok\"></span></button>"
            content += "<button class=\"btn btn-default disabled\"><span class=\"glyphicon glyphicon-remove\"></span></button>"
        }
        // If the type is a default session, check if it is only an invitation            
        else{
            // user is in the session
            if(info.sessions[i].status == 3){
                content += "<button class=\"btn btn-default disabled\"><span class=\"glyphicon glyphicon-ok\"></span></button>"
                content += "<button class=\"btn btn-danger\" onclick=\"cancelSession(" + info.sessions[i].id + ")\"><span class=\"glyphicon glyphicon-remove\"></span></button>"
            }
            // user only got an invitation so far
            else if(info.sessions[i].status == 2){
                content += "<button class=\"btn btn-success\" onclick=\"acceptSession(" + info.sessions[i].id + ")\"><span class=\"glyphicon glyphicon-ok\"></span></button>"
                content += "<button class=\"btn btn-danger\" onclick=\"declineSession(" + info.sessions[i].id + ")\"><span class=\"glyphicon glyphicon-remove\"></span></button>"
            }
            // the session got canceled -> display this too
            else if(info.sessions[i].status == 1){
                content += "<button class=\"btn btn-inverse disabled\"><span class=\"glyphicon glyphicon-ok\"></span></button>"
                content += "<button class=\"btn btn-inverse disabled\"><span class=\"glyphicon glyphicon-remove\"></span></button>"
            }
        }
        content += "</div>"
        content += "<div class=\"span1\"></div>"
        content += "</div>"
        content += "</div>"
        content += "</div>"
    }
    
    $('#upcoming_sessions').html(content);
    
}

// cancels a session
function cancelSession(session_id){
    // cancel Session
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_body').html('Canceling Session...');
    $('#status_modal').modal('show');
    var data = new Object();
    data['session_id'] = session_id;
    
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=delete_user",
        data: JSON.stringify(data),
        success: function(data){
            getSchedule();
        }
    });
}

// accept an inviation to a session
function acceptSession(session_id){
    // accept a Session
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_body').html('Accepting Session...');
    $('#status_modal').modal('show');
    
    var data = Object();
    data['session_id'] = session_id;
    
    $.ajax({
        type: "POST",
        url: "api/invitation.php?action=accept",
        data: JSON.stringify(data),
        success: function(data){
            getSchedule();
        }
    });
}

// declines a session
function declineSession(session_id){
    // decline Session
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_body').html('Declining Session...');
    $('#status_modal').modal('show');
    var data = new Object();
    data['session_id'] = session_id;
    
    $.ajax({
        type: "POST",
        url: "api/invitation.php?action=decline",
        data: JSON.stringify(data),
        success: function(data){
            getSchedule();
        }
    });
}

// gets the user's timezone or returns the default
function getTimeZone(){
    if (typeof(Storage) !== "undefined") {
        var timezone = localStorage.getItem("timezone");
        if(timezone != null){
            return timezone;
        }
        else{
            // TODO let user set his timezone
            return "Europe/Berlin";
        }
    } else {
        // Sorry! No Web Storage support..
        return "Europe/Berlin";
    }
}