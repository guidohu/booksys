import Boat from "@/api/boat";

const state = () => ({
  engineHourLog: [],
  engineHourLogLatest: null,
  fuelLog: []
});

const getters = {
  getEngineHourLog: (state) => {
    return state.engineHourLog;
  },
  getEngineHourLogLatest: (state) => {
    return state.engineHourLogLatest;
  },
  getFuelLog: (state) => {
    return state.fuelLog;
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
  queryFuelLog ({ commit }) {
    console.log("Trigger queryFuelLog");
    return new Promise((resolve, reject) => {
      Boat.getFuelLog()
      .then((response) => {
        commit('setFuelLog', response);
        resolve();
      })
      .catch(error => {
        reject(error);
      })
    })
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
  },
  addFuelEntry ({ dispatch }, fuelEntry) {
    console.log("Trigger addFuelEntry");
    return new Promise((resolve, reject) => {
      Boat.addFuelEntry(fuelEntry)
      .then(() => {
        dispatch('queryFuelLog');
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
  },
  setFuelLog (state, value){
    state.fuelLog = value;
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}

