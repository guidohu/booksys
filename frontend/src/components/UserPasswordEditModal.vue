<template>
  <modal-container name="user-password-edit-modal" :visible="visible" :inert="!visible">
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
      <button
        type="submit"
        class="btn btn-outline-info"
        :disabled="isLoading"
        @click.prevent.self="save"
      >
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

<script setup>
import { ref } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputPassword from "./forms/inputs/InputPassword.vue";

defineProps(["visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const form = ref({
  oldPassword: "",
  newPassword: "",
  newPasswordConfirm: "",
});
const errors = ref([]);
const isLoading = ref(false);

const changeUserPassword = (data) =>
  store.dispatch("user/changeUserPassword", data);

function close() {
  emit("update:visible", false);
}

function save(event) {
  event.preventDefault();

  if (form.value.oldPassword == form.value.newPassword) {
    errors.value = ["Old and new password cannot be identical."];
    return;
  }
  if (form.value.newPassword != form.value.newPasswordConfirm) {
    errors.value = [
      "The new passwords do not match. Please enter the same values for the new passwords.",
    ];
    return;
  }

  isLoading.value = true;

  changeUserPassword({
    oldPassword: form.value.oldPassword,
    newPassword: form.value.newPassword,
  })
    .then(() => {
      isLoading.value = false;
      errors.value = [];
      close();
    })
    .catch((errs) => {
      isLoading.value = false;
      errors.value = errs;
    });
}
</script>
