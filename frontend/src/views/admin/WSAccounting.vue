<template>
  <subpage-container title="Accounting">
    <card-module :nobody="true" class="mx-1 scroll">
      <accounting-tabs/>
    </card-module>
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/admin">
        <i class="bi bi-gear"></i>
        ADMIN
      </router-link>
    </template>
  </subpage-container>
</template>

<script>
import AccountingTabs from "@/components/AccountingTabs";
import CardModule from "@/components/bricks/CardModule.vue";
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";

export default {
  name: "WSAccounting",
  components: {
    AccountingTabs,
    CardModule,
    SubpageContainer
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
  overflow: hidden;
}
</style>
