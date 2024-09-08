import Configuration from "@/api/configuration";

const state = () => ({
  CONFIG_LOADED: false,
  configuration: null,
  recaptchaKey: null,
  currency: null,
  locationAddress: null,
  locationMap: null,
  timezone: "Europe/Zurich",
  maxRiders: 12,
  logoFile: null,
  logoUri: "",
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
  getTimezone: (state) => {
    return state.timezone;
  },
  getMaximumNumberOfRiders: (state) => {
    return state.maxRiders;
  },
  getLogoFile: (state) => {
    return state.logoFile;
  },
  getLogoUri: (state) => {
    return state.logoUri;
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
    state.logoUri = value;
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
