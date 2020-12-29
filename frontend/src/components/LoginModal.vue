<template>
  <b-form v-on:submit.prevent="login" role="form">
    <b-row class="text-left">
      <b-col cols="3" class="d-none d-sm-block">
      </b-col>
      <b-col cols="12" sm="6">
        <b-card no-body>
          <b-card-header>
            Please Login
          </b-card-header>
          <b-card-body>
            <b-card-text>
              <b-form-group>
                <b-col cols="12">
                    <label class="d-none d-sm-block font-weight-bold" for="username">e-mail</label>
                    <input ref="username" type="text" v-model="username" class="form-control form-control-white" placeholder="e-mail" autocomplete="username"/>
                </b-col>
              </b-form-group>
              <b-form-group>
                <b-col cols="12">
                    <label class="d-none d-sm-block font-weight-bold" for="password">password</label>
                    <input type="password" v-model="password" class="form-control form-control-white" placeholder="password" autocomplete="current-password"/>
                </b-col>
              </b-form-group>
              <b-form-group v-if="statusMessage" class="text-center">
                <b-col cols="12">
                  <b-alert show variant="danger">
                    {{ statusMessage }}
                  </b-alert>
                </b-col>
              </b-form-group>
              <b-form-group class="text-right">
                <b-col cols="12">
                  <b-button variant="outline-dark" type="submit">
                    <b-icon-caret-right-fill></b-icon-caret-right-fill>
                    Login
                  </b-button>
                </b-col>
              </b-form-group>
            </b-card-text>
          </b-card-body>
          <b-card-footer>
            <b-form-group class="text-left">
                <b-col cols="3" class="d-xs-none"></b-col>
                <b-col cols="12" xs="8">
                    <b-button to="/password/reset" variant="light">
                        Forgot password
                    </b-button>
                    <b-button to="/signup" variant="light">
                        Sign-Up
                    </b-button>
                </b-col>
                <b-col cols="1" class="d-xs-none"></b-col>
            </b-form-group>
          </b-card-footer>
        </b-card>
      </b-col>
      <b-col cols="3" class="d-none d-sm-block"></b-col>
    </b-row>
  </b-form>
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
