import { Configuration } from "../../api/configuration";

const state = () => ({
  configuration: null,
  locationAddress: null,
  locationMap: null,
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
  }
}

const mutations = {
  setConfiguration (state, value){
    state.configuration = value
    state.locationAddress = (value.location_address != null) ? value.location_address.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/\n/g, '<br>') : null;
    state.locationMap = value.location_map_iframe
    console.log('configuration set to', value)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}