import moment from 'moment';
import { sprintf } from 'sprintf-js';
import Heat from '@/api/heat';

const state = () => ({
  sessionId: null,
  userId: null,
  isRunning: false,
  isPaused: false,
  displayTime: "00:00",
  _startTime: 0,
  _timeOffset: 0,
  _isDisplayUpdaterActive: false,
  _pausedAt: null,
  _bufferedHeats: []
});

const getters = {
  getIsRunning: (state) => {
    return state.isRunning;
  },
  getIsPaused: (state) => {
    return state.isPaused;
  },
  getDisplayTime: (state) => {
    return state.displayTime;
  },
  getUserId: (state) => {
    return state.userId;
  }
};

const actions = {
  startTakingTime: ({ dispatch, commit, state }) => {
    console.log('startTakingTime called');
    if(state.isRunning == false){
      console.log("Start taking time");
      commit('setStartTime', moment().format());
      commit('setIsRunning', true);
      dispatch('displayTimeUpdate', { start: true });
    }
  },
  pauseTakingTime: ({ commit }) => {
    console.log('pauseTakingTime called');
    commit('setIsDisplayUpdaterActive', false);
    commit('setPause', moment().format());
  },
  resumeTakingTime: ({ commit, dispatch, state }) => {
    console.log('resumeTakingTime called');
    if(state.isRunning == true){
      console.log("Resume taking time");
      commit('setResume', moment().format());
      dispatch('displayTimeUpdate', { resume: true });
    }
  },
  finishTakingTime: ({ commit, state, dispatch }) => {
    console.log('finishTakingTime called');
    return new Promise((resolve, reject) => {
      if(state.isRunning == true){
        const userId = state.userId;
        const sessionId = state.sessionId;
        const comment = state.comment; // TODO comment is not given yet
        const stoppedTime = moment().format();
        const duration = moment(stoppedTime).format("X") - state._startTime - state._timeOffset;

        const newHeat = {
          user_id: userId,
          session_id: sessionId,
          comment: comment,
          duration_s: duration
        };

        // TODO collect the data to submit
        commit('setIsDisplayUpdaterActive', false);
        commit('setFinish', stoppedTime);

        Heat.addHeats([ newHeat ])
        .then((response) => {
          // TODO check if all have been added
          console.log(response)

          // trigger an update for the heats
          dispatch('heats/queryHeatsForSession', state.sessionId, { root: true});

          resolve() // or reject()
        })
        .catch((errors) => {
          // TODO auto retry?

          // trigger an update for the heats
          dispatch('heats/queryHeatsForSession', state.sessionId, { root: true});

          reject(errors);
        })
      }else{
        reject(['currently we are not taking time']);
      }
      
    })
    
  },
  setSessionId: ({ commit }, sessionId) => {
    commit('setSessionId', sessionId);
  },
  setUserId: ({ commit }, userId) => {
    console.log("triggered setUserId with", userId);
    commit('setUserId', userId);
  },
  displayTimeUpdate: ({ dispatch, commit, state }, options) => {
    console.log("displayTimeUpdate called with state:", state);
    // updates the display time every second
    if(state._isDisplayUpdaterActive == true){
      console.log("displayTimeUpdate: called from --> timer");
      commit('updateDisplayTime');
      setTimeout(() => {
        dispatch('displayTimeUpdate');
      }, 1000);
    }else if(options != null && options.start == true && state._isDisplayUpdaterActive == false){
      console.log("displayTimeUpdate: called from --> start");
      commit('updateDisplayTime');
      commit('setIsDisplayUpdaterActive', true);
      setTimeout(() => {
        dispatch('displayTimeUpdate');
      }, 1000);
    }else if(options != null && options.resume == true && state._isDisplayUpdaterActive == false){
      console.log("displayTimeUpdate: called from --> resume");
      commit('updateDisplayTime');
      commit('setIsDisplayUpdaterActive', true);
      setTimeout(() => {
        dispatch('displayTimeUpdate');
      }, 1000);
    }
  }
};

const mutations = {
  setSessionId (state, id){
    state.sessionId = id;
  },
  setUserId (state, userId){
    state.userId = userId;
  },
  setStartTime (state, value){
    const startTime = moment(value).format("X");
    console.log("set startTime to:", startTime);
    state._startTime = startTime;
  },
  updateDisplayTime (state){
    const currentTime = moment().format("X");
    const startTime = state._startTime;
    const offset = state._timeOffset;
    const timePassed = currentTime - startTime - offset;

    const seconds = timePassed % 60;
    const minutes = Math.floor(((timePassed - seconds) % 3600) / 60);
    const hours   = Math.floor((timePassed - (seconds*60) - (hours*3600)) / 3600);

    if(hours > 0){
      state.displayTime = sprintf('%02d:%02d:%02d', hours, minutes, seconds);
    }else{
      state.displayTime = sprintf('%02d:%02d', minutes, seconds);
    }
    console.log("updated display to:", state.displayTime);
  },
  setIsRunning (state, value){
    state.isRunning = value;
  },
  setIsDisplayUpdaterActive (state, value){
    console.log("set display upater active to", value);
    state._isDisplayUpdaterActive = value;
  },
  setPause (state, value) {
    // set the time we paused at in ISO format
    state.isPaused = true;
    state._pausedAt = value;
  },
  setResume (state, value) {
    // set the time we resumed in ISO format to calculate additional offset
    const pausedAt = moment(state._pausedAt).format("X");
    const resumedAt = moment(value).format("X");
    const additionalOffset = resumedAt - pausedAt;
    state._timeOffset += additionalOffset;
    state._pausedAt = null;
    state.isPaused = false;
  },
  setFinish (state) {
    state.isPaused = false;
    state.isRunning = false;
    state._timeOffset = 0;
    state._pausedAt = null;
    state._startTime = null;
    state.displayTime = "00:00";
  }
};

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}

