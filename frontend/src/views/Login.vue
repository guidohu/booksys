<template>
  <!-- <div class="login-view"> -->
  <b-overlay
    id="overlay-background"
    :show="isLoading"
    spinner-type="border"
    spinner-variant="info"
    rounded="sm"
  >
    <LoginModal 
      v-if="showLogin"
      :statusMessage="loginStatus" 
      :initialUsername="username" 
      @login="handleLogin"
    />
  </b-overlay>
  <!-- </div> -->
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import LoginModal from '@/components/LoginModal.vue';
import { BooksysBrowser } from '@/libs/browser';
import {
  BOverlay
} from 'bootstrap-vue';

export default {
  name: 'Login',
  components: {
    LoginModal,
    BOverlay
  },
  data() {
    return {
      isLoading: true,
      showLogin: false,
    }
  },
  computed: {
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    },
    ...mapGetters('login', [
      'loginStatus',
      'username'
    ])
  },
  methods: {
    ...mapActions('login', [
      'getIsLoggedIn',
      'login'
    ]),
    handleLogin: function (username, password) {
      console.log("handleLogin", username, password)
      this.login({ 
        username: username, 
        password: password
      })
    }
  },
  created () {
    this.isLoading = true;

    this.getIsLoggedIn()
    .then((loggedIn) => {
      if(loggedIn == true){
        console.log('User is logged in: redirect to content');
        if(this.$route.query != null && this.$route.query.target != null){
          this.$router.push(this.$route.query.target);
        }else{
          this.$router.push('/dashboard');
        }
      }else{
        this.showLogin = true;
      }
      this.isLoading = false;
    })
    .catch(errors => {
      this.isLoading = false;
      this.showLogin = true;
      console.error("Errors while getIsLoggedIn was called:", errors);
    })
  }
}
</script>
