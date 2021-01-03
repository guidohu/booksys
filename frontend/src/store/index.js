import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex)

export const store = new Vuex.Store({
  modules: {
    // will be added dynamically in depending on the route in router
  }
});

export const loadStoreModules = (moduleNames, callback) => {
  const loaded = [];
  moduleNames.forEach(function (moduleName) {
    const importPromise = import('@/store/modules/' + moduleName).then(module => {
      // register the imported module
      store.registerModule(moduleName, module.default);

      // initialize the modules that need initialization
      if(moduleName == "stopwatch"){
        store.dispatch('stopwatch/init');
      }

      console.log("Store module registered: " + moduleName);
    })

    loaded.push(importPromise);
  })

  if(moduleNames.length > 0){
    Promise.all(loaded).then(() => {
      callback();
    })
  }else{
    callback();
  }
}

export default store
