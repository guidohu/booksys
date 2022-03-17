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
    <hr/>
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

<script>
import InputToggle from "@/components/forms/inputs/InputToggle.vue";
import InputText from "@/components/forms/inputs/InputText.vue";
import InputPassword from "@/components/forms/inputs/InputPassword.vue";

export default {
  name: "UserSignUp",
  components: {
    InputToggle,
    InputText,
    InputPassword
  },
  props: ["userData", "showDisclaimer"],
  data() {
    return {
      signUpData: {},
      toggleWidth: 50,
      licenseToggleState: false,
      ownRiskToggleState: false,
    };
  },
  mounted() {
    if (this.signUpData.license == true) {
      this.licenseToggleState = true;
    }
    if (this.signUpData.ownRisk == true) {
      this.ownRiskToggleState = true;
    }
  },
  created() {
    // if we already get data provided upon initialization,
    // we use it
    this.signUpData = this.$props.userData;
  },
  methods: {
    save: function () {
      this.$emit("save");
    },
    update: function () {
      this.$emit("update:user", this.signUpData);
    },
    licenseToggleHandler: function () {
      this.signUpData.license = !this.signUpData.license;
      this.licenseToggleState = this.signUpData.license;
      this.update();
    },
    ownRiskToggleHandler: function () {
      this.signUpData.ownRisk = !this.signUpData.ownRisk;
      this.ownRiskToggleState = this.signUpData.ownRisk;
      this.update();
    },
  },
};
</script>
