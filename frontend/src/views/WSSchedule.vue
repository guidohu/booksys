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

<script>
import { BooksysBrowser } from "@/libs/browser";
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";
import CardModule from "@/components/bricks/CardModule.vue";
import UserSessionTabs from "@/components/UserSessionTabs.vue";

export default {
  name: "WSSchedule",
  components: {
    SubpageContainer,
    UserSessionTabs,
    CardModule,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
  },
};
</script>
