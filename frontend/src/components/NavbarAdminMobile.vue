<template>
  <b-navbar fixed="top" toggleable="lg" type="light" variant="light">
    <b-navbar-brand href="#">{{title}}</b-navbar-brand>

    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item href="#">
          <router-link to="/dashboard">
            <b-icon-house/>
            Home
          </router-link>
        </b-nav-item>
        <b-nav-item href="#">
          <router-link to="/admin">
            <b-icon-gear/>
            Administrator Menu
          </router-link>
        </b-nav-item>
        <b-nav-item href="#">
          <router-link to="/logout">
            <b-icon-box-arrow-right/>
            Logout
          </router-link>
        </b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
import {
  BNavbar,
  BNavbarBrand,
  BNavbarToggle,
  BCollapse,
  BNavbarNav,
  BNavItem,
  BIconHouse,
  BIconGear,
  BIconBoxArrowRight
} from 'bootstrap-vue';

export default {
  name: 'NavbarAdminMobile',
  components: {
    BNavbar,
    BNavbarBrand,
    BNavbarToggle,
    BCollapse,
    BNavbarNav,
    BNavItem,
    BIconHouse,
    BIconGear,
    BIconBoxArrowRight
  },
  props: ['title']
}
</script>

<style scoped>
.navbar-nav>li>a {
    height: 50px;
}
</style>