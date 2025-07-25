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
              <div v-if="getRecaptchaKey != null && getRecaptchaKey.length > 0" class="row mt-3">
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

<script setup>
import { ref, computed, onMounted } from "vue";
import { useStore, mapGetters, mapActions } from "vuex";
import { useRouter } from "vue-router";
import User from "booksys/api/user";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "../components/bricks/ModalContainer.vue";
import ModalHeader from "../components/bricks/ModalHeader.vue";
import ModalBody from "../components/bricks/ModalBody.vue";
import ModalFooter from "../components/bricks/ModalFooter.vue";
import InputText from "../components/forms/inputs/InputText.vue";
import InputPassword from "../components/forms/inputs/InputPassword.vue";
import { VueRecaptcha } from 'vue-recaptcha';

const store = useStore();
const router = useRouter();

const isLoading = ref(false);
const showEmailDialog = ref(true);
const showTokenDialog = ref(false);
const showSuccessInfo = ref(false);
const errors = ref([]);
const form = ref({
  email: null,
  recaptchaResponse: null,
});

const getRecaptchaKey = computed(
  mapGetters("configuration", ["getRecaptchaKey"]).getRecaptchaKey.bind({
    $store: store,
  })
);

const queryRecaptchaKey = mapActions("configuration", [
  "queryRecaptchaKey",
]).queryRecaptchaKey.bind({ $store: store });

const verifiedHandler = (response) => {
  form.value.recaptchaResponse = response;
};

const validatePassword = () => {
  const validationErrors = [];

  if (form.value.password !== form.value.passwordConfirm) {
    validationErrors.push("Password and Password Confirmation are not identical.");
  }
  if (form.value.password.length <= 8) {
    validationErrors.push("Please use a password longer than 8 characters.");
  }

  // password strength
  const pwUpperRegex = /[A-Z]+/;
  const pwLowerRegex = /[a-z]+/;
  const pwDigitRegex = /[0-9]+/;
  if (form.value.password.match(pwUpperRegex) == null) {
    validationErrors.push(
      "The password needs to contain at least one upper case letter (A-Z)"
    );
  }
  if (form.value.password.match(pwLowerRegex) == null) {
    validationErrors.push(
      "The password needs to contain at least one lower case letter (a-z)"
    );
  }
  if (form.value.password.match(pwDigitRegex) == null) {
    validationErrors.push("The password needs to contain at least one digit (0-9)");
  }

  return validationErrors;
};

const requestToken = (event) => {
  if (event) {
    event.preventDefault();
  }

  isLoading.value = true;

  const request = {
    email: form.value.email,
    recaptchaResponse: form.value.recaptchaResponse,
  };

  User.requestPasswordResetToken(request)
    .then(() => {
      isLoading.value = false;
      errors.value = [];
      showSuccessInfo.value = false;
      showEmailDialog.value = false;
      showTokenDialog.value = true;
    })
    .catch((error) => {
      errors.value = error;
      isLoading.value = false;
    });
};

const setPassword = (event) => {
  if (event) {
    event.preventDefault();
  }

  isLoading.value = true;

  const validationErrors = validatePassword();
  if (validationErrors.length > 0) {
    errors.value = validationErrors;
    isLoading.value = false;
    return;
  }

  const request = {
    email: form.value.email,
    password: form.value.password,
    token: form.value.token,
  };

  User.changeUserPasswordByToken(request)
    .then(() => {
      errors.value = [];
      isLoading.value = false;
      showEmailDialog.value = false;
      showTokenDialog.value = false;
      showSuccessInfo.value = true;
    })
    .catch((error) => {
      isLoading.value = false;
      errors.value = error;
    });
};

const close = () => {
  router.push("/");
};

const showEmail = () => {
  errors.value = [];
  isLoading.value = false;
  showEmailDialog.value = true;
  showTokenDialog.value = false;
  showSuccessInfo.value = false;
};

onMounted(() => {
  queryRecaptchaKey();
});
</script>
