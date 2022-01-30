<template>
  <b-form @submit.prevent="save">
    <b-row class="text-left">
      <b-col
        cols="1"
        class="d-none d-sm-block"
      />
      <b-col
        cols="12"
        sm="10"
      >
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
            v-model="signUpData.email"
            type="text"
            autocomplete="username"
            @input="update()"
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
            v-model="signUpData.password"
            type="password"
            autocomplete="new-password"
            @input="update()"
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
            v-model="signUpData.passwordConfirm"
            type="password"
            autocomplete="new-password"
            @input="update()"
          />
        </b-form-group>
        <hr>
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
            v-model="signUpData.firstName"
            type="text"
            @input="update()"
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
            v-model="signUpData.lastName"
            type="text"
            @input="update()"
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
            v-model="signUpData.street"
            type="text"
            @input="update()"
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
            v-model="signUpData.zip"
            type="text"
            @input="update()"
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
            v-model="signUpData.city"
            type="text"
            @input="update()"
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
            v-model="signUpData.phone"
            type="text"
            @input="update()"
          />
        </b-form-group>
        <hr>
        <!-- License -->
        <b-row class="text-left align-middle">
          <b-col cols="9">
            I have a driver's license for boats.
          </b-col>
          <b-col
            cols="3"
            class="text-right"
          >
            <toggle-button
              id="license-toggle"
              :value="licenseToggleState"
              sync
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{ checked: 'Yes', unchecked: 'No' }"
              @change="licenseToggleHandler"
            />
          </b-col>
        </b-row>
        <!-- Own risk -->
        <b-row
          v-if="showDisclaimer"
          class="text-left align-middle mt-3"
        >
          <b-col cols="9">
            I know the boat community and I perform the activities on my own
            risk.
          </b-col>
          <b-col
            cols="3"
            class="text-right"
          >
            <toggle-button
              id="own-risk-toggle"
              :value="ownRiskToggleState"
              sync
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{ checked: 'Yes', unchecked: 'No' }"
              @change="ownRiskToggleHandler"
            />
          </b-col>
        </b-row>
      </b-col>
      <b-col
        cols="1"
        class="d-none d-sm-block"
      />
    </b-row>
  </b-form>
</template>

<script>
import { ToggleButton } from "vue-js-toggle-button";
import { BForm, BRow, BCol, BFormGroup, BFormInput } from "bootstrap-vue";

export default {
  name: "UserSignUp",
  components: {
    ToggleButton,
    BForm,
    BRow,
    BCol,
    BFormGroup,
    BFormInput,
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
    save: function() {
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
