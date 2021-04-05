import User from "@/api/user";

const state = () => ({
  heatHistory: [],
  balance: 0,
  balanceRounded: 0,
  heatTimeMinutes: 0,
  heatTimeMinutesYTD: 0,
  heatCost: 0,
  heatCostYTD: 0,
  schedule: [],
  userList: [],
  userListDetailed: [],
  userGroups: [],
  userRoles: [],
});

const getters = {
  heatHistory: (state) => {
    return state.heatHistory;
  },
  balance: (state) => {
    return state.balance;
  },
  balanceRounded: (state) => {
    return state.balanceRounded;
  },
  heatTimeMinutes: (state) => {
    return state.heatTimeMinutes;
  },
  heatTimeMinutesYTD: (state) => {
    return state.heatTimeMinutesYTD;
  },
  heatCost: (state) => {
    return state.heatCost;
  },
  heatCostYTD: (state) => {
    return state.heatCostYTD;
  },
  userSchedule: (state) => {
    return state.schedule;
  },
  userList: (state) => {
    return state.userList;
  },
  userListDetailed: (state) => {
    return state.userListDetailed;
  },
  userGroups: (state) => {
    return state.userGroups;
  },
  userRoles: (state) => {
    return state.userRoles;
  },
};

const actions = {
  queryHeatHistory({ commit }) {
    let successCb = (response) => {
      console.log(response);
      commit("setHeatHistory", response);
    };
    let failCb = (error) => {
      console.log(error);
    };
    User.getHeats(successCb, failCb);
  },
  changeUserProfile({ dispatch }, profileData) {
    return new Promise((resolve, reject) => {
      User.changeUserProfile(profileData)
        .then(() => {
          // load the new user data
          dispatch("login/getUserInfo", {}, { root: true });
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  changeUserPassword({ dispatch }, passwordData) {
    return new Promise((resolve, reject) => {
      User.changeUserPassword(passwordData)
        .then(() => {
          // load latest user data
          dispatch("login/getUserInfo", {}, { root: true });
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryUserSchedule({ commit }) {
    return new Promise((resolve, reject) => {
      User.getUserSchedule()
        .then((data) => {
          commit("setUserSchedule", data);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  cancelSession({ dispatch }, sessionId) {
    return new Promise((resolve, reject) => {
      User.cancelSession(sessionId)
        .then(() => {
          dispatch("queryUserSchedule");
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryUserList({ commit }) {
    console.log("Trigger action queryUserList");
    return new Promise((resolve, reject) => {
      User.getUserList()
        .then((users) => {
          commit("setUserList", users);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryUserListDetailed({ commit }) {
    console.log("Trigger action queryUserListDetailed");
    return new Promise((resolve, reject) => {
      User.getDetailedUserList()
        .then((users) => {
          commit("setUserListDetailed", users);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryUserGroups({ commit }) {
    console.log("Trigger action queryUserGroups");
    return new Promise((resolve, reject) => {
      User.getUserGroups()
        .then((users) => {
          commit("setUserGroups", users);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryUserRoles({ commit }) {
    console.log("Trigger action queryUserRoles");
    return new Promise((resolve, reject) => {
      User.getUserRoles()
        .then((roles) => {
          commit("setUserRoles", roles);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  saveUserGroup({ dispatch }, userGroup) {
    console.log("Trigger action saveUserGroup");
    return new Promise((resolve, reject) => {
      User.saveUserGroup(userGroup)
        .then(() => {
          dispatch("queryUserGroups");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  deleteUserGroup({ dispatch }, userGroupId) {
    console.log("Trigger action deleteUserGroup");
    return new Promise((resolve, reject) => {
      User.deleteUserGroup(userGroupId)
        .then(() => {
          dispatch("queryUserGroups");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  lockUser({ dispatch }, userId) {
    console.log("Trigger action lockUser");
    return new Promise((resolve, reject) => {
      User.lockUser(userId)
        .then(() => {
          dispatch("queryUserListDetailed");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  unlockUser({ dispatch }, userId) {
    console.log("Trigger action unlockUser");
    return new Promise((resolve, reject) => {
      User.unlockUser(userId)
        .then(() => {
          dispatch("queryUserListDetailed");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  deleteUser({ dispatch }, userId) {
    console.log("Trigger action deleteUser");
    return new Promise((resolve, reject) => {
      User.deleteUser(userId)
        .then(() => {
          dispatch("queryUserListDetailed");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  setUserGroup({ dispatch }, data) {
    console.log("Trigger action setUserGroup with", data);
    return new Promise((resolve, reject) => {
      User.setUserGroup(data.userId, data.userGroupId)
        .then(() => {
          dispatch("queryUserListDetailed");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
};

const mutations = {
  setHeatHistory(state, value) {
    state.heatHistory = value.heats;
    state.balanceRounded = Math.round(value.balance_current * 100) / 100;
    state.balance = value.balance_current;
    state.heatTimeMinutes = value.heat_time_min;
    state.heatTimeMinutesYTD = value.heat_time_min_ytd;
    state.heatCost = value.heat_cost;
    state.heatCostYTD = value.heat_cost_ytd;
  },
  setUserSchedule(state, schedule) {
    // store:
    // sessions: [ {...}, {...}]
    // old_sessions: [ {...}, {...}]
    state.schedule = schedule;
  },
  setUserList(state, users) {
    console.log("Store: setUserList to", users);
    state.userList = users;
  },
  setUserListDetailed(state, users) {
    console.log("Store: setUserListDetailed to", users);
    state.userListDetailed = users;
  },
  setUserGroups(state, groups) {
    console.log("Store: setUserGroups to", groups);
    state.userGroups = groups;
  },
  setUserRoles(state, roles) {
    console.log("Store: setUserRoles to", roles);
    state.userRoles = roles;
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
