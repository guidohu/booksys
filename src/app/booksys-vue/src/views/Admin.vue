<template>
  <div>
    <div v-if="isDesktop" class="display">
      <b-row>
        <b-col>
          <div class="main-title">
            <h3>
              Administrator Menu
            </h3>
          </div>
        </b-col>
      </b-row>
      <b-row class="mt-5 ml-1 mr-1">
        <b-col cols="3"></b-col>
        <b-col cols="6" class="text-center">
          <b-row class="mt-4">
            <b-col cols="4" offset="2">
              <b-button to="/users" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_users"></div>
                USERS
              </b-button>
            </b-col>
            <b-col cols="4">
              <b-button to="/accounting" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_money"></div>
                ACCOUNTING
              </b-button>
            </b-col>
          </b-row>		
          <b-row class="mt-4">
            <b-col cols="4" offset="2">
              <b-button to="/settings" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_wrench"></div>
                SETTINGS
              </b-button>
            </b-col>
            <b-col cols="4">
              <b-button to="/logs" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_statistics"></div>
                LOGS
              </b-button>
            </b-col>
          </b-row>
        </b-col>
        <b-col cols="3"></b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button variant="outline-light" to="/dashboard">
          <b-icon-house></b-icon-house>
          HOME
        </b-button>
      </div>
    </div>
    <div v-else>
      <b-navbar fixed="top" toggleable="lg" type="light" variant="light">
        <b-navbar-brand href="#">Admin</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item href="#">
              <router-link to="/dashboard">
                <b-icon icon="house"></b-icon>
                Home
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

      <b-row class="mt-3"> 
        <b-col cols="1"></b-col>
        <b-col cols="10">
          <b-row>
            <b-col cols="6">
              <b-button to="/users" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_users"></div>
                USERS
              </b-button>
            </b-col>
            <b-col cols="6">
              <b-button to="/accounting" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_money"></div>
                ACCOUNTING
              </b-button>
            </b-col>
          </b-row>
          <b-row class="mt-3">
            <b-col cols="6">
              <b-button to="/settings" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_wrench"></div>
                SETTINGS
              </b-button>
            </b-col>
            <b-col cols="6">
              <b-button to="/logs" variant="light" class="btn btn-default btn-dashboard">
                <div class="bc_icon bc_icon_dashboard bc_icon_statistics"></div>
                LOGS
              </b-button>
            </b-col>
          </b-row>
        </b-col>
      </b-row>
    </div>
    
  </div>
</template>

<script>
import Vue from 'vue';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';
import { BooksysBrowser } from '@/libs/browser';

export default Vue.extend({
  name: "Admin",
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  }
})
</script>