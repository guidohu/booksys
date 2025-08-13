<template>
  <modal-container name="user-edit-modal" :visible="visible" :inert="!visible">
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
        <input-text id="last-name" label="Last Name" v-model="form.lastName" />
        <input-text id="street-nr" label="Street / Nr" v-model="form.street" />
        <input-text id="zip" label="Zip Code" v-model="form.zip" />
        <input-text id="city" label="City" v-model="form.city" />
        <input-text
          id="email"
          label="Email"
          type="email"
          v-model="form.email"
        />
        <input-text id="tel" label="Phone Nr" type="tel" v-model="form.phone" />
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
      <button
        type="submit"
        class="btn btn-outline-info"
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
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputToggle from "./forms/inputs/InputToggle.vue";

const props = defineProps(["visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const form = ref({
  license: null,
});
const errors = ref([]);

const userInfo = computed(() => store.getters["login/userInfo"]);

const changeUserProfile = (data) =>
  store.dispatch("user/changeUserProfile", data);

function close() {
  emit("update:visible", false);
}

function save() {
  if (isNaN(form.value.zip)) {
    errors.value = ["Zip Code needs to be a number."];
    return;
  }
  changeUserProfile({
    first_name: form.value.firstName,
    last_name: form.value.lastName,
    address: form.value.street,
    plz: Number(form.value.zip),
    city: form.value.city,
    email: form.value.email,
    mobile: form.value.phone,
    license: form.value.license == true,
  })
    .then(() => {
      errors.value = [];
      close();
    })
    .catch((errs) => {
      errors.value = errs;
    });
}

function setFormData(data) {
  const license = data.license != null && data.license == 1 ? true : false;
  form.value = {
    firstName: data.first_name,
    lastName: data.last_name,
    street: data.address,
    zip: data.plz,
    city: data.city,
    email: data.email,
    phone: data.mobile,
    license: license,
  };
}

watch(userInfo, (newUserInfo) => {
  console.log("userInfo changed, new value", newUserInfo);
  setFormData(newUserInfo);
});

if (userInfo.value != null) {
  console.log("Created with userInfo");
  setFormData(userInfo.value);
}
</script>
