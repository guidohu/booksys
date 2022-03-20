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

<script>
import { mapActions, mapGetters } from "vuex";
import UserProfileCard from "@/components/UserProfileCard";
import UserStatisticsCard from "@/components/UserStatisticsCard";
import UserBalanceCard from "@/components/UserBalanceCard";
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";

export default {
  name: "WSAccount",
  components: {
    UserProfileCard,
    UserStatisticsCard,
    UserBalanceCard,
    SubpageContainer,
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
  },
  methods: {
    ...mapActions("user", ["queryHeatHistory"]),
  },
  created() {
    this.queryHeatHistory();
  },
};
</script>
