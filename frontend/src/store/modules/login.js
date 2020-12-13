import { Login as ApiLogin } from "../../api/login";

// initial state
const state = () => ({
  username: '',
  role: null,
  userInfo: null,
  loginStatus: null,
  isLoggedIn: null,
})

// getters
const getters = {
  username: (state) => {
    return state.username
  },
  userInfo: (state) => {
    console.log("Getting userInfo:", state.userInfo)
    return state.userInfo
  },
  role: (state) => {
    return state.role;
  },
  loginStatus: (state) => {
    console.log("Getting loginStatus:", state.loginStatus)
    return state.loginStatus
  },
  isLoggedIn: (state) => {
    return state.isLoggedIn
  }
}

// actions
const actions = {
  login ({ commit, dispatch }, value) {
    // store the username
    dispatch('setUsername', value.username)

    console.log('do login')
    let successCb = () => {
      console.log("Login successful")
      commit('setLoginStatus', null)
      dispatch('getUserInfo')
      commit('setIsLoggedIn', true)
    }
    let failCb = (errorMsg) => {
      console.log("Login failed")
      commit('setLoginStatus', errorMsg)
    }
    console.log("try login with:", value)
    ApiLogin.login(value.username, value.password, successCb, failCb)
  },
  logout ({ commit }) {
    let successCb = () => {
      commit('setUserInfo', null)
      commit('setIsLoggedIn', false)
    }
    let failCb = (data) => {
      console.error("Logout failed:", data)
    }
    ApiLogin.logout(successCb, failCb)
  },
  getIsLoggedIn({ commit, dispatch }) {
    let successCb = () => {
      commit('setIsLoggedIn', true)
      dispatch('getUserInfo')
    }
    let failCb = () => {
      commit('setIsLoggedIn', false)
    }
    ApiLogin.isLoggedIn(successCb, failCb)
  },
  getUserInfo ({ commit }) {
    let successCb = (data) => {
      console.log("logged in userInfo:", data)
      commit('setUserInfo', data)
    }
    let failCb = (errorData) => {
      console.log("issue to get user", errorData)
      commit('setUserInfo', null)
      commit('setIsLoggedIn', false)
    }
    ApiLogin.getMyUser(successCb, failCb)
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
    state.userInfo = value;
    state.role = value.user_role_name;
    console.log('userInfo set to', value);
  },
  setLoginStatus (state, value) {
    state.loginStatus = value
    console.log('loginStatus set to', value)
  },
  setIsLoggedIn (state, value) {
    state.isLoggedIn = value
    console.log('isLoggedIn set to', value)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
