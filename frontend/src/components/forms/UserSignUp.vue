<template>
  <form @submit.prevent="save">
    <input-text
      id="email"
      label="Email"
      size="small"
      v-model="signUpData.email"
      autocomplete="username"
      @input="update()"
    />
    <input-password
      id="password"
      label="Password"
      size="small"
      v-model="signUpData.password"
      autocomplete="new-password"
      @input="update()"
    />
    <input-password
      id="password-confirm"
      label="Confirm"
      size="small"
      v-model="signUpData.passwordConfirm"
      autocomplete="new-password"
      @input="update()"
    />
    <hr />
    <input-text
      id="first-name"
      label="First Name"
      size="small"
      v-model="signUpData.firstName"
      @input="update()"
    />
    <input-text
      id="last-name"
      label="Last Name"
      size="small"
      v-model="signUpData.lastName"
      @input="update()"
    />
    <input-text
      id="street"
      label="Street / No"
      size="small"
      v-model="signUpData.street"
      @input="update()"
    />
    <input-text
      id="zip"
      label="Zip Code"
      size="small"
      v-model="signUpData.zip"
      @input="update()"
    />
    <input-text
      id="city"
      label="Place / City"
      size="small"
      v-model="signUpData.city"
      @input="update()"
    />
    <input-text
      id="phone"
      label="Phone"
      size="small"
      v-model="signUpData.phone"
      @input="update()"
    />
    <hr />
    <input-toggle
      id="license-toggle"
      label="Driver License"
      onLabel="I have a driver's license for boats."
      offLabel="I have a driver's license for boats."
      v-model="licenseToggleState"
      @change="licenseToggleHandler"
    />
    <input-toggle
      id="own-risk-toggle"
      label="GTC"
      onLabel="I know the boat community and I perform the activities on my own risk."
      offLabel="I know the boat community and I perform the activities on my own risk."
      v-model="ownRiskToggleState"
      @change="ownRiskToggleHandler"
    />
  </form>
</template>

<script setup>
import InputToggle from "@/components/forms/inputs/InputToggle.vue";
import InputText from "@/components/forms/inputs/InputText.vue";
import InputPassword from "@/components/forms/inputs/InputPassword.vue";
import { ref, watch, onMounted } from "vue";

const props = defineProps(["userData", "showDisclaimer"]);
const emit = defineEmits(["save", "update:user"]);

const signUpData = ref({});
const licenseToggleState = ref(false);
const ownRiskToggleState = ref(false);

watch(
  () => props.userData,
  (newVal) => {
    signUpData.value = newVal;
  },
  { immediate: true, deep: true }
);

onMounted(() => {
  if (signUpData.value.license == true) {
    licenseToggleState.value = true;
  }
  if (signUpData.value.ownRisk == true) {
    ownRiskToggleState.value = true;
  }
});

const save = () => {
  emit("save");
};

const update = () => {
  emit("update:user", signUpData.value);
};

const licenseToggleHandler = () => {
  signUpData.value.license = !signUpData.value.license;
  licenseToggleState.value = signUpData.value.license;
  update();
};

const ownRiskToggleHandler = () => {
  signUpData.value.ownRisk = !signUpData.value.ownRisk;
  ownRiskToggleState.value = signUpData.value.ownRisk;
  update();
};
</script>
