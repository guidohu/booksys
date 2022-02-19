<template>
  <modal-container
    name="user-password-edit-modal"
    :visible="visible"
  >
    <modal-header
      :closable="true"
      title="Change Password"
      @close="$emit('update:visible', false)"
      />
      <modal-body>
        <div class="row" v-if="errors.length">
          <warning-box :errors="errors" />
        </div>
        <form @submit.prevent.self="save">
          <input-password
            id="old-password"
            label="Current Password"
            autocomplete="current-password"
            v-model="form.oldPassword"
          />
          <input-password
            id="new-password"
            label="New Password"
            autocomplete="new-password"
            v-model="form.newPassword"
          />
          <input-password
            id="new-password-confirm"
            label="Confirm New Password"
            autocomplete="new-password"
            v-model="form.newPasswordConfirm"
          />
        </form>
      </modal-body>
      <modal-footer>
        <button type="submit" class="btn btn-outline-info" :disabled="isLoading" @click.prevent.self="save">
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
import { mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import ModalContainer from './bricks/ModalContainer.vue';
import ModalHeader from './bricks/ModalHeader.vue';
import ModalBody from './bricks/ModalBody.vue';
import ModalFooter from './bricks/ModalFooter.vue';
import InputPassword from './forms/inputs/InputPassword.vue';

export default {
  name: "UserPasswordEditModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputPassword
  },
  props: ["visible"],
    ModalContainer,
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
