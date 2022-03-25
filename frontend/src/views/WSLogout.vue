<template>
  <div class="position-absolute top-50 start-50 translate-middle">
    <div class="row">
      <div class="col-12 text-center">
        <card-module title="Goodbye">
          You have been logged out. To sign back in, please click
          <router-link to="/login">here</router-link>.
        </card-module>
      </div>
    </div>
  </div>
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
