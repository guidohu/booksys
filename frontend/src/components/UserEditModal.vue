<template>
  <b-modal
    id="userEditModal"
    title="Profile"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <b-row class="text-left">
      <b-col
        cols="1"
        class="d-none d-sm-block"
      />
      <b-col
        cols="12"
        sm="10"
      >
        <b-form @submit="save">
          <b-row v-if="errors.length">
            <b-col cols="12">
              <WarningBox :error="errors" />
            </b-col>
          </b-row>
          <b-form-group
            id="first-name"
            label="First Name"
            label-for="first-name-input"
            description=""
          >
            <b-form-input
              id="first-name-input"
              v-model="form.firstName"
              type="text"
              required
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="last-name"
            label="Last Name"
            label-for="last-name-input"
            description=""
          >
            <b-form-input
              id="last-name-input"
              v-model="form.lastName"
              type="text"
              required
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="street-nr"
            label="Street / Nr"
            label-for="street-nr-input"
            description=""
          >
            <b-form-input
              id="street-nr-input"
              v-model="form.street"
              type="text"
              required
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="plz"
            label="Zip Code"
            label-for="zip-input"
            description=""
          >
            <b-form-input
              id="zip-input"
              v-model="form.zip"
              type="text"
              required
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="city"
            label="City"
            label-for="city-input"
            description=""
          >
            <b-form-input
              id="city-input"
              v-model="form.city"
              type="text"
              required
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="email"
            label="Email"
            label-for="email-input"
            description=""
          >
            <b-form-input
              id="email-input"
              v-model="form.email"
              type="email"
              required
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="tel"
            label="Phone Nr"
            label-for="tel-input"
            description=""
          >
            <b-form-input
              id="tel-input"
              v-model="form.phone"
              type="tel"
              required
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="license"
            label="Driving License"
            label-for="license-input"
            description=""
          >
            <b-form-checkbox
              v-model="form.license"
              name="license-switch"
              switch
            >
              {{ boatLicenseText }}
            </b-form-checkbox>
          </b-form-group>
        </b-form>
      </b-col>
      <b-col
        cols="1"
        class="d-none d-sm-block"
      />
    </b-row>
    <div slot="modal-footer">
      <b-button
        type="button"
        variant="outline-info"
        :disabled="isLoading"
        @click="save"
      >
        <b-icon-person-check />
        Save
      </b-button>
      <b-button
        class="ml-1"
        type="button"
        variant="outline-danger"
        @click="close"
      >
        <b-icon-x />
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BFormCheckbox,
  BButton,
  BIconPersonCheck,
  BIconX,
} from "bootstrap-vue";

export default {
  name: "UserEditModal",
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
    BFormCheckbox,
    BButton,
    BIconPersonCheck,
    BIconX,
  },
  props: ["visible"],
  data() {
    return {
      form: {},
      errors: [],
      isLoading: false,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    boatLicenseText: function () {
      if (this.form.license != null && this.form.license == true) {
        return "I have a driving license for boats";
      } else {
        return "I do NOT have a driving license for boats";
      }
    },
  },
  methods: {
    close: function () {
      this.$emit("update:visible", false);
    },
    save: function (event) {
      event.preventDefault();

      this.isLoading = true;

      this.changeUserProfile({
        first_name: this.form.firstName,
        last_name: this.form.lastName,
        address: this.form.street,
        plz: this.form.zip,
        city: this.form.city,
        email: this.form.email,
        mobile: this.form.phone,
        license: this.form.license == true ? 1 : 0,
      })
        .then(() => {
          this.isLoading = false;
          this.errors = [];
          this.close();
        })
        .catch((errors) => {
          this.isLoading = false;
          this.errors = errors;
        });
    },
    ...mapActions("user", ["changeUserProfile"]),
  },
  beforeMount() {
    const license =
      this.userInfo.license != null && this.userInfo.license == 1
        ? true
        : false;
    this.form = {
      firstName: this.userInfo.first_name,
      lastName: this.userInfo.last_name,
      street: this.userInfo.address,
      zip: this.userInfo.plz,
      city: this.userInfo.city,
      email: this.userInfo.email,
      phone: this.userInfo.mobile,
      license: license,
    };
  },
};
</script>
