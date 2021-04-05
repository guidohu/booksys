import Vue from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";

// Bootstrap
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap-vue/dist/bootstrap-vue.css";

// Override bootstrap styling
import "@/assets/bootstrap/css/bootstrap-theme-bc.css";
import "@/assets/css/style.css";

// Deactivate normal logging in production
if (process.env.NODE_ENV === "production") {
  console.log = function () {};
}

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
