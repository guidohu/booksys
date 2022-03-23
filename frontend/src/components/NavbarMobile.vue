<template>
  <nav class="navbar fixed-top navbar-light bg-light">
    <div class="container-fluid">
      <a class="navbar-brand" href="#">{{ title }}</a>
      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarToggler"
        aria-controls="navbarToggler"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>

      <div class="collapse navbar-collapse" id="navbarToggler">
        <ul class="navbar-nav">
          <li class="nav-item mt-2">
            <router-link to="/dashboard" tag="button" class="btn btn-light">
              <i class="bi bi-house"></i>
              Dashboard
            </router-link>
          </li>
          <li class="nav-item mt-2">
            <router-link to="/account" tag="button" class="btn btn-light">
              <i class="bi bi-person-circle"></i>
              Account
            </router-link>
          </li>
          <li class="nav-item mt-2">
            <router-link to="/info" tag="button" class="btn btn-light">
              <i class="bi bi-info-circle"></i>
              Info
            </router-link>
          </li>
          <li class="nav-item mt-2">
            <router-link to="/logout" tag="button" class="btn btn-light">
              <i class="bi bi-box-arrow-right"></i>
              Logout
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: "NavbarMobile",
  components: {},
  props: ["title", "role"],
};
</script>
