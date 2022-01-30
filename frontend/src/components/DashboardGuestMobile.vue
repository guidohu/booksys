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
        <b-row class="mt-3">
          <b-col
            offset="3"
            cols="6"
          >
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
import {
  BNavbar,
  BNavbarBrand,
  BNavbarToggle,
  BCollapse,
  BNavbarNav,
  BNavItem,
  BIconBoxArrowRight,
  BRow,
  BCol,
  BButton,
} from "bootstrap-vue";

export default {
  name: "DashboardGuestMobile",
  components: {
    BNavbar,
    BNavbarBrand,
    BNavbarToggle,
    BCollapse,
    BNavbarNav,
    BNavItem,
    BIconBoxArrowRight,
    BRow,
    BCol,
    BButton,
  },
};
</script>
