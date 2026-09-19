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

// The default modules above are statically imported, so exclude them here.
// Otherwise Vite warns that they are both statically and dynamically imported
// and cannot be split into their own chunks.
const allStoreModules = import.meta.glob([
  "./modules/*.js",
  "!./modules/loginStatus.js",
  "!./modules/screenSize.js",
]);

// Function to dynamically load and register specific modules
store.loadModules = async function (moduleNames, to, from) {
  console.log(
    `Attempt to load modules ${moduleNames} while navigating to ${to.path}`,
  );
  for (const moduleName of moduleNames) {
    if (store.hasModule(moduleName)) {
      console.log(`Module '${moduleName}' already registered.`);
      continue;
    }

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
  }
};
