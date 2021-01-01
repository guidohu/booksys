<template>
  <b-form v-on:submit.prevent="save">
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
            v-model="userData.email"
            type="text"
            autocomplete="username"
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
            v-model="userData.password"
            type="password"
            autocomplete="new-password"
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
            v-model="userData.passwordConfirm"
            type="password"
            autocomplete="new-password"
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
            v-model="userData.firstName"
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
            v-model="userData.lastName"
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
            v-model="userData.street"
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
            v-model="userData.zip"
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
            v-model="userData.city"
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
            v-model="userData.phone"
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
              :value="licenseToggleState"
              sync
              @change="licenseToggleHandler"
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{checked: 'Yes', unchecked: 'No'}"
            />
          </b-col>
        </b-row>
        <!-- Own risk -->
        <b-row v-if="showDisclaimer" class="text-left align-middle mt-3"
        >
          <b-col cols="9">
            I know the boat community and I perform the activities on my own risk.
          </b-col>
          <b-col cols="3" class="text-right">
            <toggle-button 
              id="own-risk-toggle"
              :value="ownRiskToggleState"
              sync
              @change="ownRiskToggleHandler"
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{checked: 'Yes', unchecked: 'No'}"
            />
          </b-col>
        </b-row>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
  </b-form>
</template>

<script>
import Vue from 'vue';
import { ToggleButton } from 'vue-js-toggle-button';
import {
  BForm,
  BRow,
  BCol,
  BFormGroup,
  BFormInput
} from 'bootstrap-vue';

export default Vue.extend({
  name: "UserSignUp",
  components: {
    ToggleButton,
    BForm,
    BRow,
    BCol,
    BFormGroup,
    BFormInput
  },
  props: [ 'userData', 'showDisclaimer' ],
  data() {
    return {
      toggleWidth: 50,
      licenseToggleState: false,
      ownRiskToggleState: false,
    }
  },
  methods: {
    save: function() {
      this.$emit('save')
    },
    licenseToggleHandler: function(){
      console.log("licenseToggleHandler called");
      this.userData.license = !this.userData.license;
      this.licenseToggleState = this.userData.license;
    },
    ownRiskToggleHandler: function(){
      console.log("ownRiskToggleHandler called");
      this.userData.ownRisk = !this.userData.ownRisk;
      this.ownRiskToggleState = this.userData.ownRisk;
    },
  },
  mounted() {
    if(this.userData.license == true){
      this.licenseToggleState = true;
    }
    if(this.userData.ownRisk == true){
      this.ownRiskToggleState = true;
    }
  }
})
</script>