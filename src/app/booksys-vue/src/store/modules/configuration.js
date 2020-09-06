import { Configuration } from "../../api/configuration";

const state = () => ({
  configuration: null,
  currency: null,
  locationAddress: null,
  locationMap: null,
  dbUpdateStatus: null,
  dbVersionInfo: null,
  dbIsUpdating: false,
  dbUpdateResult: null,
})

const getters = {
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
  }
}

const actions = {
  queryConfiguration({ commit }) {
    let successCb = (config) => {
      commit('setConfiguration', config)
    }
    let failureCb = () => {
      commit('setConfiguration', null)
    }
    Configuration.getCustomizationParameters(successCb, failureCb)
  },
  queryDbUpdateStatus({ commit }) {
    let successCb = (response) => {
      commit('setDbUpdateStatus', response)
    }
    let failureCb = (error) => {
      console.error("store/configuration: Cannot get needsDbUpdate:", error)
      commit('setDbUpdateStatus', null)
    }
    Configuration.needsDbUpdate(successCb, failureCb)
  },
  queryDbVersionInfo({ commit }) {
    let successCb = (response) => {
      commit('setDbVersionInfo', response)
    }
    let failureCb = (error) => {
      console.error("store/configuration: Cannot get getDbVersion:", error)
      commit('setDbVersionInfo', null)
    }
    Configuration.getDbVersion(successCb, failureCb)
  },
  updateDb({ commit, dispatch }) {
    commit('setIsUpdating', true)
    let successCb = (response) => {
      console.log("updateDb response:", response);
      dispatch('queryDbVersionInfo');
      commit('setIsUpdating', false);
      commit('setUpdateResult', response); 
    };
    let failureCb = (error) => {
      console.log("updateDb error:", error);
      commit('setIsUpdating', false);
      commit('setUpdateResult ', { ok: false, msg: error});
    };
    Configuration.updateDb(successCb, failureCb);
  }
}

const mutations = {
  setConfiguration (state, value){
    state.configuration = value
    state.locationAddress = (value.location_address != null) ? value.location_address.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/\n/g, '<br>') : null;
    state.locationMap = value.location_map_iframe;
    state.currency = value.currency;
    console.log('configuration set to', value);
  },
  setDbUpdateStatus (state, value){
    console.log(value)
    state.dbUpdateStatus = value
  },
  setDbVersionInfo (state, value){
    if (value == null) {
      // set to null by default
      console.log("setDbVersionInfo to null")
      state.dbVersionInfo = null
      return
    }

    console.log("start updating")
    if (state.dbVersionInfo == null) {
      state.dbVersionInfo = {
        isUpdated: true
      };
    }
    state.dbVersionInfo.appVersion = value.app_version != null ? parseFloat(value.app_version) : null
    state.dbVersionInfo.dbVersion = value.db_version != null ? parseFloat(value.db_version) : null
    state.dbVersionInfo.dbVersionRequired = value.db_version_required != null ? parseFloat(value.db_version_required) : null

    // decide if is updated
    if (state.dbVersionInfo.dbVersion != state.dbVersionInfo.dbVersionRequired) {
      state.dbVersionInfo.isUpdated = false
    }else{
      state.dbVersionInfo.isUpdated = true
    }
    console.log("setDbVersionInfo updated to:", state.dbVersionInfo)
  },
  setIsUpdating (state, value){
    state.dbIsUpdating = value
  },
  setUpdateResult (state, value){
    if(value.ok == true){
      state.dbUpdateResult = {
        success: true,
        msg: 'Database upgrade was successful',
        queries: value.queries
      };
    }else{
      state.dbUpdateResult = {
        success: false,
        msg: value.msg
      };
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}