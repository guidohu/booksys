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
              <b-form-group
                id="input-group-username"
                label="Email"
                label-for="input-username"
                label-cols="3"
              >
                <b-form-input
                  id="input-username"
                  v-model="form.username"
                  type="text"
                  ref="username"
                  placeholder="e-mail"
                  autocomplete="username"
                />
              </b-form-group>
              <b-form-group
                id="input-group-password"
                label="Password"
                label-for="input-password"
                label-cols="3"
              >
                <b-form-input
                  id="input-password"
                  v-model="form.password"
                  type="password"
                  placeholder="password"
                  autocomplete="current-password"
                />
              </b-form-group>
              <b-row v-if="statusMessage" class="text-center">
                <b-col cols="12">
                  <b-alert show variant="danger">
                    {{ statusMessage }}
                  </b-alert>
                </b-col>
              </b-row>
              <b-row class="text-right">
                <b-col cols="12">
                  <b-button variant="outline-dark" type="submit">
                    <b-icon-caret-right-fill/>
                    Login
                  </b-button>
                </b-col>
              </b-row>
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
import {
  BForm,
  BFormGroup,
  BFormInput,
  BAlert,
  BButton,
  BIconCaretRightFill,
  BRow,
  BCol,
  BCard,
  BCardHeader,
  BCardBody,
  BCardFooter,
  BCardText,
} from 'bootstrap-vue';

export default {
  name: 'LoginModal',
  components: {
    BForm,
    BFormGroup,
    BFormInput,
    BAlert,
    BButton,
    BIconCaretRightFill,
    BRow,
    BCol,
    BCard,
    BCardHeader,
    BCardBody,
    BCardFooter,
    BCardText,
  },
  props: ['statusMessage', 'initialUsername'],
  data: function () {
    return {
      notifications: [],
      form: {
        username: null,
        password: null
      },
    }
  },
  methods: {
    login: function () {
      this.$emit("login", this.form.username, this.form.password)
    },
    focusUsername() {
      this.$refs.username.focus();
    }
  },
  mounted() {
    this.focusUsername();
    this.form.username = this.initialUsername;
  }
}
</script>
