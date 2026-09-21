<template>
  <subpage-container title="Logs">
    <log-card class="mx-1 card-height" />
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/admin">
        <i class="bi bi-gear"></i>
        ADMIN
      </router-link>
    </template>
  </subpage-container>
</template>

<script setup>
import LogCard from "booksys/components/LogCard.vue";
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
    /* The log table scrolls internally, the card itself never does */
    overflow: hidden;
  }
}
</style>
