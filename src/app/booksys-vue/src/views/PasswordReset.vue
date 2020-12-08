<template>
  <b-modal
    id="passwordResetModal"
    title="Password Reset"
    no-close-on-backdrop
    no-close-on-esc
    hide-header-close
    visible
  >
    <div v-if="showEmailDialog==true">
      <b-row class="text-left mb-4">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          Please enter your email address to get the password reset information sent to you by email.
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
      <b-form @submit="requestToken">
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
            <!-- Captcha -->
            <b-row v-if="getRecaptchaKey != null" class="mt-3">
              <b-col cols="9" offset="3">
                <vue-recaptcha :sitekey="getRecaptchaKey" :loadRecaptchaScript="true" size="compact" @verify="verifiedHandler"></vue-recaptcha>
              </b-col>
            </b-row>
          </b-col>
          <b-col cols="1" class="d-none d-sm-block"></b-col>
        </b-row>
      </b-form>
    </div>
    <div v-if="showTokenDialog==true">
      <b-row class="text-left mb-4">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          Please enter the token that has been sent to your email address and choose new password.
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
      <b-form @submit="setPassword">
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
                disabled
              />
            </b-form-group>
            <!-- Token -->
            <b-form-group
              id="token"
              label="Token"
              label-for="token-input"
              description=""
              label-cols="3"
            >
              <b-form-input
                id="token-input"
                v-model="form.token"
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
          </b-col>
          <b-col cols="1" class="d-none d-sm-block"></b-col>
        </b-row>
      </b-form>
    </div>
    <b-row class="text-center">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <p class="h4 mb-2">
        <b-icon icon="check-circle" variant="success"/>
        <br/>
        Password has been changed successfully.
        </p>
        <p>
          Please login with your new password.
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
      <b-button v-if="showEmailDialog" class="mr-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x></b-icon-x>
        Cancel
      </b-button>
      <b-button v-if="showEmailDialog" type="button" variant="outline-info" v-on:click="requestToken">
        <b-icon icon="arrow-right"/>
        Next
      </b-button>   
      <b-button v-if="showTokenDialog" class="mr-1" type="button" variant="outline-danger" v-on:click="showEmail">
        <b-icon icon="arrow-left"/>
        Back
      </b-button>
      <b-button v-if="showTokenDialog" type="button" variant="outline-info" v-on:click="setPassword">
        <b-icon icon="check"/>
        Set Password
      </b-button>
      <b-button v-if="showSuccessInfo" type="button" variant="outline-info" v-on:click="close">
        <b-icon icon="check"/>
        Done
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import { mapGetters, mapActions } from 'vuex';
import VueRecaptcha from 'vue-recaptcha';
import User from '@/api/user';
import WarningBox from '@/components/WarningBox';

export default Vue.extend({
  name: "PasswordReset",
  components: {
    VueRecaptcha,
    WarningBox
  },
  data() {
    return {
      showEmailDialog: true,
      showTokenDialog: true,
      showSuccessInfo: true,
      errors: [],
      form: {
        email: null,
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
    verifiedHandler: function(response){
      this.form.recaptchaResponse = response;
    },
    requestToken: function(event){
      if(event != null){
        event.preventDefault();
      }

      const request = {
        email:             this.form.email,
        recaptchaResponse: this.form.recaptchaResponse
      };
    
      User.requestPasswordResetToken(request)
      .then(() => {
        this.errors = [];
        this.showSuccessInfo = false;
        this.showEmailDialog = false;
        this.showTokenDialog = true;
      })
      .catch((errors) => this.errors = errors);
    },
    setPassword: function(event){
      if(event != null){
        event.preventDefault();
      }

      const request = {
        email:             this.form.email,
        password:          this.form.password,
        token:             this.form.token
      };

      User.changeUserPasswordByToken(request)
      .then(() => {
        this.errors = [];
        this.showEmailDialog = false;
        this.showTokenDialog = false;
        this.showSuccessInfo = true;
      })
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