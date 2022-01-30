<template>
  <b-modal
    id="setupModal"
    ref="setupModal"
    :title="title"
    no-close-on-backdrop
    no-close-on-esc
    hide-header-close
    visible
  >
    <b-overlay
      id="overlay-background"
      :show="isLoading"
      spinner-type="border"
      spinner-variant="info"
      rounded="sm"
    >
      <b-row v-if="errors.length > 0" class="mt-4">
        <b-col cols="1" class="d-none d-sm-block" />
        <b-col cols="12" sm="10">
          <warning-box :errors="errors" />
        </b-col>
        <b-col cols="1" class="d-none d-sm-block" />
      </b-row>
      <!-- Database setup -->
      <database-configuration
        v-if="showDbSetup"
        :db-config="dbConfig"
        @save="setDbSettings"
      />
      <!-- Administrator setup -->
      <div v-if="showUserSetup">
        <b-row class="mb-4">
          <b-col cols="1" class="d-none d-sm-block" />
          <b-col cols="12" sm="10"> Setup the Administrator Account </b-col>
          <b-col cols="1" class="d-none d-sm-block" />
        </b-row>
        <user-sign-up
          :user-data="adminUserConfig"
          :show-disclaimer="false"
          @save="addAdminUser"
          @update:user="handleUserUpdate"
        />
      </div>
      <!-- Setup Done -->
      <b-row v-if="showSetupDone" class="text-center">
        <b-col cols="1" class="d-none d-sm-block" />
        <b-col cols="12" sm="10">
          <p class="h4 mb-2">
            <b-icon-check-circle variant="success" />
            Setup Done.
          </p>
          <p>Please go back to the login page and login.</p>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block" />
      </b-row>
    </b-overlay>
    <!-- Footer -->
    <div slot="modal-footer">
      <b-button
        v-if="showDbSetup || showUserSetup"
        class="mr-1"
        type="button"
        variant="outline-danger"
        @click="close"
      >
        <b-icon-x />
        Cancel
      </b-button>
      <b-button
        v-if="showDbSetup"
        type="button"
        variant="outline-info"
        :disabled="isLoading"
        @click="setDbSettings"
      >
        <b-icon-arrow-right />
        Next
      </b-button>
      <b-button
        v-if="showUserSetup"
        type="button"
        variant="outline-info"
        :disabled="isLoading"
        @click="addAdminUser"
      >
        <b-icon-arrow-right />
        Next
      </b-button>
      <b-button
        v-if="showSetupDone"
        class="mr-1"
        type="button"
        variant="outline-info"
        @click="close"
      >
        <b-icon-check />
        Done
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { getBackendStatus } from "@/api/backend";
import Configuration from "@/api/configuration";
import User from "@/api/user";

// Lazy loaded components
const DatabaseConfiguration = () =>
  import(
    /* webpackChunkName: "database-configuration" */ "@/components/DatabaseConfiguration"
  );
const UserSignUp = () =>
  import(
    /* webpackChunkName: "user-sign-up" */ "@/components/forms/UserSignUp"
  );
const WarningBox = () =>
  import(/* webpackChunkName: "warning-box" */ "@/components/WarningBox");

import {
  BModal,
  BOverlay,
  BRow,
  BCol,
  BIconCheckCircle,
  BIconX,
  BIconArrowRight,
  BIconCheck,
  BButton,
} from "bootstrap-vue";

export default {
  name: "WSSetupPage",
  components: {
    DatabaseConfiguration,
    WarningBox,
    UserSignUp,
    BModal,
    BOverlay,
    BRow,
    BCol,
    BIconCheckCircle,
    BIconX,
    BIconArrowRight,
    BIconCheck,
    BButton,
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
    setDbSettings: function () {
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
      this.$refs["setupModal"].hide();
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
