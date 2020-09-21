import User from "../../api/user";

const state = () => ({
  heatHistory: [],
  balance: 0,
  balanceRounded: 0,
  heatTimeMinutes: 0,
  heatTimeMinutesYTD: 0,
  heatCost: 0,
  heatCostYTD: 0,
  userApiStatus: "INACTIVE",
  userApiErrors: []
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
  changeUserProfileStatus: (state) => {
    return state.changeUserProfileStatus
  },
  changeUserProfileError: (state) => {
    return state.changeUserProfileErrors
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
  changeUserProfile ( {dispatch, commit}, profileData ) {
    let successCb = (response) => {
      console.log(response);
      if(response.ok == true){
        dispatch('login/getUserInfo', {}, { root: true});
        commit('setUserApiStatus', "DONE");
      }else{
        console.error("Error reported from backend:",response);
        commit('setUserApiStatus', "ERROR");
        commit('setUserApiErrors', [ response.msg]);
      }
    };
    let failCb = (error) => {
      console.log(error)
      commit('setUserApiStatus', "ERROR");
      commit('setUserApiErrors', [ error ]);
    };

    User.changeUserProfile(profileData, successCb, failCb);
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
  setUserApiStatus (state, value) {
    state.userApiStatus = value;
  },
  setUserApiErrors (state, errors) {
    state.userApiErrors = errors;
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}