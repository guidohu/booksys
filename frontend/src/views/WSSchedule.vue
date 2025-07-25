<template>
  <subpage-container title="Schedule">
    <card-module :nobody="true" class="mx-1">
      <user-session-tabs />
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
import SubpageContainer from "booksys/components/bricks/SubpageContainer.vue";
import CardModule from "booksys/components/bricks/CardModule.vue";
import UserSessionTabs from "booksys/components/UserSessionTabs.vue";
</script>
