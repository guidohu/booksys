export default class Heat {

  /**
   * Creates heats on the backend
   * @param {*} heats an array of heats
   * [
   *    {
   *      uid: ...
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
              resolve(data.data);
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
   * Deletes a heat with a given heat_id from the database
   * @param {Number} heatId 
   */
  static deleteHeat(heatId){
    console.log("heat/deleteHeat called, with", heatId);
    return new Promise((resolve, reject) => {
      const request = {
        heat_id: heatId
      };

      fetch('/api/heat.php?action=delete_heat', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(request)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("delete heat response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("heat/deleteHeat: Cannot delete heats due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("heat/deleteHeat: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("heat/deleteHeat", error);
        reject([error]);
      })
    })
  }

  /**
   * Update a heat's duration and/or comment
   * @param {Object} heatUpdate 
   */
  static updateHeat(heatUpdate){
    console.log("heat/updateHeat called, with", heatUpdate);
    return new Promise((resolve, reject) => {
      const request = {
        heat_id: heatUpdate.id,
        user_id: heatUpdate.userId,
        duration_s: heatUpdate.duration,
        comment: heatUpdate.comment
      };

      fetch('/api/heat.php?action=update_heat', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(request)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("update heat response data:", data);
            if(data.ok){
              resolve(data);
            }else{
              console.log("heat/updateHeat: Cannot delete heats due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("heat/updateHeat: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("heat/updateHeat", error);
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