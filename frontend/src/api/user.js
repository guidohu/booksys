import values from "lodash/values";
import Login from "./login";
import Request from "@/api/common/request.js";
import { UserPointer } from "@/dataTypes/user";

export default class User {
  static getHeats() {
    console.debug("User/getHeats called");
    return Request.getRequest('/api/v2/user/my/heats');
  }

  static getHeatStats() {
    console.debug("User/getHeatStats called");
    return Request.getRequest('/api/v2/user/my/heats/statistics');
  }

  static getBalance() {
    console.debug("User/getBalance called");
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
    console.debug("User/changeUserPassword called");
    const postData = {
      password_old: passwordData.oldPassword,
      password_new: passwordData.newPassword,
    };
    return Request.postRequest('/api/v2/user/my/password/update', postData);
  }

  static changeUserPasswordByToken(tokenAndPassword) {
    const password = Login.calcHash(tokenAndPassword.password);
    const postData = {
      email: tokenAndPassword.email,
      password: password,
      token: tokenAndPassword.token,
    };
    return Request.postRequest('/api/v1/password.php?action=change_password_by_token', postData);
  }

  static getUserSchedule() {
    console.debug("User/getUserSchedule called");
    return Request.getRequest('/api/v2/user/my/sessions');
  }

  static cancelSession(sessionId) {
    const queryData = {
      session_id: sessionId,
    };
    return Request.postRequest('/api/v1/booking.php?action=delete_user', queryData);
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
        resolve(users);
      })
      .catch((error) => {
        reject(error);
      })
    });
  }

  static getDetailedUserList() {
    console.debug("User/getDetailedUserList called");
    return Request.getRequest('/api/v1/user.php?action=get_all_users_detailed');
  }

  static lockUser(userId) {
    const queryData = {
      user_id: userId,
    };
    return Request.postRequest('/api/v1/user.php?action=lock_user', queryData);
  }

  static unlockUser(userId) {
    const queryData = {
      user_id: userId,
    };
    return Request.postRequest('/api/v1/user.php?action=unlock_user', queryData);
  }

  static deleteUser(userId) {
    const queryData = {
      id: userId,
    };
    return Request.postRequest('/api/v1/user.php?action=delete_user', queryData);
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
    return Request.postRequest('/api/v1/user.php?action=change_user_group_membership', queryData);
  }

  static getUserGroups() {
    console.debug("User/getUserGroups called");
    return Request.getRequest("/api/v2/user/groups/get");
  }

  static getUserRoles() {
    console.debug("User/getUserRoles called");
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
    console.log("User/requestToken: called with userData", userData);
    const queryData = {
      email: userData.email,
      recaptcha_token: userData.recaptchaResponse,
    };
    return Request.postRequest('/api/v1/password.php?action=token_request', queryData);
  }
}
