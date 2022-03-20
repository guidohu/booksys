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

<script>
import BoatTabs from "@/components/BoatTabs";
import CardModule from "@/components/bricks/CardModule.vue";
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";

export default {
  name: "WSBoat",
  components: {
    BoatTabs,
    CardModule,
    SubpageContainer,
  },
};
</script>
