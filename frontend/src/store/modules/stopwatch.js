import moment from 'moment';
import { sprintf } from 'sprintf-js';
import _ from 'lodash';
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
  },
  getSessionId: (state) => {
    return state.sessionId;
  }
};

const actions = {
  init: ({ commit, state, dispatch }) => {
    console.log("Before initState");
    commit('initState');

    // Start watch again
    if(state._isDisplayUpdaterActive){
      console.log("Display updater was active:", state._isDisplayUpdaterActive, dispatch);
      dispatch('displayTimeUpdate');
    }
    console.log("After initState");
  },
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

        // Use current time, or the time when we paused
        const stoppedTime = (state.isPaused) 
          ? moment(state._pausedAt).format("X")
          : moment().format("X");
        
        const startTime = state._startTime;
        console.log("state at time of finish:", state);
        const offset = state._timeOffset;
        console.log("stoppedTime", moment(stoppedTime, "X").format(), stoppedTime);
        console.log("startTime", moment(state._startTime, "X").format(), state._startTime);
        console.log("offset", moment(state._timeOffset, "X").format(), state._timeOffset);
        const duration = stoppedTime - startTime - offset;

        // create a unique ID
        const uid = "u"+ userId + "s" + sessionId + "r" + Math.floor(Math.random() * 2**16);

        const newHeat = {
          uid: uid,
          user_id: userId,
          session_id: sessionId,
          comment: comment,
          duration_s: duration
        };

        commit('setIsDisplayUpdaterActive', false);

        Heat.addHeats([ newHeat ])
        .then((response) => {
          // TODO check if all have been added
          console.log(response);

          // check if the heat that was just added, was added successfully
          const allUids = _.keys(response);
          const failedHeatUids = allUids.filter(id => response[id].ok != true);
          if(failedHeatUids.length != 0){
            dispatch('heats/queryHeatsForSession', state.sessionId, { root: true});
            reject(['Heat could not be added to the backend.']);
          }else{
            commit('setFinish', stoppedTime);
            dispatch('heats/queryHeatsForSession', state.sessionId, { root: true});
            resolve();
          }
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
    console.log("setSessionId with", id);
    console.log("current state is", state);
    state.sessionId = id;
    storeState(state);
  },
  setUserId (state, userId){
    console.log("setUserId with", userId);
    console.log("current state is", state);
    state.userId = userId;
    storeState(state);
  },
  setStartTime (state, value){
    const startTime = moment(value).format("X");
    console.log("set startTime to:", startTime);
    state._startTime = startTime;
    storeState(state);
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
    storeState(state);
  },
  setIsRunning (state, value){
    console.log("setIsRunning with", value);
    state.isRunning = value;
    storeState(state);
  },
  setIsDisplayUpdaterActive (state, value){
    console.log("set display upater active to", value);
    state._isDisplayUpdaterActive = value;
    storeState(state);
  },
  setPause (state, value) {
    console.log("setPause with", value);
    // set the time we paused at in ISO format
    state.isPaused = true;
    state._pausedAt = value;
    state._isDisplayUpdaterActive = false;
    storeState(state);
  },
  setResume (state, value) {
    console.log("setResume with", value);
    // set the time we resumed in ISO format to calculate additional offset
    const pausedAt = moment(state._pausedAt).format("X");
    const resumedAt = moment(value).format("X");
    const additionalOffset = resumedAt - pausedAt;
    state._timeOffset += additionalOffset;
    state._pausedAt = null;
    state.isPaused = false;
    storeState(state);
  },
  setFinish (state) {
    console.log("setFinish with state before", state);
    state.isPaused = false;
    state.isRunning = false;
    state._timeOffset = 0;
    state._pausedAt = null;
    state._startTime = null;
    state.displayTime = "00:00";
    console.log("setFinish with state after", state);
    storeState(state);
  },
  initState (state) {
    const storedState = loadState();
    if(storedState != null){
      // Load state from localStorage
      state.displayTime = storedState.displayTime;
      state.isPaused = storedState.isPaused;
      state.isRunning = storedState.isRunning;
      state.sessionId = storedState.sessionId;
      state.userId = storedState.userId;
      state._bufferedHeats = storedState._bufferedHeats;
      state._isDisplayUpdaterActive = storedState._isDisplayUpdaterActive;
      state._startTime = storedState._startTime;
      state._pausedAt = storedState._pausedAt;
      state._timeOffset = storedState._timeOffset;
    }
  }
};

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}

function storeState(state){
  if(typeof(Storage) !== "undefined"){
    localStorage.setItem("stopwatchState", JSON.stringify(state));
  }else{
    // No permanent session functionality supported
    console.warn("Cannot use permanent store for heats.");
  }
}

function loadState(){
  if(typeof(Storage) !== "undefined"){
    const localState = localStorage.getItem("stopwatchState");
    return JSON.parse(localState);
  }else{
    // No permanent session functionality supported
    console.warn("Cannot use permanent store for heats.");
    return null;
  }
}

