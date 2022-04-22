import values from "lodash/values";
import Login from "./login";
import Request from "@/api/common/request.js";
import { UserPointer } from "@/dataTypes/user";

export default class User {
  static getHeats() {
    console.debug("User/getHeats called");
    return Request.getRequest('/api/user.php?action=get_my_user_heats');
  }

  static makeAdmin(userId) {
    const requestData = {
      user_id: userId,
    };
    return Request.postRequest('/api/configuration.php?action=make_user_admin', requestData);
  }

  static changeUserProfile(profileData) {
    return Request.postRequest('/api/user.php?action=change_my_user_data', profileData);
  }

  static changeUserPassword(passwordData) {
    const password_old = Login.calcHash(passwordData.oldPassword);
    const password_new = Login.calcHash(passwordData.newPassword);
    const postData = {
      password_old,
      password_new,
    };
    return Request.postRequest('/api/password.php?action=change_password_by_password', postData);
  }

  static changeUserPasswordByToken(tokenAndPassword) {
    const password = Login.calcHash(tokenAndPassword.password);
    const postData = {
      email: tokenAndPassword.email,
      password: password,
      token: tokenAndPassword.token,
    };
    return Request.postRequest('/api/password.php?action=change_password_by_token', postData);
  }

  static getUserSchedule() {
    console.debug("User/getUserSchedule called");
    return Request.getRequest('/api/user.php?action=get_my_user_sessions');
  }

  static cancelSession(sessionId) {
    const queryData = {
      session_id: sessionId,
    };
    return Request.postRequest('/api/booking.php?action=delete_user', queryData);
  }

  static getUserList() {
    return new Promise((resolve, reject) => {
      Request.getRequest('/api/user.php?action=get_all_users')
      .then((response) => {
        const usersResponse = response;
        let users = [];
        response.forEach((u) => {
          users.push(new UserPointer(u.id, u.first_name, u.last_name));
        });
        resolve(users);
      })
      .catch((error) => {
        reject(error);
      })
    });
  }

  static getDetailedUserList() {
    console.debug("User/getDetailedUserList called");
    return Request.getRequest('/api/user.php?action=get_all_users_detailed');
  }

  static lockUser(userId) {
    const queryData = {
      user_id: userId,
    };
    return Request.postRequest('/api/user.php?action=lock_user', queryData);
  }

  static unlockUser(userId) {
    const queryData = {
      user_id: userId,
    };
    return Request.postRequest('/api/user.php?action=unlock_user', queryData);
  }

  static deleteUser(userId) {
    const queryData = {
      id: userId,
    };
    return Request.postRequest('/api/user.php?action=delete_user', queryData);
  }

  static setUserGroup(userId, userGroupId) {
    console.log(
      "User/setUserGroup: called with userId",
      userId,
      "userGroupId",
      userGroupId
    );
    const queryData = {
      user_id: userId,
      status_id: userGroupId,
    };
    return Request.postRequest('/api/user.php?action=change_user_group_membership', queryData);
  }

  static getUserGroups() {
    console.debug("User/getUserGroups called");
    return Request.getRequest("/api/user.php?action=get_user_groups");
  }

  static getUserRoles() {
    console.debug("User/getUserRoles called");
    return Request.getRequest("/api/user.php?action=get_user_roles");
  }

  static saveUserGroup(userGroup) {
    console.log("User/saveUserGroup: called with userGroup", userGroup);
    const queryData = userGroup;
    return Request.postRequest('/api/user.php?action=save_user_group', queryData);
  }

  static deleteUserGroup(userGroupId) {
    console.log("User/deleteUserGroup: called with userGroupId", userGroupId);
    const queryData = {
      user_group_id: userGroupId,
    };
    return Request.postRequest('/api/user.php?action=delete_user_group', queryData);
  }

  static signUp(userData) {
    console.log("User/signUp: called with userGroupId", userData);
    const password = Login.calcHash(userData.password);
    const queryData = {
      username: userData.email,
      password: password,
      first_name: userData.firstName,
      last_name: userData.lastName,
      address: userData.street,
      mobile: userData.phone,
      plz: userData.zip,
      city: userData.city,
      email: userData.email,
      license: userData.license,
      recaptcha_token: userData.recaptchaResponse,
    };
    return Request.postRequest('/api/sign_up.php?action=sign_up', queryData);
  }

  static requestPasswordResetToken(userData) {
    console.log("User/requestToken: called with userData", userData);
    const queryData = {
      email: userData.email,
      recaptcha_token: userData.recaptchaResponse,
    };
    return Request.postRequest('/api/password.php?action=token_request', queryData);
  }
}
