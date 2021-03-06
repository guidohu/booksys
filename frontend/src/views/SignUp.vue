<template>
  <b-modal
    id="signUpModal"
    title="Sign Up"
    no-close-on-backdrop
    no-close-on-esc
    hide-header-close
    visible
  >
    <div v-if="isSignedUp == false">
      <b-row class="text-left mb-4">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          Please fill the form below to create your own personal account.
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
      <user-sign-up :user-data="form" :show-disclaimer="true" />
      <b-row v-if="getRecaptchaKey != null" class="mt-3">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <vue-recaptcha
            :sitekey="getRecaptchaKey"
            :loadRecaptchaScript="true"
            @verify="recaptchaVerifiedHandler"
          ></vue-recaptcha>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </div>
    <!-- Sign up confirmation -->
    <b-row v-if="isSignedUp == true" class="text-center">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <p class="h4 mb-2">
          <b-icon-check-circle variant="success" />
          <br />
          You have been signed up successfully.
        </p>
        <p>Please wait until we activate your new account.</p>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-row v-if="errors.length > 0" class="mt-4">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <warning-box :errors="errors" />
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <div slot="modal-footer">
      <b-button
        v-if="isSignedUp == true"
        type="button"
        variant="outline-info"
        v-on:click="close"
      >
        <b-icon-person-check />
        Done
      </b-button>
      <b-button
        v-if="isSignedUp == false"
        type="button"
        variant="outline-info"
        v-on:click="save"
      >
        <b-icon-person-check />
        Sign-Up
      </b-button>
      <b-button
        v-if="isSignedUp == false"
        class="ml-1"
        type="button"
        variant="outline-danger"
        v-on:click="close"
      >
        <b-icon-x />
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import User from "@/api/user";
import {
  BModal,
  BRow,
  BCol,
  BIconCheckCircle,
  BIconPersonCheck,
  BIconX,
  BButton,
} from "bootstrap-vue";

const VueRecaptcha = () =>
  import(/* webpackChunkName: "vue-recaptcha" */ "vue-recaptcha");
const UserSignUp = () =>
  import(
    /* webpackChunkName: "user-sign-up" */ "@/components/forms/UserSignUp"
  );

export default {
  name: "SignUp",
  components: {
    WarningBox,
    UserSignUp,
    VueRecaptcha,
    BModal,
    BRow,
    BCol,
    BIconCheckCircle,
    BIconPersonCheck,
    BIconX,
    BButton,
  },
  data() {
    return {
      errors: [],
      isSignedUp: false,
      toggleWidth: 50,
      form: {
        license: false,
        ownRisk: false,
        recaptchaResponse: null,
      },
    };
  },
  computed: {
    ...mapGetters("configuration", ["getRecaptchaKey"]),
  },
  methods: {
    ...mapActions("configuration", ["queryRecaptchaKey"]),
    recaptchaVerifiedHandler: function (response) {
      this.form.recaptchaResponse = response;
    },
    validateForm: function () {
      const errors = [];

      if (this.form.ownRisk == false) {
        errors.push(
          "To sign-up you need to agree that the activities are performed at your own risk."
        );
      }
      if (this.form.password != this.form.passwordConfirm) {
        errors.push("Password and Password Confirmation are not identical.");
      }
      if (this.form.password.length <= 8) {
        errors.push("Please use a password longer than 8 characters.");
      }
      if (this.form.recaptchaResponse == null && this.getRecaptchaKey) {
        errors.push("Please tick `I'm not a robot`.");
      }

      // password strength
      const pwUpperRegex = /[A-Z]+/;
      const pwLowerRegex = /[a-z]+/;
      const pwDigitRegex = /[0-9]+/;
      if (this.form.password.match(pwUpperRegex) == null) {
        errors.push(
          "The password needs to contain at least one upper case letter (A-Z)"
        );
      }
      if (this.form.password.match(pwLowerRegex) == null) {
        errors.push(
          "The password needs to contain at least one lower case letter (a-z)"
        );
      }
      if (this.form.password.match(pwDigitRegex) == null) {
        errors.push("The password needs to contain at least one digit (0-9)");
      }

      return errors;
    },
    save: function (event) {
      // prevent default submit event
      if (event != null) {
        event.preventDefault();
      }

      this.errors = this.validateForm();
      if (this.errors.length > 0) {
        return;
      }

      User.signUp(this.form)
        .then(() => (this.isSignedUp = true))
        .catch((errors) => (this.errors = errors));
    },
    close: function () {
      this.$router.push("/");
    },
  },
  created() {
    this.queryRecaptchaKey();
  },
};
</script>
