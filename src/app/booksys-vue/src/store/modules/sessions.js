import Sessions from "../../api/sessions"

const state = () => ({
  sessions: null
})

const getters = {
  getSessions: (state) => {
    return state.sessions
  }
}

const actions = {
  querySessions({ commit }, time) {
    console.log("Trigger action with", time.start, time.end)
    let successCb = (sessions) => {
      console.log("got sessions:", sessions)
      commit('setSessions', sessions)
    }
    let failureCb = () => {
      commit('setSessions', null)
    }
    Sessions.getSessions(time.start, time.end, successCb, failureCb)
  }
}

const mutations = {
  setSessions (state, value){
    state.sessions = value
    console.log('sessions set to', value)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}