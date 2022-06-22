import Request from "@/api/common/request.js";

export default class Log {
  /**
   * Get all log entries from the backend
   */
  static getLogs() {
    console.debug("Log/getLogs called");
    return Request.getRequest("/api/log.php?action=get_logs");
  }
}
