<template>
  <subpage-container title="Account">
    <user-profile-card class="m-1" />
    <user-statistics-card class="m-1" />
    <user-balance-card class="m-1" />
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/dashboard">
        <i class="bi bi-house"></i>
        HOME
      </router-link>
    </template>
  </subpage-container>
</template>

<script setup>
import { useStore } from "vuex";
import { onMounted } from "vue";
import UserProfileCard from "../components/UserProfileCard.vue";
import UserStatisticsCard from "../components/UserStatisticsCard.vue";
import UserBalanceCard from "../components/UserBalanceCard.vue";
import SubpageContainer from "../components/bricks/SubpageContainer.vue";

const store = useStore();

onMounted(() => {
  store.dispatch("user/queryHeatHistory");
});
</script>
