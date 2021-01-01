<template>
  <b-card v-if="userInfo != null" title="Goodbye">
    <b-card-text>
      Logging you out.
    </b-card-text>
  </b-card>
  <b-card v-else title="Goodbye">
    <b-card-text>
      You have been logged out. To sign back in, please click <b-link to="/login">here</b-link>.
    </b-card-text>
  </b-card>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import {
  BCard,
  BCardText,
  BLink
} from 'bootstrap-vue';

export default Vue.extend({
  name: 'Logout',
  components: {
    BCard,
    BCardText,
    BLink
  },
  created () {
    this.logout()
  },
  computed: {
    ...mapGetters('login', [
      'userInfo'
    ])
  },
  methods: {
    ...mapActions('login', [
      'logout'
    ])
  }
})
</script>