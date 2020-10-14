import Sessions from "../../api/sessions";
import moment from 'moment';

const state = () => ({
  sessions: null
})

const getters = {
  getSessions: (state) => {
    return state.sessions
  }
}

const actions = {
  querySessions({ commit, state }, time) {
    console.log("Trigger querySessions with timespan", time);

    if(time == null && state.sessions == null){
      console.log("vuex/querySessions: no time window provided, do nothing");
      return;
    }

    // do a reload in case we already have sessions loaded
    if(time == null && state.sessions != null){
      time = {
        start: moment(state.sessions.window_start, 'X').format(),
        end: moment(state.sessions.window_end, 'X').format()
      }
    }

    let successCb = (sessions) => {
      console.log("got sessions:", sessions)
      commit('setSessions', { sessions: sessions, timeWindow: time });
    }
    let failureCb = () => {
      commit('setSessions', null)
    }
    Sessions.getSessions(time.start, time.end, successCb, failureCb)
  },
  createSession({ dispatch }, sessionObj) {
    console.log("Trigger createSession action with", sessionObj);
    return new Promise((resolve, reject) => {
      Sessions.createSession(sessionObj)
        .then((response) => {
          dispatch('querySessions');
          resolve(response.data);
        })
        .catch(error => {
          reject(error);
        })
    });
  },
  deleteSession({ dispatch }, sessionObj) {
    console.log("Trigger deleteSession action with", sessionObj);
    return new Promise((resolve ,reject) => {
      Sessions.deleteSession(sessionObj)
        .then(() => {
          dispatch('querySessions');
          resolve();
        })
        .catch(error => {
          reject(error);
        })
    });
  },
  addUsersToSession({ dispatch }, data) {
    console.log("Trigger addUsersToSession action with", data);
    return new Promise((resolve, reject) => {
      Sessions.addUsersToSession(data.sessionId, data.users)
        .then(() => {
          dispatch('querySessions');
          resolve();
        })
        .catch(error => {
          reject(error);
        })
    });
  },
  deleteUserFromSession({ dispatch }, data) {
    console.log("Trigger deleteUserFromSession action with", data);
    return new Promise((resolve, reject) => {
      Sessions.deleteUserFromSession(data.sessionId, data.user)
        .then(() => {
          dispatch('querySessions');
          resolve();
        })
        .catch(error => {
          reject(error);
        })
    });
  },
}

const mutations = {
  setLoadedTime (state, value){
    state.loadedTime = {
      start: value.start,
      end: value.end
    };
  },
  setSessions (state, value){
    state.sessions = value.sessions;
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