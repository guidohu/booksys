<template>
  <div v-if="getDbUpdateStatus != null && getDbUpdateStatus.updateAvailable == true">
    <DatabaseUpdateModal/>
  </div>
  <div v-else>
    <div v-if="isDesktop" class="display">
      <header v-if="userInfo!=null" class="welcome">
        Welcome, 
        <b-link to="/account" class="header-desktop">
          {{ userInfo.first_name }} {{ userInfo.last_name}}
        </b-link>
      </header>
      <DashboardAdmin v-if="getSessions!=null" v-bind:sessionData="getSessions"/>
    </div>
    <div v-else>
      <DashboardAdminMobile v-if="getSessions!=null" v-bind:sessionData="getSessions"/>
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapGetters, mapActions } from 'vuex'
import { BooksysBrowser } from '@/libs/browser'
import DashboardAdmin from '../components/DashboardAdmin'
import DashboardAdminMobile from '../components/DashboardAdminMobile'
import DatabaseUpdateModal from '../components/DatabaseUpdate'
import moment from 'moment-timezone'


export default Vue.extend({
  name: 'Dashboard',
  components: {
    DashboardAdmin,
    DashboardAdminMobile,
    DatabaseUpdateModal,
  },
  computed: {
    ...mapGetters('login', [
      'userInfo'
    ]),
    ...mapGetters('sessions', [
      'getSessions'
    ]),
    ...mapGetters('configuration', [
      'getDbUpdateStatus'
    ]),
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  },
  methods: {
    ...mapActions('sessions', [
      'querySessions'
    ]),
    ...mapActions('configuration', [
      'queryDbUpdateStatus'
    ]),
    getTimeZone: function() {
      return "Europe/Berlin";
    }
  },
  created() {
    var dateStart = moment().tz(this.getTimeZone()).startOf('day');
    var dateEnd   = moment().tz(this.getTimeZone()).endOf('day');
    console.log("Query sessions from", dateStart, "to", dateEnd);
    this.querySessions({start: dateStart, end: dateEnd});

    this.queryDbUpdateStatus()
  }
})
</script>
