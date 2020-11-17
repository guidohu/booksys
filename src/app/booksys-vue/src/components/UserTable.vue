<template>
  <div class="text-left">
    <WarningBox v-if="errors.length > 0" :errors="errors"/>
    <div v-else>
      <b-button :disabled="selectedItems.length==0" size="sm" variant="outline-info" class="mr-1 mb-2">View</b-button>
      <b-button :disabled="selectedItems.length==0" size="sm" variant="outline-danger" class="mb-2">Delete</b-button>
      <b-table 
        striped
        responsive
        hover
        small
        sort-by="first_name"
        :items="items"
        :fields="fields"
        :selectable="true"
        select-mode="single"
        @row-selected="rowSelected"
      >
        <template #cell(license)="data">
          <div class="text-center">
            <b-icon v-if="data.item.license==1" icon="patch-check" variant="success"/>        
          </div>
        </template>
        <template #cell(balance)="data">
          <div class="text-right"> 
            {{ getBalance(data.item) }}       
          </div>
        </template>
        <template #cell(status)="data">
          <b-form-select v-if="userGroupList.length>1" v-model="data.item.status" @change="groupChangeHandler($event, data.item)" :options="userGroupList" size="sm" :data-index="data.item"></b-form-select>
          <div v-if="userGroupList.length==1">
            {{ userGroupList[0].text }}
          </div>
          <div v-if="userGroupList.length==0">
            -
          </div>
        </template>
        <template #cell(locked)="data">
          <div class="text-center">
            <b-button size="sm" style="font-size: 0.8em;" variant="light">
              <b-icon v-if="data.item.locked==1" v-on:click="unlock(data.item)" icon="lock-fill" variant="danger"/>
              <b-icon v-if="data.item.locked==0" v-on:click="lock(data.item)" icon="unlock-fill" variant="success"/>
            </b-button>
          </div>
        </template>
      </b-table>
      </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapGetters, mapActions } from 'vuex';
import { sprintf } from 'sprintf-js';
import WarningBox from '@/components/WarningBox';

export default Vue.extend({
  name: 'UserTable',
  components: {
    WarningBox
  },
  data() {
    return {
      errors: [],
      currency: 'TODO',
      userGroupList: [
      ],
      items: [],
      selectedItems: [],
      fields: [
        { 
          key: "first_name",
          label: "First Name",
          sortable: true
        },
        { 
          key: "last_name",
          label: "Surname",
          sortable: true
        },
        { 
          key: "mobile",
          label: "Mobile"
        },
        { 
          key: "email",
          label: "Email"
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
          sortByFormatted: true,
          formatter: (value, key, item) => { return sprintf("%.2f%s", (item.total_payment - item.total_heat_cost), this.currency) }
        },
        { 
          key: "status",
          sortable: true,
          label: "User Group"
        },
        { 
          key: "locked",
          sortable: true,
          label: "Locked"
        }
      ],
    };
  },
  computed: {
    ...mapGetters("user", [
      'userListDetailed',
      'userGroups'
    ])
  },
  watch: {
    userGroups: function(newUserGroups) {
      console.log("userGroups changed to", newUserGroups)
      this.userGroupToList(newUserGroups)
    },
    userListDetailed: function(newUserListDetailed) {
      this.items = newUserListDetailed;
    }
  },
  methods: {
    ...mapActions("user", [
      'queryUserListDetailed',
      'queryUserGroups',
      'lockUser',
      'unlockUser',
      'setUserGroup'
    ]),
    setRows: function() {
      this.items = this.userListDetailed;
    },
    getBalance: function(row) {
      return sprintf("%.2f %s", row.total_payment - row.total_heat_cost, this.currency)
    },
    rowSelected: function(rows) {
      this.selectedItems = rows;
    },
    isSelected: function() {
      return this.selectedItems.length > 0;
    },
    lock: function(user){
      this.lockUser(user.id)
      .catch((errors) => this.errors = errors);
    },
    unlock: function(user){
      this.unlockUser(user.id)
      .catch((errors) => this.errors = errors);
    },
    groupChangeHandler: function(userGroupId, user){
      const userGroupUpdate = {
        userId: user.id,
        userGroupId: userGroupId
      };
      this.setUserGroup(userGroupUpdate)
      .then(() => this.errors = [])
      .catch((errors) => this.errors = errors);
    },
    userGroupToList: function(userGroups){
      this.userGroupList = userGroups.map(ug => {
        return {
          value: ug.user_group_id,
          text: ug.user_group_name
        }
      });
    }
  },
  created() {
    // TODO get user groups
    this.queryUserGroups()
    .then(() => this.userGroupToList())
    .catch((errors) => this.errors = errors);

    this.queryUserListDetailed()
    .then(() => this.setRows())
    .catch((errors) => this.errors = errors);
  }
})
</script>