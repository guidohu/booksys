import { Login } from './login'

export default class User {

  static getHeats(cbSuccess, cbFailure){
    fetch('/api/user.php?action=get_my_user_heats', {
      method: 'GET',
      cache: 'no-cache'
    })
    .then(response => {
      return response.json()
    })
    .then(data => {
      cbSuccess(data)
    })
    .catch(error => {
      console.log('User/getHeats', error)
      cbFailure(
        error
      )
    })
  }

  static changeUserProfile(profileData){
    return new Promise((resolve, reject) => {
      fetch('/api/user.php?action=change_my_user_data', {
        method: 'POST',
        cache: 'no-cache',
        body: JSON.stringify(profileData)
      })
      .then(response => {
        response.json()
          .then( data => {
            if(data == null){
              console.error("User/changeUserProfile: Cannot parse server response", data);
              reject(["Cannot parse server response"]);
            }else if(data.ok == false){
              console.log("Cannot change user profile:", data.msg);
              reject([data.msg]);
            }else{
              console.log("User profile successfully changed");
              resolve(data.msg);
            }
          })
          .catch( error => {
            console.error("User/changeUserProfile: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error('User/changeUserProfile', error)
        reject([error]);
      })
    });
  }

  static changeUserPassword(passwordData){
    return new Promise((resolve, reject) => {
      const password_old = Login.calcHash(passwordData.oldPassword)
      const password_new = Login.calcHash(passwordData.newPassword)
      const postData = {
        password_old,
        password_new
      }

      fetch('/api/password.php?action=change_password_by_password', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(postData)
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data == null){
              console.error("User/changeUserPassword: Cannot parse server response", data);
              reject(["Cannot parse server response"]);
            }else if(data.ok == false){
              console.log("Cannot change user password:", data.message);
              reject([data.message]);
            }else{
              console.log("User password successfully changed");
              resolve(data.message);
            }
          })
          .catch( error => {
            console.error("User/changeUserPassword: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error('User/changeUserPassword', error)
        reject([error]);
      })
    })
  }

  static getUserSchedule(){
    return new Promise((resolve, reject) => {
      fetch('/api/user.php?action=get_my_user_sessions', {
        method: "GET",
        cache: 'no-cache'
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data == null){
              console.error("User/getUserSchedule: Cannot parse server response", data);
              reject(["Cannot parse server response"]);
            }else if(data.sessions == null || data.sessions_old == null){
              console.log("Cannot get user's sessions:", response);
              reject(["Cannot get user's sessions"]);
            }else{
              console.log("User's sessions retrieved");
              resolve(data);
            }
          })
          .catch( error => {
            console.error("User/getUserSchedule: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error('User/getUserSchedule', error)
        reject([error]);
      })
    })
  }

  static cancelSession(sessionId){
    return new Promise((resolve, reject) => {
      const queryData = {
        session_id: sessionId
      };

      fetch('/api/booking.php?action=delete_user', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(queryData)
      })
      .then(response => {
        if(!response.ok){
          throw Error(response.statusText);
        }else{
          resolve();
        }
      })
      .catch(error => {
        console.error("Cannot remove user from session");
        reject([error])
      })
    })
  }
}