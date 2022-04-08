<template>
  <sectioned-card-module>
    <template v-slot:body>
      <warning-box v-if="errors.length > 0" :errrors="errors" />
      <logo-upload @logoChanged="logoChangeHandler" />
      <form @submit.stop="save">
        <div class="alert alert-info">
          Settings related to the boat you are using.
        </div>
        <input-select
          :options="engineHourLogFormats"
          v-model="form.engineHourFormat"
          label="Engine Hour Format"
          size="small"
          description="select the format your engine hour display uses."
        />
        <div class="alert alert-info">Settings for your daily operations.</div>
        <input-select
          :options="fuelPaymentTypes"
          v-model="form.fuelPaymentType"
          label="Fuel Payments"
          size="small"
          description="choose between direct payments at the time of consumption and credited payment where you add payments for fuel in the payments section separately after paying a bill."
        />
        <div class="alert alert-info">
          Timezone and location settings. Users will see the times in the
          timezone defined by these settings, also sunrise and sunset will be
          calculated based on these settings.
        </div>
        <input-text
          id="timezone"
          label="Timezone"
          description="e.g., Europe/Zurich"
          size="small"
          v-model="form.timezone"
        />
        <input-text
          id="longitude"
          label="Longitude"
          description="Define the longitude in degrees (example 8.5432)"
          size="small"
          v-model="form.longitude"
        />
        <input-text
          id="latitude"
          label="Latitude"
          description="Define the latitude in degrees (example 47.3456)"
          size="small"
          v-model="form.latitude"
        />
        <input-text
          id="map-iframe"
          label="Google Maps Iframe"
          description="Location using an Iframe link from Google Maps."
          size="small"
          v-model="form.mapIframe"
        />
        <input-text-multiline
          id="address"
          label="Address"
          size="small"
          v-model="form.address"
          rows="3"
        />
        <div class="alert alert-info">
          Payment Information. This information is displayed to users in case
          they want to top-up their balance.
        </div>
        <input-text
          id="currency"
          label="Currency"
          size="small"
          v-model="form.currency"
        />
        <input-text
          id="account-owner"
          label="Account Owner"
          size="small"
          v-model="form.accountOwner"
        />
        <input-text
          id="account-iban"
          label="IBAN"
          size="small"
          v-model="form.iban"
        />
        <input-text
          id="account-bic"
          label="BIC / SWIFT"
          size="small"
          v-model="form.bic"
        />
        <input-text
          id="account-comment"
          label="Comment"
          size="small"
          v-model="form.comment"
          description="Additional advice shown to users as payments instructions."
        />
        <div class="alert alert-info">
          Notification Settings. The application will send emails to users,
          therefore it needs to have access to an email account for doing so.
        </div>
        <input-text
          id="smtp-sender"
          label="Sender Address"
          size="small"
          v-model="form.smtpSender"
        />
        <input-text
          id="smtp-server"
          label="SMTP Server"
          size="small"
          v-model="form.smtpServer"
        />
        <input-text
          id="smtp-username"
          label="Username"
          size="small"
          v-model="form.smtpUsername"
        />
        <input-password
          id="smtp-password"
          label="Password"
          size="small"
          v-model="form.smtpPassword"
        />
        <div class="alert alert-info">
          This application supports myNautique and can integrate information
          such as fuel level or engine hours from myNautique directly into this
          app.
        </div>
        <input-toggle
          id="my-nautique-toggle"
          label="myNautique"
          on-label="On"
          off-label="Off"
          v-model="form.myNautiqueEnabled"
        />
        <input-text
          v-if="form.myNautiqueEnabled"
          id="my-nautique-user"
          label="User"
          size="small"
          v-model="form.myNautiqueUser"
        />
        <input-password
          v-if="form.myNautiqueEnabled"
          id="my-nautique-password"
          label="Password"
          size="small"
          v-model="form.myNautiquePassword"
        />
        <input-text
          v-if="form.myNautiqueEnabled"
          id="my-nautique-boat-id"
          label="Boat ID"
          size="small"
          description="The ID you can find in the myNautique app settings. Typically a 5-6 digit number."
          v-model="form.myNautiqueBoatId"
        />
        <input-fuel
          v-if="form.myNautiqueEnabled"
          id="my-nautique-fuel-capacity"
          label="Fuel Capacity"
          size="small"
          v-model="form.myNautiqueFuelCapacity"
        />
        <div class="alert alert-info">
          ReCAPTCHA protects the application and its users from SPAMers. Thus it
          is recommended to use ReCAPTCHA.
        </div>
        <input-text
          id="recaptcha-private-key"
          label="Private Key"
          size="small"
          v-model="form.recaptchaPrivateKey"
        />
        <input-text
          id="recaptcha-public-key"
          label="Public Key"
          size="small"
          v-model="form.recaptchaPublicKey"
        />
      </form>
    </template>
    <template v-if="showControls" v-slot:footer>
      <div class="row">
        <div class="col-12 text-end">
          <button class="btn btn-outline-danger" @click="cancel">
            <i class="bi bi-x" />
            Cancel
          </button>
          <button class="btn btn-outline-info ms-1" @click="save">
            <i class="bi bi-check" />
            Save
          </button>
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import LogoUpload from "@/components/LogoUpload";
import SectionedCardModule from "@/components/bricks/SectionedCardModule.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputPassword from "./forms/inputs/InputPassword.vue";
import InputTextMultiline from "./forms/inputs/InputTextMultiline.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import InputToggle from "./forms/inputs/InputToggle.vue";
import InputFuel from "./forms/inputs/InputFuel.vue";

