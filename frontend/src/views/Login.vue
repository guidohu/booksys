<template>
  <div class="login-view">
    <LoginModal :isMobile="isMobile" :statusMessage="loginStatus" @login="handleLogin"/>
  </div>
</template>

<script>
import { mapActions, mapGetters, mapState } from 'vuex';
import LoginModal from '@/components/LoginModal.vue';
import { BooksysBrowser } from '@/libs/browser';

console.log("Loaded Login.vue");

export default {
  name: 'Login',
  computed: {
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    },
    ...mapGetters('login', [
      'userInfo',
      'loginStatus'
    ]),
    ...mapState('login', [
      'isLoggedIn'
    ])
  },
  watch: {
    isLoggedIn(newValue){
      // we are only properly logged in if we know both,
      // that we are logged in and when we know some basic user info
      if(newValue == true && this.userInfo != null){
        this.$router.push("/dashboard")
      }
    },
    userInfo(newValue){
      // we are only properly logged in if we know both,
      // that we are logged in and when we know some basic user info
      if(newValue != null && this.isLoggedIn == true){
        this.$router.push("/dashboard")
      }
    }
  },
  methods: {
    ...mapActions('login', [
      'getIsLoggedIn',
      'setUserInfo',
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
  components: {
    LoginModal
  },
  created () {
    console.log("Try to get all information if logged in already")
    this.getIsLoggedIn()
  }
}
</script>
