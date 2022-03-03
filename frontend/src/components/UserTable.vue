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

<script>
import { mapGetters, mapActions } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "@/components/WarningBox";
import TableModule from "./bricks/TableModule.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import { confirm, info } from "@/components/bricks/DialogModal";

export default {
  name: "UserTable",
  components: {
    WarningBox,
    TableModule,
    InputSelect,
  },
  data() {
    return {
      errors: [],
      userGroupList: [],
      items: [],
      selectedItems: [],
      fields: [
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
          //formatter: (value, key, item) => { return sprintf("%.2f%s", (item.total_payment - item.total_heat_cost), this.getCurrency) }
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
      ],
    };
  },
  computed: {
    ...mapGetters("user", ["userListDetailed", "userGroups"]),
    ...mapGetters("configuration", ["getCurrency"]),
  },
  watch: {
    userGroups: function (newUserGroups) {
      this.userGroupToList(newUserGroups);
    },
    userListDetailed: function (newUserListDetailed) {
      // precalculate the balance as a separate field
      // to simplify sorting
      const userList = newUserListDetailed.map((e) => {
        e.balance = e.total_payment - e.total_heat_cost;
        return e;
      });
      this.items = userList;
    },
  },
  methods: {
    ...mapActions("user", [
      "queryUserListDetailed",
      "queryUserGroups",
      "lockUser",
      "unlockUser",
      "deleteUser",
      "setUserGroup",
    ]),
    ...mapActions("configuration", ["queryConfiguration"]),
    dismissedHandler: function () {
      this.errors = [];
    },
    setRows: function () {
      this.items = this.userListDetailed;
    },
    getBalance: function (row) {
      return sprintf(
        "%.2f %s",
        row.total_payment - row.total_heat_cost,
        this.getCurrency
      );
    },
    rowSelected: function (rows) {
      this.selectedItems = rows;
    },
    isSelected: function () {
      return this.selectedItems.length > 0;
    },
    lock: function (user) {
      this.lockUser(user.id).catch((errors) => (this.errors = errors));
    },
    unlock: function (user) {
      this.unlockUser(user.id).catch((errors) => (this.errors = errors));
    },
    groupChangeHandler: function (userGroupId, userId) {
      console.log("Change to: userGroupId", userGroupId, "for user", userId);
      const userGroupUpdate = {
        userId: userId,
        userGroupId: userGroupId,
      };
      this.setUserGroup(userGroupUpdate)
        .then(() => (this.errors = []))
        .catch((errors) => (this.errors = errors));
    },
    showDeleteUserDialog: function () {
      const name =
        this.selectedItems[0].first_name +
        " " +
        this.selectedItems[0].last_name;
      confirm({
        title: "Delete User",
        message: "Do you really want to delete user " + name + "?",
      })
        .then((value) => {
          console.log("Returned with", value);
          // delete user
          if (value == true) {
            this.deleteUser(this.selectedItems[0].id).catch(
              (errors) => (this.errors = errors)
            );
          }
        })
        .catch((error) => {
          console.warn("Dialog returned error", error);
          this.errors = [error];
        });
    },
    showDetails: function () {
      info({
        title: "Not implemented.",
        message: "This functionality is still missing.",
      });
    },
    userGroupToList: function (userGroups) {
      if (userGroups == null) {
        userGroups = this.userGroups;
      }
      this.userGroupList = userGroups.map((ug) => {
        return {
          value: ug.user_group_id,
          text: ug.user_group_name,
        };
      });
    },
  },
  created() {
    this.queryConfiguration().catch((errors) => this.errors.push(...errors));

    // TODO get user groups
    this.queryUserGroups()
      .then(() => this.userGroupToList())
      .catch((errors) => {
        console.log("UserTable (queryUserGroups) got errors:", this.errors);
        this.errors.push(...errors);
      });

    this.queryUserListDetailed()
      .then(() => this.setRows())
      .catch((errors) => {
        console.log(
          "UserTable (queryUserListDetailed) got errors:",
          this.errors
        );
        this.errors.push(...errors);
      });
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