<template>
  <form role="form" @submit="login">
    <div class="card">
      <div class="card-header">Welcome back</div>
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
          <div class="col-12 text-end login-submit">
            <button class="btn btn-outline-dark btn-lg" type="submit">
              <i class="bi bi-chevron-compact-right"></i>
              Login
            </button>
          </div>
        </div>
      </div>
      <div class="card-footer login-footer">
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

/* Mobile only: the login card is the first thing a user sees, so it gets the
   full treatment here while the desktop layout stays as it is. */
@media (max-width: 992px) {
  .card-header {
    padding: 1.25rem 1.25rem 0.25rem;
    border-bottom: 0;
    color: var(--bk-on-dark);
    font-size: 1.35rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    text-transform: none;
  }

  .card-body {
    padding: 1rem 1.25rem 1.25rem;
  }

  /* A single full-width target: there is only one thing to do on this
     screen, and it should be reachable with the thumb that is holding the
     phone rather than tucked into the far corner. */
  /* The one action on the screen, so it is the only filled aqua surface in
     the card rather than another outline. */
  .login-submit .btn {
    width: 100%;
    border-color: transparent;
    background-color: var(--bk-aqua);
    color: #04222b;
    font-size: 1rem;
  }

  .login-submit .btn:active {
    background-color: var(--bk-aqua-bright);
    border-color: transparent;
    color: #04222b;
  }

  /* The chevron was decoration next to a right-aligned button; centred in a
     full-width one it just crowds the label. */
  .login-submit .btn .bi {
    display: none;
  }

  .login-footer {
    padding: 0.9rem 1.25rem;
    background-color: rgba(0, 0, 0, 0.14);
    text-align: center;
  }

  .login-footer a {
    color: var(--bk-aqua-bright);
    font-size: 0.85rem;
    font-weight: 600;
    text-decoration: none;
  }
}
</style>
