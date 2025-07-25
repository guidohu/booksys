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

<script setup>
import { useStore } from "vuex";
import { onMounted } from "vue";
import CardModule from "booksys/components/bricks/CardModule.vue";

const store = useStore();

onMounted(() => {
  store.dispatch("login/logout").catch(() => console.error("Logout failed."));
});
</script>
