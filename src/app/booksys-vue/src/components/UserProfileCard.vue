<template>
  <b-card no-body class="text-left">
    <b-card-header>
      Profile
    </b-card-header>
    <b-card-body>
      <b-row>
        <b-col cols="6">
          {{userInfo.first_name}} {{userInfo.last_name}}
        </b-col>
        <b-col cols="6">
          {{userInfo.email}}
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="6">
          {{userInfo.address}}
        </b-col>
        <b-col cols="6">
          {{userInfo.mobile}}
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="6">
          {{userInfo.plz}} {{userInfo.city}}
        </b-col>
        <b-col cols="6">
          Driving License {{userInfo.license == true ? 'Yes' : 'No' }}
        </b-col>
      </b-row>
    </b-card-body>
  </b-card>
</template>

<script>
import Vue from 'vue'
import { mapGetters } from 'vuex'

export default Vue.extend({
  name: 'UserProfileCard',
  computed: {
    ...mapGetters('login', [
      'userInfo'
    ]),
  }  
})
</script>