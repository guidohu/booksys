<template>
  <div class="d-flex flex-column">    
    <div class="container">
      <div class="row text-center mb-3">
        <div class="col-12">
          <div class="ratio ratio-21x9">
            <img
              class="img-fluid"
              v-if="getLogoFile != null"
              :src="getLogoFile"
              alt="Logo"
              :height="100"
              width="auto"
            />
          </div>
        </div>
      </div>
      <div class="row">
        <div class="col-3 d-none d-sm-block"/>
        <div class="col-12 col-sm-6">
          <overlay-spinner v-model:active="isLoading" :fullPage="false">
            <login-form
              v-if="showLogin"
              :statusMessage="status"
              :initialUsername="username"
              @login="handleLogin"
            />
          </overlay-spinner>
        <div class="col-3 d-none d-sm-block"/>
      </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import LoginForm from "@/components/LoginForm.vue";
import { BooksysBrowser } from "@/libs/browser";
import OverlaySpinner from "@/components/styling/OverlaySpinner.vue";

export default {
  name: "WSLogin",
  components: {
    LoginForm,
    OverlaySpinner
  },
  data() {
    return {
      isLoading: true,
      showLogin: false,
      status: null,
    };
  },
  computed: {
    isMobile: function () {
      return BooksysBrowser.isMobile();
    },
    ...mapGetters("login", ["username"]),
    ...mapGetters("configuration", ["getLogoFile"]),
  },
  methods: {
    ...mapActions("login", ["getIsLoggedIn", "login"]),
    ...mapActions("configuration", ["queryLogoFile"]),
    handleLogin: function (username, password) {
      this.isLoading = true;
      this.login({
        username: username,
        password: password,
      })
        .then(() => {
          if (this.$route.query != null && this.$route.query.target != null) {
            this.$router.push(this.$route.query.target);
          } else {
            this.$router.push({name:"Dashboard"});
          }
          this.isLoading = false;
        })
        .catch((errors) => {
          if (errors.length > 0) {
            this.status = errors[0];
          }
          this.isLoading = false;
        });
    },
  },
  created() {
    this.isLoading = true;

    this.queryLogoFile().catch((errors) => console.log(errors));

    this.getIsLoggedIn()
      .then((loggedIn) => {
        if (loggedIn == true) {
          console.log("User is logged in: redirect to content");
          console.log(this.$route);
          if (this.$route.query != null && this.$route.query.target != null) {
            console.log("Redirect to:", this.$route.query.target);
            this.$router.push(this.$route.query.target);
          } else {
            console.log("Redirect to: /dashboard");
            this.$router.push("/dashboard");
          }
        } else {
          this.showLogin = true;
        }
        this.isLoading = false;
      })
      .catch((errors) => {
        this.isLoading = false;
        this.showLogin = true;
        console.error("Errors while getIsLoggedIn was called:", errors);
      });
  },
};
</script>
