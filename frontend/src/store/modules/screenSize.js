import { BooksysBrowser } from "@/libs/browser.js";

// initial state
const state = () => ({
  isMobile: false,
  _handlerPresent: false,
});

// getters
const getters = {
  isMobile: (state) => {
    return state.isMobile;
  },
};

// actions
const actions = {
  initScreenSize({ commit, state, dispatch }) {
    if(! state._handlerPresent){
      state._handlerPresent = true;
      window.addEventListener("resize", function(){
        dispatch("refreshScreenSize");
      });
      commit("setIsMobile", BooksysBrowser.isMobileResponsive());
    }
  },
  refreshScreenSize({ commit }) {
    console.log("Refresh Screen Size");
    commit("setIsMobile", BooksysBrowser.isMobileResponsive());
  },
};

// mutations
const mutations = {
  setIsMobile(state, isMobile) {
    state.isMobile = isMobile;
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
