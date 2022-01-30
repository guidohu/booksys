<template>
  <b-modal
    id="userPasswordEditModal"
    title="Change Password"
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
            id="old-password"
            label="Current Password"
            label-for="old-password-input"
            description=""
          >
            <b-form-input
              id="old-password-input"
              v-model="form.oldPassword"
              type="password"
              required
              placeholder=""
              autocomplete="current-password"
            />
          </b-form-group>
          <b-form-group
            id="new-password"
            label="New Password"
            label-for="new-password-input"
            description=""
          >
            <b-form-input
              id="new-password-input"
              v-model="form.newPassword"
              type="password"
              required
              placeholder=""
              autocomplete="new-password"
            />
          </b-form-group>
          <b-form-group
            id="new-password-confirm"
            label="Confirm New Password"
            label-for="new-password-confirm-input"
            description=""
          >
            <b-form-input
              id="new-password-confirm-input"
              v-model="form.newPasswordConfirm"
              type="password"
              required
              placeholder=""
              autocomplete="new-password"
            />
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
import { mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BButton,
  BIconPersonCheck,
  BIconX,
} from "bootstrap-vue";

export default {
  name: "UserPasswordEditModal",
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
    BButton,
    BIconPersonCheck,
    BIconX,
  },
  props: ["visible"],
  data() {
    return {
      form: {
        oldPassword: "",
        newPassword: "",
        newPasswordConfirm: "",
      },
      errors: [],
      isLoading: false,
    };
  },
  methods: {
    close: function () {
      this.$emit("update:visible", false);
    },
    save: function (event) {
      event.preventDefault();

      if (this.form.oldPassword == this.form.newPassword) {
        this.errors = ["Old and new password cannot be identical."];
        return;
      }
      if (this.form.newPassword != this.form.newPasswordConfirm) {
        this.errors = [
          "The new passwords do not match. Please enter the same values for the new passwords.",
        ];
        return;
      }

      this.isLoading = true;
      console.log(
        "change password from",
        this.form.oldPassword,
        "to",
        this.form.newPassword
      );

      this.changeUserPassword({
        oldPassword: this.form.oldPassword,
        newPassword: this.form.newPassword,
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
    ...mapActions("user", ["changeUserPassword"]),
  },
};
</script>
