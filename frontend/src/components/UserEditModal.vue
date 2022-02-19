<template>
  <modal-container
    name="user-edit-modal"
    :visible="visible"
  >
    <modal-header
      :closable="true"
      title="Profile"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <div class="row" v-if="errors.length">
        <warning-box :errors="errors" />
      </div>
      <form @submit.prevent.self="save">
        <input-text
          id="first-name"
          label="First Name"
          v-model="form.firstName"
        />
        <input-text
          id="last-name"
          label="Last Name"
          v-model="form.lastName"
        />
        <input-text
          id="street-nr"
          label="Street / Nr"
          v-model="form.street"
        />
        <input-text
          id="zip"
          label="Zip Code"
          v-model="form.zip"
        />
        <input-text
          id="city"
          label="City"
          v-model="form.city"
        />
        <input-text
          id="email"
          label="Email"
          type="email"
          v-model="form.email"
        />
        <input-text
          id="tel"
          label="Phone Nr"
          type="tel"
          v-model="form.phone"
        />
        <input-toggle
          id="license"
          label="Driving License"
          v-model="form.license"
          offLabel="I'm NOT authorized to drive boats"
          onLabel="I'm authorized to drive boats"
        />
      </form>
    </modal-body>
    <modal-footer>
      <button type="submit" class="btn btn-outline-info" @click.prevent.self="save">
        <i class="bi bi-check"></i>
        Save
      </button>
      <button type="button" class="btn btn-outline-danger" @click="close">
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import ModalContainer from './bricks/ModalContainer.vue';
import ModalHeader from './bricks/ModalHeader.vue';
import ModalBody from './bricks/ModalBody.vue';
import ModalFooter from './bricks/ModalFooter.vue';
import InputText from './forms/inputs/InputText.vue';
import InputToggle from './forms/inputs/InputToggle.vue';

export default {
  name: "UserEditModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputText,
    InputToggle
  },
  props: ["visible"],
  data() {
    return {
      form: {
        license: null,
      },
      errors: [],
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
  },
  methods: {
    close: function () {
      this.$emit("update:visible", false);
    },
    save: function () {

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
          this.errors = [];
          this.close();
        })
        .catch((errors) => {
          this.errors = errors;
        });
    },
    setFormData: function(data){
      const license =
        data.license != null && data.license == 1
          ? true
          : false;
      this.form = {
        firstName: data.first_name,
        lastName: data.last_name,
        street: data.address,
        zip: data.plz,
        city: data.city,
        email: data.email,
        phone: data.mobile,
        license: license,
      };
    },
    ...mapActions("user", ["changeUserProfile"]),
  },
  watch: {
    userInfo: function(newUserInfo){
      console.log("userInfo changed, new value", newUserInfo);
      this.setFormData(newUserInfo);
    }
  },
  created() {
    if(this.userInfo != null){
      console.log("Created with userInfo");
      this.setFormData(this.userInfo);
      return;
    }
    console.log("Created without userInfo");
    return;
  },
};
</script>
