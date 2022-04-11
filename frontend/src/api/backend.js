import Request from "@/api/common/request.js";

export const getBackendStatus = () => {
  console.log("backend/getBackendStatus called");
  return Request.getRequest('/api/backend.php?action=get_status');
};
