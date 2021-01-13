const state = () => ({
  isLoggedIn: null
})

const getters = {
  isLoggedIn: (state) => {
    return state.isLoggedIn;
  }
}

const actions = {
  setIsLoggedIn ({ commit }, value) {
    commit('setIsLoggedIn', value);
  }
}

const mutations = {
  setIsLoggedIn (state, value) {
    state.isLoggedIn = value;
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}