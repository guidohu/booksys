<template>
  <div class="text-left">
    <WarningBox
      v-if="errors.length > 0"
      :errors="errors"
      dismissible="true"
      @dismissed="dismissedHandler"
    />
    <div v-else>
      <b-row>
        <b-col cols="12">
          <b-button
            :disabled="selectedItems.length == 0"
            size="sm"
            variant="outline-info"
            class="mr-1 mb-2"
            @click="showDetails"
          >
            View
          </b-button>
          <b-button
            :disabled="selectedItems.length == 0"
            size="sm"
            variant="outline-danger"
            class="mb-2"
            @click="showDeleteUserDialog"
          >
            Delete
          </b-button>
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="12">
          <b-table
            striped
            hover
            response
            small
            sticky-header
            borderless
            sort-by="first_name"
            :items="items"
            :fields="fields"
            :selectable="true"
            select-mode="single"
            @row-selected="rowSelected"
          >
            <template #cell(license)="data">
              <div class="text-center">
                <b-icon-patch-check
                  v-if="data.item.license == 1"
                  variant="success"
                />
              </div>
            </template>
            <template #cell(balance)="data">
              <div class="text-right">
                {{ getBalance(data.item) }}
              </div>
            </template>
            <template #cell(status)="data">
              <b-form-select
                v-if="userGroupList.length > 1"
                v-model="data.item.status"
                :options="userGroupList"
                size="sm"
                :data-index="data.item"
                @change="groupChangeHandler($event, data.item)"
              />
              <div v-if="userGroupList.length == 1">
                {{ userGroupList[0].text }}
              </div>
              <div v-if="userGroupList.length == 0">-</div>
            </template>
            <template #cell(locked)="data">
              <div class="text-center">
                <b-button size="sm" style="font-size: 0.8em" variant="light">
                  <b-icon-lock-fill
                    v-if="data.item.locked == 1"
                    variant="danger"
                    @click="unlock(data.item)"
                  />
                  <b-icon-unlock-fill
                    v-if="data.item.locked == 0"
                    variant="success"
                    @click="lock(data.item)"
                  />
                </b-button>
              </div>
            </template>
          </b-table>
        </b-col>
      </b-row>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "@/components/WarningBox";
import {
  BRow,
  BCol,
  BButton,
  BTable,
  BIconPatchCheck,
  BIconLockFill,
  BIconUnlockFill,
  BFormSelect,
} from "bootstrap-vue";

export default {
  name: "UserTable",
  components: {
    WarningBox,
    BRow,
    BCol,
    BButton,
    BTable,
    BIconPatchCheck,
    BIconLockFill,
    BIconUnlockFill,
    BFormSelect,
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
    groupChangeHandler: function (userGroupId, user) {
      const userGroupUpdate = {
        userId: user.id,
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
      this.boxTwo = "";
      this.$bvModal
        .msgBoxConfirm("Do you really want to delete user " + name + "?", {
          title: "Delete User",
          size: "sm",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "Delete",
          cancelTitle: "Cancel",
          footerClass: "p-2",
          hideHeaderClose: false,
          centered: true,
        })
        .then((value) => {
          // delete user
          if (value == true) {
            this.deleteUser(this.selectedItems[0].id).catch(
              (errors) => (this.errors = errors)
            );
          }
        })
        .catch((err) => {
          this.errors = [err];
        });
    },
    showDetails: function () {
      console.log("TODO show details of user");
      this.$bvModal
        .msgBoxOk("This functionality is still missing.", {
          title: "Not implemented",
          size: "sm",
          buttonSize: "sm",
          okVariant: "success",
          headerClass: "p-2 border-bottom-0",
          footerClass: "p-2 border-top-0",
          centered: true,
        })
        .then((value) => {
          this.errors = value;
        })
        .catch((err) => {
          // An error occurred
          this.errors = [err];
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
.b-table-sticky-header {
  max-height: 340px;
}
</style>
