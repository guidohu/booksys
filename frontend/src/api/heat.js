import Request from "@/api/common/request";

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
  static addHeats(heats) {
    console.log("heat/addHeats called, with", heats);
    return Request.postRequest("/api/heat.php?action=add_heats", heats);
  }

  /**
   * Deletes a heat with a given heat_id from the database
   * @param {Number} heatId
   */
  static deleteHeat(heatId) {
    console.log("heat/deleteHeat called, with", heatId);
    const request = {
      heat_id: heatId,
    };
    return Request.postRequest("/api/heat.php?action=delete_heat", request);
  }

  /**
   * Update a heat's duration and/or comment
   * @param {Object} heatUpdate
   */
  static updateHeat(heatUpdate) {
    console.log("heat/updateHeat called, with", heatUpdate);
    const request = {
      heat_id: heatUpdate.id,
      user_id: heatUpdate.userId,
      duration_s: heatUpdate.duration,
      comment: heatUpdate.comment,
    };
    return Request.postRequest("/api/heat.php?action=update_heat", request);
  }

  /**
   * Get all heats from a session
   * @param {*} sessionId the ID of the session to get the heats from
   */
  static getHeatsBySession(sessionId) {
    console.log("heat/getHeatsBySession called, with", sessionId);
    const request = {
      session_id: sessionId,
    };
    return Request.postRequest("/api/heat.php?action=get_session_heats", request);
  }
}
