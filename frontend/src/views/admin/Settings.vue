<template>
  <div>
    <div v-if="isDesktop" class="display">
      <b-row>
        <b-col>
          <div class="main-title">
            <h3>
              Settings
            </h3>
          </div>
        </b-col>
      </b-row>
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <SettingsCard class="content-max-height"/>
        </b-col>
      </b-row>
    </div>
    <div v-else>
      <NavbarAdminMobile title="Settings"/>
      <SettingsCard/>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarAdminMobile from '@/components/NavbarAdminMobile';
import SettingsCard from '@/components/SettingsCard';
import {
  BRow,
  BCol
} from 'bootstrap-vue';

export default Vue.extend({
  name: 'Settings',
  components: {
    NavbarAdminMobile,
    SettingsCard,
    BRow,
    BCol
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
    max-height: 520px;
    overflow-y: scroll;
  }
</style>