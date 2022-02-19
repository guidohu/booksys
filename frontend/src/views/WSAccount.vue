<template>
  <div v-if="isDesktop">
    <subpage-desktop title="Account">
      <template v-slot:content>
        <user-profile-card class="m-1"/>
        <user-statistics-card class="m-1"/>
        <user-balance-card class="m-1"/>
      </template>
      <template v-slot:bottom>
        <router-link tag="button" class="btn btn-outline-light" to="/dashboard">
          <i class="bi bi-house"></i>
          HOME
        </router-link>
      </template>
    </subpage-desktop>
  </div>
  <div v-else>
    <navbar-mobile title="Account" />
    <user-profile-card />
    <user-statistics-card class="mt-2" />
    <user-balance-card class="mt-2" />
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import UserProfileCard from "@/components/UserProfileCard";
import UserStatisticsCard from "@/components/UserStatisticsCard";
import UserBalanceCard from "@/components/UserBalanceCard";
import SubpageDesktop from '@/components/bricks/SubpageDesktop.vue';


export default {
  name: "WSAccount",
  components: {
    NavbarMobile,
    UserProfileCard,
    UserStatisticsCard,
    UserBalanceCard,
    SubpageDesktop,
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    isMobile: function () {
      return BooksysBrowser.isMobile();
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
    mapHeight: function () {
      if (this.isMobile) {
        return 350;
      } else {
        return 300;
      }
    },
    mapWidth: function () {
      if (this.isMobile) {
        return 340;
      } else {
        return 600;
      }
    },
  },
  methods: {
    ...mapActions("user", ["queryHeatHistory"]),
  },
  created() {
    console.log("Info.vue: Try to get all information required");
    this.queryHeatHistory();
  },
};
</script>
