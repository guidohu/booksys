<template>
  <b-modal
    id="userGroupModal"
    :title="title"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <b-form @submit="save">
      <b-row v-if="errors.length > 0">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <WarningBox :errors="errors" />
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <!-- User Group -->
          <h6 class="mt-2 mb-3">User Group</h6>
          <b-form-group
            id="user-group-name"
            label="Name"
            label-for="user-group-name"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="user-group-name"
              v-model="form.user_group_name"
              type="text"
              placeholder=""
              :disabled="!isEditMode"
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="user-group-description"
            label="Description"
            label-for="user-group-description"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="user-group-description"
              v-model="form.user_group_description"
              type="text"
              placeholder=""
              :disabled="!isEditMode"
            ></b-form-input>
          </b-form-group>
          <!-- User Role -->
          <h6 class="mt-5 mb-3">User Role and Permissions</h6>
          <b-form-group
            id="user-role-name"
            label="Role"
            label-for="user-role-name-select"
            description=""
            label-cols="3"
          >
            <b-form-select
              id="user-role-name-select"
              v-model="form.user_role_id"
              @change="roleChangeHandler($event)"
              :options="userRoleList"
              :disabled="!isEditMode"
            ></b-form-select>
          </b-form-group>
          <b-form-group
            id="user-role-description"
            label="Description"
            label-for="user-role-description-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="user-role-description-input"
              v-model="form.user_role_description"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <!-- Pricing -->
          <h6 class="mt-5 mb-3">Pricing</h6>
          <b-form-group
            id="price"
            label="Price"
            label-for="price-input"
            description=""
            label-cols="3"
          >
            <b-input-group>
              <b-form-input
                id="price-input"
                v-model="form.price_min"
                type="text"
                placeholder=""
                :disabled="!isEditMode"
              ></b-form-input>
              <b-input-group-append is-text>
                {{ getCurrency }}/min
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="price-description"
            label="Description"
            label-for="price-description-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="price-description-input"
              v-model="form.price_description"
              type="text"
              placeholder=""
              :disabled="!isEditMode"
            ></b-form-input>
          </b-form-group>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <div class="text-left d-inline">
        <b-button
          v-if="isEditMode == true"
          class="mr-1"
          type="button"
          variant="outline-danger"
          v-on:click="remove"
        >
          <b-icon-trash />
          Delete
        </b-button>
      </div>
      <div class="text-right d-inline">
        <b-button
          v-if="isEditMode == true"
          class="ml-4"
          type="button"
          variant="outline-info"
          v-on:click="save"
        >
          <b-icon-check />
          Save
        </b-button>
        <b-button
          v-if="isEditMode == false"
          class="ml-1"
          type="button"
          variant="outline-info"
          v-on:click="enableEditMode"
        >
          <b-icon-pencil />
          Edit
        </b-button>
        <b-button
          class="ml-1"
          type="button"
          variant="outline-danger"
          v-on:click="close"
        >
          <b-icon-x />
          Cancel
        </b-button>
      </div>
    </div>
  </b-modal>
</template>

<script>
import Vue from "vue";
import { mapGetters, mapActions } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "@/components/WarningBox";
import {
  BForm,
  BRow,
  BCol,
  BFormGroup,
  BFormInput,
  BFormSelect,
  BInputGroup,
  BInputGroupAppend,
  BButton,
  BIconTrash,
  BIconCheck,
  BIconPencil,
  BIconX,
  ModalPlugin,
} from "bootstrap-vue";

Vue.use(ModalPlugin);

export default {
  name: "UserGroupModal",
  props: ["userGroup", "editMode", "visible"],
  components: {
    WarningBox,
    BForm,
    BRow,
    BCol,
    BFormGroup,
    BFormInput,
    BFormSelect,
    BInputGroup,
    BInputGroupAppend,
    BButton,
    BIconTrash,
    BIconCheck,
    BIconPencil,
    BIconX,
  },
  data() {
    return {
      errors: [],
      userRoleList: [],
      form: {
        user_role_description: "",
      },
      isEditMode: false,
      title: "User Group",
    };
  },
  computed: {
    ...mapGetters("configuration", ["getCurrency"]),
    ...mapGetters("user", ["userRoles"]),
  },
  watch: {
    userGroup: function (newValue) {
      if (newValue != null) {
        this.form = { ...newValue };
        this.form.price_min = sprintf("%.2f", newValue.price_min);
        console.log("price_min", this.form);
      } else {
        this.form = {};
      }
    },
    editMode: function (newValue) {
      console.log("editMode changed to:", newValue);
      this.isEditMode = newValue;
      this.setTitle();
    },
    userRoles: function (newValue) {
      this.userRolesToList(newValue);
    },
  },
  methods: {
    save: function () {
      this.saveUserGroup(this.form)
        .then(() => {
          this.errors = [];
          this.close();
        })
        .catch((errors) => (this.errors = errors));
    },
    close: function () {
      this.setTitle();
      this.$emit("update:visible", false);
    },
    remove: function () {
      const name = this.form.user_group_name;
      const id = this.form.user_group_id;
      this.$bvModal
        .msgBoxConfirm(
          "Do you really want to delete user group " + name + "?",
          {
            title: "Delete User Group",
            size: "sm",
            buttonSize: "sm",
            okVariant: "danger",
            okTitle: "Delete",
            cancelTitle: "Cancel",
            footerClass: "p-2",
            hideHeaderClose: false,
            centered: true,
          }
        )
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
      console.log("Enable edit mode");
      this.isEditMode = true;
      this.setTitle();
    },
    setTitle: function () {
      console.log("Set title with:", this.isEditMode);
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
    roleChangeHandler: function (user_role_id) {
      const role = this.userRoles.find((ur) => ur.user_role_id == user_role_id);
      this.form.user_role_name = role.user_role_name;
      this.form.user_role_description = role.user_role_description;
    },
  },
  created() {
    this.queryConfiguration();

    this.queryUserRoles()
      .then(() => this.userRolesToList())
      .catch((errors) => this.errors.append(...errors));

    if (this.userGroup != null) {
      this.form = { ...this.userGroup };
    }

    this.isEditMode = this.editMode;
    this.setTitle();
  },
  mounted() {
    this.$root.$on("bv::modal::show", () => {
      this.isEditMode = this.editMode;
      this.setTitle();
    });
  },
};
</script>
