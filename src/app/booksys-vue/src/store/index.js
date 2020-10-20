import Vue from 'vue';
import Vuex from 'vuex';
import login from './modules/login';
import configuration from './modules/configuration';
import sessions from './modules/sessions';
import user from './modules/user';
import boat from './modules/boat';

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    login: login,
    configuration: configuration,
    sessions: sessions,
    user: user,
    boat: boat
  }
})

export default store
