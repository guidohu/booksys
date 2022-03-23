<template>
  <div class="text-left">
    <warning-box
      v-if="errors.length > 0"
      :errors="errors"
      dismissible="true"
      @dismissed="dismissedHandler"
    />
    <div class="box box-fix-height">
      <div class="row box-fix-content">
        <div class="col-12 ms-1">
          <button
            type="button"
            class="btn btn-outline-info btn-sm me-1 mb-2"
            @click="newGroup"
          >
            <i class="bi bi-plus"></i>
            New
          </button>
          <button
            v-if="selectedItems.length > 0"
            type="button"
            class="btn btn-outline-info btn-sm me-1 mb-2"
            @click="showDetails"
          >
            <i class="bi bi-eye"></i>
            View
          </button>
          <button
            v-if="selectedItems.length > 0"
            type="button"
            class="btn btn-outline-info btn-sm me-1 mb-2"
            @click="editGroup"
          >
            <i class="bi bi-pencil"></i>
            Edit
          </button>
          <button
            v-if="selectedItems.length > 0"
            type="button"
            class="btn btn-outline-danger btn-sm me-1 mb-2"
            @click="showDeleteUserGroupDialog"
          >
            <i class="bi bi-trash"></i>
            Delete
          </button>
        </div>
      </div>
      <div class="row box-flex-content">
        <div class="col-12">
          <table-module
            :columns="fields"
            :rows="items"
            :selectable="true"
            select-mode="single"
            @select-row="rowSelected"
          />
        </div>
      </div>
    </div>
    <div>
      <user-group-modal
        v-model:visible="showUserGroupModal"
        :user-group="selectedItems[0]"
        :edit-mode="userGroupEditMode"
        @save="resetSelection()"
      />
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "@/components/WarningBox";
import TableModule from "./bricks/TableModule.vue";
import { confirm } from "@/components/bricks/DialogModal";
import UserGroupModal from "@/components/UserGroupModal";

export default {
  name: "UserGroupTable",
  components: {
    WarningBox,
    TableModule,
    UserGroupModal,
  },
  data() {
    return {
      errors: [],
      userGroupList: [],
      items: [],
      selectedItems: [],
      fields: [
        {
          key: "user_group_name",
          label: "User Group",
          sortable: true,
        },
        {
          key: "user_role_description",
          label: "User Role",
          sortable: true,
        },
        {
          key: "price_min",
          label: "Price " + this.getCurrency + "/min",
          sortable: true,
          sortKey: "price_min",
          formatter: (value) => {
            return sprintf("%.2f", value);
          },
        },
      ],
      userGroupEditMode: false,
      showUserGroupModal: false,
    };
  },
  computed: {
    ...mapGetters("user", ["userListDetailed", "userGroups"]),
    ...mapGetters("configuration", ["getCurrency"]),
  },
  watch: {
    userGroups: function (newUserGroups) {
      this.setRows(newUserGroups);
    },
    getCurrency: function (currency) {
      this.fields[2].label = "Price " + currency + "/min";
    },
  },
  methods: {
    ...mapActions("user", [
      "queryUserListDetailed",
      "queryUserGroups",
      "deleteUserGroup",
      "setUserGroup",
    ]),
    ...mapActions("configuration", ["queryConfiguration"]),
    dismissedHandler: function () {
      this.errors = [];
    },
    setRows: function () {
      this.items = this.userGroups;
    },
    rowSelected: function (rows) {
      this.selectedItems = rows;
    },
    isSelected: function () {
      return this.selectedItems.length > 0;
    },
    showDeleteUserGroupDialog: function () {
      const name = this.selectedItems[0].user_group_name;
      const id = this.selectedItems[0].user_group_id;
      confirm({
        title: "Delete User Group",
        message: "Do you really want to delete user group " + name + "?",
      })
        .then((value) => {
          // delete user group
          if (value == true) {
            this.deleteUserGroup(id).catch((errors) => (this.errors = errors));
          }
        })
        .catch((err) => {
          this.errors = [err];
        });
    },
    showDetails: function () {
      this.userGroupEditMode = false;
      this.showUserGroupModal = true;
    },
    editGroup: function () {
      this.userGroupEditMode = true;
      this.showUserGroupModal = true;
    },
    newGroup: function () {
      this.userGroupEditMode = true;
      this.selectedItems = [];
      this.showUserGroupModal = true;
    },
    resetSelection: function () {
      this.selectedItems = [];
    },
  },
  created() {
    this.queryConfiguration().catch((errors) => this.errors.push(...errors));

    this.queryUserGroups()
      .then(() => this.setRows())
      .catch((errors) => this.errors.push(...errors));
  },
};
</script>

<style scoped>
.box {
  display: flex;
  flex-flow: column;
  height: 100%;
}

.box-fix-height {
  max-height: 430px;
}

@media (max-width: 992px) {
  .box-fix-height {
    max-height: 90vh;
  }
}

.box-fix-content {
  flex: 0 0 auto;
}

.box-flex-content {
  flex: 1 1 auto;
  overflow: scroll;
}
</style>
