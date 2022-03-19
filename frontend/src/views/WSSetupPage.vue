<template>
  <modal-container name="setup-modal" :visible="true">
    <modal-header :title="title"/>
    <modal-body>
      <warning-box v-if="errors.length > 0" class="mt-4" :errors="errors"/>
      <!-- Database setup -->
      <database-configuration
        v-if="showDbSetup"
        :dbconfig="dbConfig"
        @config-change="dbConfigInputHandler"
        @save="setDbSettings"
      />
      <!-- Administrator Account Setup -->
      <div v-if="showUserSetup">
        <div class="row mb-4">
          <div class="col-12">
            Setup the Administrato Account
          </div>
        </div>
        <user-sign-up
          :user-data="adminUserConfig"
          :show-disclaimer="false"
          @save="addAdminUser"
          @update:user="handleUserUpdate"
        />
      </div>
      <!-- Setup Done -->
      <div v-if="showSetupDone" class="row text-center">
        <div class="col-12">
          <p class="h4 mb-2">
            <i class="bi bi-check-circle text-success" />
            Setup Done.
          </p>
          <p>Please go back to the login page and login.</p>
        </div>
      </div>
    </modal-body>
    <modal-footer>
      <button
        v-if="showDbSetup || showUserSetup"
        class="btn btn-outline-danger me-1"
        type="button"
        @click="close"
      >
        <i class="bi bi-x" />
        Cancel
      </button>
      <button
        v-if="showDbSetup"
        class="btn btn-outline-info"
        type="button"
        :disabled="isLoading"
        @click="setDbSettings"
      >
        <i class="bi bi-arrow-right" />
        Next
      </button>
      <button
        v-if="showUserSetup"
        type="button"
        class="btn btn-outline-info"
        :disabled="isLoading"
        @click="addAdminUser"
      >
        <i class="bi bi-arrow-right" />
        Next
      </button>
      <button
        v-if="showSetupDone"
        class="btn btn-outline-info me-1"
        type="button"
        variant="outline-info"
        @click="close"
      >
        <i class="bi bi-check" />
        Done
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { defineAsyncComponent } from "vue"; 
import { getBackendStatus } from "@/api/backend";
import Configuration from "@/api/configuration";
import User from "@/api/user";
import ModalContainer from "../components/bricks/ModalContainer.vue";
import ModalHeader from "../components/bricks/ModalHeader.vue";
import ModalBody from "../components/bricks/ModalBody.vue";
import ModalFooter from "../components/bricks/ModalFooter.vue";

// Lazy loaded components
const DatabaseConfiguration = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "database-configuration" */ "@/components/DatabaseConfiguration"
  )
);
const UserSignUp = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "user-sign-up" */ "@/components/forms/UserSignUp"
  )
);
const WarningBox = defineAsyncComponent(() =>
  import(/* webpackChunkName: "warning-box" */ "@/components/WarningBox")
);

