<template>
<div v-if="isDesktop" class="display">
  <main-title title-name="Schedule"/>
  <b-row class="ml-1 mr-1">
    <b-col>
      <div class="accordion" role="tablist">
        <b-card no-body class="mb-1">
          <b-card-header header-tag="header" class="p-1" role="tab">
            <b-button block v-b-toggle.accordion-1 variant="outline-dark">Past Sessions</b-button>
          </b-card-header>
          <b-collapse id="accordion-1" accordion="my-accordion" role="tabpanel">
            <b-card-body>
              <UserSessionsTable v-if="pastSessions.length > 0" :userSessions="pastSessions"/>
              <b-card-text v-else class="text-center">No sessions to display</b-card-text>
            </b-card-body>
          </b-collapse>
        </b-card>

        <b-card no-body class="mb-1">
          <b-card-header header-tag="header" class="p-1" role="tab">
            <b-button block v-b-toggle.accordion-2 variant="outline-dark">Upcoming Sessions</b-button>
          </b-card-header>
          <b-collapse id="accordion-2" visible accordion="my-accordion" role="tabpanel">
            <b-card-body>
              <UserSessionsTable v-if="upcomingSessions.length > 0" :userSessions="upcomingSessions" :showCancel="true" @cancel="cancelSessionHandler"/>
              <b-card-text v-else class="text-center">No sessions to display</b-card-text>
            </b-card-body>
          </b-collapse>
        </b-card>
      </div>
    </b-col>
  </b-row>
  <div class="bottom mr-2">
    <b-button variant="outline-light" to="/dashboard">
      <b-icon-house></b-icon-house>
      HOME
    </b-button>
  </div>
</div>
<div v-else>
  <NavbarMobile title="Schedule"/>
  <div>
    <b-card>
      <b-row>
        <b-col>
          <b-tabs content-class="mt-3 text-left">
            <b-tab title="Upcoming Sessions" active>
              <UserSessionsTable v-if="upcomingSessions.length > 0" :userSessions="upcomingSessions" :showCancel="true" @cancel="cancelSessionHandler"/>
              <div v-else class="text-center">No sessions to display</div>
            </b-tab>
            <b-tab title="Past Sessions">
              <UserSessionsTable v-if="pastSessions.length > 0" :userSessions="pastSessions"/>
              <div v-else class="text-center">No sessions to display</div>
            </b-tab>
          </b-tabs>
        </b-col>
      </b-row>
    </b-card>
  </div>
</div>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import { BooksysBrowser } from '@/libs/browser';
import NavbarMobile from '@/components/NavbarMobile';
import UserSessionsTable from '@/components/UserSessionsTable';
import MainTitle from '@/components/MainTitle';
import {
  BRow,
  BCol,
  BCard,
  BCardHeader,
  BCollapse,
  BCardBody,
  BCardText,
  BButton,
  BIconHouse,
  BTabs,
  BTab
} from 'bootstrap-vue';


export default Vue.extend({
  name: 'Schedule',
  components: {
    NavbarMobile,
    MainTitle,
    UserSessionsTable,
    BRow,
    BCol,
    BCard,
    BCardHeader,
    BCollapse,
    BCardBody,
    BCardText,
    BButton,
    BIconHouse,
    BTabs,
    BTab
  },
  computed: {
    isDesktop: function() {
      return !BooksysBrowser.isMobile()
    },
    ...mapGetters('user', [
      'userSchedule'
    ]),
    upcomingSessions: function() {
      return this.userSchedule.sessions;
    },
    pastSessions: function() {
      return this.userSchedule.sessions_old;
    }
  },
  methods: {
    ...mapActions('user', [
      'queryUserSchedule',
      'cancelSession'
    ]),
    cancelSessionHandler(value){
      console.log("Cancel session with ID:", value, "in parent");
      this.cancelSession(value);
    }
  },
  created () {
    this.queryUserSchedule()
      .then(() => {
        console.log("Retrieved userSchedule");
        console.log(this.userSchedule);
      })
      .catch(error => {
        console.error(error)
      })
  }
})
</script>

<style scoped>

  .card-body {
    max-height: 380px;
    overflow-x: scroll;
  }
</style>