const Log = () => import('@/api/log');

const state = () => ({
  logLines: []
});

const getters = {
  getLogLines: (state) => {
    return state.logLines;
  }
}

const actions = {
  queryLogLines({ commit }) {
    console.log("Query logLines");
    return new Promise((resolve, reject) => {
      Log.getLogs()
      .then((lines) => {
        commit('setLogLines', lines);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  }
}

const mutations = {
  setLogLines (state, lines){
    state.logLines = lines;
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}

