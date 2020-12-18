<template>
  <b-modal
    id="setupModal"
    title="Setup"
    no-close-on-backdrop
    no-close-on-esc
    hide-header-close
    visible
  >
    <b-row v-if="errors.length > 0" class="mt-4">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <warning-box :errors="errors"/>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <database-configuration v-if="showDbSetup" :dbConfig="dbConfig" @save="setDbSettings"/>
    <div v-if="showUserSetup">Admin User</div>
    <div slot="modal-footer">
      <b-button v-if="showDbSetup" class="mr-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x></b-icon-x>
        Cancel
      </b-button>
      <b-button v-if="showDbSetup" type="button" variant="outline-info" v-on:click="setDbSettings" :disabled="isLoading">
        <b-icon icon="arrow-right"/>
        Next
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import Backend from '@/api/backend';
import { Configuration } from '@/api/configuration';
import DatabaseConfiguration from '@/components/DatabaseConfiguration';
import WarningBox from '@/components/WarningBox';

export default Vue.extend({
  name: "Setup",
  components: {
    DatabaseConfiguration,
    WarningBox
  },
  data(){
    return {
      errors: [],
      showDbSetup: false,
      showUserSetup: false,
      isLoading: false,
      dbConfig: {}
    }
  },
  methods: {
    setDbSettings: function(){
      Configuration.setDbConfig(this.dbConfig)
      .then(() => {
        this.showDbSetup = false;
        this.showUserSetup = true;
        this.errors = [];
      })
      .catch((errors) => this.errors = errors);
    },
    close: function() {
      this.$bvModal.hide("setupModal");
      this.$router.push("/login");
    },
    getDbConfig: function() {

    }
  },
  mounted(){
    Backend.getStatus()
    .then(status => {
      if(status.configFile == false){
        this.showDbSetup = true;
        this.showUserSetup = false;
      }
    })
    .catch(errors => this.errors = errors);
  }

})
</script>