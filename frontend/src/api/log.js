import Request from "booksys/api/common/request.js";

export default class Log {
  /**
   * Get all log entries from the backend
   */
  static getLogs() {
    return Request.getRequest("/api/v2/admin/logs");
  }
}
