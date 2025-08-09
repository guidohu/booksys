import { createStore } from "vuex";
import { nextTick } from "vue";
import loginStatus from "booksys/store/modules/loginStatus";
import screenSize from "booksys/store/modules/screenSize";

export const store = createStore({
  modules: {
    // Modules will be added dynamically depending on the route
    // in router. Only a few default modules are defined here.
    loginStatus: loginStatus,
    screenSize: screenSize,
  },
});

const allStoreModules = import.meta.glob("./modules/*.js"); // Adjust path to be relative to index.js

// Function to dynamically load and register specific modules
store.loadModules = async function (moduleNames, to, from) {
  console.log(
    `Attempt to load modules ${moduleNames} while navigating to ${to.path}`,
  );
  for (const moduleName of moduleNames) {
    // Construct the expected path within the glob result
    // The key will be relative to the project root for aliases, or relative to the current file for relative paths
    let modulePath;
    if (allStoreModules[`./modules/${moduleName}.js`]) {
      // If using relative path glob
      modulePath = `./modules/${moduleName}.js`;
    } else if (allStoreModules[`booksys/store/modules/${moduleName}.js`]) {
      // If using 'booksys' alias glob
      modulePath = `booksys/store/modules/${moduleName}.js`;
    } else {
      console.warn(`Module '${moduleName}' not found in glob pattern.`);
      continue;
    }

    if (!store.hasModule(moduleName) && allStoreModules[modulePath]) {
      try {
        const importModule = allStoreModules[modulePath];
        const module = await importModule(); // Call the dynamic import function
        store.registerModule(moduleName, module.default || module);
        console.log(`Module '${moduleName}' registered.`);
        // Ensure reactivity updates if needed
        await nextTick(); // Allows Vue to process reactivity updates
      } catch (error) {
        console.error(`Failed to load module '${moduleName}':`, error);
      }
    } else if (store.hasModule(moduleName)) {
      console.log(`Module '${moduleName}' already registered.`);
    }
  }
};

export default store;
