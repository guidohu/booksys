<template>
  <div class="login-view">
      <LoginModal :isMobile="isMobile" :statusMessage="loginStatusMessage" @login="handleLogin"/>
  </div>
</template>

<script>
import Vue from 'vue'
import { mapActions } from 'vuex'
import LoginModal from '@/components/LoginModal.vue'
import { BooksysBrowser } from '@/libs/browser'
import { Login as ApiLogin} from '@/api/login'

export default Vue.extend({
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
    ...mapActions('login', [
      'setUserInfo'
    ]),
    handleLogin: function (username, password) {
      console.log(username, password)
      this.loginStatusMessage = "Cannot login"
      let successCb = () => {
        console.log("Login Successful")
        this.loginStatusMessage = null
        // send login state to store
        let loadDashboard = (userInfo) => {
          // store in store
          console.log("userinfo", userInfo)

          this.setUserInfo(userInfo)

          // route to /dashboard
          this.$router.push('/dashboard')
        }
        let errorCase = (error) => {
          console.log("Error", error)
        }
        ApiLogin.getMyUser(loadDashboard, errorCase) 
      }
      let failCb = (message) => {
        if(message.login_status_code == -1){
          this.loginStatusMessage = "incorrect username or password"
        } else if(message.login_status_code == -2){
          this.loginStatusMessage = "please be patient until we activate your account"
        } else {
          this.loginStatusMessage = message
        }
        console.log("Login Failed", message)
      }

      ApiLogin.login(username, password, successCb, failCb)
    }
  },
  components: {
    LoginModal
  },
  beforeCreate () {}
})
</script>
