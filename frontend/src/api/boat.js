import Request from "@/api/common/request.js";

export default class Boat {
  static getEngineHourLog() {
    console.log("api/v2/getEngineHourLog called");
    return Request.getRequest('/api/v2/boat/engine-hours/get');
  }

  static getEngineHourLogLatest() {
    console.log("api/v2/getEngineHourLogLatest called");
    return Request.getRequest('/api/v2/boat/engine-hour/latest/get');
  }

  /**
   * Adds a new engine hour entry.
   * @param {} engineHourEntry consisting of:
   * {
   *   user_id: ...
   *   engine_hours_before: ...
   *   engine_hours_after: ...
   *   type: ...                 // currently 0 for private, 1 for course
   * }
   */
  static addEngineHours(engineHourEntry) {
    console.log("api/v2/addEngineHours called, with", engineHourEntry);
    return Request.postRequest('/api/v2/boat/engine-hour/update', engineHourEntry);
  }

  /**
   * Updates the type of an existing engineHourEntry.
   * @param {} engineHourEntry consisting of:
   * {
   *   type: ...   (0: private, 1: course)
   *   id: ...
   * }
   */
  static updateEngineHours(engineHourEntryUpdate) {
    console.log("api/v2/updateEngineHours called, with", engineHourEntryUpdate);
    return Request.postRequest('/api/v2/boat/engine-hour/entry/update', engineHourEntryUpdate);
  }

  static getFuelLog() {
    console.log("api/v2/getFuelLog called");
    return Request.getRequest('/api/v2/boat/fuel-entries/get');
  }

  /**
   * Adds a new fuel entry.
   * @param {} fuelEntry consisting of:
   * {
   *   user_id: ...
   *   engine_hours: ...
   *   liters: ...
   *   cost: ...
   * }
   */
  static addFuelEntry(fuelEntry) {
    console.log("api/v2/addFuelEntry called");
    return Request.postRequest("/api/v2/boat/fuel-entry/add", fuelEntry);
  }

  /**
   * Updates an existing fuel entry.
   * @param {} fuelEntry consisting of:
   * {
   *   id: ...            id of the entry
   *   engine_hours: ...
   *   liters: ...
   *   cost: ...          net costs
   *   cost_brutto: ...   gross costs
   * }
   */
  static updateFuelEntry(fuelEntry) {
    console.log("api/v1/updateFuelEntry called, with", fuelEntry);
    return Request.postRequest("/api/v2/boat/fuel-entry/edit", fuelEntry);
  }

  static getMaintenanceLog() {
    console.log("api/v1/getMaintenanceLog called");
    return Request.getRequest('/api/v2/boat/maintenance-entries/get');
  }

  /**
   * Adds a new maintenance entry.
   * @param {} maintenanceEntry consisting of:
   * {
   *   user_id: ...
   *   engine_hours: ...
   *   description: ...
   * }
   */
  static addMaintenanceEntry(maintenanceEntry) {
    console.log("api/v2/addMaintenanceEntry called, with", maintenanceEntry);
    return Request.postRequest("/api/v2/boat/maintenance-entry/add", maintenanceEntry);
  }

  static getMyNautiqueInfo(boatId, token, tokenExpiry) {
    console.log("api/v1/getMyNautiqueInfo called");
    const request = {
      boat_id: boatId,
      token: token,
      token_expiry: tokenExpiry,
    };
    return Request.postRequest("/api/v2/boat/mynautique/telemetry/get", request);
  }
}
