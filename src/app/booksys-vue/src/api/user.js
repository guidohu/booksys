import _ from 'lodash';
import { Login } from './login';
import { UserPointer } from '@/dataTypes/user';

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

  static getUserList(){
    return new Promise((resolve, reject) => {
      fetch('/api/user.php?action=get_all_users', {
        method: "GET",
        cache: 'no-cache',
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data == null){
              console.error("User/getUserList: Cannot parse server response", data);
              reject(["Cannot parse server response"]);
            }else if(data.data == null){
              console.log("User/getUserList: Cannot get user list:", response);
              reject(["Cannot get user list"]);
            }else{
              console.log("User/getUserList: User's sessions retrieved");
              const usersResponse = data.data
              let users = [];
              usersResponse.forEach(u => {
                users.push(new UserPointer(
                  u.id,
                  u.first_name,
                  u.last_name
                ))
              })
              resolve(users);
            }
          })
          .catch( error => {
            console.error("User/getUserList: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error('User/getUserList:', error)
        reject([error]);
      })
    })
  }

  static getDetailedUserList(){
    return new Promise((resolve, reject) => {
      fetch('/api/user.php?action=get_all_users_detailed', {
        method: "GET",
        cache: 'no-cache',
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data == null){
              console.error("User/getDetailedUserList: Cannot parse server response", data);
              reject(["Cannot parse server response"]);
            }else if(data.data.users == null){
              console.log("User/getDetailedUserList: Cannot get user list:", response);
              reject(["Cannot get user list"]);
            }else{
              console.log("User/getDetailedUserList: User list retrieved");
              const usersResponse = _.values(data.data.users);
              resolve(usersResponse);
            }
          })
          .catch( error => {
            console.error("User/getDetailedUserList: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error('User/getDetailedUserList:', error)
        reject([error]);
      })
    })
  }

  static lockUser(userId){
    return new Promise((resolve, reject) => {
      const queryData = {
        user_id: userId
      };

      fetch('/api/user.php?action=lock_user', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(queryData)
      })
      .then(response => {
        response.json()
        .then(data => {
          if(!data.ok){
            console.error("User/lockUser: Cannot lock user, got error", data.msg);
            reject([data.msg]);
          }else{
            resolve();
          }
        })
        .catch(error => {
          console.error("User/lockUser: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("User/lockUser: Cannot lock user with ID", error);
        reject([error])
      })
    })
  }

  static unlockUser(userId){
    return new Promise((resolve, reject) => {
      const queryData = {
        user_id: userId
      };

      fetch('/api/user.php?action=unlock_user', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(queryData)
      })
      .then(response => {
        response.json()
        .then(data => {
          if(!data.ok){
            console.error("User/unlockUser: Cannot unlock user, got error", data.msg);
            reject([data.msg]);
          }else{
            resolve();
          }
        })
        .catch(error => {
          console.error("User/unlockUser: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("User/unlockUser: Cannot unlock user with ID", error);
        reject([error])
      })
    })
  }

  static deleteUser(userId){
    return new Promise((resolve, reject) => {
      const queryData = {
        id: userId
      };

      fetch('/api/user.php?action=delete_user', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(queryData)
      })
      .then(response => {
        response.json()
        .then(data => {
          if(!data.ok){
            console.error("User/deleteUser: Cannot delete user, got error", data.msg);
            reject([data.msg]);
          }else{
            resolve();
          }
        })
        .catch(error => {
          console.error("User/deleteUser: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("User/deleteUser: Cannot delete user with ID", error);
        reject([error])
      })
    })
  }

  static setUserGroup(userId, userGroupId){
    console.log("User/setUserGroup: called with userId", userId, "userGroupId", userGroupId);
    return new Promise((resolve, reject) => {
      const queryData = {
        user_id: userId,
        status_id: userGroupId
      };

      fetch('/api/user.php?action=change_user_group_membership', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(queryData)
      })
      .then(response => {
        response.json()
        .then(data => {
          if(!data.ok){
            console.error("User/setUserGroup: Cannot set user group, got error", data.msg);
            reject([data.msg]);
          }else{
            resolve();
          }
        })
        .catch(error => {
          console.error("User/setUserGroup: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("User/setUserGroup: Cannot unlock user with ID", error);
        reject([error])
      })
    })
  }

  static getUserGroups(){
    return new Promise((resolve, reject) => {
      fetch('/api/user.php?action=get_user_groups', {
        method: "GET",
        cache: 'no-cache',
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data == null){
              console.error("User/getUserGroups: Cannot parse server response", data);
              reject(["Cannot parse server response"]);
            }else{
              console.log("User/getUserGroups: User group list retrieved");
              const usersResponse = _.values(data.data);
              resolve(usersResponse);
            }
          })
          .catch( error => {
            console.error("User/getUserGroups: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error('User/getUserGroups:', error)
        reject([error]);
      })
    })
  }

  static getUserRoles(){
    return new Promise((resolve, reject) => {
      fetch('/api/user.php?action=get_user_roles', {
        method: "GET",
        cache: 'no-cache',
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data == null){
              console.error("User/getUserRoles: Cannot parse server response", data);
              reject(["Cannot parse server response"]);
            }else{
              console.log("User/getUserRoles: User role list retrieved");
              const usersResponse = _.values(data.data);
              resolve(usersResponse);
            }
          })
          .catch( error => {
            console.error("User/getUserRoles: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error('User/getUserRoles:', error)
        reject([error]);
      })
    })
  }

  static saveUserGroup(userGroup){
    console.log("User/saveUserGroup: called with userGroup", userGroup);
    return new Promise((resolve, reject) => {
      const queryData = userGroup;

      fetch('/api/user.php?action=save_user_group', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(queryData)
      })
      .then(response => {
        response.json()
        .then(data => {
          if(!data.ok){
            console.error("User/saveUserGroup: Cannot save user group, got error", data.msg);
            reject([data.msg]);
          }else{
            resolve();
          }
        })
        .catch(error => {
          console.error("User/saveUserGroup: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("User/saveUserGroup: Cannot unlock user with ID", error);
        reject([error])
      })
    })
  }

  static deleteUserGroup(userGroupId){
    console.log("User/deleteUserGroup: called with userGroupId", userGroupId);
    return new Promise((resolve, reject) => {
      const queryData = {
        user_group_id: userGroupId
      };

      fetch('/api/user.php?action=delete_user_group', {
        method: "POST",
        cache: 'no-cache',
        body: JSON.stringify(queryData)
      })
      .then(response => {
        response.json()
        .then(data => {
          if(!data.ok){
            console.error("User/deleteUserGroup: Cannot delete user group, got error", data.msg);
            reject([data.msg]);
          }else{
            resolve();
          }
        })
        .catch(error => {
          console.error("User/deleteUserGroup: Cannot parse server response", error);
          reject([error]);
        })
      })
      .catch(error => {
        console.error("User/deleteUserGroup: Cannot unlock user with ID", error);
        reject([error])
      })
    })
  }

}