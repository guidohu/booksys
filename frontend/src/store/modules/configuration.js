import Configuration from "@/api/configuration";

const state = () => ({
  CONFIG_LOADED: false,
  configuration: null,
  recaptchaKey: null,
  currency: null,
  locationAddress: null,
  locationMap: null,
  dbUpdateStatus: null,
  dbVersionInfo: null,
  dbIsUpdating: false,
  dbUpdateResult: null,
  timezone: "Europe/Zurich",
  maxRiders: 12,
  logoFile: null,
  engineHourFormat: "hh.h",
  fuelPaymentType: null,
  myNautiqueEnabled: null,
  myNautiqueUser: null,
  myNautiquePassword: null,
});

const getters = {
  getRecaptchaKey: (state) => {
    return state.recaptchaKey;
  },
  getConfiguration: (state) => {
    return state.configuration;
  },
  getCurrency: (state) => {
    return state.currency;
  },
  getLocationAddress: (state) => {
    return state.locationAddress;
  },
  getLocationMap: (state) => {
    return state.locationMap;
  },
  getDbUpdateStatus: (state) => {
    return state.dbUpdateStatus;
  },
  getDbVersionInfo: (state) => {
    console.log("getDbVersionInfo set to", state.dbVersionInfo);
    return state.dbVersionInfo;
  },
  getDbIsUpdating: (state) => {
    return state.dbIsUpdating;
  },
  getDbUpdateResult: (state) => {
    return state.dbUpdateResult;
  },
  getTimezone: (state) => {
    return state.timezone;
  },
  getMaximumNumberOfRiders: (state) => {
    return state.maxRiders;
  },
  getLogoFile: (state) => {
    return state.logoFile;
  },
  getEngineHourFormat: (state) => {
    return state.engineHourFormat;
  },
  getFuelPaymentType: (state) => {
    return state.fuelPaymentType;
  },
  getMyNautiqueEnabled: (state) => {
    return state.myNautiqueEnabled;
  },
  getMyNautiqueBoatId: (state) => {
    return state.myNautiqueBoatId;
  },
};

const actions = {
  invalidateConfiguration({ commit }) {
    commit("invalidateConfiguration", {});
  },
  queryConfiguration({ commit, state }) {
    if (state.CONFIG_LOADED == true) {
      return;
    }

    return new Promise((resolve, reject) => {
      Configuration.getConfiguration()
      .then((config) => {
        commit("setConfiguration", config);
        resolve(config);
      })
      .catch((errors) => {
        commit("setConfiguration", null);
        reject(errors);
      })
    })
  },
  queryDbUpdateStatus({ commit }) {
    return new Promise((resolve, reject) => {
      Configuration.needsDbUpdate()
      .then((updateInfo) => {
        commit("setDbUpdateStatus", updateInfo);
        resolve(updateInfo);
      })
      .catch((errors) => {
        commit("setDbUpdateStatus", null);
        reject(errors);
      })
    });
  },
  queryDbVersionInfo({ commit }) {
    return new Promise((resolve, reject) => {
      Configuration.getDbVersion()
        .then((info) => {
          commit("setDbVersionInfo", info);
          resolve(info);
        })
        .catch((errors) => {
          commit("setDbVersionInfo", null);
          reject(errors);
        });
    });
  },
  queryRecaptchaKey({ commit }) {
    return new Promise((resolve, reject) => {
      Configuration.getRecaptchaKey()
        .then((response) => {
          // load latest configuration
          commit("setRecaptchaKey", response);
          resolve();
        })
        .catch((errors) => {
          reject(errors);
        });
    });
  },
  queryLogoFile({ commit }) {
    return new Promise((resolve, reject) => {
      Configuration.getLogoFile()
        .then((response) => {
          // load latest configuration
          commit("setLogoFile", response.uri);
          resolve(response.uri);
        })
        .catch((errors) => {
          reject(errors);
        });
    });
  },
  updateDb({ commit, dispatch }) {
    return new Promise((resolve, reject) => {
      commit("setIsUpdating", true);
      Configuration.updateDb()
      .then(() => {
        dispatch("queryDbVersionInfo");
        commit("setIsUpdating", false);
        commit("setUpdateResult", true);
      })
      .catch((errors) => {
        console.log("updateDb error:", errors);
        commit("setIsUpdating", false);
        commit("setUpdateResult ", { ok: false, msg: errors[0] });
      })
    })
  },
  setConfiguration({ dispatch, commit }, configurationValues) {
    return new Promise((resolve, reject) => {
      Configuration.setConfiguration(configurationValues)
        .then(() => {
          // load latest configuration
          commit("invalidateConfiguration", {});
          dispatch("queryConfiguration", {});
          resolve();
        })
        .catch((errors) => {
          reject(errors);
        });
    });
  },
};

const mutations = {
  invalidateConfiguration(state) {
    state.CONFIG_LOADED = false;
  },
  setConfiguration(state, value) {
    state.configuration = value;
    state.locationAddress =
      value.location_address != null
        ? value.location_address
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/\n/g, "<br>")
        : null;
    state.locationMap = value.location_map;
    state.currency = value.currency;
    state.logoFile = value.logo_file;
    state.engineHourFormat = value.engine_hour_format;
    state.fuelPaymentType = value.fuel_payment_type;
    state.myNautiqueEnabled = value.mynautique_enabled;
    state.myNautiqueBoatId = value.mynautique_boat_id;

    state.CONFIG_LOADED = true;
    console.log("configuration set to", value);
  },
  setRecaptchaKey(state, response) {
    state.recaptchaKey = response.key;
  },
  setLogoFile(state, value) {
    state.logoFile = value;
  },
  setDbUpdateStatus(state, value) {
    console.log(value);
    if(value != null){
      state.dbUpdateStatus = value.updateAvailable;
    } else {
      state.dbUpdateStatus = null;
    }
  },
  setDbVersionInfo(state, value) {
    if (value == null) {
      // set to null by default
      console.log("setDbVersionInfo to null");
      state.dbVersionInfo = null;
      return;
    }

    console.log("start updating");
    if (state.dbVersionInfo == null) {
      state.dbVersionInfo = {
        isUpdated: true,
      };
    }
    state.dbVersionInfo.appVersion =
      value.app_version != null ? parseFloat(value.app_version) : null;
    state.dbVersionInfo.dbVersion =
      value.db_version != null ? parseFloat(value.db_version) : null;
    state.dbVersionInfo.dbVersionRequired =
      value.db_version_required != null
        ? parseFloat(value.db_version_required)
        : null;

    // decide if is updated
    if (
      state.dbVersionInfo.dbVersion != state.dbVersionInfo.dbVersionRequired
    ) {
      state.dbVersionInfo.isUpdated = false;
    } else {
      state.dbVersionInfo.isUpdated = true;
    }
    console.log("setDbVersionInfo updated to:", state.dbVersionInfo);
  },
  setIsUpdating(state, value) {
    state.dbIsUpdating = value;
  },
  setUpdateResult(state, success, errorMessage="") {
    if (success == true) {
      state.dbUpdateResult = {
        success: true,
        msg: "Database upgrade was successful",
        queries: value.queries,
      };
    } else {
      state.dbUpdateResult = {
        success: false,
        msg: errorMessage,
      };
    }
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
