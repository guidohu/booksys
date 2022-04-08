import { createApp } from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import { store } from "./store";

// Bootstrap
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.esm";
import "bootstrap-icons/font/bootstrap-icons.css";

// Override bootstrap styling
import "@/assets/bootstrap/css/bootstrap-theme-bc.css";
import "@/assets/css/style.css";

// Deactivate normal logging in production
if (process.env.NODE_ENV === "production") {
  console.log = function () {};
}

// Setup App
const app = createApp(App);
app.use(router);
app.use(store);
app.mount("#app");
