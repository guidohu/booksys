import { createStore } from "vuex";
import loginStatus from "@/store/modules/loginStatus";

export const store = createStore({
  modules: {
    // will be added dynamically in depending on the route in router
    loginStatus: loginStatus,
  },
});

export const loadStoreModules = (moduleNames, callback) => {
  const loaded = [];
  moduleNames.forEach(function (moduleName) {
    if (!store.hasModule(moduleName)) {
      const importPromise = import("@/store/modules/" + moduleName).then(
        (module) => {
          // register the imported module
          if (!store.hasModule(moduleName)) {
            store.registerModule(moduleName, module.default);
          }

          // initialize the modules that need initialization
          if (moduleName == "stopwatch") {
            store.dispatch("stopwatch/init");
          }

          // console.log("Store module registered: " + moduleName);
        }
      );

      loaded.push(importPromise);
    }
  });

  if (moduleNames.length > 0) {
    Promise.all(loaded).then(() => {
      callback();
    });
  } else {
    callback();
  }
};

export default store;
