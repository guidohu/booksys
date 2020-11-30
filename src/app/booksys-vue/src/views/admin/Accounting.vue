<template>
  <div>
    <div v-if="isDesktop" class="display">
      <b-row>
        <b-col>
          <div class="main-title">
            <h3>
              Accounting
            </h3>
          </div>
        </b-col>
      </b-row>
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <b-card no-body>
            <AccountingTabs class="content-max-height"/>
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
      <NavbarAdminMobile title="Accounting"/>
      <b-card no-body>
        <AccountingTabs/>
      </b-card>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarAdminMobile from '@/components/NavbarAdminMobile';
import AccountingTabs from '@/components/AccountingTabs';

export default Vue.extend({
  name: 'Accounting',
  components: {
    NavbarAdminMobile,
    AccountingTabs
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