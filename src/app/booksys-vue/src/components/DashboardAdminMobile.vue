<template>
  <div>
    <b-navbar fixed="top" toggleable="lg" type="light" variant="light">
      <b-navbar-brand href="#">Dashboard</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item href="#">
            <router-link to="/account">
              <b-icon icon="person"></b-icon>
              Account
              </router-link>
            </b-nav-item>
          <b-nav-item href="#">
            <router-link to="/logout">
              <b-icon icon="box-arrow-right"></b-icon>
              Logout
            </router-link>
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-row> 
      <b-col cols="1"></b-col>
      <b-col cols="10">
        <b-row>
          <b-col cols="6">
            <b-button to="/boat" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_boat" id="boat"></div>
              BOAT
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button to="/ride" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_wakeboard"></div>
              RIDE
            </b-button>
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col cols="6">
            <b-button to="/today" variant="light" class="btn btn-default btn-dashboard">
              <!-- <div class="bc_icon bc_icon_dashboard bc_icon_boat" id="boat"></div> -->
              TODAY
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button to="/calendar" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_calendar_add"></div>
              BOOK
            </b-button>
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col cols="6">
            <b-button to="/admin" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_setting"></div>
              ADMIN
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button to="/schedule" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_calendar"></div>
              MY CAL
            </b-button>
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col cols="6">
            <b-button to="/account" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_user"></div>
              ACCOUNT
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button to="/info" variant="light" class="btn btn-default btn-dashboard">
              <div class="bc_icon bc_icon_dashboard bc_icon_info"></div>
              INFO
            </b-button>
          </b-col>
        </b-row>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import Vue from 'vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

export default Vue.extend({
  name: "DashboardAdminMobile"
})
</script>