<template>
  <div>
    <div v-if="isDesktop" class="display">
      <b-row>
        <b-col>
          <div class="main-title">
            <h3>
              User Management
            </h3>
          </div>
        </b-col>
      </b-row>
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <b-card no-body>
            <UserManagementTabs class="content-max-height"/>
          </b-card>
        </b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button variant="outline-light" to="/admin">
          <b-icon icon="gear"></b-icon>
          ADMIN
        </b-button>
      </div>
    </div>
    <div v-else>
      <NavbarAdminMobile title="User Management"/>
      <b-card no-body>
        <UserManagementTabs/>
      </b-card>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarAdminMobile from '@/components/NavbarAdminMobile';
import UserManagementTabs from '@/components/UserManagementTabs';

export default Vue.extend({
  name: 'Users',
  components: {
    NavbarAdminMobile,
    UserManagementTabs
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  }
})
</script>

<style>
  .content-max-height {
    min-height: 480px;
    max-height: 480px;
    overflow-y: scroll;
  }
</style>