<template>
  <div v-if="isDesktop">
    <subpage-desktop title="Schedule">
      <template v-slot:content>
        <card-module :nobody="true" class="m-1">
          <user-session-tabs />
        </card-module>
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
    <navbar-mobile title="Schedule" />
    <card-module :nobody="true">
      <user-session-tabs />
    </card-module>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import SubpageDesktop from "@/components/bricks/SubpageDesktop.vue";
import CardModule from "@/components/bricks/CardModule.vue";
import UserSessionTabs from "@/components/UserSessionTabs.vue";

export default {
  name: "WSSchedule",
  components: {
    NavbarMobile,
    SubpageDesktop,
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
