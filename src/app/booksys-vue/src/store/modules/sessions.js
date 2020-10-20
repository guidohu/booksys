import Sessions from "../../api/sessions";
import moment from 'moment';

const state = () => ({
  sessions: null,
  sessionsCalendar: null
})

const getters = {
  getSessions: (state) => {
    return state.sessions;
  },
  getSessionsCalendar: (state) => {
    return state.sessionsCalendar;
  }
}

const actions = {
  querySessionsCalendar({ commit }, month) {
    console.log("Trigger querySessionsCalendar for month", month);

    return new Promise((resolve, reject) => {
      Sessions.getSessionsCalendar(month)
      .then((sessions) => {
        commit('setSessionsCalendar', { sessions: sessions });
        resolve();
      })
      .catch(error => {
        commit('setSessionsCalendar', null);
        reject(error);
      })
    })    
  },
  querySessions({ commit, state, rootState }, time) {
    console.log("Trigger querySessions with timespan", time);
    if(time == null && state.sessions == null){
      console.log("vuex/querySessions: no time window provided, do nothing");
      return;
    }

    // do a reload in case we already have sessions loaded
    if(time == null && state.sessions != null){
      time = {
        start: moment(state.sessions.window_start).format(),
        end: moment(state.sessions.window_end).format()
      }
    }

    return new Promise((resolve, reject) => {
      Sessions.getSessions(time.start, time.end, rootState.configuration.timezone)
      .then((sessions) => {
        commit('setSessions', { sessions: sessions, timeWindow: time });
        resolve();
      })
      .catch(error => {
        commit('setSessions', null);
        reject(error);
      })
    })
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
  editSession({ dispatch }, sessionObj) {
    console.log("Trigger editSession action with", sessionObj);
    return new Promise((resolve, reject) => {
      Sessions.editSession(sessionObj)
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
  setSessionsCalendar (state, value){
    state.sessionsCalendar = value.sessions;
    console.log('sessionsCalendar set to', value.sessions);
  },
  setSessions (state, value){
    state.sessions = value.sessions;
    console.log('sessions set to', value);
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}