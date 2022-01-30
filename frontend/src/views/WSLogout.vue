<template>
  <card-module title="Goodbye">
    You have been logged out. To sign back in, please click
    <router-link to="/login">here</router-link>.
  </card-module>
</template>

<script>
import { mapActions } from "vuex";
import CardModule from "@/components/bricks/CardModule.vue";

export default {
  name: "WSLogout",
  components: {
    CardModule,
  },
  created() {
    this.logout();
  },
  methods: {
    ...mapActions("login", ["logout"]),
  },
};
</script>
