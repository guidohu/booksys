import { Configuration } from "../../api/configuration";

const state = () => ({
  configuration: null,
  locationAddress: null,
  locationMap: null,
  dbUpdateStatus: null,
  dbVersionInfo: null,
})

const getters = {
  getConfiguration: (state) => {
    return state.configuration
  },
  getLocationAddress: (state) => {
    return state.locationAddress
  },
  getLocationMap: (state) => {
    return state.locationMap
  },
  getDbUpdateStatus: (state) => {
    return state.dbUpdateStatus
  },
  getDbVersionInfo: (state) => {
    return state.dbVersionInfo
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
  }
}

const mutations = {
  setConfiguration (state, value){
    state.configuration = value
    state.locationAddress = (value.location_address != null) ? value.location_address.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/\n/g, '<br>') : null;
    state.locationMap = value.location_map_iframe
    console.log('configuration set to', value)
  },
  setDbUpdateStatus (state, value){
    console.log(value)
    state.dbUpdateStatus = value
  },
  setDbVersionInfo (state, value){
    state.dbVersionInfo = value
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}