<template>
  <div v-if="isDesktop" class="display">
    <header v-if="userInfo!=null" class="welcome">
      Welcome, 
      <router-link to="/account">
        {{ userInfo.first_name }} {{ userInfo.last_name}}
      </router-link>
    </header>
    <!-- <div class="login-view">
      <div class="row">
        <div class="col-sm-3"></div>
        <div class="col-sm-6 text-center">
          <div class="row row-dashboard">
            <div class="col-sm-4">
              <a href="today.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_dashboard" id="today"></div>
                TODAY
              </a>
            </div>
            <div class="col-sm-4">
              <a href="boat.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_boat" id="boat"></div>
                BOAT
              </a>
            </div>
            <div class="col-sm-4">
              <a href="ride.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_wakeboard"></div>
                RIDE
              </a>
            </div>
          </div>		
          <div class="row row-dashboard">
            <div class="col-sm-4">
              <a href="calendar.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_calendar_add"></div>
                BOOK
              </a>
            </div>
            <div class="col-sm-4">
              <a href="schedule.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_calendar"></div>
                MY CAL
              </a>
            </div>
            <div class="col-sm-4">
              <a href="admin.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_setting"></div>
                ADMIN
              </a>
            </div>
          </div>
          <div class="row row-dashboard">
            <div class="col-sm-4">
              <a href="account.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_user"></div>
                ACCOUNT
              </a>
            </div>
            <div class="col-sm-4">
              <a href="info.html" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_info"></div>
                INFO
              </a>
            </div>
            <div class="col-sm-4">
              <a href="logout.php" type="button" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_logout"></div>
                LOGOUT
              </a>
            </div>
          </div>
        </div>
        <div class="col-sm-3"></div>
      </div>
    </div> -->
    <DashboardAdmin/>
  </div>
  <div v-else>
    <DashboardAdminMobile/>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { BooksysBrowser } from '@/libs/browser'
import DashboardAdmin from '../components/DashboardAdmin'
import DashboardAdminMobile from '../components/DashboardAdminMobile'


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
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  },
  methods: {

  }
})
</script>
