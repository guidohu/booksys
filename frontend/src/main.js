import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
// import { BootstrapVue, IconsPlugin } from 'bootstrap-vue';
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap-vue/dist/bootstrap-vue.css";

Vue.config.productionTip = false;

// Install BootstrapVue
// Vue.use(BootstrapVue)
// Vue.use(IconsPlugin)

// Override bootstrap styling
import "./assets/bootstrap/css/bootstrap-theme-bc.css";
import "./assets/css/style.css";

new Vue({
  router,
  store,
  render: h => h(App),
  beforeCreate() {
    this.$store.dispatch('initialize');
  }
}).$mount("#app");
