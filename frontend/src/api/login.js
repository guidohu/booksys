import { sha256 } from 'js-sha256';

export default class Login {

  // perform backend login
  static login(username, password){
    console.log('Login/login called');
    return new Promise((resolve, reject) => {
      if(username == null || username == ""){
        reject(["no username provided"]);
      }
      if(password == null || password == ""){
        reject(["no password provided"]);
      }

      const pwHash = Login.calcHash(password);
      const request = {
        username: username,
        password: pwHash
      };

      fetch('/api/login.php?action=login', {
        method: 'POST',
        cache: 'no-cache',
        body: JSON.stringify(request)
      })
      .then(response => {
        response.json()
        .then(data => {
          if(data.ok){
            resolve();
          }else{
            reject([data.msg]);
          }
        })
        .catch(error => {
          reject([error]);
        })
      })
      .catch(error => {
        reject([error]);
      })
    })
  }

  static logout(cbSuccess, cbFailed){
      fetch('/api/logout.php', {
          method: 'GET',
          cache: 'no-cache'
      })
      .then(response => {
          return response.json()})
      .then(data => {
          if(data.ok){
            cbSuccess(data)
          }else{
            cbFailed("Error while logging out")
          }
      })
      .catch(error => {
          console.log('action: login/logout', error)
          cbFailed({
              ok: false,
              message: error
          })
      })
  }

  static getMyUser(){
    console.log('Login/getMyUser called');
    return new Promise((resolve, reject) => {
      fetch('/api/user.php?action=get_my_user', {
        method: "GET",
        cache: "no-cache"
      })
      .then(response => {
        response.json()
        .then(data => {
          if(data.ok){
            resolve(data.data);
          }else{
            console.error("Login/getMyUser: Cannot get user info", data);
            reject([data.msg]);
          }
        })
        .catch(error => {
          console.error("Login/getMyUser: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("Login/getMyUser:", error);
        reject([error]);
      })
    })
  }

  static isLoggedIn(){    
    console.log('Login/isLoggedIn called');
    return new Promise((resolve, reject) => {
      fetch('/api/login.php?action=isLoggedIn', {
        method: "GET",
        cache: "no-cache"
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data.ok){
              resolve(data.data.loggedIn);
            }else{
              console.error("Login/isLoggedIn: Cannot get login status", data);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Login/isLoggedIn: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Login/isLoggedIn:", error);
        reject([error]);
      })
    })
  }

  static async postData(data) {
      const result = await fetch('/api/login.php?action=login', {
          method: 'POST',
          cache: 'no-cache',
          body: JSON.stringify(data)
      })
      return result.json()
  }

  static async getData(url) {
      const result = await fetch(url, {
          method: 'GET',
          cache: 'no-cache'
      })
      return result.json()
  }

  static calcHash(password) {
    try {
      var hash = sha256.create();
      hash.update(password);
      return hash.hex();
    } catch(e) {
      console.error("Cannot calculate sha256 of password", e);
      return null;
    }
  }
}