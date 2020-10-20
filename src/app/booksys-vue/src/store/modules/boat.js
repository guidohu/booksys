import Boat from "@/api/boat";

const state = () => ({
  engineHourLog: [],
  engineHourLogLatest: null
});

const getters = {
  getEngineHourLog: (state) => {
    return state.engineHourLog;
  },
  getEngineHourLogLatest: (state) => {
    return state.engineHourLogLatest;
  }
};

const actions = {
  queryEngineHourLog ({ commit }) {
    console.log("Trigger queryEngineHourLog");
    return new Promise((resolve, reject) => {
      Boat.getEngineHourLog()
      .then((response) => {
        commit('setEngineHourLog', response);
        resolve();
      })
      .catch(error => {
        reject(error);
      })
    });
  },
  queryEngineHourLogLatest ({ commit }) {
    console.log("Trigger queryEngineHourLogLatest");
    return new Promise((resolve, reject) => {
      Boat.getEngineHourLogLatest()
      .then((response) => {
        commit('setEngineHourLogLatest', response);
        resolve();
      })
      .catch(error => {
        reject(error)
      })
    });
  },
  addEngineHours ({ dispatch }, engineHourEntry) {
    console.log("Trigger addEngineHourLogEntry");
    return new Promise((resolve, reject) => {
      Boat.addEngineHours(engineHourEntry)
      .then(() => {
        dispatch('queryEngineHourLog');
        dispatch('queryEngineHourLogLatest');
        resolve();
      })
      .catch(error => {
        reject(error)
      })
    });
  }
};

const mutations = {
  setEngineHourLog (state, value){
    state.engineHourLog = value;
  },
  setEngineHourLogLatest (state, value){
    state.engineHourLogLatest = value;
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}

