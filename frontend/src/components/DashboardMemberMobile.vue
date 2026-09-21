<template>
  <div>
    <navbar-mobile title="Dashboard" role="member-home" />
    <div class="bk-page bk-bento">
      <dashboard-tile
        to="/schedule"
        label="My calendar"
        hint="Your upcoming sessions"
        icon="bi-calendar-week"
        wide
      />
      <dashboard-tile to="/account" label="Account" icon="bi-person" />
      <dashboard-tile to="/info" label="Info" icon="bi-geo-alt" />
    </div>
  </div>
</template>

<script setup>
import NavbarMobile from "booksys/components/NavbarMobile.vue";
import DashboardTile from "booksys/components/bricks/DashboardTile.vue";
</script>
