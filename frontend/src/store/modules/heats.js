const Heat = () => import('@/api/heat');

const state = () => ({
  sessionId: null,
  heatsForSession: []
});

const getters = {
  getSessionId: (state) => {
    return state.sessionId;
  },
  getHeatsForSession: (state) => {
    return state.heatsForSession;
  }
};

const actions = {
  queryHeatsForSession({ commit }, sessionId) {
    return new Promise((resolve, reject) => {
      Heat.getHeatsBySession(sessionId)
      .then((heats) => {
        console.log(heats);
        commit('setSessionId', sessionId);
        commit('setHeatsForSession', heats);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  removeHeat({ dispatch, state }, heatId) {
    return new Promise((resolve, reject) => {
      Heat.deleteHeat(heatId)
      .then(() => {
        dispatch('queryHeatsForSession', state.sessionId);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  updateHeat({ dispatch, state }, heatUpdate) {
    return new Promise((resolve, reject) => {
      Heat.updateHeat(heatUpdate)
      .then(() => {
        dispatch('queryHeatsForSession', state.sessionId);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  }
};

const mutations = {
  setSessionId (state, sessionId){
    state.sessionId = sessionId;
  },
  setHeatsForSession (state, heats){
    state.heatsForSession = heats;
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}