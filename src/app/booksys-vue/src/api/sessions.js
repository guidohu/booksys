import moment from 'moment';

export default class Sessions {

  static getSessions(dateStart, dateEnd, cbSuccess, cbFailure){
    let query = {
      start: moment(dateStart).format("X"),
      end: moment(dateEnd).format("X")
    }
    fetch('/api/booking.php?action=get_booking_day', {
      method: 'POST',
      cache: 'no-cache',
      body: JSON.stringify(query)
    })
    .then(response => {
      return response.json()
    })
    .then(data => {
      cbSuccess(data)
    })
    .catch(error => {
      console.log('Sessions/getSessions', error)
      cbFailure(
        error
      )
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