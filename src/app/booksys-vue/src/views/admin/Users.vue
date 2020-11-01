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
            User Management here
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
        User management here
      </b-card>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarAdminMobile from '@/components/NavbarAdminMobile';

export default Vue.extend({
  name: 'Users',
  components: {
    NavbarAdminMobile
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  }
})
</script>