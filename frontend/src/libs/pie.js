import * as dayjs from 'dayjs';
import * as dayjsUTC from 'dayjs/plugin/utc';
import * as dayjsTimezone from 'dayjs/plugin/timezone';
import * as dayjsAdvancedFormat from 'dayjs/plugin/advancedFormat';

const Raphael = require("raphael");

export default class BooksysPie {

  static addPiePlugin = () => {
    Raphael.fn.pieChart = function (cx, cy, r, values, labels, colors, stroke, animate, clickCallBack) {
      var paper = this,
        rad = Math.PI / 180,
        chart = this.set();

      var sectors = new Object();
      var sector_active = 0;
      var sectors_active = new Object();

      // Allows to select a certain sector by ID
      this.selectSector = function (i) {
        resetAllSectors();
        sectors[i].stop().animate({ transform: "s1.15 1.15 " + cx + " " + cy }, 0, "elastic");
        sector_active = 1;
        sectors_active[i] = 1;
      }

      // Resets all sectors
      this.reset = function () {
        resetAllSectors();
      }

      // Creates a single sector
      // center (cx, cy), radius, startAngle, endAngle, design-parameters
      function sector(cx, cy, r, startAngle, endAngle, params) {
        var x1 = cx + r * Math.cos(-startAngle * rad),
          x2 = cx + r * Math.cos(-endAngle * rad),
          y1 = cy + r * Math.sin(-startAngle * rad),
          y2 = cy + r * Math.sin(-endAngle * rad);

        // we have to round values because the sector does not get drawn if start and endpoint are identical
        x1 = Math.round(x1 * 10000) / 10000;
        x2 = Math.round(x2 * 10000) / 10000;
        y1 = Math.round(y1 * 10000) / 10000;
        y2 = Math.round(y2 * 10000) / 10000;

        if (x1 == x2 && y1 == y2) {
          x2 = x2 - 0.001;
        }

        return paper.path(["M", cx, cy, "L", x1, y1, "A", r, r, 0, +(endAngle - startAngle < -180), 1, x2, y2, "z"]).attr(params);
      }

      function resetAllSectors() {
        for (var i = 0, ii = values.length; i < ii; i++) {
          sectors[i].stop().animate({ transform: "" }, 500, "elastic");
          sectors_active[i] = 0;
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
        var p = sector(cx, cy, r, angle, angle - angleplus, { fill: colors[j], stroke: stroke, "stroke-width": 1 });
        sectors[j] = p;
        sectors_active[j] = 0;

        // Check if labels are defined, if yes we create the labels
        if (labels[j]) {
          var txt = paper.text(cx + (r + delta + 55) * Math.cos(-popangle * rad), cy + (r + delta + 25) * Math.sin(-popangle * rad),
            labels[j]).attr({ fill: colors[j], stroke: "none", opacity: 0, "font-size": 18 });
        }

        // If the pie has some functionality/animation we add events
        if (animate) {
          p.mouseover(function () {
            p.stop().animate({ transform: "s1.15 1.15 " + cx + " " + cy }, ms, "elastic");
            if (labels[j]) {
              txt.stop().animate({ opacity: 1 }, ms, "elastic");
            }
          }).mouseout(function () {
            if (sectors_active[j] == 0) {
              p.stop().animate({ transform: "" }, ms, "elastic");
            }
            if (labels[j]) {
              txt.stop().animate({ opacity: 0 }, ms);
            }
          })

          if (clickCallBack) {
            p.click(function () {
              if (sector_active == 0) {
                sector_active = 1;
                sectors_active[j] = 1;
              } else {
                resetAllSectors();
                sectors_active[j] = 1;
                sector_active = 1;
              }
              clickCallBack(j);
              p.stop().animate({ transform: "s1.15 1.15 " + cx + " " + cy }, ms, "elastic");
            });
          }
        }
        angle -= angleplus;
        chart.push(p);
      };

      // get the total hours that are available
      for (var i = 0, ii = values.length; i < ii; i++) {
        total += values[i];
      }
      // process every single sector
      for (i = 0; i < ii; i++) {
        process(i);
      }
      return this;
    };
  }

  // draws the pie chart
  // - location       where to draw the pie chart
  // - data           the data we got from the server (get_booking_day)
  // - callback       function that gets called upon selection of a pie element
  //
  // Returns an object containing all the data of the pie.
  drawPie(location, data, callback, properties) {
    BooksysPie.addPiePlugin();

    // get timezone
    dayjs.extend(dayjsUTC);
    dayjs.extend(dayjsTimezone);
    dayjs.extend(dayjsAdvancedFormat);
    let timezone = 'UTC';
    if (properties.timezone != null) {
      timezone = properties.timezone;
    }

    // define colors (hardcode for now)
    const colorNoSlot = "#424242";
    const colorCourse = "#d9534e"; //"#3AAFA9"; // "#FC4445";
    const colorSlot = "#5cb85b";
    const colorSlotFull = "#d9534e";
    const colorOffHour = "#212121";

    // reset current content
    location.html = "";

    // values / labels / colors for the pie
    var sessions = data.sessions;
    var values = [];
    var labels = [];
    var colors = [];
    var pieSessions = [];

    // specific times required for proper display
    const sunrise = dayjs(data.sunrise);
    const sunset = dayjs(data.sunset);
    const business_day_start = dayjs(data.business_day_start);
    const business_day_end = dayjs(data.business_day_end);
    // get the start/end of the day (e.g. either sunrise or first allowed time)
    const dayStart = dayjs.unix(Math.max(sunrise.format("X"), business_day_start.format("X"))).format();
    const dayEnd = dayjs.unix(Math.min(sunset.format("X"), business_day_end.format("X"))).format();

    // reminder variables to fill the gaps between sessions
    let lastEnd = dayStart;

    // add all sessions to the pie (including gaps)
    for (var i = 0; i < sessions.length; i++) {
      // add a gap-session before this session if needed
      if (dayjs(lastEnd).format("X") < dayjs(sessions[i].start).format("X")) {
        const session = sessions[i];
        const duration = dayjs(session.start).format("X") - 1 - dayjs(lastEnd).format("X");
        let label = "";
        let start = "";
        let end = "";
        if (lastEnd == dayStart) {
          start = dayjs(lastEnd).format();
          end = dayjs(session.start).add(-1, "minute").format();
          label = dayjs(start).format("HH:mm") + " - " + dayjs(end).format("HH:mm");
        } else {
          start = dayjs(lastEnd).add(1, 'minute').format();
          end = dayjs(session.start).add(-1, "minute").format();
          label = dayjs(start).format("HH:mm") + " - " + dayjs(end).format("HH:mm");
        }
        values.push(duration);
        labels.push(label);
        colors.push(colorNoSlot);
        pieSessions.push({
          id: null,
          start: start,
          end: end
        });
      }

      // calculate duration offsets for sessions longer than a day
      // as we only display a single day
      var duration_offset = 0;
      if (dayjs(sessions[i].start).format("X") < dayjs(data.window_start).format("X")) {
        duration_offset = dayjs(data.window_start).format("X") - dayjs(sessions[i].start).format("X");
      }
      if (dayjs(sessions[i].end).format("X") > dayjs(data.window_end).format("X")) {
        duration_offset = dayjs(sessions[i].end).format("X") - dayjs(data.window_end).format("X");
      }

      const sessionDurationSeconds = dayjs(sessions[i].end).format("X") - dayjs(sessions[i].start).format("X");

      // get color for this slot
      var color = "";
      // check if free and is not course
      if (sessions[i].maximumRiders > 0 && sessions[i].type != 1) {
        color = colorSlot;
      } else if (sessions[i].type == 1) {
        color = colorCourse;
      } else {
        color = colorSlotFull;
      }

      values.push(sessionDurationSeconds - duration_offset);
      labels.push(dayjs(sessions[i].start).format("HH:mm")
        + " - "
        + dayjs(sessions[i].end).format("HH:mm"));
      colors.push(color);

      pieSessions.push({
        id: sessions[i].id,
        start: sessions[i].start,
        end: sessions[i].end,
        title: sessions[i].title,
        comment: sessions[i].comment,
        free: sessions[i].free,
        riders: sessions[i].riders
      });
      lastEnd = dayjs(sessions[i].end).format();
    }

    // in case there are no sessions yet
    if (pieSessions.length == 0) {
      values.push(
        dayjs(dayEnd).diff(dayjs(dayStart), 'seconds')
      );
      labels.push(
        dayjs(dayStart).tz(timezone).format("HH:mm")
        + " - "
        + dayjs(dayEnd).tz(timezone).format("HH:mm")
      );
      colors.push(colorNoSlot);

      pieSessions.push({
        id: null,
        start: dayStart,
        end: dayEnd
      });
    }

    // add first pie(s) of the day (if needed)
    if (dayjs(pieSessions[0].start).format("X") > dayjs(data.window_start).format("X")) {
      let session = pieSessions[0];
      let duration = dayjs(session.start).format("X") - 1 - dayjs(data.window_start).format("X");

      duration = duration / 6;
      values.unshift(duration);
      labels.unshift(
        dayjs(data.window_start).format("HH:mm")
        + " - "
        + dayjs(session.start).add(-1, 'minute').format("HH:mm")
      );
      colors.unshift(colorOffHour);

      pieSessions.unshift({
        id: null,
        start: data.window_start,
        end: dayjs(session.start).add(-1, 'minute').format()
      });
    }

    // add gap filler session until end of day
    if (dayjs(pieSessions[pieSessions.length - 1].end).format("X") < dayjs(dayEnd).format("X")) {
      let session = pieSessions[pieSessions.length - 1]
      var duration = dayjs(dayEnd).format("X") - dayjs(session.end).format("X");
      values.push(duration);
      labels.push(
        dayjs(session.end).add(1, 'minute').format("HH:mm")
        + " - "
        + dayjs(dayEnd).format("HH:mm"));
      colors.push(colorNoSlot);

      pieSessions.push({
        id: null,
        start: dayjs(session.end).add(1, 'minute').format(),
        end: dayEnd
      });
    }

    // add last pie(s) of the day (if needed)
    if (dayjs(pieSessions[pieSessions.length - 1].end).format("X") < dayjs(data.window_end).format("X")) {
      let session = pieSessions[pieSessions.length - 1];
      let duration = dayjs(data.window_end).format("X") - dayjs(session.end).add(1, "minute").format("X");
      duration = duration / 6;
      values.push(duration);
      labels.push(
        dayjs(session.end).add(1, "minute").format("HH:mm")
        + " - "
        + dayjs(data.window_end).format("HH:mm"));
      colors.push(colorOffHour);

      pieSessions.push({
        id: null,
        start: dayjs(session.end).add(1, "minute").format(),
        end: data.window_end
      });
    }

    // accept sizing properties
    let containerHeight;
    let containerWidth;
    let circleX;
    let circleY;
    let circleRadius;
    let animate;

    if (properties != null && properties.containerHeight != null) {
      containerHeight = properties.containerHeight;
    } else {
      containerHeight = Math.min(Math.min(location.offsetHeight, window.innerWidth) * 0.75, 400);
    }
    if (properties != null && properties.containerWidth != null) {
      containerWidth = properties.containerWidth;
    } else {
      containerWidth = containerHeight;
    }
    if (properties != null && properties.circleX != null) {
      circleX = properties.circleX;
    } else {
      circleX = containerWidth / 2;
    }
    if (properties != null && properties.circleY != null) {
      circleY = properties.circleY;
    } else {
      circleY = containerHeight / 2;
    }
    if (properties != null && properties.circleRadius != null) {
      circleRadius = properties.circleRadius;
    } else {
      circleRadius = Math.min(containerWidth, containerHeight) / 2 * 0.8;
    }
    if (properties != null && properties.animate != null) {
      animate = properties.animate;
    } else {
      animate = true;
    }
    if (properties != null && properties.labels != null) {
      if (properties.labels == false) {
        labels = [];
      }
    }

    // create the pie-chart
    this.pie = Raphael(location, containerWidth, containerHeight).pieChart(
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
  selectSector(id) {
    this.pie.selectSector(id);
  }

  resetSelection() {
    this.pie.reset();
  }
}

BooksysPie.pie = null;