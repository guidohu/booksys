<template>
  <b-tabs content-class="mt-3">
    <b-tab active>
      <div slot="title">
        <p class="h6 mt-2 mb-2"><b-icon icon="person"></b-icon>Users</p>
      </div>
      <div :class="tabClass">
        <UserTable/>
      </div>
    </b-tab>
    <b-tab>
      <div slot="title">
        <p class="h6 mt-2 mb-2"><b-icon icon="people"></b-icon>User Groups</p>
      </div>
      <div :class="tabClass">
        <UserGroupTable/>
      </div>
    </b-tab>           
  </b-tabs>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import UserTable from '@/components/UserTable';
import UserGroupTable from '@/components/UserGroupTable';

export default Vue.extend({
  name: 'UserManagementTabs',
  components: {
    UserTable,
    UserGroupTable,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    },
    tabClass: function (){
      if(!BooksysBrowser.isMobile()){
        return "user-management-tab-limited-height";
      }else{
        return "";
      }
    }
  }
})
</script>

<style scoped>
  .user-management-tab-limited-height {
    max-height: 400px;
    overflow: scroll;
  }
</style>