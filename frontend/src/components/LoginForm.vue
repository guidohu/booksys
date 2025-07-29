<template>
  <form role="form" @submit="login">
    <div class="card">
      <div class="card-header">Please login</div>
      <div class="card-body">
        <div class="form-floating mb-3">
          <input
            v-model="form.username"
            type="text"
            class="form-control"
            id="username"
            ref="username"
            placeholder="name@example.com"
            autocomplete="username"
          />
          <label for="floatingInput">Email address</label>
        </div>
        <div class="form-floating mb-3">
          <input
            v-model="form.password"
            type="password"
            class="form-control"
            id="password"
            placeholder="Password"
            autocomplete="current-password"
          />
          <label for="floatingPassword">Password</label>
        </div>
        <div class="row text-center mb-3" v-if="statusMessage">
          <div class="col-12">
            <div class="alert alert-danger" role="alert">
              {{ statusMessage }}
            </div>
          </div>
        </div>
        <div class="row justify-content-end my-2">
          <div class="col-12 text-end">
            <button class="btn btn-outline-dark btn-lg" type="submit">
              <i class="bi bi-chevron-compact-right"></i>
              Login
            </button>
          </div>
        </div>
      </div>
      <div class="card-footer">
        <div class="row">
          <div class="col-12">
            <router-link tag="button" to="/password/reset">
              Forgot password
            </router-link>
            <router-link class="ms-2" to="/signup"> Sign-Up </router-link>
          </div>
        </div>
      </div>
    </div>
  </form>
</template>

<script setup>
import { ref, onMounted } from "vue";

const props = defineProps(["statusMessage", "initialUsername"]);
const emit = defineEmits(["login"]);

const form = ref({
  username: null,
  password: null,
});

const username = ref(null);

onMounted(() => {
  focusUsername();
  form.value.username = props.initialUsername;
});

function login(e) {
  emit("login", form.value.username, form.value.password);
  e.preventDefault();
}

function focusUsername() {
  username.value.focus();
}
</script>

<style scoped>
a {
  color: rgb(39, 39, 39);
}
</style>
