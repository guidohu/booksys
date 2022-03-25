<template>
  <subpage-container title="User Management">
    <card-module :nobody="true" class="mx-1 scroll">
      <user-management-tabs />
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
import UserManagementTabs from "@/components/UserManagementTabs";
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";
import CardModule from "@/components/bricks/CardModule.vue";

export default {
  name: "WSUsers",
  components: {
    UserManagementTabs,
    SubpageContainer,
    CardModule,
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
