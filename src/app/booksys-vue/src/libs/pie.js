// import {Raphael} from "raphael"
//  from "./raphael.pie"
// import moment from "moment"
// import from "timezone"
import moment from 'moment-timezone/moment-timezone';

const Raphael = require("raphael");

export default class BooksysPie {

    static getTimeZone = () => {
        return 'UTC'
    }

    static addPiePlugin = () => {
        Raphael.fn.pieChart = function (cx, cy, r, values, labels, colors, stroke, animate, clickCallBack) {
            var paper = this,
                rad = Math.PI / 180,
                chart = this.set();
            
            var sectors = new Object();
            var sector_active = 0;
            var sectors_active = new Object();
            
            // Allows to select a certain sector by ID
            this.selectSector = function(i){
                resetAllSectors();
                sectors[i].stop().animate({transform: "s1.15 1.15 " + cx + " " + cy}, 0, "elastic");
                sector_active = 1;
                sectors_active[i] = 1;
            }
            
            // Resets all sectors
            this.reset = function(){
                resetAllSectors();
            }
                
            // Creates a single sector
            // center (cx, cy), radius, startAngle, endAngle, design-parameters
            function sector(cx, cy, r, startAngle, endAngle, params) {
                var x1 = cx + r * Math.cos(-startAngle * rad),
                    x2 = cx + r * Math.cos(-endAngle * rad),
                    y1 = cy + r * Math.sin(-startAngle * rad),
                    y2 = cy + r * Math.sin(-endAngle * rad);
                // draw a patch starting at (cx, cy), line to (x1, y1), ellipse with (rx ry x-axis-rotation large-arc-flag sweep-flag (x2, y2))+
                // and close the path
                //alert("End-Angle: " + endAngle + "\n" +
                //	  "startAngle: " + startAngle + "\n" +
                //	  "diff: " + (endAngle-startAngle) + "\n" + 
                //	  "Raw: 1: (" + x1 + ", " + y1 + ")\t2: ( " + x2 + ", " + y2 + ")");
                
                // we have to round values because the sector does not get drawn if start and endpoint are identical
                x1 = Math.round(x1*10000)/10000;
                x2 = Math.round(x2*10000)/10000;
                y1 = Math.round(y1*10000)/10000;
                y2 = Math.round(y2*10000)/10000;
                
                //alert("Rounded: 1: (" + x1 + ", " + y1 + ")\t2: ( " + x2 + ", " + y2 + ")");
                
                if(x1==x2 && y1==y2){
                    x2=x2-0.001;	
                    //alert("Corrected: 1: (" + x1 + ", " + y1 + ")\t2: ( " + x2 + ", " + y2 + ")");
                }
                
                //alert(+(endAngle - startAngle > 180));
                //return paper.path(["M", 50, 50, "L", 50, 10, "A", r, r, 0, 1, 1, 49.999, 10, "z"]).attr(params);
                //return paper.path(["M", cx, cy, "L", x1, y1, "A", r, r, 0, 0, 0, x2, y2, "z"]).attr(params);
                return paper.path(["M", cx, cy, "L", x1, y1, "A", r, r, 0, +(endAngle - startAngle < -180), 1, x2, y2, "z"]).attr(params);
                //return paper.path(["M", cx, cy, "L", x1, y1, "A", r, r, 0, 50, 1, x2, y2, "z"]).attr(params);
            }
            
            function resetAllSectors(){
                for (var i = 0, ii = values.length; i < ii; i++) {
                    sectors[i].stop().animate({transform: ""}, 500, "elastic");
                    sectors_active[i]=0;
                }
            }
            
            // draws the sector and adds the gestures and the text
            let angle = 90
            let total = 0
            // let start = 0
            let process = function (j) {
                    var value = values[j],
                        angleplus = 360 * value / total,
                        popangle = angle - (angleplus / 2),
                        ms = 700,
                        delta = 25;
                        // center (cx, cy), radius, startAngle, endAngle, design-parameters
                    //alert("Value: " + value + "\n" +
                    //	  "Angleplus: " + angleplus + "\n" +
                    //	  "cx, cy, r, angle, angle-angleplus:" + cx + ", " + cy + ", " + r + ", " + angle + ", " + (angle-angleplus));
                    var p = sector(cx, cy, r, angle, angle - angleplus, {fill: colors[j], stroke: stroke, "stroke-width": 1});
                    sectors[j]=p;
                    sectors_active[j]=0;
        
                    // Check if labels are defined, if yes we create the labels
                    if(labels[j]){
                        var txt = paper.text(cx + (r + delta + 55) * Math.cos(-popangle * rad), cy + (r + delta + 25) * Math.sin(-popangle * rad), 
                                  labels[j]).attr({fill: colors[j], stroke: "none", opacity: 0, "font-size": 18});
                    }
                    
                    // If the pie has some functionality/animation we add events
                    if(animate){
                        p.mouseover(function () {
                            p.stop().animate({transform: "s1.15 1.15 " + cx + " " + cy}, ms, "elastic");
                            if(labels[j]){
                                txt.stop().animate({opacity: 1}, ms, "elastic");
                            }
                        }).mouseout(function () {
                            if(sectors_active[j] == 0){
                                p.stop().animate({transform: ""}, ms, "elastic");
                            }
                            if(labels[j]){
                                txt.stop().animate({opacity: 0}, ms);
                            }
                        })
                        
                        if(clickCallBack){
                            p.click(function () {
                                if(sector_active == 0 ){
                                    sector_active = 1;
                                    sectors_active[j] = 1;
                                }else{
                                    resetAllSectors();
                                    sectors_active[j] = 1;
                                    sector_active = 1;
                                }
                                clickCallBack(j);
                                p.stop().animate({transform: "s1.15 1.15 " + cx + " " + cy}, ms, "elastic");
                            });
                        }
                    }
                    angle -= angleplus;
                    chart.push(p);
                    //chart.push(txt);
                    //start += .1;
                };
            
            // get the total hours that are available
            for (var i = 0, ii = values.length; i < ii; i++) {
                total += values[i];
            }
            // process every single sector
            for (i = 0; i < ii; i++) {
                process(i);
            }
            //return chart;
            return this;
        };
    }
    // draws the pie chart
    // - location       where to draw the pie chart
    // - data           the data we got from the server (get_booking_day)
    // - callback       function that gets called upon selection of a pie element
    //
    // Returns an object containing all the data of the pie.
    static drawPie(location, data, callback, properties){
        BooksysPie.addPiePlugin()

        // define colors (hardcode for now)
        let colorNoSlot = "#424242";
        let colorCourse = "#d9534e"; //"#3AAFA9"; // "#FC4445";
        let colorSlot   = "#5cb85b";
        let colorSlotFull = "#d9534e";
        let colorOffHour = "#212121";

        // reset current content
        location.html = "";

        // values / labels / colors for the pie
        var sessions = data.sessions;
        var values = [];
        var labels = [];
        var colors = [];
        var pieSessions = [];

        // specific times required for proper display
        var sunrise = moment.utc(data.sunrise, "X");
        var sunset  = moment.utc(data.sunset, "X");
        var business_day_start = moment(Number(data.window_start), "X").tz(this.getTimeZone());
        var business_day_end   = moment(Number(data.window_end), "X").tz(this.getTimeZone());
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
                labels.push(lastEnd.tz(this.getTimeZone()).format("HH:mm") 
                    + " - " 
                    + moment(sessions[i].start, "X").tz(this.getTimeZone()).format("HH:mm"));
                colors.push(colorNoSlot);
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
            var color = "";
            console.log(sessions[i]);
            // check if free and is not course
            if(sessions[i].free > 0 && sessions[i].type != 1){
                color = colorSlot;
            } else if(sessions[i].type == 1) {
                color = colorCourse;
            } else {
                color = colorSlotFull;
            }


            // console.log("Add regular session");
            // console.log(sessions[i].duration - duration_offset);
            values.push(sessions[i].duration - duration_offset);
            labels.push(moment.utc(sessions[i].start, "X").tz(this.getTimeZone()).format("HH:mm")
                + " - " 
                + moment.utc(sessions[i].end, "X").tz(this.getTimeZone()).format("HH:mm"));
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
            labels.push(dayStart.tz(this.getTimeZone()).format("HH:mm")
                + " - "
                + dayEnd.tz(this.getTimeZone()).format("HH:mm"));
            colors.push(colorNoSlot);

            pieSessions.push({
                id:     null,
                start:  dayStart,
                end:    dayEnd
            });
        }

