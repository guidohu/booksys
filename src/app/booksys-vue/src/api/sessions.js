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


  // in case we do not have an ID, create a new session
//   if(session.id == null){
//     var req = new Object();
//     req.title      = session.title;
//     req.comment    = session.description;
//     req.max_riders = session.maxRiders;
//     req.start      = session.start.format('X');
//     req.end        = session.end.format('X');
//     req.type       = session.type;
//     $.ajax({
//         type: "POST",
//         url:  "api/booking.php?action=add_session",
//         data: JSON.stringify(req),
//         success: function(response){
//             var json = $.parseJSON(response);
//             cb.success(json.session_id);
//         },
//         error: function(xhr,status,error){
//             console.log(xhr);
//             console.log(status);
//             console.log(error);
//             var errorMsg = $.parseJSON(xhr.responseText);
//             cb.failure([errorMsg.error]);
//         },
//     });
// }
}