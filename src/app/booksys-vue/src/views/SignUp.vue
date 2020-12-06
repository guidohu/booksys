<template>
  <b-modal
    id="signUpModal"
    title="Sign Up"
    no-close-on-backdrop
    no-close-on-esc
    hide-header-close
    visible
  >
    <b-row v-if="isSignedUp==false" class="text-left mb-4">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        Please fill the form below to create your own personal account.
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-form v-if="isSignedUp==false" @submit="save">
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <!-- Email -->
          <b-form-group
            id="email"
            label="Email"
            label-for="email-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="email-input"
              v-model="form.email"
              type="text"
            />
          </b-form-group>
          <!-- Password -->
          <b-form-group
            id="password"
            label="Password"
            label-for="password-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="password-input"
              v-model="form.password"
              type="password"
            />
          </b-form-group>
          <!-- Password Confirm -->
          <b-form-group
            id="password-confirm"
            label="Confirm"
            label-for="password-confirm-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="password-confirm-input"
              v-model="form.passwordConfirm"
              type="password"
            />
          </b-form-group>
          <hr/>
          <!-- First Name -->
          <b-form-group
            id="first-name"
            label="First Name"
            label-for="first-name-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="first-name-input"
              v-model="form.firstName"
              type="text"
            />
          </b-form-group>
          <!-- Surname -->
          <b-form-group
            id="last-name"
            label="Surname"
            label-for="last-name-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="last-name-input"
              v-model="form.lastName"
              type="text"
            />
          </b-form-group>
          <!-- Street / No -->
          <b-form-group
            id="street"
            label="Street / No"
            label-for="street-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="street-input"
              v-model="form.street"
              type="text"
            />
          </b-form-group>
          <!-- Zip Code -->
          <b-form-group
            id="zip-name"
            label="Zip Code"
            label-for="zip-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="zip-input"
              v-model="form.zip"
              type="text"
            />
          </b-form-group>
          <!-- Place / City -->
          <b-form-group
            id="city"
            label="Place / City"
            label-for="city-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="city-input"
              v-model="form.city"
              type="text"
            />
          </b-form-group>
          <!-- Phone -->
          <b-form-group
            id="phone"
            label="Phone"
            label-for="phone-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="phone-input"
              v-model="form.phone"
              type="text"
            />
          </b-form-group>
          <hr/>
          <!-- License -->
          <b-row class="text-left align-middle"
          >
            <b-col cols="9">
              I have a driver's license for boats.
            </b-col>
            <b-col cols="3" class="text-right">
              <toggle-button 
                id="license-toggle"
                :value="form.license"
                sync
                @change="licenseToggleHandler"
                color="#17a2b8"
                :width="toggleWidth"
                :labels="{checked: 'Yes', unchecked: 'No'}"
              />
            </b-col>
          </b-row>
          <!-- Own risk -->
          <b-row class="text-left align-middle mt-3"
          >
            <b-col cols="9">
              I know the boat community and I perform the activities on my own risk.
            </b-col>
            <b-col cols="3" class="text-right">
              <toggle-button 
                id="own-risk-toggle"
                :value="form.ownRisk"
                sync
                @change="ownRiskToggleHandler"
                color="#17a2b8"
                :width="toggleWidth"
                :labels="{checked: 'Yes', unchecked: 'No'}"
              />
            </b-col>
          </b-row>
          <!-- Captcha -->
          <b-row v-if="getRecaptchaKey != null" class="mt-3">
            <b-col cols="12">
              <vue-recaptcha :sitekey="getRecaptchaKey" :loadRecaptchaScript="true" @verify="verifiedHandler"></vue-recaptcha>
            </b-col>
          </b-row>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <b-row class="text-center">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <p class="h4 mb-2">
        <b-icon icon="check-circle" variant="success"/>
        <br/>
        You have been signed up successfully.
        </p>
        <p>
          Please wait until we activate your new account.
        </p>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-row v-if="errors.length > 0" class="mt-4">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <warning-box :errors="errors"/>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <div slot="modal-footer">
      <b-button v-if="isSignedUp==true" type="button" variant="outline-info" v-on:click="close">
        <b-icon-person-check></b-icon-person-check>
        Done
      </b-button>
      <b-button v-if="isSignedUp==false" type="button" variant="outline-info" v-on:click="save">
        <b-icon-person-check></b-icon-person-check>
        Sign-Up
      </b-button>
      <b-button v-if="isSignedUp==false" class="ml-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x></b-icon-x>
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import { mapGetters, mapActions } from 'vuex';
import { ToggleButton } from 'vue-js-toggle-button';
import WarningBox from '@/components/WarningBox';
import VueRecaptcha from 'vue-recaptcha';
import User from '@/api/user';

export default Vue.extend({
  name: "SignUp",
  components: {
    ToggleButton,
    WarningBox,
    VueRecaptcha
  },
  data() {
    return {
      errors: [],
      isSignedUp: false,
      toggleWidth: 50,
      form: {
        license: false,
        ownRisk: false,
        recaptchaResponse: null
      }
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getRecaptchaKey'
    ])
  },
  methods: {
    ...mapActions('configuration', [
      'queryRecaptchaKey'
    ]),
    licenseToggleHandler: function(){
      this.form.license = !this.form.license;
    },
    ownRiskToggleHandler: function(){
      this.form.ownRisk = !this.form.ownRisk;
    },
    verifiedHandler: function(response){
      this.form.recaptchaResponse = response;
    },
    validateForm: function() {
      const errors = [];

      if(this.form.ownRisk == false){
        errors.push("To sign-up you need to agree that the activities are performed at your own risk.");
      }
      if(this.form.password != this.form.passwordConfirm){
        errors.push("Password and Password Confirmation are not identical.");
      }
      if(this.form.password.length <= 8){
        errors.push("Please use a password longer than 8 characters.");
      }
      if(this.form.recaptchaResponse == null && this.getRecaptchaKey){
        errors.push("Please tick `I'm not a robot`.");
      }

      // password strength
      const pwUpperRegex = /[A-Z]+/;
      const pwLowerRegex = /[a-z]+/;
      const pwDigitRegex = /[0-9]+/;
      if(this.form.password.match(pwUpperRegex) == null){
          errors.push("The password needs to contain at least one upper case letter (A-Z)");
      }
      if(this.form.password.match(pwLowerRegex) == null){
          errors.push("The password needs to contain at least one lower case letter (a-z)");
      }
      if(this.form.password.match(pwDigitRegex) == null){
          errors.push("The password needs to contain at least one digit (0-9)");
      }

      return errors;
    },
    save: function(event) {
      // prevent default submit event
      if(event != null){
        event.preventDefault();
      }

      this.errors = this.validateForm();
      if(this.errors.length > 0){
        return;
      }

      User.signUp(this.form)
      .then(() => this.isSignedUp = true)
      .catch((errors) => this.errors = errors);
    },
    close: function() {
      this.$router.push("/");
    }
  },
  created() {
    this.queryRecaptchaKey();
  }
})
</script>