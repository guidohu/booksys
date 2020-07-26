<template>
  <div class="login-view">
      <LoginModal :isMobile="isMobile" :statusMessage="loginStatusMessage" @login="handleLogin"/>
  </div>
</template>

<script>
import LoginModal from '@/components/LoginModal.vue'
import { BooksysBrowser } from '@/libs/browser'
import { Login as ApiLogin} from '@/api/login'

export default {
  name: 'Login',
  data: function() {
    return {
      // initialize login message to be empty
      loginStatusMessage: null
    }
  },
  computed: {
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  },
  methods: {
    handleLogin: function (username, password) {
      console.log(username, password)
      this.loginStatusMessage = "Cannot login"
      let successCb = () => {
        console.log("Login Successful")
      }
      let failCb = () => {
        console.log("Login Failed")
      }

      ApiLogin.login(username, password, successCb, failCb)
    }
  },
  components: {
    LoginModal
  },
  beforeCreate () {}
}
</script>
