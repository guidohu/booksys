<template>
  <div v-if="isDesktop">
    <subpage-desktop title="Boat">
      <template v-slot:content>
        <card-module :nobody="true" class="m-1">
          <BoatTabs />
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
    <NavbarMobile title="Boat" />
    <card-module :nobody="true">
      <BoatTabs />
    </card-module>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import BoatTabs from "@/components/BoatTabs";
import SubpageDesktop from "@/components/bricks/SubpageDesktop.vue";
import CardModule from "@/components/bricks/CardModule.vue";

export default {
  name: "WSBoat",
  components: {
    NavbarMobile,
    BoatTabs,
    SubpageDesktop,
    CardModule,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
  },
};
</script>
