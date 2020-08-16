<template>
  <div class="login-view">
    <b-row>
      <b-col cols="3"></b-col>
      <b-col cols="6" class="text-center">
        <b-row class="row-dashboard">
          <b-col cols="4">
            <b-button to="/today" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_dashboard" id="today"></div>
              TODAY
            </b-button>
          </b-col>
          <b-col cols="4">
            <b-button to="/boat" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_boat"></div>
              BOAT
            </b-button>
          </b-col>
          <b-col cols="4">
            <b-button to="/ride" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_wakeboard"></div>
              RIDE
            </b-button>
          </b-col>
        </b-row>		
        <b-row class="row-dashboard">
          <b-col cols="4">
            <b-button to="/calendar" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_calendar_add"></div>
              BOOK
            </b-button>
          </b-col>
          <b-col cols="4">
            <b-button to="/schedule" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_calendar"></div>
              MY CAL
            </b-button>
          </b-col>
          <b-col cols="4">
            <b-button to="/admin" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_setting"></div>
              ADMIN
            </b-button>
          </b-col>
        </b-row>
        <b-row class="row-dashboard">
          <b-col cols="4">
            <b-button to="/account" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_user"></div>
              ACCOUNT
            </b-button>
          </b-col>
          <b-col cols="4">
            <b-button to="/info" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_info"></div>
              INFO
            </b-button>
          </b-col>
          <b-col cols="4">
            <b-button to="/logout" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_logout"></div>
              LOGOUT
            </b-button>
          </b-col>
        </b-row>
      </b-col>
      <b-col cols="3"></b-col>
    </b-row>
  </div>
</template>

<script>
import Vue from 'vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

export default Vue.extend({
  name: "DashboardAdmin",
  props: ['isMobile']
})
</script>