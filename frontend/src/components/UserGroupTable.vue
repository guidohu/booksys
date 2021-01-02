<template>
  <div class="text-left">
    <WarningBox v-if="errors.length > 0" :errors="errors" dismissible="true" @dismissed="dismissedHandler"/>
    <div v-else>
      <UserGroupModal
        :userGroup="selectedItems[0]"
        :editMode="userGroupEditMode"
      />
      <b-row>
        <b-col cols="12">
          <b-button v-on:click="newGroup" size="sm" variant="outline-info" class="mr-1 mb-2">
            <b-icon-plus/>
            New
          </b-button>
          <b-button v-if="selectedItems.length>0" v-on:click="showDetails" size="sm" variant="outline-info" class="mr-1 mb-2">
            <b-icon-eye/>
            View
          </b-button>
          <b-button v-if="selectedItems.length>0" v-on:click="editGroup" size="sm" variant="outline-info" class="mr-1 mb-2">
            <b-icon-pencil/>
            Edit
          </b-button>
          <b-button v-if="selectedItems.length>0" v-on:click="showDeleteUserGroupDialog" size="sm" variant="outline-danger" class="mb-2">
            <b-icon-trash/>
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
          </b-table>
        </b-col>
      </b-row>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import { sprintf } from 'sprintf-js';
import WarningBox from '@/components/WarningBox';
import UserGroupModal from '@/components/UserGroupModal';
import {
  BRow,
  BButton,
  BIconPlus,
  BIconEye,
  BIconPencil,
  BIconTrash,
  BCol,
  BTable
} from "bootstrap-vue";

export default {
  name: 'UserGroupTable',
  components: {
    WarningBox,
    UserGroupModal,
    BRow,
    BButton,
    BIconPlus,
    BIconEye,
    BIconPencil,
    BIconTrash,
    BCol,
    BTable
  },
  data() {
    return {
      errors: [],
      userGroupList: [
      ],
      items: [],
      selectedItems: [],
      fields: [
        { 
          key: "user_group_name",
          label: "User Group",
          sortable: true
        },
        { 
          key: "user_role_description",
          label: "User Role",
          sortable: true
        },
        { 
          key: "price_min",
          label: "Price " + this.getCurrency + "/min",
          sortable: true,
          sortKey: "price_min",
          formatter: (value) => { return sprintf("%.2f", value) }
        }
      ],
      userGroupEditMode: false
    };
  },
  computed: {
    ...mapGetters("user", [
      'userListDetailed',
      'userGroups'
    ]),
    ...mapGetters("configuration", [
      'getCurrency'
    ])
  },
  watch: {
    userGroups: function(newUserGroups) {
      this.setRows(newUserGroups);
    },
    getCurrency: function(currency){
      this.fields[2].label = "Price " + currency + "/min";
    }
  },
  methods: {
    ...mapActions("user", [
      'queryUserListDetailed',
      'queryUserGroups',
      'lockUser',
      'unlockUser',
      'deleteUserGroup',
      'setUserGroup'
    ]),
    ...mapActions("configuration", [
      'queryConfiguration'
    ]),
    dismissedHandler: function() {
      this.errors = [];
    },
    setRows: function() {
      this.items = this.userGroups;
    },
    getBalance: function(row) {
      return sprintf("%.2f %s", row.total_payment - row.total_heat_cost, this.getCurrency)
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
    showDeleteUserGroupDialog: function(){
      const name = this.selectedItems[0].user_group_name;
      const id   = this.selectedItems[0].user_group_id;
      this.$bvModal.msgBoxConfirm('Do you really want to delete user group '+name+'?', {
        title: 'Delete User Group',
        size: 'sm',
        buttonSize: 'sm',
        okVariant: 'danger',
        okTitle: 'Delete',
        cancelTitle: 'Cancel',
        footerClass: 'p-2',
        hideHeaderClose: false,
        centered: true
      })
      .then(value => {
        // delete user group
        if(value == true){
          this.deleteUserGroup(id)
          .catch((errors) => this.errors = errors);
        }
      })
      .catch(err => {
        this.errors = [ err ]
      })
    },
    showDetails: function() {
      this.userGroupEditMode = false;
      this.$bvModal.show('userGroupModal');
    },
    editGroup: function() {
      console.log("TODO implement editGroup");
      this.userGroupEditMode = true;
      this.$bvModal.show('userGroupModal');
    },
    newGroup: function() {
      console.log("TODO implement newGroup");
      this.userGroupEditMode = true;
      this.selectedItems = [];
      this.$bvModal.show('userGroupModal');
    }
  },
  created() {
    this.queryConfiguration()
    .catch((errors) => this.errors.push(...errors));

    this.queryUserGroups()
    .then(() => this.setRows())
    .catch((errors) => this.errors.push(...errors));
  }
}
</script>

<style scoped>
  .b-table-sticky-header {
    max-height: 340px;
  }
</style>