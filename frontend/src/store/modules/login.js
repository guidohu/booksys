import ApiLogin from "@/api/login";

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
      commit('setLoginStatus', null);
      dispatch('getUserInfo');
      commit('setIsLoggedIn', true);
      dispatch('loginStatus/setIsLoggedIn', true, { root: true });
    }
    let failCb = (errorMsg) => {
      console.log("Login failed")
      commit('setLoginStatus', errorMsg)
    }
    console.log("try login with:", value)
    ApiLogin.login(value.username, value.password, successCb, failCb)
  },
  logout ({ commit, dispatch }) {
    let successCb = () => {
      commit('setUserInfo', null);
      commit('setIsLoggedIn', false);
      dispatch('loginStatus/setIsLoggedIn', false, { root: true });
    }
    let failCb = (data) => {
      console.error("Logout failed:", data)
    }
    ApiLogin.logout(successCb, failCb)
  },
  getIsLoggedIn({ commit, dispatch }) {
    return new Promise((resolve, reject) => {
      ApiLogin.isLoggedIn()
      .then((status) => {
        commit('setIsLoggedIn', status);
        dispatch('loginStatus/setIsLoggedIn', status, { root: true });
        if(status == true){
          dispatch('getUserInfo');
        }
        resolve(status);
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  getUserInfo ({ commit, dispatch }) {
    return new Promise((resolve, reject) => {
      ApiLogin.getMyUser()
      .then((userInfo) => {
        commit('setUserInfo', userInfo);
        resolve();
      })
      .catch((errors) => {
        console.log("Cannot get user info, will reject with errors", errors);
        commit('setUserInfo', null);
        commit('setIsLoggedIn', false);
        dispatch('loginStatus/setIsLoggedIn', false, { root: true });
        reject(errors);
      })
    })
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
    state.username = message;
  },
  setUserInfo (state, value) {
    state.userInfo = value;
    if(value != null){
      state.role = value.user_role_name;
    }else{
      state.role = null;
    }
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
