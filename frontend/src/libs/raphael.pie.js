/*
	Create to create a pie-chart. This function is based on the example on the website
	of the Raphael library.
*/

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
    var angle = 90,
        total = 0,
        start = 0,
        process = function (j) {
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
                });
				
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
            start += .1;
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