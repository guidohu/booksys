import Boat from "@/api/boat";
import { sortBy, reverse, forEach, round } from 'lodash';

const state = () => ({
  avgFuelConsumption: null,
  engineHourLog: [],
  engineHourLogLatest: null,
  fuelLog: [],
  maintenanceLog: [],
  myNautique: {
    token: "",
    tokenExpiry: 0,
    boat: {
      fuelLevel: 0,
      fuelCapacity: 1,
      engineHours: null,
    }
  }
});

const getters = {
  getAvgFuelConsumption: (state) => {
    return state.avgFuelConsumption;
  },
  getEngineHourLog: (state) => {
    return state.engineHourLog;
  },
  getEngineHourLogLatest: (state) => {
    return state.engineHourLogLatest;
  },
  getFuelLog: (state) => {
    return state.fuelLog;
  },
  getMaintenanceLog: (state) => {
    return state.maintenanceLog;
  },
  getMyNautiqueFuelLevel: (state) => {
    return state.myNautique.boat.fuelLevel;
  },
  getMyNautiqueFuelCapacity: (state) => {
    return state.myNautique.boat.fuelCapacity;
  },
  getMyNautiqueEngineHours: (state) => {
    return state.myNautique.boat.engineHours;
  }
};

const actions = {
  queryEngineHourLog({ commit }) {
    console.log("Trigger queryEngineHourLog");
    return new Promise((resolve, reject) => {
      Boat.getEngineHourLog()
        .then((response) => {
          commit("setEngineHourLog", response);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryEngineHourLogLatest({ commit }) {
    console.log("Trigger queryEngineHourLogLatest");
    return new Promise((resolve, reject) => {
      Boat.getEngineHourLogLatest()
        .then((response) => {
          commit("setEngineHourLogLatest", response);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryFuelLog({ commit }) {
    console.log("Trigger queryFuelLog");
    return new Promise((resolve, reject) => {
      Boat.getFuelLog()
        .then((response) => {
          commit("setFuelLog", response);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryMaintenanceLog({ commit }) {
    console.log("Trigger queryMaintenanceLog");
    return new Promise((resolve, reject) => {
      Boat.getMaintenanceLog()
        .then((response) => {
          commit("setMaintenanceLog", response);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  queryMyNautiqueInfo({ commit, state }, boatId) {
    console.log("Trigger queryMyNautiqueInfo");
    return new Promise((resolve, reject) => {
      Boat.getMyNautiqueInfo(boatId, state.myNautique.token, state.myNautique.tokenExpiry)
        .then((response) => {
          commit("setMyNautiqueInfo", response);
          resolve();
        })
        .catch((error) => {
          reject(error);
        })
    })
  },
  addEngineHours({ dispatch }, engineHourEntry) {
    console.log("Trigger addEngineHourLogEntry");
    return new Promise((resolve, reject) => {
      Boat.addEngineHours(engineHourEntry)
        .then(() => {
          dispatch("queryEngineHourLog");
          dispatch("queryEngineHourLogLatest");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  addFuelEntry({ dispatch }, fuelEntry) {
    console.log("Trigger addFuelEntry");
    return new Promise((resolve, reject) => {
      Boat.addFuelEntry(fuelEntry)
        .then(() => {
          dispatch("queryFuelLog");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  addMaintenanceEntry({ dispatch }, maintenanceEntry) {
    console.log("Trigger addMaintenanceEntry");
    return new Promise((resolve, reject) => {
      Boat.addMaintenanceEntry(maintenanceEntry)
        .then(() => {
          dispatch("queryMaintenanceLog");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  updateEngineHours({ dispatch }, engineHourUpdate) {
    console.log("Trigger updateEngineHours, with", engineHourUpdate);
    return new Promise((resolve, reject) => {
      Boat.updateEngineHours(engineHourUpdate)
        .then(() => {
          dispatch("queryEngineHourLog");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  updateFuelEntry({ dispatch }, fuelEntry) {
    console.log("Trigger updateFuelEntry, with", fuelEntry);
    return new Promise((resolve, reject) => {
      Boat.updateFuelEntry(fuelEntry)
        .then(() => {
          dispatch("queryFuelLog");
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
};

const mutations = {
  setEngineHourLog(state, value) {
    state.engineHourLog = value;
  },
  setEngineHourLogLatest(state, value) {
    state.engineHourLogLatest = value;
  },
  setFuelLog(state, value) {
    state.fuelLog = value;

    // calculate the average fuel consumption
    // of the last 5 pitstops as a reference
    const sortedLog = reverse(
      sortBy(
        value, 
        function(v){ return v.timestamp }
      )
    );
    
    let i = 0;
    let totalDiffHours = 0;
    let totalFuel = 0;
    forEach(sortedLog, function(v){
      if(i < 5 && v.avg_liters_per_hour != null && v.diff_hours != null){
        if(v.diff_hours > 0 && v.avg_liters_per_hour > 0){
          totalFuel += v.diff_hours * v.avg_liters_per_hour;
          totalDiffHours += v.diff_hours;
          i++;
        }
      }
    });

    // calculate average fuel consumption per hour
    if(totalDiffHours > 0){
      state.avgFuelConsumption = round(
        totalFuel / totalDiffHours,
        1
      );
    }
  },
  setMaintenanceLog(state, value) {
    state.maintenanceLog = value;
  },
  setMyNautiqueInfo(state, value) {
    if(value == null){
      console.log("myNautique not available");
      return;
    }
    
    state.myNautique.token = value.token;
    state.myNautique.tokenExpiry = value.token_expiry;
    state.myNautique.boat.fuelLevel = value.boat.telemetry.FUEL_LEVEL_LINC;
    state.myNautique.boat.fuelCapacity = value.boat.metainfo.fuel_capacity;
    state.myNautique.boat.engineHours = value.boat.telemetry.EngineTotalHoursOfOperation;
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
