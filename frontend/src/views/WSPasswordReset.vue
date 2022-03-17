<template>
  <modal-container name="password-reset-modal" :visible="true">
    <modal-header title="Password Reset" />
    <modal-body>
      <div v-if="showEmailDialog == true">
        <div class="row text-left mb-4">
          <div class="col-12">
            Please enter your email address to get the password reset
            information sent to you by email.
          </div>
        </div>
        <form @submit.prevent="requestToken">
          <div class="row text-left">
            <div class="col-12">
              <!-- Email -->
              <input-text
                id="email"
                label="Email"
                v-model="form.email"
                size="small"
                autocomplete="username"
              />
              <!-- Captcha -->
              <div v-if="getRecaptchaKey != null" class="row mt-3">
                <div class="col-12">
                  <vue-recaptcha
                    :sitekey="getRecaptchaKey"
                    :load-recaptcha-script="true"
                    size="compact"
                    @verify="verifiedHandler"
                  />
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>
      <div v-if="showTokenDialog == true">
        <div class="row text-left mb-4">
          <div class="col-12">
            Please enter the token that has been sent to your email address and
            choose new password.
          </div>
        </div>
        <form @submit.prevent="setPassword">
          <div class="row text-left">
            <div class="col-12">
              <!-- Email -->
              <input-text
                id="email"
                label="Email"
                v-model="form.email"
                size="small"
                autocomplete="username"
              />
              <!-- Token -->
              <input-text
                id="token"
                label="Token"
                v-model="form.token"
                size="small"
              />
              <!-- Password -->
              <input-password
                id="password"
                label="Password"
                v-model="form.password"
                size="small"
                autocomplete="new-password"
              />
              <input-password
                id="password-confirm"
                label="Confirm"
                v-model="form.passwordConfirm"
                size="small"
                autocomplete="new-password"
              />
            </div>
          </div>
        </form>
      </div>
      <div v-if="showSuccessInfo" class="row text-center">
        <div class="col-12">
          <p class="h4 mb-2">
            <i class="bi bi-check-circle text-success" />
            <br />
            Password has been changed successfully.
          </p>
          <p>Please login with your new password.</p>
        </div>
      </div>
      <warning-box v-if="errors.length > 0" :errors="errors" />
    </modal-body>
    <modal-footer>
      <button
        v-if="showEmailDialog"
        class="btn btn-outline-danger me-1"
        type="button"
        @click="close"
      >
        <i class="bi bi-x" />
        Cancel
      </button>
      <button
        v-if="showEmailDialog"
        type="button"
        class="btn btn-outline-info"
        :disabled="isLoading"
        @click="requestToken"
      >
        <i class="bi bi-arrow-right" />
        Next
      </button>
      <button
        v-if="showTokenDialog"
        class="btn btn-outline-danger me-1"
        type="button"
        @click="showEmail"
      >
        <i class="bi bi-arrow-left" />
        Back
      </button>
      <button
        v-if="showTokenDialog"
        type="button"
        class="btn btn-outline-info"
        :disabled="isLoading"
        @click="setPassword"
      >
        <i class="bi bi-arrow-check" />
        Set Password
      </button>
      <button
        v-if="showSuccessInfo"
        type="button"
        class="btn btn-outline-info"
        @click="close"
      >
        <i class="bi bi-arrow-check" />
        Done
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import User from "@/api/user";
import WarningBox from "@/components/WarningBox";
import ModalContainer from "../components/bricks/ModalContainer.vue";
import ModalHeader from "../components/bricks/ModalHeader.vue";
import ModalBody from "../components/bricks/ModalBody.vue";
import ModalFooter from "../components/bricks/ModalFooter.vue";
import InputText from "../components/forms/inputs/InputText.vue";
import InputPassword from "../components/forms/inputs/InputPassword.vue";

const VueRecaptcha = () =>
  import(/* webpackChunkName: "vue-recaptcha" */ "vue-recaptcha");

export default {
  name: "WSPasswordReset",
  components: {
    VueRecaptcha,
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputText,
    InputPassword,
  },
  data() {
    return {
      isLoading: false,
      showEmailDialog: true,
      showTokenDialog: false,
      showSuccessInfo: false,
      errors: [],
      form: {
        email: null,
        recaptchaResponse: null,
      },
    };
  },
  computed: {
    ...mapGetters("configuration", ["getRecaptchaKey"]),
  },
  methods: {
    ...mapActions("configuration", ["queryRecaptchaKey"]),
    verifiedHandler: function (response) {
      this.form.recaptchaResponse = response;
    },
    validatePassword: function () {
      const errors = [];

      if (this.form.password != this.form.passwordConfirm) {
        errors.push("Password and Password Confirmation are not identical.");
      }
      if (this.form.password.length <= 8) {
        errors.push("Please use a password longer than 8 characters.");
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
    requestToken: function (event) {
      if (event != null) {
        event.preventDefault();
      }

      this.isLoading = true;

      const request = {
        email: this.form.email,
        recaptchaResponse: this.form.recaptchaResponse,
      };

      User.requestPasswordResetToken(request)
        .then(() => {
          this.isLoading = false;
          this.errors = [];
          this.showSuccessInfo = false;
          this.showEmailDialog = false;
          this.showTokenDialog = true;
        })
        .catch((errors) => {
          this.errors = errors;
          this.isLoading = false;
        });
    },
    setPassword: function (event) {
      if (event != null) {
        event.preventDefault();
      }

      this.isLoading = true;

      const errors = this.validatePassword();
      if (errors != null && errors.length != 0) {
        this.errors = errors;
        this.isLoading = false;
        return;
      }

      const request = {
        email: this.form.email,
        password: this.form.password,
        token: this.form.token,
      };

      User.changeUserPasswordByToken(request)
        .then(() => {
          this.errors = [];
          this.isLoading = false;
          this.showEmailDialog = false;
          this.showTokenDialog = false;
          this.showSuccessInfo = true;
        })
        .catch((errors) => {
          this.isLoading = false;
          this.errors = errors;
        });
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
