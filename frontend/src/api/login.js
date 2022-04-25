import { sha256 } from "js-sha256";
import Request from "@/api/common/request.js";

export default class Login {
  // perform backend login
  static login(username, password) {
    console.log("Login/login called");
    return new Promise((resolve, reject) => {
      if (username == null || username == "") {
        reject(["no username provided"]);
      }
      if (password == null || password == "") {
        reject(["no password provided"]);
      }

      const pwHash = Login.calcHash(password);
      const request = {
        username: username,
        password: pwHash,
      };

      Request.postRequest("/api/login.php?action=login")
      .then((response) => resolve(response))
      .catch((error) => reject(error));
    });
  }

  static logout() {
    console.log("Login/logout called");
    return Request.getRequest('/api/logout.php');
  }

  static getMyUser() {
    console.log("Login/getMyUser called");
    return Request.getRequest('/api/user.php?action=get_my_user');
  }

  static isLoggedIn() {
    console.log("Login/isLoggedIn called");
    return Request.getRequest('/api/login.php?action=isLoggedIn');
  }

  static calcHash(password) {
    try {
      var hash = sha256.create();
      hash.update(password);
      return hash.hex();
    } catch (e) {
      console.error("Cannot calculate sha256 of password", e);
      return null;
    }
  }
}
