// initial state
// shape: [{ id, quantity }]
const state = () => ({
  username: '',
  role: '',
  userInfo: {}
})

// getters
const getters = {
  username: (state) => {
    return state.username
  },
  userInfo: (state) => {
    console.log("Getting userInfo:", state.userInfo)
    return state.userInfo
  }
}

// actions
const actions = {
  doLogin () {
    console.log('do login')
  },
  setUsername ({ commit }, value) {
    commit('setUsername', value)
  },
  setUserInfo ({ commit }, value) {
    commit('setUserInfo', value)
  }
}

// mutations
const mutations = {
  setUsername (state, message) {
    state.username = message
    console.log('username: ', state.username)
  },
  setUserInfo (state, value) {
    state.userInfo = value
    console.log('userInfo set to', value)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
