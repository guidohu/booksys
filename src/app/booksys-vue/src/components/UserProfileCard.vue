<template>
  <b-card no-body class="text-left">
    <b-card-header>
      <b-row>
        <b-col cols="6">
          Profile
        </b-col>
        <b-col cols="6" class="text-right">
          <b-dropdown variant="outline-info" size="sm" no-caret dropleft>
            <template v-slot:button-content>
              <b-icon-pencil-square></b-icon-pencil-square><span class="sr-only">Search</span>
            </template>
            <b-dropdown-item href="#" v-on:click="showUserEdit">
              Edit Profile
            </b-dropdown-item>
            <b-dropdown-item href="#" v-on:click="showPasswordEdit">
              Change Password
              </b-dropdown-item>
          </b-dropdown>
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body>
      <UserEditModal/>
      <UserPasswordEditModal/>
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
import UserEditModal from './UserEditModal'
import UserPasswordEditModal from './UserPasswordEditModal'

export default Vue.extend({
  name: 'UserProfileCard',
  components: {
    UserEditModal,
    UserPasswordEditModal
  },
  computed: {
    ...mapGetters('login', [
      'userInfo'
    ]),
  },
  methods: {
    showUserEdit: function() {
      this.$bvModal.show('userEditModal')
    },
    showPasswordEdit: function() {
      this.$bvModal.show('userPasswordEditModal')
    }
  }
})
</script>