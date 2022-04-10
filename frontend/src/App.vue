<template>
  <div id="app-content">
    <alert-message
      v-if="backendReachable == false"
      :alert-message="backendNotReachableAlertMsg"
    />
    <router-view id="router-view" v-else />
    <footer class="legal-footer d-none d-lg-block">
      developed 2013-2022 by Guido Hungerbuehler
      <a href="https://github.com/guidohu/booksys">Find me on Github</a>
    </footer>
  </div>
</template>

<script>
import { defineAsyncComponent } from "vue";
import { getBackendStatus } from "@/api/backend";
import { mapActions, mapGetters } from "vuex";

// Lazy imports
const AlertMessage = defineAsyncComponent(() =>
  import("@/components/AlertMessage")
);

export default {
  name: "App",
  components: {
    AlertMessage,
  },
  data: function () {
    return {
      errors: [],
      backendReachable: true,
      backendNotReachableAlertMsg:
        "The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible. You might try to simply refresh the page if you feel lucky.",
    };
  },
  computed: {
    ...mapGetters("loginStatus", ["isLoggedIn"]),
  },
  watch: {
    isLoggedIn: function (newStatus, oldStatus) {
      if (newStatus == false && oldStatus == true) {
        console.log("Logout detected by App.");
        if (
          this.$route.path != "/logout" &&
          this.$route.path != "/setup" &&
          this.$route.path != "/signup"
        ) {
          console.log("Redirect to login page.");
          this.$router.push("/login");
        }
      }
    },
  },
  methods: {
    ...mapActions("screenSize", ["initScreenSize"]),
  },
  created() {
    this.initScreenSize();
  },
  mounted() {
    // short pulse check on the
    // backend to verify whether
    // it is up and configured
    getBackendStatus()
      .then((status) => {
        // if we do not have a configFile for the app -> go to setup
        if (!status.configFile || !status.configDb) {
          console.log("App Setup Done: no");
          if (this.$route.path !== "/setup") {
            this.$router.push("/setup");
          }
        } else if (!status.dbReachable) {
          this.backendReachable = false;
        } else {
          console.log("App Setup Done: yes (no automated forward to setup)");

          // TODO check if user is logged in, if not -> redirect to login page
        }
      })
      .catch((error) => {
        this.errors = error;
        this.backendReachable = false;
        if (this.errors.length > 0) {
          this.backendNotReachableAlertMsg = this.errors[0];
        } else {
          this.backendNotReachableAlertMsg =
            "Unknown error when connecting to the backend.";
        }
      });
  },
};
</script>
