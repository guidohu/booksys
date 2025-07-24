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

<script setup>
import { ref, computed, onMounted, defineAsyncComponent } from "vue";
import { useStore, mapGetters, mapActions } from "vuex";
import { useRouter } from "vue-router";
import WarningBox from "@/components/WarningBox";
import User from "@/api/user";
import ModalContainer from "@/components/bricks/ModalContainer.vue";
import ModalHeader from "@/components/bricks/ModalHeader.vue";
import ModalBody from "@/components/bricks/ModalBody.vue";
import ModalFooter from "@/components/bricks/ModalFooter.vue";
import { VueRecaptcha } from 'vue-recaptcha';

const UserSignUp = defineAsyncComponent(() =>
  import(/* webpackChunkName: "user-sign-up" */ "@/components/forms/UserSignUp")
);

const store = useStore();
const router = useRouter();

const errors = ref([]);
const isSignedUp = ref(false);
const form = ref({
  license: false,
  ownRisk: false,
  recaptchaResponse: null,
});
const visible = ref(true);

const getRecaptchaKey = computed(
  mapGetters("configuration", ["getRecaptchaKey"]).getRecaptchaKey.bind({
    $store: store,
  })
);

const queryRecaptchaKey = mapActions("configuration", [
  "queryRecaptchaKey",
]).queryRecaptchaKey.bind({ $store: store });

const recaptchaVerifiedHandler = (response) => {
  form.value.recaptchaResponse = response;
};

const handleUserUpdate = (userData) => {
  form.value.email = userData.email;
  form.value.password = userData.password;
  form.value.passwordConfirm = userData.passwordConfirm;
  form.value.firstName = userData.firstName;
  form.value.lastName = userData.lastName;
  form.value.street = userData.street;
  form.value.zip = userData.zip;
  form.value.city = userData.city;
  form.value.phone = userData.phone;
  form.value.ownRisk = userData.ownRisk;
  form.value.license = userData.license;
};

const validateForm = () => {
  const validationErrors = [];

  if (form.value.ownRisk == false) {
    validationErrors.push(
      "To sign-up you need to agree that the activities are performed at your own risk."
    );
  }
  if (form.value.password != form.value.passwordConfirm) {
    validationErrors.push("Password and Password Confirmation are not identical.");
  }
  if (form.value.password.length <= 8) {
    validationErrors.push("Please use a password longer than 8 characters.");
  }
  if (form.value.recaptchaResponse == null && getRecaptchaKey.value) {
    validationErrors.push("Please tick `I'm not a robot`.");
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

const save = () => {
  errors.value = validateForm();
  if (errors.value.length > 0) {
    return;
  }

  User.signUp(form.value)
    .then(() => (isSignedUp.value = true))
    .catch((err) => (errors.value = err));
};

const close = () => {
  router.push("/");
};

onMounted(() => {
  queryRecaptchaKey();
});
</script>
