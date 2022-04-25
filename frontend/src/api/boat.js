import Request from "@/api/common/request.js";

export default class Boat {
  static getEngineHourLog() {
    console.log("api/getEngineHourLog called");
    return Request.getRequest('/api/boat.php?action=get_engine_hours_log');
  }

  static getEngineHourLogLatest() {
    console.log("api/getEngineHourLogLatest called");
    return Request.getRequest('/api/boat.php?action=get_engine_hours_latest');
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
    console.log("api/addEngineHours called");
    return Request.getRequest('/api/boat.php?action=update_engine_hours', engineHourEntry);
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
    console.log("api/updateEngineHours called, with", engineHourEntryUpdate);
    return Request.postRequest('/api/boat.php?action=update_engine_hours_entry', engineHourEntryUpdate);
  }

  static getFuelLog() {
    console.log("api/getFuelLog called");
    return Request.getRequest('/api/boat.php?action=get_fuel_log');
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
    console.log("api/addFuelEntry called");
    return Request.postRequest("/api/boat.php?action=update_fuel", fuelEntry);
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
    console.log("api/updateFuelEntry called, with", fuelEntry);
    return Request.postRequest("/api/boat.php?action=update_fuel_entry", fuelEntry);
  }

  static getMaintenanceLog() {
    console.log("api/getMaintenanceLog called");
    return Request.getRequest('/api/boat.php?action=get_maintenance_log');
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
    console.log("api/addMaintenanceEntry called, with", maintenanceEntry);
    return Request.postRequest("/api/boat.php?action=update_maintenance_log", maintenanceEntry);
  }

  static getMyNautiqueInfo(boatId, token, tokenExpiry) {
    console.log("api/getMyNautiqueInfo called");
    const request = {
      boat_id: boatId,
      token: token,
      token_expiry: tokenExpiry,
    };
    return Request.postRequest("/api/mynautique.php?action=get_boat_info", request);
  }
}
