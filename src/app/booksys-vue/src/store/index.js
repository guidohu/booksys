import Vue from 'vue'
import Vuex from 'vuex'
import login from './modules/login'
import configuration from './modules/configuration'
import sessions from './modules/sessions'
import user from './modules/user'

Vue.use(Vuex)

const store = new Vuex.Store({
  // state: {
  //   // username: 'username from store'
  // },
  // mutations: {
  //   // setUsername (state, message) {
  //   //   state.username = message
  //   //   console.log('username: ', state.username)
  //   // }
  // },
  // actions: {
  // },
  modules: {
    login,
    configuration,
    sessions,
    user
  }
})

export default store
