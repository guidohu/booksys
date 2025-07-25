<template>
  <subpage-container title="Boat">
    <card-module :nobody="true" class="mx-1">
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
