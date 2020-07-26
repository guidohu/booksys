<template>
  <form class="form-horizontal" v-on:submit.prevent="login" role="form">
    <div class="row">
      <div class="col-sm-3 hidden-xs">
      </div>
      <div class="col-sm-6 col-xs-12">
        <div class="panel panel-default">
          <div class="panel-heading">
            Please Login
          </div>
          <div class="panel-body text-left">
            <div class="form-group">
                <div class="col-sm-3 hidden-xs"></div>
                <div class="col-sm-8 col-xs-12">
                    <label class="hidden-xs" for="username">e-mail</label>
                    <input ref="username" type="text" v-model="username" class="form-control form-control-white" placeholder="e-mail" autocomplete="username"/>
                </div>
                <div class="col-sm-1 hidden-xs"></div>
            </div>
            <div class="form-group">
                <div class="col-sm-3 hidden-xs"></div>
                <div class="col-sm-8 col-xs-12">
                    <label class="hidden-xs" for="password">password</label>
                    <input type="password" v-model="password" class="form-control form-control-white" placeholder="password" autocomplete="current-password"/>
                </div>
                <div class="col-sm-1 hidden-xs"></div>
            </div>
            <div v-if="statusMessage" class="form-group text-center">
                <div class="col-sm-3 hidden-xs"></div>
                <div class="col-sm-8 col-xs-12">
                  <div class="label label-danger">
                    {{ statusMessage }}
                  </div>
                </div>
                <div class="col-sm-1 hidden-xs"></div>
            </div>
            <div class="form-group text-right">
                <div class="col-sm-3 hidden-xs"></div>
                <div class="col-sm-8 col-xs-12">
                    <button type="submit" id="loginBtn" class="btn btn-default inline">
                        <span class="glyphicon glyphicon-log-in"></span>
                        Login
                    </button>
                </div>
                <div class="col-sm-1 hidden-xs"></div>
            </div>
          </div>
          <div class="panel-footer">
            <div class="form-group text-left">
                <div class="col-sm-3 hidden-xs"></div>
                <div class="col-sm-8 col-xs-12">
                    <a href="forgot_password.html" type="button" class="btn btn-default btn-xs inline">
                        Forgot password
                    </a>
                    <a href="sign_up.html" type="button" class="btn btn-default btn-xs inline">
                        Sign-Up
                    </a>
                </div>
                <div class="col-sm-1 hidden-xs"></div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-sm-3 hidden-xs"></div>
    </div>
  </form>
</template>

<script>
import Vue from 'vue'
import { mapActions } from 'vuex'

export default Vue.extend({
  name: 'LoginModal',
  components: {
  },
  props: ['isMobile', 'statusMessage'],
  data: function () {
    return {
      hasNotifications: true,
      notifications: [],
      password: ''
    }
  },
  computed: {
    // ...mapGetters('login', [
    //   'username'
    // ]),
    // ...mapState('login', {
    //   username: state => state.login.username
    // }),
    username: {
      get () {
        return this.$store.state.login.username
      },
      set (value) {
        // this.$store.commit('setUsername', value)
        this.setUsername(value)
      }
    }
  },
  mounted() {
    this.focusUsername()
  },
  methods: {
    login: function () {
      console.log("emit login")
      this.$emit("login", this.username, this.password)
    },
    ...mapActions('login', [
      'doLogin',
      'setUsername'
    ]),
    goSignUp: function() {
      this.$router.push('SignUp')
    },
    focusUsername() {
      this.$refs.username.focus();
    }
  }
})
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
