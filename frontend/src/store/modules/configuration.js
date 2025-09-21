import Configuration from "booksys/api/configuration";

const state = () => ({
  ADMIN_CONFIG_LOADED: false,
  CONFIG_LOADED: false,
  configuration: null,
  adminConfiguration: null,
  recaptchaKey: null,
  currency: null,
  locationAddress: null,
  locationMap: null,
  timezone: "Europe/Zurich",
  maxRiders: 12,
  logoFile: null,
  logoUri: null,
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
  getAdminConfiguration: (state) => {
    return state.adminConfiguration;
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
      console.debug(
        "Configuration is already loaded. Skip query for configuration on backend.",
      );
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
        });
    });
  },
  queryAdminConfiguration({ commit, state }) {
    if (state.ADMIN_CONFIG_LOADED == true) {
      console.debug(
        "Admin Configuration is already loaded. Skip query for configuration on backend.",
      );
      return;
    }

    return new Promise((resolve, reject) => {
      Configuration.getAdminConfiguration()
        .then((config) => {
          commit("setAdminConfiguration", config);
          resolve(config);
        })
        .catch((errors) => {
          commit("setAdminConfiguration", null);
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
  setConfiguration({ dispatch, commit }, configurationValues) {
    return new Promise((resolve, reject) => {
      Configuration.setConfiguration(configurationValues)
        .then(() => {
          // load latest configuration
          commit("invalidateConfiguration", {});
          dispatch("queryAdminConfiguration", {});
          dispatch("queryConfiguration", {});
          dispatch("queryLogoFile", {});
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
    state.ADMIN_CONFIG_LOADED = false;
    state.CONFIG_LOADED = false;
    state.configuration = null;
    state.adminConfiguration = null;
  },
  setConfiguration(state, value) {
    state.CONFIG_LOADED = true;
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
    console.log("configuration set to", value);
  },
  setAdminConfiguration(state, value) {
    state.ADMIN_CONFIG_LOADED = true;
    state.adminConfiguration = value;
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
    state.url = value.url;
    console.log("admin configuration set to", value);
  },
  setRecaptchaKey(state, response) {
    state.recaptchaKey = response.key;
  },
  setLogoFile(state, value) {
    if (value == "") {
      state.logoUri = null;
      return;
    }
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
