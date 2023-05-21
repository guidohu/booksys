import Request from "@/api/common/request.js";

export const getBackendStatus = () => {
  console.debug("backend/getBackendStatus called");
  return Request.getRequest('/api/v1/backend.php?action=get_status');
};

export const getAPIV2Status = () => {
  console.debug("getAPIV2Status: /api/v2/ping called");
  return Request.getRequest('/api/v2/ping');
}
