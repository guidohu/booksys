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
  /* Fill what is left of the phone screen once the top bar, the tab bar and
     the page padding (.bk-page, 0.65rem top and bottom) have taken theirs. */
  .card-height {
    min-height: calc(
      100vh - var(--bk-shell-top) - var(--bk-shell-bottom) - 1.3rem
    );
    max-height: calc(
      100vh - var(--bk-shell-top) - var(--bk-shell-bottom) - 1.3rem
    );
    height: calc(100vh - var(--bk-shell-top) - var(--bk-shell-bottom) - 1.3rem);
    min-height: calc(
      100dvh - var(--bk-shell-top) - var(--bk-shell-bottom) - 1.3rem
    );
    max-height: calc(
      100dvh - var(--bk-shell-top) - var(--bk-shell-bottom) - 1.3rem
    );
    height: calc(
      100dvh - var(--bk-shell-top) - var(--bk-shell-bottom) - 1.3rem
    );
    /* The tab content scrolls internally, the card itself never does */
    overflow: hidden;
  }
}
</style>
