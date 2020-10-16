import moment from 'moment';
import Session from '@/dataTypes/session';

export default class Sessions {

  static getSessions(dateStart, dateEnd){
    console.log("getSessions:", dateStart);
    return new Promise((resolve, reject) => {
      const query = {
        start: moment(dateStart).format("X"),
        end: moment(dateEnd).format("X")
      }
      fetch('/api/booking.php?action=get_booking_day', {
        method: 'POST',
        cache: 'no-cache',
        body: JSON.stringify(query)
      })
      .then(response => {
        response.json()
        .then(data => {
          console.log("getSessions response data:", data);
          if(data.ok){
            let res = data.data;
            let timezone = res.timezone;
            const winStart   = res.window_start;
            // const winEnd     = res.window_end;
            res.window_start = moment(res.window_start, "X").tz(timezone).format();
            res.window_end   = moment(res.window_end, "X").tz(timezone).format();
            res.sunrise      = moment(res.sunrise, "X").tz(timezone).format();
            res.sunset       = moment(res.sunset, "X").tz(timezone).format();
            const busDayStart = res.business_day_start.split(":");
            const busDayEnd   = res.business_day_end.split(":");
            res.business_day_start = moment(winStart, "X")
              .tz(timezone)
              .set('hour', busDayStart[0])
              .set('minutes', busDayStart[1])
              .set('seconds', busDayStart[2])
              .format();
            res.business_day_end = moment(winStart, "X")
              .tz(timezone)
              .set('hour', busDayEnd[0])
              .set('minutes', busDayEnd[1])
              .set('seconds', busDayEnd[2])
              .format();

            // convert session times to ISO time
            for(let i=0; i<res.sessions.length; i++){
              const s = res.sessions[i];
              const session = new Session(
                s.id,
                s.title,
                s.comment,
                moment(s.start, "X").tz(timezone).format(),
                moment(s.end, "X").tz(timezone).format(),
                s.free,
                s.type
              );
              session.addRiders(s.riders);
              res.sessions[i] = session;
            }

            resolve(res);
          }else{
            console.log("Sessions/getSessions: Cannot get sessions, due to:", data.msg);
            reject([data.msg]);
          }
        })
        .catch(error => {
          console.error("Sessions/getSessions: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("Sessions/getSessions", error);
        reject([error]);
      })
    })
  }

  static createSession(sessionData){
    return new Promise((resolve, reject) => {
      // build request body
      const session = {
        title: sessionData.title,
        comment: sessionData.description,
        max_riders: sessionData.maximumRiders,
        start: moment(sessionData.start).format('X'),
        end: moment(sessionData.end).format('X'),
        type: sessionData.type
      };

      fetch('/api/booking.php?action=add_session', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(session)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("create session response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("Sessions/createSession: Cannot create session, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Sessions/createSession: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Sessions/createSession", error);
        reject([error]);
      })
    })
  }

  static editSession(sessionData){
    return new Promise((resolve, reject) => {
      // build request body
      const session = {
        id: sessionData.id,
        title: sessionData.title,
        comment: sessionData.description,
        max_riders: sessionData.maximumRiders,
        start: moment(sessionData.start).format('X'),
        end: moment(sessionData.end).format('X'),
        type: sessionData.type
      };

      fetch('/api/booking.php?action=edit_session', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(session)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("edit session response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("Sessions/editSession: Cannot create session, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Sessions/editSession: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Sessions/editSession", error);
        reject([error]);
      })
    })
  }

  static deleteSession(sessionData){
    console.log("api/deleteSession called for session:", sessionData);
    return new Promise((resolve, reject) => {
      // build request body
      const session = {
        session_id: sessionData.id
      };

      fetch('/api/booking.php?action=delete_session', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(session)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("Sessions/deleteSession response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("Sessions/deleteSession: Cannot create session, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Sessions/deleteSession: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Sessions/deleteSession", error);
        reject([error]);
      })
    })
  }

  static addUsersToSession(sessionId, users){
    console.log("api/addUsersToSession called for session:", sessionId, "and users", users);
    return new Promise((resolve, reject) => {
      // build request body
      const requestBody = {
        user_ids: users.map(u => u.id),
        session_id: sessionId
      };

      fetch('/api/booking.php?action=add_users', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestBody)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("Sessions/addUsersToSession response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("Sessions/addUsersToSession: Cannot add user to session, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Sessions/addUsersToSession: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Sessions/addUsersToSession", error);
        reject([error]);
      })
    })
  }

  static deleteUserFromSession(sessionId, user){
    console.log("api/deleteUserFromSession called for session:", sessionId, "and user", user);
    return new Promise((resolve, reject) => {
      // build request body
      const requestBody = {
        user_id: user.id,
        session_id: sessionId
      };

      fetch('/api/booking.php?action=delete_user', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestBody)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("Sessions/deleteUserFromSession response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("Sessions/deleteUserFromSession: Cannot delete user from session, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Sessions/deleteUserFromSession: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Sessions/deleteUserFromSession", error);
        reject([error]);
      })
    })
  }
}