<template>
  <div v-if="isDesktop" class="display">
    <header v-if="userInfo!=null" class="welcome">
      Welcome, 
      <router-link to="/account">
        {{ userInfo.first_name }} {{ userInfo.last_name}}
      </router-link>
    </header>
    <DashboardAdmin v-if="getSessions!=null" v-bind:sessionData="getSessions"/>
  </div>
  <div v-else>
    <DashboardAdminMobile v-if="getSessions!=null" v-bind:sessionData="getSessions"/>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapGetters, mapActions } from 'vuex'
import { BooksysBrowser } from '@/libs/browser'
import DashboardAdmin from '../components/DashboardAdmin'
import DashboardAdminMobile from '../components/DashboardAdminMobile'
import moment from 'moment-timezone'


export default Vue.extend({
  name: 'Dashboard',
  components: {
    DashboardAdmin,
    DashboardAdminMobile
  },
  //props: ['isMobile'],
  // data: function () {

  // },
  computed: {
    // userInfo: {
    //   get () {
    //     return this.$store.getters.login.userInfo
    //   }
    // }
    ...mapGetters('login', [
      'userInfo'
    ]),
    ...mapGetters('sessions', [
      'getSessions'
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
    getTimeZone: function() {
      return "Europe/Berlin";
    }
  },
  created() {
    var dateStart = moment().tz(this.getTimeZone()).startOf('day');
    var dateEnd   = moment().tz(this.getTimeZone()).endOf('day');
    console.log("Query sessions from", dateStart, "to", dateEnd);
    this.querySessions({start: dateStart, end: dateEnd});
  }
})
</script>