export default {
  name: "SettingsCard",
  components: {
    WarningBox,
    LogoUpload,
    SectionedCardModule,
    InputToggle,
    InputText,
    InputPassword,
    InputTextMultiline,
    InputSelect,
    InputFuel,
  },
  props: ["showControls"],
  emits: ["save", "saved", "cancelled", "change"],
  data() {
    return {
      errors: [],
      form: {
        logoFile: null,
        fuelPaymentType: "instant",
      },
      fuelPaymentTypes: [
        { value: "instant", text: "pay directly" },
        { value: "billed", text: "pay by bill" },
      ],
      engineHourLogFormats: [
        { value: "hh.h", text: "hh.h  - such as 9.7" },
        { value: "hh:mm", text: "hh:mm - such as 9:42" },
      ],
    };
  },
  computed: {
    ...mapGetters("configuration", ["getConfiguration"]),
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration", "setConfiguration"]),
    logoChangeHandler: function (logoUri) {
      this.form.logoFile = logoUri;
    },
    configChangeHandler: function () {
      const configuration = this.getSanitizedConfiguration();
      if (configuration != null) {
        console.log("emit change", configuration);
        this.$emit("change", configuration);
      }
    },
    getSanitizedConfiguration: function () {
      // extract correct URL from mapIframe
      const v = this.form;
      let mapUrl = [];
      if (v.mapIframe != null && v.mapIframe != "") {
        const mapUrlRegex =
          /https:\/\/www\.google\.com\/maps\/embed\?pb=[^"\s]+/;
        mapUrl = v.mapIframe.match(mapUrlRegex);
        if (mapUrl == null || mapUrl.length != 1) {
          this.errors = ["Google Maps Iframe URL is not in a valid format."];
          return;
        }
      }

      const newConfiguration = {
        logo_file: v.logoFile,
        engine_hour_format: v.engineHourFormat,
        fuel_payment_type: v.fuelPaymentType,
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
        recaptcha_publickey: v.recaptchaPublicKey,
        mynautique_enabled: v.myNautiqueEnabled,
        mynautique_user: v.myNautiqueUser,
        mynautique_password: v.myNautiquePassword,
        mynautique_boat_id: v.myNautiqueBoatId,
        mynautique_fuel_capacity: v.myNautiqueFuelCapacity,
      };

      return newConfiguration;
    },
    save: function () {
      const newConfiguration = this.getSanitizedConfiguration();
      if (newConfiguration == null) {
        return;
      }

      // if not responsible for data controlling -> just emit signal
      if (this.showControls == false) {
        this.$emit("save", newConfiguration);
        return;
      }

      // handle save otherwise
      this.setConfiguration(newConfiguration)
        .then(() => {
          this.errors = [];
          this.$emit("saved", newConfiguration);
        })
        .catch((errors) => (this.errors = errors));
    },
    cancel: function () {
      this.$emit("cancelled");
    },
    setFormDefaults: function (defaultValues) {
      if (defaultValues == null) {
        this.form = {};
        return;
      }

      console.log("set form to defaults:", defaultValues);
      this.form = {
        logoFile: defaultValues.logo_file,
        engineHourFormat: defaultValues.engine_hour_format,
        fuelPaymentType: defaultValues.fuel_payment_type,
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
        recaptchaPublicKey: defaultValues.recaptcha_publickey,
        myNautiqueEnabled: defaultValues.mynautique_enabled,
        myNautiqueUser: defaultValues.mynautique_user,
        myNautiquePassword: defaultValues.mynautique_password,
        myNautiqueBoatId: defaultValues.mynautique_boat_id,
        myNautiqueFuelCapacity: defaultValues.mynautique_fuel_capacity,
      };
    },
  },
  watch: {
    getConfiguration: function (newValues) {
      this.setFormDefaults(newValues);
    },
    form: {
      deep: true,
      handler() {
        console.log("config changed");
        this.configChangeHandler();
      },
    },
  },
  created() {
    this.queryConfiguration();
  },
  mounted() {
    this.setFormDefaults(this.getConfiguration);
  },
};
</script>

<style>
.overflow-scroll {
  overflow-y: scroll;
}
</style>
