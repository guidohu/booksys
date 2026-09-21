<template>
  <div>
    <navbar-mobile title="Dashboard" role="guest-home" />
    <div class="bk-page bk-bento">
      <dashboard-tile
        to="/account"
        label="Account"
        hint="Your profile and balance"
        icon="bi-person"
        wide
      />
      <dashboard-tile
        to="/info"
        label="Info"
        hint="Location and map"
        icon="bi-geo-alt"
        wide
      />
    </div>
  </div>
</template>

<script setup>
import NavbarMobile from "booksys/components/NavbarMobile.vue";
import DashboardTile from "booksys/components/bricks/DashboardTile.vue";
</script>
