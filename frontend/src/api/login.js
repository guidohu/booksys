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

      fetch("/api/login.php?action=login", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(request),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              if (data.ok) {
                resolve();
              } else {
                reject([data.msg]);
              }
            })
            .catch((error) => {
              reject([error]);
            });
        })
        .catch((error) => {
          reject([error]);
        });
    });
  }

  static logout(cbSuccess, cbFailed) {
    fetch("/api/logout.php", {
      method: "GET",
      cache: "no-cache",
    })
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        if (data.ok) {
          cbSuccess(data);
        } else {
          cbFailed("Error while logging out");
        }
      })
      .catch((error) => {
        console.log("action: login/logout", error);
        cbFailed({
          ok: false,
          message: error,
        });
      });
  }

  static getMyUser() {
    console.log("Login/getMyUser called");
    return Request.getRequest('/api/user.php?action=get_my_user');
  }

  static isLoggedIn() {
    console.log("Login/isLoggedIn called");
    return Request.getRequest('/api/login.php?action=isLoggedIn');
  }

  static async postData(data) {
    const result = await fetch("/api/login.php?action=login", {
      method: "POST",
      cache: "no-cache",
      body: JSON.stringify(data),
    });
    return result.json();
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
