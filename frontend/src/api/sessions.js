import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";
import * as dayjsAdvancedFormat from "dayjs/plugin/advancedFormat";
import Session from "@/dataTypes/session";
import Request from "@/api/common/request";

export default class Sessions {
  /**
   * Returns 6 weeks of sessions
   * @param {ISO date} month a day within the month
   */
  static getSessionsCalendar(month) {
    dayjs.extend(dayjsUTC);
    dayjs.extend(dayjsTimezone);
    dayjs.extend(dayjsAdvancedFormat);

    let dateIterator = dayjs(month).startOf("month");

    // start with a monday
    while (dateIterator.day() != 1) {
      dateIterator = dateIterator.subtract(1, "day");
    }

    // build 42 time windows (7 days for 6 weeks)
    let query = {
      time_windows: [],
    };

    for (let i = 0; i < 42; i++) {
      // create new time window
      const start = parseInt(dateIterator.startOf("day").format("X"));
      const end = parseInt(dateIterator.endOf("day").format("X"));
      const window = {
        start: start,
        end: end,
      };
      query.time_windows.push(window);
      dateIterator = dateIterator.add(1, "day");
    }

    console.log("getSessionsCalendar:", query);
    return new Promise((resolve, reject) => {
      Request.postRequest("/api/v2/booking/series/list", query)
      .then((days) => {
        let monthSessions = [];
        for (let i = 0; i < days.length; i++) {
          let res = days[i];
          let timezone = res.timezone;
          const winStart = res.window_start;
          // const winEnd     = res.window_end;
          res.window_start = dayjs
            .unix(res.window_start)
            .tz(timezone)
            .format();
          res.window_end = dayjs
            .unix(res.window_end)
            .tz(timezone)
            .format();
          res.sunrise = dayjs.unix(res.sunrise).tz(timezone).format();
          res.sunset = dayjs.unix(res.sunset).tz(timezone).format();
          const busDayStart = res.business_day_start.split(":");
          const busDayEnd = res.business_day_end.split(":");
          res.business_day_start = dayjs
            .unix(winStart)
            .tz(timezone)
            .set("hour", busDayStart[0])
            .set("minutes", busDayStart[1])
            .set("seconds", busDayStart[2])
            .format();
          res.business_day_end = dayjs
            .unix(winStart)
            .tz(timezone)
            .set("hour", busDayEnd[0])
            .set("minutes", busDayEnd[1])
            .set("seconds", busDayEnd[2])
            .format();

          // convert session times to ISO time
          for (let i = 0; i < res.sessions.length; i++) {
            const s = res.sessions[i];
            const session = new Session(
              s.id,
              s.title,
              s.comment,
              dayjs.unix(s.start).tz(timezone).format(),
              dayjs.unix(s.end).tz(timezone).format(),
              s.free,
              s.type,
              s.creator_id,
              s.creator_first_name,
              s.creator_last_name
            );
            session.addRiders(s.riders);
            res.sessions[i] = session;
          }
          monthSessions.push(res);
        }
        resolve(monthSessions);
      })
      .catch((errors) => {
        reject(errors);
      })
    });
  }

