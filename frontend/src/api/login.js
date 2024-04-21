import Request from "@/api/common/request.js";

export default class Login {
  // perform backend login
  static login(username, password) {
    console.log("Login/login called");
    return new Promise((resolve, reject) => {
      if (username == null || username == "") {
        reject(["no username provided"]);
        return;
      }
      if (password == null || password == "") {
        reject(["no password provided"]);
        return;
      }

      const request = {
        username: username,
        password: password,
      };

      Request.postRequest("/api/v2/auth/login", request)
      .then((response) => resolve(response))
      .catch((error) => reject(error));
    });
  }

  static logout() {
    console.log("Login/logout called");
    return Request.getRequest('/api/v2/auth/logout');
  }

  static getMyUser() {
    return Request.getRequest('/api/v2/auth/user');
  }

  static isLoggedIn() {
    console.log("Login/isLoggedIn called");
    return Request.getRequest('/api/v2/auth/isloggedin');
  }
}
