import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'

Vue.config.productionTip = false

// Before you create app
Vue.config.devtools = process.env.NODE_ENV === 'development'

new Vue({ // eslint-disable-line no-new
  el: '#app',
  router: router,
  store: store,
  render: h => h(App)
})
// .$mount('#app')

// After you create app
// window.__VUE_DEVTOOLS_GLOBAL_HOOK__.Vue = app.constructor

// then had to add in ./store.js as well.
// Vue.config.devtools = process.env.NODE_ENV === 'development'
