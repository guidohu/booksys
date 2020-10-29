export default class Heat {

  /**
   * Creates heats on the backend
   * @param {*} heats an array of heats
   * [
   *    {
   *      user_id: ...
   *      session_id: ...
   *      duration_s: ...
   *      comment: ...
   *    }
   * ]
   */
  static addHeats(heats){
    console.log("heat/addHeats called, with", heats);
    return new Promise((resolve, reject) => {
      fetch('/api/heat.php?action=add_heats', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(heats)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("add heats response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("heat/addHeats: Cannot create heats, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("heat/addHeats: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("heat/addHeats", error);
        reject([error]);
      })
    })
  }

  /**
   * Get all heats from a session
   * @param {*} sessionId the ID of the session to get the heats from
   */
  static getHeatsBySession(sessionId){
    console.log('heat/getHeatsBySession called, with', sessionId);
    return new Promise((resolve, reject) => {
      const request = {
        session_id: sessionId
      };

      fetch('/api/heat.php?action=get_session_heats', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(request)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("get heats by session response data:", data);
            if(data.ok){
              resolve(data.data);
            }else{
              console.log("heat/getHeatsBySession: Cannot get heats, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("heat/getHeatsBySession: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("heat/getHeatsBySession", error);
        reject([error]);
      })
    })
  }
}