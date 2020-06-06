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

        $("#body").load("res/calendar_mobile.html", function(){
            loadContent();
        });
    } else {
        $("#body").load("res/calendar.html", function(){
            loadContent();
        });
    }
});

// the month we currently display
var displayedMonth = moment().tz(getTimeZone()).startOf('month');

// load all the dynamic content
function loadContent(){
    $('#calendar').hide();
    updateBookings(displayedMonth);

    // stay in a webapp
    $.stayInWebApp();
}

// Requests the new data from the server and
// updates the view
// - date		startOfMonth
function updateBookings(date){
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_modal').modal('show');

    // get the session-information from the database
    var data = new Object();
    data.timeWindows = new Array();

    // dateIterator is used to loop through all dates
    // to build the calendar view 
    var dateIterator;
    if(date == null){
        dateIterator = moment().tz(getTimeZone()).startOf('month');
    }else{
        dateIterator = date.clone();
    }
    var monthText = date.format("MMMM") + " " + date.format("YYYY");

    // go back to a Monday (our calendar starts with Mondays)
    while(dateIterator.tz(getTimeZone()).day() != 1){
        dateIterator.tz(getTimeZone()).subtract(1, 'day');
    }

    // get 42 time windows
    for(var i = 0; i < 42; i++){
        // add a time window
        var window = {
            start: dateIterator.tz(getTimeZone()).startOf('day').format("X"),
            end:   dateIterator.tz(getTimeZone()).endOf('day').format("X"),
        };
        data.timeWindows.push(window);

        dateIterator.tz(getTimeZone()).add(1, 'day');
    }

    $.ajax({
        type: "POST",
        url: "api/booking.php?action=get_booking_month",
        data: JSON.stringify(data),
        success: function(data){
            var json = $.parseJSON(data);

            // in case we are not logged in --> redirect to login page
            if(typeof json.redirect != "undefined"){
                console.log(window);
                location.href = json.redirect;
                return;
            }
                
            $("#titleDate").text(monthText);
            $("#detail_rider").html("-");
            $("#detail_date").html("please select a session");
            // console.log(displayedMonth.format());
            drawPies(json, date);
            $("#calendar").fadeIn();
            $('#status_modal').modal('hide');				
        },
        error: function(data, text, errorCode){
            $("#status_body").html("Error" + data.status + ": " + data.responseText);
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

// Draws the pie with the information given
function drawPies(data, displayedMonth){
    // prepare placeholders
    $('.pie-small.pie-small-inactive').each(function (index) {
        $(this).removeClass('pie-small-inactive');
    });

    // empty placeholders
    $("[class=pie-small]").each(function (index) {
        $(this).html("");
    });

    // get all placeholders
    $("[class=pie-small]").each(function (index) {
        var session = data[index];
        var start   = moment(session.window_start, "X").tz(getTimeZone());
        var end     = moment(session.window_end, "X").tz(getTimeZone());
        var day     = moment(session.window_start, "X").tz(getTimeZone()).format("D");
        var html    = "<div id='pieChart"+index+"''></div>\n";
        html       += "<div class='calendar-day'>" + day +"<div>";
        $(this).html(html)
        //$(this).html("<div class='calendar-day'>" + day +"<div>");
        $(this).bind('mouseenter', function(){
            showDetails(session);
        });
        
        $(this).bind('click', function(){
            clickDay(session);
        });

        // different shade whether it is in the current month or not
        if(displayedMonth.tz(getTimeZone()).startOf('month') > start.tz(getTimeZone())
            || displayedMonth.tz(getTimeZone()).endOf('month') < start.tz(getTimeZone()))
        {
            // console.log("Session outside of month");
            $(this).addClass('pie-small-inactive');
        }
        else
        {
            $(this).bind('click', function(){
                clickDay(session);
            });			
        }
    });
    
    var properties = {
        containerHeight:	64,
        containerWidth:     90,
        circleX:			45,
        circleY:            34,
        circleRadius:       23,
        animate:            false,
    }
    if(BooksysBrowser.isMobile()){
        console.log(document.getElementById('calendar_container').getBoundingClientRect());
        let tileWidth = document.getElementById('calendar_container').getBoundingClientRect().width / 7;
        properties = {
            containerHeight:	66,
            containerWidth:     tileWidth,
            circleX:			tileWidth/2,
            circleY:            34,
            circleRadius:       20,
            animate:            false,
        }
    }
    for(var j=0; j<data.length; j++){
        BooksysPie.drawPie("pieChart"+j, data[j], null, properties);
    }
}

// Switch to next month
function nextMonth(){
    displayedMonth.add(1, 'month');
    displayedMonth = displayedMonth.startOf('month');
    updateBookings(displayedMonth);
}

// Switch to previous month
function prevMonth(){
    displayedMonth.subtract(1, 'month');
    displayedMonth = displayedMonth.startOf('month');
    updateBookings(displayedMonth);
}

// display the details to that day and its sessions
function showDetails(day){
    var date_text = moment(day.window_start, "X").tz(getTimeZone()).format("dddd DD.MM.YYYY");
    $("#detail_date").text(date_text);
    
    // display sessions
    var text = '';
    for(var i=0; i<day.sessions.length; i++){
        session = day.sessions[i];
        
        var start = moment(session.start, "X").tz(getTimeZone());
        var end   = moment(session.end, "X").tz(getTimeZone());
        text += start.tz(getTimeZone()).format('HH:mm') + ' - ' + end.tz(getTimeZone()).format('HH:mm') + '</br>';
        for(var j=0; j<session.riders.length; j++){
            text += session.riders[j].name + '</br>';
        }
        text += '</br>';
    }
    
    if(text == ''){
        text='no sessions';
    }
    
    $('#detail_sessions').html(text);

}

// click action for tiles
function clickDay(day){
    var startOfDay = moment(day.window_start, "X").tz(getTimeZone());
    var dateStr    = startOfDay.format("YYYY-MM-DD");
    goToURL("today.html?getdate="+dateStr);
}

// go to a given URL
function goToURL(url){
    window.location.href = url;
}