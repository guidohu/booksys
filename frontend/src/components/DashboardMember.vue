<template>
  <b-row class="mt-5">
    <b-col cols="3"></b-col>
    <b-col cols="6" class="text-center">
      <b-row class="mt-4">
        <b-col offset="2" cols="4">
          <b-button to="/schedule" variant="light" class="btn btn-default btn-dashboard">
            <div class="bc_icon bc_icon_dashboard bc_icon_calendar"></div>
            MY CAL
          </b-button>
        </b-col>
        <b-col cols="4">
          <b-button to="/account" variant="light" class="btn btn-default btn-dashboard">
            <div class="bc_icon bc_icon_dashboard bc_icon_user"></div>
            ACCOUNT
          </b-button>
        </b-col>
      </b-row>
      <b-row class="mt-4">
        <b-col offset="2" cols="4">
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
</template>

<script>
import {
  BRow,
  BCol,
  BButton
} from 'bootstrap-vue';

export default {
  name: "DashboardMember",
  components: {
    BRow,
    BCol,
    BButton
  }
}
</script>