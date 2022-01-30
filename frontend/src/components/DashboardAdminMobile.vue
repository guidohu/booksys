<template>
  <div>
    <b-navbar
      fixed="top"
      toggleable="lg"
      type="light"
      variant="light"
    >
      <b-navbar-brand href="#">
        Dashboard
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse" />

      <b-collapse
        id="nav-collapse"
        is-nav
      >
        <b-navbar-nav>
          <b-nav-item href="#">
            <router-link to="/account">
              <b-icon-person />
              Account
            </router-link>
          </b-nav-item>
          <b-nav-item href="#">
            <router-link to="/logout">
              <b-icon-box-arrow-right />
              Logout
            </router-link>
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-row class="mt-3">
      <b-col cols="1" />
      <b-col cols="10">
        <b-row>
          <b-col cols="6">
            <b-button
              to="/boat"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <div
                id="boat"
                class="bc_icon bc_icon_dashboard bc_icon_boat"
              />
              BOAT
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button
              to="/ride"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <div class="bc_icon bc_icon_dashboard bc_icon_wakeboard" />
              RIDE
            </b-button>
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col cols="6">
            <b-button
              to="/today"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <BooksysPie
                :session-data="sessionData"
                :properties="properties"
              />
              TODAY
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button
              to="/calendar"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <div class="bc_icon bc_icon_dashboard bc_icon_calendar_add" />
              BOOK
            </b-button>
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col cols="6">
            <b-button
              to="/admin"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <div class="bc_icon bc_icon_dashboard bc_icon_setting" />
              ADMIN
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button
              to="/schedule"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <div class="bc_icon bc_icon_dashboard bc_icon_calendar" />
              MY CAL
            </b-button>
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col cols="6">
            <b-button
              to="/account"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <div class="bc_icon bc_icon_dashboard bc_icon_user" />
              ACCOUNT
            </b-button>
          </b-col>
          <b-col cols="6">
            <b-button
              to="/info"
              variant="light"
              class="btn btn-default btn-dashboard"
            >
              <div class="bc_icon bc_icon_dashboard bc_icon_info" />
              INFO
            </b-button>
          </b-col>
        </b-row>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import BooksysPie from "./Pie.vue";
import {
  BNavbar,
  BNavbarBrand,
  BNavbarToggle,
  BCollapse,
  BNavbarNav,
  BNavItem,
  BIconPerson,
  BIconBoxArrowRight,
  BRow,
  BCol,
  BButton,
} from "bootstrap-vue";

export default {
  name: "DashboardAdminMobile",
  components: {
    BooksysPie,
    BNavbar,
    BNavbarBrand,
    BNavbarToggle,
    BCollapse,
    BNavbarNav,
    BNavItem,
    BIconPerson,
    BIconBoxArrowRight,
    BRow,
    BCol,
    BButton,
  },
  props: ["sessionData"],
  data() {
    return {
      properties: {
        containerWidth: 80,
        containerHeight: 85,
        animate: false,
      },
    };
  },
};
</script>
