<template>
  <modal-container name="sign-up-modal" :visible="true">
    <modal-header title="Sign Up" />
    <modal-body>
      <div v-if="isSignedUp == false">
        <div class="row">
          <div class="col-12">
            Please fill the form below to create your own personal account.
          </div>
        </div>
        <user-sign-up
          class="mt-2"
          :user-data="form"
          :show-disclaimer="true"
          @update:user="handleUserUpdate"
        />
        <div class="row mt-3" v-if="getRecaptchaKey != null">
          <vue-recaptcha
            :sitekey="getRecaptchaKey"
            :load-recaptcha-script="true"
            @verify="recaptchaVerifiedHandler"
          />
        </div>
      </div>
      <div class="row text-center" v-if="isSignedUp == true">
        <div class="col-12">
          <p class="h4 mb-2">
            <i class="bi bi-check-circle text-success" />
            <br />
            You have been signed up successfully.
          </p>
          <p>Please wait until we activate your new account.</p>
        </div>
      </div>
      <warning-box v-if="errors.length > 0" :errors="errors" />
    </modal-body>
    <modal-footer>
      <button
        v-if="isSignedUp == true"
        type="button"
        class="btn btn-outline-info"
        @click="close"
      >
        <i class="bi bi-person-check" />
        Done
      </button>
      <button
        v-if="isSignedUp == false"
        type="button"
        class="btn btn-outline-info"
        @click="save"
      >
        <i class="bi bi-person-check" />
        Sign-Up
      </button>
      <button
        v-if="isSignedUp == false"
        type="button"
        class="btn btn-outline-dark ml-1"
        @click="close"
      >
        <i class="bi bi-x" />
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import User from "@/api/user";
import { defineAsyncComponent } from "vue";
import ModalContainer from "@/components/bricks/ModalContainer.vue";
import ModalHeader from "@/components/bricks/ModalHeader.vue";
import ModalBody from "@/components/bricks/ModalBody.vue";
import ModalFooter from "@/components/bricks/ModalFooter.vue";
import { VueRecaptcha } from 'vue-recaptcha';

const UserSignUp = defineAsyncComponent(() =>
  import(/* webpackChunkName: "user-sign-up" */ "@/components/forms/UserSignUp")
);

export default {
  name: "WSSignUp",
  components: {
    WarningBox,
    UserSignUp,
    VueRecaptcha,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
  },
  data() {
    return {
      errors: [],
      isSignedUp: false,
      form: {
        license: false,
        ownRisk: false,
        recaptchaResponse: null,
      },
      visible: true,
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
    handleUserUpdate: function (userData) {
      this.form.email = userData.email;
      this.form.password = userData.password;
      this.form.passwordConfirm = userData.passwordConfirm;
      this.form.firstName = userData.firstName;
      this.form.lastName = userData.lastName;
      this.form.street = userData.street;
      this.form.zip = userData.zip;
      this.form.city = userData.city;
      this.form.phone = userData.phone;
      this.form.ownRisk = userData.ownRisk;
      this.form.license = userData.license;
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
    save: function () {
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
