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

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "booksys/components/WarningBox.vue";
import TableModule from "./bricks/TableModule.vue";
import { confirm } from "booksys/components/bricks/DialogModal.js";
import UserGroupModal from "booksys/components/UserGroupModal.vue";

const store = useStore();

const errors = ref([]);
const items = ref([]);
const selectedItems = ref([]);
const userGroupEditMode = ref(false);
const showUserGroupModal = ref(false);

const userGroups = computed(() => store.getters["user/userGroups"]);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);

const fields = computed(() => [
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
    label: `Price ${getCurrency.value}/min`,
    sortable: true,
    sortKey: "price_min",
    formatter: (value) => {
      return sprintf("%.2f", value);
    },
  },
]);

watch(userGroups, () => {
  setRows();
});

const queryUserGroups = () => store.dispatch("user/queryUserGroups");
const deleteUserGroup = (id) => store.dispatch("user/deleteUserGroup", id);
const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function dismissedHandler() {
  errors.value = [];
}

function setRows() {
  items.value = userGroups.value;
}

function rowSelected(rows) {
  selectedItems.value = rows;
}

function showDeleteUserGroupDialog() {
  const name = selectedItems.value[0].user_group_name;
  const id = selectedItems.value[0].user_group_id;
  confirm({
    title: "Delete User Group",
    message: `Do you really want to delete user group ${name}?`,
  })
    .then((value) => {
      if (value == true) {
        deleteUserGroup(id).catch((errs) => (errors.value = errs));
      }
    })
    .catch((err) => {
      errors.value = [err];
    });
}

function showDetails() {
  userGroupEditMode.value = false;
  showUserGroupModal.value = true;
}

function editGroup() {
  userGroupEditMode.value = true;
  showUserGroupModal.value = true;
}

function newGroup() {
  userGroupEditMode.value = true;
  selectedItems.value = [];
  showUserGroupModal.value = true;
}

function resetSelection() {
  selectedItems.value = [];
}

queryConfiguration().catch((errs) => errors.value.push(...errs));

queryUserGroups()
  .then(() => setRows())
  .catch((errs) => errors.value.push(...errs));
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
