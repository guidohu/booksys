<template>
  <modal-container name="user-group-modal" :visible="visible">
    <modal-header
      :closable="true"
      :title="title"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length" :errors="errors" />
      <form @submit.prevent="save">
        <h6 class="mt-2 mb-3">User Group</h6>
        <input-text
          id="user-group-name"
          label="Name"
          v-model="form.user_group_name"
          size="small"
          :disabled="!isEditMode"
        />
        <input-text
          id="user-group-description"
          label="Description"
          v-model="form.user_group_description"
          size="small"
          :disabled="!isEditMode"
        />
        <!-- User Role -->
        <h6 class="mt-5 mb-3">User Role and Permissions</h6>
        <input-select
          id="user-role-name"
          label="Role"
          :options="userRoleList"
          v-model="form.user_role_id"
          :disabled="!isEditMode"
          size="small"
        />
        <input-text
          id="user-role-description"
          label="Description"
          v-model="userRoleDescription"
          size="small"
          :disabled="true"
        />
        <!-- Pricing -->
        <h6 class="mt-5 mb-3">Pricing</h6>
        <input-currency
          id="price"
          label="Price"
          v-model="form.price_min"
          :disabled="!isEditMode"
          size="small"
          :currency="getCurrency + '/min'"
        />
        <input-text
          id="price-description"
          label="Description"
          v-model="form.price_description"
          size="small"
          :disabled="!isEditMode"
        />
      </form>
    </modal-body>
    <modal-footer class="block-footer">
      <div class="row">
        <div class="col-6 text-start">
          <button
            v-if="isEditMode && form.user_group_id != null"
            type="button"
            class="btn btn-outline-danger"
            @click.stop="remove"
          >
            <i class="bi bi-trash"></i>
            Delete
          </button>
        </div>
        <div class="col-6 text-end">
          <button
            v-if="isEditMode == true"
            type="submit"
            class="btn btn-outline-info ms-2"
            @click.prevent.self="save"
          >
            <i class="bi bi-check"></i>
            Save
          </button>
          <button
            v-if="isEditMode == false"
            type="submit"
            class="btn btn-outline-info ms-2"
            @click.prevent.self="enableEditMode"
          >
            <i class="bi bi-pencil"></i>
            Edit
          </button>
          <button
            type="button"
            class="btn btn-outline-danger ms-2"
            @click.prevent.self="close"
          >
            <i class="bi bi-x"></i>
            Cancel
          </button>
        </div>
      </div>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "booksys/components/WarningBox.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import InputCurrency from "./forms/inputs/InputCurrency.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import { confirm } from "booksys/components/bricks/DialogModal.js";

const props = defineProps(["userGroup", "editMode", "visible"]);
const emit = defineEmits(["update:visible", "save"]);

const store = useStore();

const errors = ref([]);
const userRoleList = ref([]);
const form = ref({});
const isEditMode = ref(false);
const title = ref("User Group");

const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const userRoles = computed(() => store.getters["user/userRoles"]);

const userRoleDescription = computed(() => {
  if (form.value.user_role_id == null) {
    return "";
  }
  const role = userRoles.value.find(
    (ur) => ur.user_role_id == form.value.user_role_id
  );
  return role.user_role_description;
});

watch(
  () => props.userGroup,
  () => {
    reloadProps();
  }
);

watch(
  () => props.editMode,
  (newValue) => {
    isEditMode.value = newValue;
    setTitle();
  }
);

watch(userRoles, (newValue) => {
  userRolesToList(newValue);
});

watch(
  () => props.visible,
  () => {
    reloadProps();
  }
);

const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");
const queryUserRoles = () => store.dispatch("user/queryUserRoles");
const saveUserGroup = (group) => store.dispatch("user/saveUserGroup", group);
const deleteUserGroup = (id) => store.dispatch("user/deleteUserGroup", id);

function reloadProps() {
  isEditMode.value = props.editMode;
  if (props.userGroup != null) {
    form.value = { ...props.userGroup };
    form.value.price_min = sprintf("%.2f", props.userGroup.price_min);
  } else {
    form.value = {};
  }
  setTitle();
}

function save() {
  saveUserGroup(form.value)
    .then(() => {
      errors.value = [];
      emit("save", form.value);
      form.value = {};
      close();
    })
    .catch((errs) => (errors.value = errs));
}

function close() {
  emit("update:visible", false);
}

function remove() {
  const name = form.value.user_group_name;
  const id = form.value.user_group_id;
  confirm({
    title: "Delete User Group",
    message: "Do you really want to delete user group " + name + "?",
  })
    .then((value) => {
      // delete user group
      if (value == true) {
        deleteUserGroup(id)
          .then(() => close())
          .catch((errs) => (errors.value = errs));
      }
    })
    .catch((err) => {
      errors.value = [err];
    });
}

function enableEditMode() {
  isEditMode.value = true;
  setTitle();
}

function setTitle() {
  if (isEditMode.value == true && form.value.user_group_id == null) {
    title.value = "New User Group";
  } else if (isEditMode.value == true) {
    title.value = "Edit User Group";
  } else {
    title.value = "User Group";
  }
}

function userRolesToList(roles) {
  if (roles == null) {
    roles = userRoles.value;
  }
  userRoleList.value = roles.map((ur) => {
    return {
      value: ur.user_role_id,
      text: ur.user_role_name,
    };
  });
}

queryConfiguration();

queryUserRoles()
  .then(() => userRolesToList())
  .catch((errs) => errors.value.push(...errs));

reloadProps();
</script>

<style scoped>
.block-footer {
  display: block;
}
</style>
