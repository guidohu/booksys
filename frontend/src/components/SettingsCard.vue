<template>
  <b-card no-body class="text-left">
    <b-card-body class="overflow-scroll">
      <WarningBox v-if="errors.length>0" :errors="errors"/>
      <b-form @submit="save">
        <b-alert variant="info" show>
          Timezone and location settings. Users will see the times in the timezone defined by these settings, also sunrise and sunset will be calculated based on these settings.
        </b-alert>
        <!-- Timezone -->
        <b-form-group
          id="input-group-timezone"
          label="Timezone*"
          label-for="input-timezone"
          label-cols="3"
          description="e.g. Europe/Zurich"
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-timezone"
              v-model="form.timezone"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Longitude -->
        <b-form-group
          id="input-group-longitude"
          label="Longitude*"
          label-for="input-longitude"
          label-cols="3"
          description="Define the longitude in degrees (example 8.5432)"
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-longitude"
              v-model="form.longitude"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Latitude -->
        <b-form-group
          id="input-group-latitude"
          label="Latitude*"
          label-for="input-latitude"
          label-cols="3"
          description="Define the latitude in degrees (example 47.3456)"
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-latitude"
              v-model="form.latitude"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Google Iframe -->
        <b-form-group
          id="input-group-map-iframe"
          label="Google Maps Iframe"
          label-for="input-timezone"
          label-cols="3"
          description=""
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-map-iframe"
              v-model="form.mapIframe"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Address -->
        <b-form-group
          id="input-group-address"
          label="Address"
          label-for="input-address"
          label-cols="3"
          description=""
        >
          <b-form-textarea size="sm"
              id="input-address"
              v-model="form.address"
              rows="3"
              max-rows="3"
          />
        </b-form-group>


        <b-alert variant="info" show>
          Payment Information. This information is displayed to users in case they want to top-up their balance.
        </b-alert>
        <!-- Currency -->
        <b-form-group
          id="input-group-currency"
          label="Currency"
          label-for="input-currency"
          label-cols="3"
          description=""
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-currency"
              v-model="form.currency"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Payment Account Owner -->
        <b-form-group
          id="input-group-account-owner"
          label="Account Owner"
          label-for="input-account-owner"
          label-cols="3"
          description=""
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-account-owner"
              v-model="form.accountOwner"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Account IBAN -->
        <b-form-group
          id="input-group-iban"
          label="IBAN"
          label-for="input-iban"
          label-cols="3"
          description=""
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-iban"
              v-model="form.iban"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Account BIC -->
        <b-form-group
          id="input-group-bic"
          label="BIC / SWIFT"
          label-for="input-bic"
          label-cols="3"
          description=""
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-bic"
              v-model="form.bic"
              type="text"
            />
          </b-input-group>
        </b-form-group>
        <!-- Account Comment -->
        <b-form-group
          id="input-group-comment"
          label="Comment"
          label-for="input-comment"
          label-cols="3"
          description=""
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-comment"
              v-model="form.comment"
              type="text"
            />
          </b-input-group>
        </b-form-group>


        <b-alert variant="info" show>
          Notification Settings. The application will send emails to users, therefore it needs to have access to an email account for doing so.
        </b-alert>
        <!-- SMTP Sender -->
        <b-form-group
          id="input-group-smtp-sender"
          label="Sender Address"
          label-for="input-smtp-sender"
          label-cols="3"
          description=""
        >
          <b-input-group size="sm">
            <b-form-input
              id="input-smtp-sender"
              v-model="form.smtpSender"
              type="text"
            />
          </b-input-group>
        </b-form-group>
      </b-form>
      <!-- SMTP Server -->
      <b-form-group
        id="input-group-smtp-server"
        label="SMTP Server"
        label-for="input-smtp-server"
        label-cols="3"
        description=""
      >
        <b-input-group size="sm">
          <b-form-input
            id="input-smtp-server"
            v-model="form.smtpServer"
            type="text"
          />
        </b-input-group>
      </b-form-group>
      <!-- SMTP Username -->
      <b-form-group
        id="input-group-smtp-username"
        label="Username"
        label-for="input-smtp-username"
        label-cols="3"
        description=""
      >
        <b-input-group size="sm">
          <b-form-input
            id="input-smtp-username"
            v-model="form.smtpUsername"
            type="text"
          />
        </b-input-group>
      </b-form-group>
      <!-- SMTP Password -->
      <b-form-group
        id="input-group-smtp-password"
        label="Password"
        label-for="input-smtp-password"
        label-cols="3"
        description=""
      >
        <b-input-group size="sm">
          <b-form-input
            id="input-smtp-password"
            v-model="form.smtpPassword"
            type="password"
          />
        </b-input-group>
      </b-form-group>


      <b-alert variant="info" show>
        ReCAPTCHA protects the application and its users from SPAMers. Thus it is recommended to use ReCAPTCHA.
      </b-alert>
      <!-- ReCAPTCHA private key -->
      <b-form-group
        id="input-group-recaptcha-private-key"
        label="Private Key"
        label-for="input-recaptcha-private-key"
        label-cols="3"
        description=""
      >
        <b-input-group size="sm">
          <b-form-input
            id="input-recaptcha-private-key"
            v-model="form.recaptchaPrivateKey"
            type="text"
          />
        </b-input-group>
      </b-form-group>
      <!-- ReCAPTCHA Public Key -->
      <b-form-group
        id="input-group-recaptcha-public-key"
        label="Public Key"
        label-for="input-recaptcha-public-key"
        label-cols="3"
        description=""
      >
        <b-input-group size="sm">
          <b-form-input
            id="input-recaptcha-public-key"
            v-model="form.recaptchaPublicKey"
            type="text"
          />
        </b-input-group>
      </b-form-group>
    </b-card-body>
    <b-card-footer>
        <b-row class="text-right">
          <b-col cols="9" offset="3">
            <b-button class="mr-2" variant="outline-danger" to="/admin">
              <b-icon-x/>
              Cancel
            </b-button>
            <b-button variant="outline-info" v-on:click="save">
              <b-icon-check/>
              Save
            </b-button>
          </b-col>
        </b-row>
    </b-card-footer>
  </b-card>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import WarningBox from '@/components/WarningBox';
