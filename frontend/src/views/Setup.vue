<template>
  <b-modal
    id="setupModal"
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
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <warning-box :errors="errors"/>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
      <!-- Database setup -->
      <database-configuration v-if="showDbSetup" :dbConfig="dbConfig" @save="setDbSettings"/>
      <!-- Administrator setup -->
      <div v-if="showUserSetup">
        <b-row class="mb-4">
          <b-col cols="1" class="d-none d-sm-block"></b-col>
          <b-col cols="12" sm="10">
            Setup the Administrator Account
          </b-col>
          <b-col cols="1" class="d-none d-sm-block"></b-col>
        </b-row>
        <user-sign-up
          :userData="adminUserConfig"
          :showDisclaimer="false"
          @save="addAdminUser"
        />
      </div>
      <!-- Setup Done -->
      <b-row v-if="showSetupDone" class="text-center">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <p class="h4 mb-2">
            <b-icon icon="check-circle" variant="success"/>
            Setup Done.
          </p>
          <p>
            Please go back to the login page and login.
          </p>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-overlay>
    <!-- Footer -->
    <div slot="modal-footer">
      <b-button v-if="showDbSetup || showUserSetup" class="mr-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x></b-icon-x>
        Cancel
      </b-button>
      <b-button v-if="showDbSetup" type="button" variant="outline-info" v-on:click="setDbSettings" :disabled="isLoading">
        <b-icon icon="arrow-right"/>
        Next
      </b-button>
      <b-button v-if="showUserSetup" type="button" variant="outline-info" v-on:click="addAdminUser" :disabled="isLoading">
        <b-icon icon="arrow-right"/>
        Next
      </b-button>
      <b-button v-if="showSetupDone" class="mr-1" type="button" variant="outline-info" v-on:click="close">
        <b-icon icon="check"/>
        Done
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import Backend from '@/api/backend';
import { Configuration } from '@/api/configuration';
import DatabaseConfiguration from '@/components/DatabaseConfiguration';
import User from '@/api/user';
import WarningBox from '@/components/WarningBox';
import UserSignUp from '@/components/forms/UserSignUp';

export default Vue.extend({
  name: "Setup",
  components: {
    DatabaseConfiguration,
    WarningBox,
    UserSignUp
  },
  data(){
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
      setupDoneTitle: "Setup Done"
    }
  },
  methods: {
    setDbSettings: function(){
      this.isLoading = true;

      Configuration.setDbConfig(this.dbConfig)
      .then(() => {
        this.showDbSetup = false;
        this.showUserSetup = true;
        this.title = this.userSetupTitle;
        this.errors = [];
        this.isLoading = false;
      })
      .catch((errors) => {
        this.errors = errors;
        this.isLoading = false;
      });
    },
    addAdminUser: function(){
      console.log("addAdminUser called");
      
      // check for obvious validation errors
      const errors = this.validateAdminUser();
      if(errors.length > 0){
        this.errors = errors;
        return;
      }

      this.isLoading = true;

      User.signUp(this.adminUserConfig)
      .then((user) => {
        this.errors = [];
        console.log("userCreated:", user);
        this.makeUserAdmin(user.user_id)
      })
      .catch((errors) => {
        this.errors = errors;
        this.isLoading = false;
      });
    },
    makeUserAdmin: function(userId){
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
    close: function() {
      this.$bvModal.hide("setupModal");
      this.$router.push("/login");
    },
    validateAdminUser: function() {
      const errors = [];

      if(this.adminUserConfig.password != this.adminUserConfig.passwordConfirm){
        errors.push("Password and Password Confirmation are not identical.");
      }
      if(this.adminUserConfig.password.length <= 8){
        errors.push("Please use a password longer than 8 characters.");
      }
      if(this.adminUserConfig.recaptchaResponse == null && this.getRecaptchaKey){
        errors.push("Please tick `I'm not a robot`.");
      }

      // password strength
      const pwUpperRegex = /[A-Z]+/;
      const pwLowerRegex = /[a-z]+/;
      const pwDigitRegex = /[0-9]+/;
      if(this.adminUserConfig.password.match(pwUpperRegex) == null){
          errors.push("The password needs to contain at least one upper case letter (A-Z)");
      }
      if(this.adminUserConfig.password.match(pwLowerRegex) == null){
          errors.push("The password needs to contain at least one lower case letter (a-z)");
      }
      if(this.adminUserConfig.password.match(pwDigitRegex) == null){
          errors.push("The password needs to contain at least one digit (0-9)");
      }

      return errors;      
    }
  },
  mounted(){
    this.isLoading = true;

    Backend.getStatus()
    .then(status => {
      if(status.configFile == false){
        // no configuration at all yet
        this.showDbSetup = true;
        this.showUserSetup = false;
        this.title = this.dbSetupTitle;
      }else if(status.configDb == false){
        // no database configuration
        this.showDbSetup = true;
        this.showUserSetup = false;
        this.title = this.dbSetupTitle;
      }else if(status.dbReachable == false){
        // cannot reach database, thus allow to change settings
        this.showDbSetup = true;
        this.showUserSetup = false;
        this.title = this.dbSetupTitle;
      }else if(status.adminExists == false){
        // database is up, but there is no admin user yet
        this.showDbSetup = false;
        this.showUserSetup = true;
        this.title = this.userSetupTitle;
      }else{
        // setup is done
        this.showSetupDone = true;
        this.title = this.setupDoneTitle;
      }
      this.isLoading = false;
    })
    .catch(errors => {
      this.errors = errors;
      this.isLoading = false;
    });
  }

})
</script>