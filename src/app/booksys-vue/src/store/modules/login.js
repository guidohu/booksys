// initial state
// shape: [{ id, quantity }]
const state = () => ({
  username: '',
  role: ''
})

// getters
const getters = {
  username: (state) => {
    return state.username
  }
}

// actions
const actions = {
  doLogin () {
    console.log('do login')
  },
  setUsername ({ commit }, value) {
    commit('setUsername', value)
  }
}

// mutations
const mutations = {
  setUsername (state, message) {
    state.username = message
    console.log('username: ', state.username)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
