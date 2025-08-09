<template>
  <div class="d-flex flex-column justify-content-center">
    <div class="container">
      <div class="row text-center mb-3 mt-5">
        <div class="col-12">
          <img
            class="img-fluid custom-height"
            v-if="getLogoUri != null && getLogoUri != ''"
            :src="getLogoUri"
            alt="Logo"
          />
        </div>
      </div>
      <div class="row">
        <div class="col-4 d-none d-lg-block" />
        <div class="col-12 col-lg-4">
          <overlay-spinner v-model:active="isLoading" :fullPage="false">
            <login-form
              v-if="showLogin"
              :statusMessage="status"
              :initialUsername="username"
              @login="handleLogin"
            />
          </overlay-spinner>
          <div class="col-4 d-none d-lg-block" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useStore, mapGetters } from "vuex";
import { useRouter, useRoute } from "vue-router";
import LoginForm from "booksys/components/LoginForm.vue";
import OverlaySpinner from "booksys/components/styling/OverlaySpinner.vue";

const store = useStore();
const router = useRouter();
const route = useRoute();

const isLoading = ref(true);
const showLogin = ref(false);
const status = ref(null);

const username = computed(
  mapGetters("login", ["username"]).username.bind({ $store: store }),
);
const getLogoUri = computed(
  mapGetters("configuration", ["getLogoUri"]).getLogoUri.bind({
    $store: store,
  }),
);

const handleLogin = (username, password) => {
  isLoading.value = true;
  store
    .dispatch("login/login", {
      username: username,
      password: password,
    })
    .then(() => {
      if (route.query != null && route.query.target != null) {
        router.push(route.query.target);
      } else {
        router.push({ name: "Dashboard" });
      }
      isLoading.value = false;
    })
    .catch((errors) => {
      console.error("login failed:", errors);
      if (errors.length > 0) {
        status.value = errors[0];
      }
      isLoading.value = false;
    });
};

onMounted(() => {
  isLoading.value = true;

  store
    .dispatch("configuration/queryLogoFile")
    .catch((errors) => console.log(errors));

  store
    .dispatch("login/getIsLoggedIn")
    .then((loggedIn) => {
      if (loggedIn == true) {
        console.log("User is logged in: redirect to content");
        console.log(route);
        if (route.query != null && route.query.target != null) {
          console.log("Redirect to:", route.query.target);
          router.push(route.query.target);
        } else {
          console.log("Redirect to: /dashboard");
          router.push("/dashboard");
        }
      } else {
        console.log("User not logged in: show login");
        showLogin.value = true;
      }
      isLoading.value = false;
    })
    .catch((errors) => {
      isLoading.value = false;
      showLogin.value = true;
      console.error("Errors while getIsLoggedIn was called:", errors);
    });
});
</script>

<style scoped>
.custom-height {
  max-height: 100px;
}
</style>