export default {
  name: "WSSetupPage",
  components: {
    DatabaseConfiguration,
    WarningBox,
    UserSignUp,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
  },
  data() {
    return {
      errors: [],
      title: "Setup",
      showDbSetup: false,
      showUserSetup: false,
      showSetupDone: false,
      isLoading: false,
      dbConfig: {},
      adminUserConfig: {},
      dbSetupTitle: "Setup 1/2",
      userSetupTitle: "Setup 2/2",
      setupDoneTitle: "Setup Done",
    };
  },
  mounted() {
    this.isLoading = true;
    this.getBackendStatus();
  },
  methods: {
    dbConfigInputHandler: function(config) {
      this.dbConfig = config;
    },
    setDbSettings: function () {
      console.log("setDB called:", this.dbConfig);
      this.isLoading = true;

      Configuration.setDbConfig(this.dbConfig)
        .then(() => {
          this.isLoading = false;
          this.getBackendStatus();
        })
        .catch((errors) => {
          this.errors = errors;
          this.isLoading = false;
        });
    },
    handleUserUpdate: function (u) {
      this.adminUserConfig.email = u.email;
      this.adminUserConfig.password = u.password;
      this.adminUserConfig.passwordConfirm = u.passwordConfirm;
      this.adminUserConfig.firstName = u.firstName;
      this.adminUserConfig.lastName = u.lastName;
      this.adminUserConfig.street = u.street;
      this.adminUserConfig.zip = u.zip;
      this.adminUserConfig.city = u.city;
      this.adminUserConfig.phone = u.phone;
      this.adminUserConfig.ownRisk = u.ownRisk;
      this.adminUserConfig.license = u.license;
    },
    addAdminUser: function () {
      // check for obvious validation errors
      const errors = this.validateAdminUser();
      if (errors.length > 0) {
        this.errors = errors;
        return;
      }

      this.isLoading = true;

      User.signUp(this.adminUserConfig)
        .then((user) => {
          this.errors = [];
          this.makeUserAdmin(user.user_id);
        })
        .catch((errors) => {
          this.errors = errors;
          this.isLoading = false;
        });
    },
    makeUserAdmin: function (userId) {
      User.makeAdmin(userId)
        .then(() => {
          this.errors = [];
          this.showDbSetup = false;
          this.showUserSetup = false;
          this.showSetupDone = true;
          this.title = this.setupDoneTitle;
          this.isLoading = false;
        })
        .catch((errors) => {
          this.errors = errors;
          this.isLoading = false;
        });
    },
    close: function () {
      console.log("Navigate to login");
      console.log(this.$router);
      this.$router.push("/login");
    },
    validateAdminUser: function () {
      const errors = [];

      if (
        this.adminUserConfig.password != this.adminUserConfig.passwordConfirm
      ) {
        errors.push("Password and Password Confirmation are not identical.");
      }
      if (this.adminUserConfig.password.length <= 8) {
        errors.push("Please use a password longer than 8 characters.");
      }
      if (
        this.adminUserConfig.recaptchaResponse == null &&
        this.getRecaptchaKey
      ) {
        errors.push("Please tick `I'm not a robot`.");
      }

      // password strength
      const pwUpperRegex = /[A-Z]+/;
      const pwLowerRegex = /[a-z]+/;
      const pwDigitRegex = /[0-9]+/;
      if (this.adminUserConfig.password.match(pwUpperRegex) == null) {
        errors.push(
          "The password needs to contain at least one upper case letter (A-Z)"
        );
      }
      if (this.adminUserConfig.password.match(pwLowerRegex) == null) {
        errors.push(
          "The password needs to contain at least one lower case letter (a-z)"
        );
      }
      if (this.adminUserConfig.password.match(pwDigitRegex) == null) {
        errors.push("The password needs to contain at least one digit (0-9)");
      }

      return errors;
    },
    getBackendStatus: function () {
      // reset errors
      this.errors = [];

      getBackendStatus()
        .then((status) => {
          if (status.configFile == false) {
            // no configuration at all yet
            this.showDbSetup = true;
            this.showUserSetup = false;
            this.title = this.dbSetupTitle;
          } else if (status.configDb == false) {
            // no database configuration
            this.showDbSetup = true;
            this.showUserSetup = false;
            this.title = this.dbSetupTitle;
          } else if (status.dbReachable == false) {
            // cannot reach database, thus allow to change settings
            this.showDbSetup = true;
            this.showUserSetup = false;
            this.title = this.dbSetupTitle;
          } else if (status.adminExists == false) {
            // database is up, but there is no admin user yet
            this.showDbSetup = false;
            this.showUserSetup = false;
            this.showUserSetup = true;
            this.title = this.userSetupTitle;
          } else {
            // setup is done

            this.showDbSetup = false;
            this.showSetupDone = true;
            this.title = this.setupDoneTitle;
          }
          this.isLoading = false;
        })
        .catch((errors) => {
          this.errors = errors;
          this.isLoading = false;
        });
    },
  },
};
</script>
