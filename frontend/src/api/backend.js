import Request from "api/common/request.js";

export const getBackendStatus = () => {
  console.debug("backend/getBackendStatus called");
  return Request.getRequest('/api/v2/health/status');
};
