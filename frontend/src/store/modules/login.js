import ApiLogin from "@/api/login";

// initial state
const state = () => ({
  username: "",
  role: null,
  userInfo: null,
  isLoggedIn: null,
});

// getters
const getters = {
  username: (state) => {
    return state.username;
  },
  userInfo: (state) => {
    console.log("Getting userInfo:", state.userInfo);
    return state.userInfo;
  },
  role: (state) => {
    return state.role;
  },
  isLoggedIn: (state) => {
    return state.isLoggedIn;
  },
};

// actions
const actions = {
  login({ commit, dispatch }, value) {
    dispatch("setUsername", value.username);

    return new Promise((resolve, reject) => {
      ApiLogin.login(value.username, value.password)
        .then(() => {
          dispatch("getUserInfo");
          commit("setIsLoggedIn", true);
          dispatch("loginStatus/setIsLoggedIn", true, { root: true });
          dispatch("configuration/invalidateConfiguration", true, {
            root: true,
          });
          resolve();
        })
        .catch((errors) => {
          reject(errors);
        });
    });
  },
  logout({ commit, dispatch }) {
    return new Promise((resolve, reject) => {
      ApiLogin.logout()
      .then(() => {
        commit("setUserInfo", null);
        commit("setIsLoggedIn", false);
        dispatch("loginStatus/setIsLoggedIn", false, { root: true });
        resolve();
        console.log("Logout successful");
      })
      .catch((errors) => {
        console.error("Logout failed:", errors);
        reject(errors);
      })
    })
  },
  getIsLoggedIn({ commit, dispatch }) {
    return new Promise((resolve, reject) => {
      ApiLogin.isLoggedIn()
        .then((status) => {
          commit("setIsLoggedIn", status.loggedIn);
          dispatch("loginStatus/setIsLoggedIn", status.loggedIn, { root: true });
          if (status.loggedIn == true) {
            dispatch("getUserInfo");
          }
          resolve(status.loggedIn);
        })
        .catch((errors) => {
          reject(errors);
        });
    });
  },
  getUserInfo({ commit, dispatch }) {
    return new Promise((resolve, reject) => {
      ApiLogin.getMyUser()
        .then((userInfo) => {
          commit("setUserInfo", userInfo);
          resolve();
        })
        .catch((errors) => {
          console.log("Cannot get user info, will reject with errors", errors);
          commit("setUserInfo", null);
          commit("setIsLoggedIn", false);
          dispatch("loginStatus/setIsLoggedIn", false, { root: true });
          reject(errors);
        });
    });
  },
  setUsername({ commit }, value) {
    commit("setUsername", value);
  },
  setUserInfo({ commit }, value) {
    commit("setUserInfo", value);
  },
};

// mutations
const mutations = {
  setUsername(state, message) {
    state.username = message;
  },
  setUserInfo(state, value) {
    state.userInfo = value;
    if (value != null) {
      state.role = value.user_role_name;
    } else {
      state.role = null;
    }
    console.log("userInfo set to", value);
  },
  setIsLoggedIn(state, value) {
    state.isLoggedIn = value;
    console.log("isLoggedIn set to", value);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
