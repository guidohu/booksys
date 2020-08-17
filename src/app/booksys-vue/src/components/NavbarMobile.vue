<template>
  <b-navbar fixed="top" toggleable="lg" type="light" variant="light">
    <b-navbar-brand href="#">{{title}}</b-navbar-brand>

    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item href="#">
          <router-link to="/dashboard">
            <b-icon icon="house"></b-icon>
            Dashboard
          </router-link>
        </b-nav-item>
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
</template>

<script>
import Vue from 'vue'

export default Vue.extend({
  name: 'NavbarMobile',
  props: ['title']
})
</script>