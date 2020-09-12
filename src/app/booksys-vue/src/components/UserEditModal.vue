<template>
  <b-modal
    id="userEditModal"
    title="Profile"
  >
    <b-row class="text-left">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <!-- here it comes -->
      </b-col>
      <b-col cols="1" class="d-none d-sm-block">
      </b-col>
    </b-row>    
    <div slot="modal-footer">
      <b-button type="button" variant="outline-info" v-on:click="save">
        <b-icon-person-check></b-icon-person-check>
        Save
      </b-button>
      <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x></b-icon-x>
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue'
import { mapGetters } from 'vuex'

export default Vue.extend({
  name: 'UserEditModal',
  computed: {
    ...mapGetters('configuration', [
      'getConfiguration'
    ])
  },
  methods: {
    close: function(){
      this.$bvModal.hide('userEditModal');
    },
    save: function(){
      
    }
  }
})
</script>