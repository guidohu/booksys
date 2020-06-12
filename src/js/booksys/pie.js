class BooksysPie {

    // draws the pie chart
    // - location       where to draw the pie chart
    // - data           the data we got from the server (get_booking_day)
    // - callback       function that gets called upon selection of a pie element
    //
    // Returns an object containing all the data of the pie.
    static drawPie(location, data, callback, properties){
        // reset current content
        $('#'+location).html("");

        // values / labels / colors for the pie
        var sessions = data.sessions;
        var values = [];
        var labels = [];
        var colors = [];
        var pieSessions = [];

        // specific times required for proper display
        var sunrise = moment.utc(data.sunrise, "X");
        var sunset  = moment.utc(data.sunset, "X");
        var business_day_start = moment(Number(data.window_start), "X").tz(getTimeZone());
        var business_day_end   = moment(Number(data.window_end), "X").tz(getTimeZone());
        // console.log(business_day_start.format() + " " + business_day_end.format());
        business_day_start.set('hour', data.business_day_start.substring(0,2));
        business_day_start.set('minute', data.business_day_start.substring(3,5));
        business_day_start.set('second', 0);
        business_day_end.set('hour', data.business_day_end.substring(0,2));
        business_day_end.set('minute', data.business_day_end.substring(3,5));
        business_day_end.set('second', 0);
        // get the start/end of the day (e.g. either sunrise or first allowed time)
        // console.log("sunrise: " + sunrise.format());
        // console.log("sunset: " + sunset.format());
        // console.log("business_day_start: " + business_day_start.format());
        // console.log("business_day_end: " + business_day_end.format());
        var dayStart = moment.utc(Math.max(sunrise.format("X"), business_day_start.format("X")), "X");
        var dayEnd   = moment.utc(Math.min(sunset.format("X"), business_day_end.format("X")), "X");

        // reminder variables to fill the gaps between sessions
        var lastEnd = dayStart;
        
        // add all sessions to the pie (including gaps)
        for(var i=0; i<sessions.length; i++){
            // add a gap-session before this session if needed
            if(lastEnd.format("X") < sessions[i].start){
                // console.log("Add gap session before first session");
                // console.log(sessions[i].start - lastEnd.format("X"));
                values.push(sessions[i].start - lastEnd.format("X"));
                labels.push(lastEnd.tz(getTimeZone()).format("HH:mm") 
                    + " - " 
                    + moment(sessions[i].start, "X").tz(getTimeZone()).format("HH:mm"));
                colors.push("#333");
                pieSessions.push({
                    id:    null,
                    start: lastEnd,
                    end:   moment(sessions[i].start, "X")
                });
            }

            // calculate duration offsets for sessions longer than a day
            // as we only display a single day
            var duration_offset = 0;
            if(sessions[i].start < data.window_start){
                duration_offset = data.window_start - sessions[i].start;
            }
            if(sessions[i].end > data.window_end){
                duration_offset = sessions[i].end - data.window_end;
            }

            // get color for this slot
            var color = "red";
            if(sessions[i].free > 0){
                color = "green";
            } else {
                color = "red";
            }


            // console.log("Add regular session");
            // console.log(sessions[i].duration - duration_offset);
            values.push(sessions[i].duration - duration_offset);
            labels.push(moment.utc(sessions[i].start, "X").tz(getTimeZone()).format("HH:mm")
                + " - " 
                + moment.utc(sessions[i].end, "X").tz(getTimeZone()).format("HH:mm"));
            colors.push(color);

            pieSessions.push({
                id:     sessions[i].id,
                start:  moment.utc(sessions[i].start, "X"),
                end:    moment.utc(sessions[i].end, "X"),
                title:  sessions[i].title,
                comment: sessions[i].comment,
                free:   sessions[i].free,
                riders: sessions[i].riders
            });
            lastEnd = moment.utc(sessions[i].end, "X");
        }

        // in case there are no sessions yet
        if(pieSessions.length == 0){
            // console.log("Add empty day session");
            // console.log(dayEnd.diff(dayStart, 'seconds'));
            // console.log(dayStart.format());
            // console.log(dayEnd.format());
            values.push(dayEnd.diff(dayStart, 'seconds'));
            labels.push(dayStart.tz(getTimeZone()).format("HH:mm")
                + " - "
                + dayEnd.tz(getTimeZone()).format("HH:mm"));
            colors.push("#333");

            pieSessions.push({
                id:     null,
                start:  dayStart,
                end:    dayEnd
            });
        }

        // add first pie(s) of the day (if needed)
        if(pieSessions[0].start.format("X") > data.window_start){
            var session = pieSessions[0];
            var duration = session.start.format("X") - data.window_start;

            // console.log("Add pre day session");
            // console.log(duration);
            duration = duration / 6;
            values.unshift(duration);
            labels.unshift(
                moment(data.window_start, "X").tz(getTimeZone()).format("HH:mm")
                + " - "
                + session.start.tz(getTimeZone()).format("HH:mm")
            );
            colors.unshift("#111");

            pieSessions.unshift({
                id:     null,
                start:  moment(data.window_start, "X"),
                end:    pieSessions[0].start
            });
        }

        // add gap filler session until end of day
        if(pieSessions[pieSessions.length-1].end.format("X") < dayEnd.format("X")){
            session = pieSessions[pieSessions.length-1]
            var duration = dayEnd.format("X") - session.end.format("X");
            // console.log("Add end of day gap filler");
            // console.log(duration);
            values.push(duration);
            labels.push(
                session.end.tz(getTimeZone()).format("HH:mm") 
                + " - " 
                + dayEnd.tz(getTimeZone()).format("HH:mm"));
            colors.push("#333");

            pieSessions.push({
                id:     null,
                start:  session.end,
                end:    dayEnd
            });
        }

        // add last pie(s) of the day (if needed)
        if(pieSessions[pieSessions.length-1].end.format("X") < data.window_end){
            session = pieSessions[pieSessions.length-1];
            var duration = data.window_end - session.end.format("X");
            // console.log("Add after-day");
            // console.log(duration);
            duration = duration / 6;
            values.push(duration);
            labels.push(
                session.end.tz(getTimeZone()).format("HH:mm") 
                + " - " 
                + moment(data.window_end, "X").tz(getTimeZone()).format("HH:mm"));
            colors.push("#111");

            pieSessions.push({
                id:     null,
                start:  session.end,
                end:    moment(data.window_end, "X")
            });
        }
        
        // accept sizing properties
        var containerHeight;
        var containerWidth;
        var circleX;
        var circleY;
        var circleRadius;
        var animate;

        if(properties != null && properties.containerHeight != null){
            containerHeight = properties.containerHeight;
        }else{
            containerHeight = Math.min(Math.min($('#'+location).height(), $(window).width())*0.75,400);
        }
        if(properties != null && properties.containerWidth != null){
            containerWidth = properties.containerWidth;
        }else{
            containerWidth = containerHeight;
        }
        if(properties != null && properties.circleX != null){
            circleX = properties.circleX;
        }else{
            circleX = containerWidth / 2;
        }
        if(properties != null && properties.circleY != null){
            circleY = properties.circleY;
        }else{
            circleY = containerHeight / 2;
        }
        if(properties != null && properties.circleRadius != null){
            circleRadius = properties.circleRadius;
        }else{
            circleRadius = Math.min(containerWidth, containerHeight) / 2 * 0.8;
        }
        if(properties != null && properties.animate != null){
            animate = properties.animate;
        }else{
            animate = true;
        }
        if(properties != null && properties.labels != null){
            if(properties.labels == false){
                labels = [];
            }
        }
        
        
        // create the pie-chart
        // console.log(values);
        BooksysPie.pie = Raphael(location, containerWidth, containerHeight).pieChart(  
            circleX, 
            circleY, 
            circleRadius, 
            values,
            labels,
            colors, 
            "#fff", 
            animate, 
            callback
        );

        return pieSessions;
    }

    // selects a specific sector of the pie
    static selectSector(id){
        BooksysPie.pie.selectSector(id);
    }

    static resetSelection(){
        BooksysPie.pie.reset();
    }
}

BooksysPie.pie = null;