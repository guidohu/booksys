import Request from "booksys/api/common/request.js";

export default class Boat {
  static getEngineHourLog() {
    return Request.getRequest("/api/v2/boat/engine-hours/list");
  }

  static getEngineHourLogLatest() {
    return Request.getRequest("/api/v2/boat/engine-hour/latest/get");
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
    return Request.postRequest(
      "/api/v2/boat/engine-hour/update",
      engineHourEntry,
    );
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
    return Request.postRequest(
      "/api/v2/boat/engine-hour/entry/update",
      engineHourEntryUpdate,
    );
  }

  static getFuelLog() {
    return Request.getRequest("/api/v2/boat/fuel-entries/get");
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
    return Request.postRequest("/api/v2/boat/fuel-entry/edit", fuelEntry);
  }

  static removeFuelEntry(id) {
    return Request.postRequest("/api/v2/boat/fuel-entry/remove", {
      id: id,
    });
  }

  static getMaintenanceLog() {
    return Request.getRequest("/api/v2/boat/maintenance-entries/get");
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
    return Request.postRequest(
      "/api/v2/boat/maintenance-entry/add",
      maintenanceEntry,
    );
  }

  static getMyNautiqueInfo(boatId, token, tokenExpiry) {
    const request = {
      boat_id: boatId,
    };
    return Request.postRequest(
      "/api/v2/boat/mynautique/telemetry/get",
      request,
    );
  }
}
