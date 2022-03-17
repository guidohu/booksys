<template>
  <subpage-container title="Logs">
    <log-card class="mx-1 scroll" />
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/admin">
        <i class="bi bi-gear"></i>
        ADMIN
      </router-link>
    </template>
  </subpage-container>
</template>

<script>
import LogCard from "@/components/LogCard";
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";

export default {
  name: "WSLogs",
  components: {
    LogCard,
    SubpageContainer,
  },
};
</script>

<style scoped>
/* responsive, for small screens, scrolling */
@media (max-width: 992px) {
  .scroll {
    overflow-x: scroll;
  }
}

.scroll {
  max-height: 480px;
  overflow: scroll;
}
</style>
