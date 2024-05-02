import values from "lodash/values";
import sortBy from "lodash/sortBy";
import Login from "./login";
import Request from "@/api/common/request.js";
import { UserPointer } from "@/dataTypes/user";

export default class User {
  static getHeats() {
    return Request.getRequest('/api/v2/user/my/heats');
  }

  static getHeatStats() {
    return Request.getRequest('/api/v2/user/my/heats/statistics');
  }

  static getBalance() {
    return Request.getRequest('/api/v2/user/my/balance');
  }

  static makeAdmin(userId) {
    const requestData = {
      user_id: userId,
    };
    return Request.postRequest('/api/v2/user/create-admin', requestData);
  }

  static changeUserProfile(profileData) {
    return Request.postRequest('/api/v2/user/my/update', profileData);
  }

  static changeUserPassword(passwordData) {
    const postData = {
      password_old: passwordData.oldPassword,
      password_new: passwordData.newPassword,
    };
    return Request.postRequest('/api/v2/user/my/password/update', postData);
  }

  static changeUserPasswordByToken(tokenAndPassword) {
    console.debug("User/changeUserPasswordByToken called");
    const postData = {
      email: tokenAndPassword.email,
      password: tokenAndPassword.password,
      token: tokenAndPassword.token,
    };
    return Request.postRequest('/api/v2/user/password/reset-by-token', postData);
  }

  static getUserSchedule() {
    return Request.getRequest('/api/v2/user/my/sessions');
  }

  static cancelSession(sessionId) {
    const queryData = {
      session_id: sessionId,
    };
    return Request.postRequest('/api/v2/user/my/session/delete', queryData);
  }

  static getUserList() {
    return new Promise((resolve, reject) => {
      Request.getRequest('/api/v2/user/list-short')
      .then((response) => {
        const usersResponse = response;
        let users = [];
        response.forEach((u) => {
          users.push(new UserPointer(u.id, u.first_name, u.last_name));
        });
        users = sortBy(users, ['firstName', 'lastName']);
        resolve(users);
      })
      .catch((error) => {
        reject(error);
      })
    });
  }

  static getDetailedUserList() {
    return Request.getRequest('/api/v2/user/list-detailed');
  }

  static lockUser(userId, locked) {
    const queryData = {
      user_id: userId,
      locked: locked,
    };
    return Request.postRequest('/api/v2/user/lock/set', queryData);
  }

  static deleteUser(userId) {
    const queryData = {
      user_id: userId,
    };
    return Request.postRequest('/api/v2/user/delete', queryData);
  }

  static setUserGroup(userId, userGroupId) {
    const queryData = {
      user_id: parseInt(userId),
      status_id: parseInt(userGroupId),
    };
    return Request.postRequest('/api/v2/user/group/set', queryData);
  }

  static getUserGroups() {
    return Request.getRequest("/api/v2/user/groups/get");
  }

  static getUserRoles() {
    return Request.getRequest("/api/v2/user/roles/get");
  }

  static saveUserGroup(userGroup) {
    console.log("User/saveUserGroup: called with userGroup", userGroup);
    let url = '/api/v2/user/group/edit'
    if (userGroup.user_group_id == null) {
      url = '/api/v2/user/group/create'
    }
    const queryData = {
      price_id: parseInt(userGroup.price_id),
      price_description: userGroup.price_description,
      price_min: parseFloat(userGroup.price_min),
      user_group_id: parseInt(userGroup.user_group_id),
      user_group_description: userGroup.user_group_description,
      user_group_name: userGroup.user_group_name,
      user_role_id: parseInt(userGroup.user_role_id),
    }
    return Request.postRequest(url, queryData);
  }

  static deleteUserGroup(userGroupId) {
    console.log("User/deleteUserGroup: called with userGroupId", userGroupId);
    const queryData = {
      user_group_id: userGroupId,
    };
    return Request.postRequest('/api/v2/user/group/delete', queryData);
  }

  static signUp(userData) {
    console.log("User/signUp: called with userGroupId", userData);
    const queryData = {
      username: userData.email,
      password: userData.password,
      first_name: userData.firstName,
      last_name: userData.lastName,
      address: userData.street,
      mobile: userData.phone,
      plz: parseInt(userData.zip),
      city: userData.city,
      email: userData.email,
      license: userData.license,
      ownRisk: userData.ownRisk,
      recaptcha_token: userData.recaptchaResponse,
    };
    return Request.postRequest('/api/v2/user/signup', queryData);
  }

  static requestPasswordResetToken(userData) {
    console.log("/api/v2/user/password/token-request: called with userData", userData);
    const queryData = {
      email: userData.email,
      recaptcha_token: userData.recaptchaResponse,
    };
    return Request.postRequest('/api/v2/user/password/token-request', queryData);
  }
}
