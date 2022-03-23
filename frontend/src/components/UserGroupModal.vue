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

<script>
import { mapGetters, mapActions } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "@/components/WarningBox";
import InputText from "./forms/inputs/InputText.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import InputCurrency from "./forms/inputs/InputCurrency.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import { confirm } from "@/components/bricks/DialogModal";

export default {
  name: "UserGroupModal",
  components: {
    WarningBox,
    InputText,
    InputSelect,
    ModalContainer,
    ModalHeader,
    ModalBody,
    InputCurrency,
    ModalFooter,
  },
  props: ["userGroup", "editMode", "visible"],
  emits: ["update:visible", "save"],
  data() {
    return {
      errors: [],
      userRoleList: [],
      form: {},
      isEditMode: false,
      title: "User Group",
    };
  },
  computed: {
    ...mapGetters("configuration", ["getCurrency"]),
    ...mapGetters("user", ["userRoles"]),
    userRoleDescription: function () {
      if (this.form.user_role_id == null) {
        return "";
      }
      const role = this.userRoles.find(
        (ur) => ur.user_role_id == this.form.user_role_id
      );
      return role.user_role_description;
    },
  },
  watch: {
    userGroup: function () {
      this.reloadProps();
    },
    editMode: function (newValue) {
      this.isEditMode = newValue;
      this.setTitle();
    },
    userRoles: function (newValue) {
      this.userRolesToList(newValue);
    },
    visible: function () {
      this.reloadProps();
    },
  },
  methods: {
    reloadProps: function () {
      this.isEditMode = this.editMode;
      if (this.userGroup != null) {
        this.form = { ...this.userGroup };
        this.form.price_min = sprintf("%.2f", this.userGroup.price_min);
      } else {
        this.form = {};
      }
      this.setTitle();
    },
    save: function () {
      this.saveUserGroup(this.form)
        .then(() => {
          this.errors = [];
          this.$emit("save", this.form);
          this.form = {};
          this.close();
        })
        .catch((errors) => (this.errors = errors));
    },
    close: function () {
      this.$emit("update:visible", false);
    },
    remove: function () {
      const name = this.form.user_group_name;
      const id = this.form.user_group_id;
      confirm({
        title: "Delete User Group",
        message: "Do you really want to delete user group " + name + "?",
      })
        .then((value) => {
          // delete user group
          if (value == true) {
            this.deleteUserGroup(id)
              .then(() => this.close())
              .catch((errors) => (this.errors = errors));
          }
        })
        .catch((err) => {
          this.errors = [err];
        });
    },
    ...mapActions("configuration", ["queryConfiguration"]),
    ...mapActions("user", [
      "queryUserRoles",
      "saveUserGroup",
      "deleteUserGroup",
    ]),
    enableEditMode: function () {
      this.isEditMode = true;
      this.setTitle();
    },
    setTitle: function () {
      if (this.isEditMode == true && this.form.user_group_id == null) {
        this.title = "New User Group";
      } else if (this.isEditMode == true) {
        this.title = "Edit User Group";
      } else {
        this.title = "User Group";
      }
    },
    userRolesToList: function (userRoles) {
      if (userRoles == null) {
        userRoles = this.userRoles;
      }
      this.userRoleList = userRoles.map((ur) => {
        return {
          value: ur.user_role_id,
          text: ur.user_role_name,
        };
      });
    },
  },
  created() {
    this.queryConfiguration();

    this.queryUserRoles()
      .then(() => this.userRolesToList())
      .catch((errors) => this.errors.append(...errors));

    this.reloadProps();
  },
};
</script>

<style scoped>
.block-footer {
  display: block;
}
</style>
