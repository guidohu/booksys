<template>
  <subpage-container title="Boat">
    <card-module :nobody="true" fill class="mx-1 card-height">
      <boat-tabs />
    </card-module>
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/dashboard">
        <i class="bi bi-house"></i>
        HOME
      </router-link>
    </template>
  </subpage-container>
</template>

<script setup>
import BoatTabs from "booksys/components/BoatTabs.vue";
import CardModule from "booksys/components/bricks/CardModule.vue";
import SubpageContainer from "booksys/components/bricks/SubpageContainer.vue";
</script>

<style scoped>
@media (max-width: 992px) {
  /* Fill the whole phone screen below the fixed navbar (60px padding-top) */
  .card-height {
    min-height: calc(100vh - 66px);
    max-height: calc(100vh - 66px);
    height: calc(100vh - 66px);
    min-height: calc(100dvh - 66px - env(safe-area-inset-bottom, 0px));
    max-height: calc(100dvh - 66px - env(safe-area-inset-bottom, 0px));
    height: calc(100dvh - 66px - env(safe-area-inset-bottom, 0px));
    /* The tab content scrolls internally, the card itself never does */
    overflow: hidden;
  }
}
</style>