        // add first pie(s) of the day (if needed)
        if(pieSessions[0].start.format("X") > data.window_start){
            let session = pieSessions[0];
            let duration = session.start.format("X") - data.window_start;

            // console.log("Add pre day session");
            // console.log(duration);
            duration = duration / 6;
            values.unshift(duration);
            labels.unshift(
                moment(data.window_start, "X").tz(this.getTimeZone()).format("HH:mm")
                + " - "
                + session.start.tz(this.getTimeZone()).format("HH:mm")
            );
            colors.unshift(colorOffHour);

            pieSessions.unshift({
                id:     null,
                start:  moment(data.window_start, "X"),
                end:    pieSessions[0].start
            });
        }

        // add gap filler session until end of day
        if(pieSessions[pieSessions.length-1].end.format("X") < dayEnd.format("X")){
            let session = pieSessions[pieSessions.length-1]
            var duration = dayEnd.format("X") - session.end.format("X");
            // console.log("Add end of day gap filler");
            // console.log(duration);
            values.push(duration);
            labels.push(
                session.end.tz(this.getTimeZone()).format("HH:mm") 
                + " - " 
                + dayEnd.tz(this.getTimeZone()).format("HH:mm"));
            colors.push(colorNoSlot);

            pieSessions.push({
                id:     null,
                start:  session.end,
                end:    dayEnd
            });
        }

        // add last pie(s) of the day (if needed)
        if(pieSessions[pieSessions.length-1].end.format("X") < data.window_end){
            let session = pieSessions[pieSessions.length-1];
            let duration = data.window_end - session.end.format("X");
            // console.log("Add after-day");
            // console.log(duration);
            duration = duration / 6;
            values.push(duration);
            labels.push(
                session.end.tz(this.getTimeZone()).format("HH:mm") 
                + " - " 
                + moment(data.window_end, "X").tz(this.getTimeZone()).format("HH:mm"));
            colors.push(colorOffHour);

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
            containerHeight = Math.min(Math.min(location.offsetHeight, window.innerWidth)*0.75,400);
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