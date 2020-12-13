<template>
  <div>
    <div v-if="isDesktop" class="display">
      <b-row>
        <b-col>
          <div class="main-title">
            <h3>
              Logs
            </h3>
          </div>
        </b-col>
      </b-row>
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <LogCard class="content-max-height"/>
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
      <NavbarAdminMobile title="Logs"/>
      <LogCard/>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarAdminMobile from '@/components/NavbarAdminMobile';
import LogCard from '@/components/LogCard';

export default Vue.extend({
  name: 'Logs',
  components: {
    NavbarAdminMobile,
    LogCard
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
    max-height: 480px;
    overflow-y: scroll;
  }
</style>