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
            :disabled="selectedItems.length == 0"
            class="btn btn-outline-info btn-sm me-1 mb-2"
            @click="showDetails"
          >
            View
          </button>
          <button
            :disabled="selectedItems.length == 0"
            class="btn btn-outline-danger btn-sm mb-2"
            @click="showDeleteUserDialog"
          >
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
          >
            <template #cell(license)="data">
              <div class="text-center">
                <i
                  class="bi bi-patch-check text-success"
                  v-if="data.cell == 1"
                ></i>
              </div>
            </template>
            <template #cell(balance)="data">
              <div class="text-end">
                {{ getBalance(data.row) }}
              </div>
            </template>
            <template #cell(status)="data">
              <input-select
                v-if="userGroupList.length > 1"
                v-model="data.row.status"
                :options="userGroupList"
                size="small"
                @change="groupChangeHandler(data.row.status, data.row.id)"
              />
              <div v-if="userGroupList.length == 1">
                {{ userGroupList[0].text }}
              </div>
              <div v-if="userGroupList.length == 0">-</div>
            </template>
            <template #cell(locked)="data">
              <div class="text-center">
                <button class="btn btn-light btn-sm" style="font-size: 0.8em">
                  <i
                    class="bi bi-lock-fill text-danger"
                    v-if="data.row.locked == 1"
                    @click.stop="unlock(data.row)"
                  />
                  <i
                    class="bi bi-unlock-fill text-success"
                    v-if="data.row.locked == 0"
                    @click.stop="lock(data.row)"
                  />
                </button>
              </div>
            </template>
          </table-module>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "booksys/components/WarningBox.vue";
import TableModule from "./bricks/TableModule.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import { confirm, info } from "booksys/components/bricks/DialogModal.js";

const store = useStore();

const errors = ref([]);
const userGroupList = ref([]);
const items = ref([]);
const selectedItems = ref([]);

const fields = ref([
  {
    key: "first_name",
    label: "First Name",
    sortable: true,
  },
  {
    key: "last_name",
    label: "Surname",
    sortable: true,
  },
  {
    key: "mobile",
    label: "Mobile",
  },
  {
    key: "email",
    label: "Email",
  },
  {
    key: "license",
    label: "License",
    sortable: true,
  },
  {
    key: "balance",
    label: "Balance",
    sortable: true,
    sortKey: "balance",
  },
  {
    key: "status",
    sortable: true,
    label: "User Group",
  },
  {
    key: "locked",
    sortable: true,
    label: "Locked",
  },
]);

const userListDetailed = computed(() => store.getters["user/userListDetailed"]);
const userGroups = computed(() => store.getters["user/userGroups"]);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);

watch(userGroups, (newUserGroups) => {
  userGroupToList(newUserGroups);
});

watch(userListDetailed, (newUserListDetailed) => {
  const userList = newUserListDetailed.map((e) => {
    e.balance = e.total_payment - e.total_heat_cost;
    return e;
  });
  items.value = userList;
});

const queryUserListDetailed = () => store.dispatch("user/queryUserListDetailed");
const queryUserGroups = () => store.dispatch("user/queryUserGroups");
const lockUser = (data) => store.dispatch("user/lockUser", data);
const deleteUser = (id) => store.dispatch("user/deleteUser", id);
const setUserGroup = (data) => store.dispatch("user/setUserGroup", data);
const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function dismissedHandler() {
  errors.value = [];
}

function getBalance(row) {
  return sprintf(
    "%.2f %s",
    row.total_payment - row.total_heat_cost,
    getCurrency.value
  );
}

function rowSelected(rows) {
  selectedItems.value = rows;
}

function lock(user) {
  lockUser({ user: user.id, locked: true }).catch(
    (errs) => (errors.value = errs)
  );
}

function unlock(user) {
  lockUser({ user: user.id, locked: false }).catch(
    (errs) => (errors.value = errs)
  );
}

function groupChangeHandler(userGroupId, userId) {
  const userGroupUpdate = {
    userId: userId,
    userGroupId: userGroupId,
  };
  setUserGroup(userGroupUpdate)
    .then(() => (errors.value = []))
    .catch((errs) => (errors.value = errs));
}

function showDeleteUserDialog() {
  const name =
    selectedItems.value[0].first_name + " " + selectedItems.value[0].last_name;
  confirm({
    title: "Delete User",
    message: `Do you really want to delete user ${name}?`,
  })
    .then((value) => {
      if (value == true) {
        deleteUser(selectedItems.value[0].id).catch(
          (errs) => (errors.value = errs)
        );
      }
    })
    .catch((error) => {
      errors.value = [error];
    });
}

function showDetails() {
  info({
    title: "Not implemented.",
    message: "This functionality is still missing.",
  });
}

function userGroupToList(groups) {
  if (groups == null) {
    groups = userGroups.value;
  }
  userGroupList.value = groups.map((ug) => {
    return {
      value: ug.user_group_id,
      text: ug.user_group_name,
    };
  });
}

queryConfiguration().catch((errs) => errors.value.push(...errs));

queryUserGroups()
  .then(() => userGroupToList())
  .catch((errs) => {
    errors.value.push(...errs);
  });

queryUserListDetailed()
  .then(() => (items.value = userListDetailed.value))
  .catch((errs) => {
    errors.value.push(...errs);
  });
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
