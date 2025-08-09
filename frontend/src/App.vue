<template>
  <div id="app-content">
    <AlertMessage
      v-if="backendReachable == false"
      :alert-message="backendNotReachableAlertMsg"
    />
    <router-view id="router-view" v-else />
    <footer class="legal-footer d-none d-lg-block">
      developed 2013-2025 by Guido Hungerbuehler
      <a href="https://github.com/guidohu/booksys">Find me on Github</a>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import { getBackendStatus } from "booksys/api/backend";
import AlertMessage from "booksys/components/AlertMessage.vue";

const store = useStore();
const router = useRouter();

const errors = ref([]);
const backendReachable = ref(true);
const backendNotReachableAlertMsg = ref(
  "The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible. You might try to simply refresh the page if you feel lucky.",
);

const isLoggedIn = computed(() => store.getters["loginStatus/isLoggedIn"]);

watch(isLoggedIn, (newStatus, oldStatus) => {
  if (newStatus == false && oldStatus == true) {
    console.log("Logout detected by App.");
    if (
      router.currentRoute.value.path != "/logout" &&
      router.currentRoute.value.path != "/setup" &&
      router.currentRoute.value.path != "/signup"
    ) {
      console.log("Redirect to login page.");
      router.push("/login");
    }
  }
});

const initScreenSize = () => store.dispatch("screenSize/initScreenSize");

initScreenSize();

onMounted(() => {
  // short pulse check on the
  // backend to verify whether
  // it is up and configured
  getBackendStatus()
    .then((status) => {
      if (!status.configDb || !status.usersExist) {
        console.log("App Setup Done: no (automated forward to setup)");
        if (router.currentRoute.value.path !== "/setup") {
          router.push("/setup");
        }
      } else if (!status.dbReachable) {
        backendReachable.value = false;
      } else {
        console.log("App Setup Done: yes (no automated forward to setup)");
      }
    })
    .catch((errorArr) => {
      errors.value = errorArr;
      backendReachable.value = false;
      if (errors.value.length > 0) {
        backendNotReachableAlertMsg.value = errors.value[0];
      } else {
        backendNotReachableAlertMsg.value =
          "Unknown error when connecting to the backend.";
      }
    });
});
</script>