import {
  BCard,
  BCardBody,
  BCardFooter,
  BForm,
  BAlert,
  BFormGroup,
  BInputGroup,
  BFormInput,
  BFormTextarea,
  BRow,
  BCol,
  BButton,
  BIconX,
  BIconCheck  
} from 'bootstrap-vue';

export default {
  name: "SettingsCard",
  components: {
    WarningBox,
    BCard,
    BCardBody,
    BCardFooter,
    BForm,
    BAlert,
    BFormGroup,
    BInputGroup,
    BFormInput,
    BFormTextarea,
    BRow,
    BCol,
    BButton,
    BIconX,
    BIconCheck
  },
  data() {
    return {
      errors: [],
      form: {}
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getConfiguration'
    ])
  },
  methods: {
    ...mapActions('configuration', [
      'queryConfiguration',
      'setConfiguration'
    ]),
    save: function() {
      // extract correct URL from mapIframe
      const v = this.form;
      let mapUrl = [];
      if(v.mapIframe != null && v.mapIframe != ''){
        const mapUrlRegex = /https:\/\/www\.google\.com\/maps\/embed\?pb=[^"\s]+/;
        mapUrl = v.mapIframe.match(mapUrlRegex);
        if(mapUrl == null || mapUrl.length != 1){
          this.errors = ['Google Maps Iframe URL is not in a valid format.'];
          return;
        }
      }

      const newConfiguration = {
        location_time_zone: v.timezone,
        location_longitude: Number(v.longitude),
        location_latitude: Number(v.latitude),
        location_map: mapUrl[0],
        location_address: v.address,
        currency: v.currency,
        payment_account_owner: v.accountOwner,
        payment_account_iban: v.iban,
        payment_account_bic: v.bic,
        payment_account_comment: v.comment,
        smtp_sender: v.smtpSender,
        smtp_server: v.smtpServer,
        smtp_username: v.smtpUsername,
        smtp_password: v.smtpPassword,
        recaptcha_privatekey: v.recaptchaPrivateKey,
        recaptcha_publickey: v.recaptchaPublicKey
      };

      this.setConfiguration(newConfiguration)
      .then(() => {
        this.errors = [];
        this.$router.go(-1);
      })
      .catch((errors) => this.errors = errors);
    },
    setFormDefaults: function(defaultValues) {
      if(defaultValues == null){
        this.form = {};
        return;
      }

      console.log("set form to defaults:", defaultValues);
      this.form = {
        timezone: defaultValues.location_time_zone,
        longitude: defaultValues.location_longitude,
        latitude: defaultValues.location_latitude,
        mapIframe: defaultValues.location_map,
        address: defaultValues.location_address,
        currency: defaultValues.currency,
        accountOwner: defaultValues.payment_account_owner,
        iban: defaultValues.payment_account_iban,
        bic: defaultValues.payment_account_bic,
        comment: defaultValues.payment_account_comment,
        smtpSender: defaultValues.smtp_sender,
        smtpServer: defaultValues.smtp_server,
        smtpUsername: defaultValues.smtp_username,
        smtpPassword: defaultValues.smtp_password,
        recaptchaPrivateKey: defaultValues.recaptcha_privatekey,
        recaptchaPublicKey: defaultValues.recaptcha_publickey
      }
    }
  },
  watch: {
    getConfiguration: function(newValues){
      this.setFormDefaults(newValues)
    }
  },
  created() {
    this.queryConfiguration();
  },
  mounted() {
    this.setFormDefaults(this.getConfiguration)
  }
}
</script>

<style>
  .overflow-scroll {
    overflow-y: scroll;
  }
</style>