  static getSessions(dateStart, dateEnd) {
    dayjs.extend(dayjsUTC);
    dayjs.extend(dayjsTimezone);
    dayjs.extend(dayjsAdvancedFormat);

    const query = {
      start: parseInt(dayjs(dateStart).format("X")),
      end: parseInt(dayjs(dateEnd).format("X")),
    };

    return new Promise((resolve, reject) => {
      Request.postRequest("/api/v2/booking/day/list", query)
      .then((res) => {
        let timezone = res.timezone;
        const winStart = res.window_start;
        // const winEnd     = res.window_end;
        res.window_start = dayjs
          .unix(res.window_start)
          .tz(timezone)
          .format();
        res.window_end = dayjs
          .unix(res.window_end)
          .tz(timezone)
          .format();
        res.sunrise = dayjs.unix(res.sunrise).tz(timezone).format();
        res.sunset = dayjs.unix(res.sunset).tz(timezone).format();
        const busDayStart = res.business_day_start.split(":");
        const busDayEnd = res.business_day_end.split(":");
        res.business_day_start = dayjs
          .unix(winStart)
          .tz(timezone)
          .set("hour", busDayStart[0])
          .set("minutes", busDayStart[1])
          .set("seconds", busDayStart[2])
          .format();
        res.business_day_end = dayjs
          .unix(winStart)
          .tz(timezone)
          .set("hour", busDayEnd[0])
          .set("minutes", busDayEnd[1])
          .set("seconds", busDayEnd[2])
          .format();

        // convert session times to ISO time
        for (let i = 0; i < res.sessions.length; i++) {
          const s = res.sessions[i];
          const session = new Session(
            s.id,
            s.title,
            s.comment,
            dayjs.unix(s.start).tz(timezone).format(),
            dayjs.unix(s.end).tz(timezone).format(),
            s.free,
            s.type,
            s.creator_id,
            s.creator_first_name,
            s.creator_last_name
          );
          session.addRiders(s.riders);
          res.sessions[i] = session;
        }

        resolve(res);
      })
      .catch((errors) => reject(errors));
    });
  }

  static getSession(sessionId) {
    const requestData = {
      id: sessionId,
    };
    return new Promise((resolve, reject) => {
      Request.postRequest("/api/v2/session/get", requestData)
        .then((sR) => {
          const session = new Session(
            sR.id,
            sR.title,
            sR.comment,
            dayjs.unix(sR.start_time).format(),
            dayjs.unix(sR.end_time).format(),
            sR.riders_max,
            sR.type
          );
          session.addRiders(sR.riders);
          resolve({
            session: session,
          });
        })
        .catch((error) => {
          reject(error);
        });
    });
  }

  static getSessionMetadata(sessionId) {
    console.log("sessions/getSessionMetadata called, with", sessionId);
    const requestData = {
      id: sessionId,
    };
    return new Promise((resolve, reject) => {
      Request.postRequest("/api/v2/session/metadata/get", requestData)
        .then((resp) => {
          resolve({
            metaInfo: resp,
          });
        })
        .catch((error) => {
          reject(error);
        });
    });
  }

  static createSession(sessionData) {
    dayjs.extend(dayjsAdvancedFormat);

    // build request body
    const session = {
      title: sessionData.title,
      comment: sessionData.description,
      max_riders: sessionData.maximumRiders,
      start: parseInt(dayjs(sessionData.start).format("X")),
      end: parseInt(dayjs(sessionData.end).format("X")),
      type: sessionData.type,
    };

    return Request.postRequest("/api/v2/session/create", session);
  }

  static editSession(sessionData) {
    dayjs.extend(dayjsAdvancedFormat);

    // build request body
    const session = {
      id: sessionData.id,
      title: sessionData.title,
      comment: sessionData.description,
      max_riders: parseInt(sessionData.maximumRiders),
      start: parseInt(dayjs(sessionData.start).format("X")),
      end: parseInt(dayjs(sessionData.end).format("X")),
      type: sessionData.type,
    };

    return Request.postRequest("/api/v2/session/edit", session);
  }

  static deleteSession(sessionData) {
    // build request body
    const session = {
      session_id: sessionData.id,
    };

    return Request.postRequest("/api/v2/session/delete", session);
  }

  static addUsersToSession(sessionId, users) {
    // build request body
    const requestBody = {
      user_ids: users.map((u) => u.id),
      session_id: sessionId,
    };

    return Request.postRequest("/api/v2/session/user/add", requestBody);
  }

  /**
   * Deletes a user from a specific session.
   * @param {*} sessionId ID of the session
   * @param {*} user ID of the user
   * @returns
   */
  static deleteUserFromSession(sessionId, user) {
    // build request body
    const requestBody = {
      user_id: user.id,
      session_id: sessionId,
    };

    return Request.postRequest("/api/v2/session/user/remove", requestBody);
  }
}
