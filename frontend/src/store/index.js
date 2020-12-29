import Vue from 'vue';
import Vuex from 'vuex';
import login from './modules/login';
import configuration from './modules/configuration';
import sessions from './modules/sessions';
import user from './modules/user';
import boat from './modules/boat';
import stopwatch from './modules/stopwatch';
import heats from './modules/heats';
import log from './modules/log';
import accounting from './modules/accounting';

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    login: login,
    configuration: configuration,
    sessions: sessions,
    user: user,
    boat: boat,
    stopwatch: stopwatch,
    heats: heats,
    log: log,
    accounting: accounting
  }
})

export default store