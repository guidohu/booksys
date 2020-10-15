// Go to the next day
function nextDay(){
    selectedSessionIdx = -1;
    BooksysPie.resetSelection();
    startDay = startDay.tz(getTimeZone()).add(1, 'day').startOf('day');
    endDay   = endDay.tz(getTimeZone()).add(1, 'day').endOf('day');
    updateBookings(startDay, endDay);
}

// Go to the previous day
function prevDay(){
    selectedSessionIdx = -1;
    BooksysPie.resetSelection();
    startDay = startDay.tz(getTimeZone()).subtract(1, 'day').startOf('day');
    endDay   = endDay.tz(getTimeZone()).subtract(1, 'day').endOf('day');
    updateBookings(startDay, endDay);
}

// Requests the new data from the server and
// updates the view
function updateBookings(start, end){
    // hide continue button upon reload
    $("#btn_stop_watch").hide();

    // get the session-information from the database
    var data = {
        start:	start.tz(getTimeZone()).format("X"),
        end:    end.tz(getTimeZone()).format("X")
    };
    $.ajax({
        type: "POST",
        url: "api/booking.php?action=get_booking_day",
        data: JSON.stringify(data),
        success: function(resp){
            var json = $.parseJSON(resp);

            // in case we are not logged in --> redirect to login page
            if(typeof json.redirect != "undefined"){
                window.location.replace = json.redirect;
                return;
            }
            
            // Some general info which has to be adjusted
            $("#titleDate").text(start.tz(getTimeZone()).format("dddd DD.MM.YYYY"));
            var sunrise = moment(json.sunrise, "X").tz(getTimeZone());
            var sunset  = moment(json.sunset, "X").tz(getTimeZone());
            $("#detail_sunrise").html(sunrise.format('HH:mm'));
            $("#detail_sunset").html(sunset.format('HH:mm'));

            // add event listener for titleDate
            document.getElementById("titleDate").addEventListener("click", function(){
                // alert("clicked");
                window.location.href = "/calendar.html?";
            });
            
            // Draw the pie and get the pie content returned
            let properties = {
                containerHeight: 350,
                containerWidth:  700,
                circleX:         300,
                circleY:         170,
                circleRadius:    100,
                animation:       true,
                labels:          true,
            }
            if(BooksysBrowser.isMobile()){
                properties = {
                    containerHeight: 300,
                    containerWidth:  350,
                    circleX:         175,
                    circleY:         150,
                    circleRadius:    100,
                    animation:       false,
                    labels:          true,
                }
                BooksysViewSessionDetails.display("session_details", null, null, function(){});
            }
            
            pieSessions = BooksysPie.drawPie("pie", json, function(id){
                updateDetail(id);
            }, properties);

            // we might need to display a specific session
            if(selectedSessionIdx != -1){
                BooksysPie.selectSector(selectedSessionIdx);
                updateDetail(selectedSessionIdx);
            }else{
                // show the session details but without a selected session
                BooksysViewSessionDetails.display("session_details", null, null, function(){});
            }
        }
    });	
}

function continueToWatch(){
    window.location.href = '/watch.html?sessionId='+selectedSession.id;
}