import User from "../../api/user";

const state = () => ({
  heatHistory: [],
  balance: 0,
  balanceRounded: 0,
  heatTimeMinutes: 0,
  heatTimeMinutesYTD: 0,
  heatCost: 0,
  heatCostYTD: 0,
  schedule: [],
  userList: []
})

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
  }
}

const actions = {
  queryHeatHistory ( {commit} ) {
    let successCb = (response) => {
      console.log(response);
      commit('setHeatHistory', response);
    };
    let failCb = (error) => {
      console.log(error)
    };
    User.getHeats(successCb, failCb);
  },
  changeUserProfile ( {dispatch}, profileData ) {
    return new Promise((resolve, reject) => {
      User.changeUserProfile(profileData)
        .then(() => {
          // load the new user data
          dispatch('login/getUserInfo', {}, { root: true });
          resolve();
        })
        .catch(error => {
          reject(error);
        })
    });
  },
  changeUserPassword ( {dispatch}, passwordData ) {
    return new Promise((resolve, reject) => {
      console.log(passwordData);
      User.changeUserPassword(passwordData)
        .then(() => {
          // load latest user data
          dispatch('login/getUserInfo', {}, { root: true });
          resolve();
        })
        .catch(error => {
          reject(error);
        })
    })
  },
  queryUserSchedule ( {commit} ) {
    return new Promise((resolve, reject) => {
      User.getUserSchedule()
        .then((data) => {
          commit('setUserSchedule', data);
          resolve();
        })
        .catch(error => {
          reject(error);
        })
    })
  },
  cancelSession ( { dispatch }, sessionId ) {
    return new Promise((resolve, reject) => {
      User.cancelSession(sessionId)
        .then( () => {
          dispatch('queryUserSchedule')
        })
        .catch( (error) => {
          reject(error);
        })
    })
  },
  queryUserList ( { commit }) {
    console.log("Trigger action queryUserList")
    return new Promise((resolve, reject) => {
      User.getUserList()
        .then( (users) => {
          commit('setUserList', users);
          resolve();
        })
        .catch( (error) => {
          reject(error);
        })
    })
  }
}

const mutations = {
  setHeatHistory (state, value) {
    state.heatHistory = value.heats;
    state.balanceRounded = Math.round(value.balance_current*100)/100;
    state.balance = value.balance_current;
    state.heatTimeMinutes = value.heat_time_min;
    state.heatTimeMinutesYTD = value.heat_time_min_ytd;
    state.heatCost = value.heat_cost;
    state.heatCostYTD = value.heat_cost_ytd;
  },
  setUserSchedule (state, schedule) {
    // store:
    // sessions: [ {...}, {...}]
    // old_sessions: [ {...}, {...}]
    state.schedule = schedule;
  },
  setUserList (state, users) {
    console.log("Store: setUserList to", users);
    state.userList = users;
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}