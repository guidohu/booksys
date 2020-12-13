<template>
  <b-modal 
    id="view_database_update_modal"
    title='Database Update'
  >
    <div v-if="getDbVersionInfo == null" class="text-center">
      <b-spinner variant="primary" label="Spinning"></b-spinner>
    </div>
    <div v-else>
      <form>
        <b-row v-if="errors.length">
          <b-col offset="1" cols="10">
            <b-alert variant="warning" show>
              <b>Please correct the following error(s):</b>
              <ul>
                <li v-for="error in errors" :key="error">{{ error }}</li>
              </ul>
            </b-alert>
          </b-col>
        </b-row>
        <b-row v-if="getDbVersionInfo != null">
          <b-col offset="1" cols="4" class="font-weight-bold">
            App Version
          </b-col>
          <b-col cols="6">
            v{{ getDbVersionInfo.appVersion }}
          </b-col>
        </b-row>
        <b-row v-if="getDbVersionInfo != null">
          <b-col offset="1" cols="4" class="font-weight-bold">
            DB Version
          </b-col>
          <b-col cols="6">
            v{{ getDbVersionInfo.dbVersion}}
            {{dbVersionText}}
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col offset="1" cols="10">
            <b-alert variant="info" v-if="getDbVersionInfo != null && getDbVersionInfo.isUpdated == false && getDbUpdateResult == null" show>
              The database version is outdated and needs to be upgraded. If the upgrade is not done, there might occur data inconsisstencies and the app will stop working properly. Thus it is highly recommended to upgrade now.
            </b-alert>
            <b-alert variant="success" v-if="getDbVersionInfo.isUpdated == false && getDbIsUpdating == true" show>
              Updating database. Please wait, this might take a few minutes.
            </b-alert>
            <b-alert variant="success" v-if="getDbVersionInfo.isUpdated == true && getDbIsUpdating == false && (getDbUpdateResult == null || getDbUpdateResult.success == true)" show>
              Update successful.
            </b-alert>
            <b-alert variant="danger" v-if="getDbIsUpdating == false && getDbUpdateResult != null && getDbUpdateResult.success == false" show>
              Update failed. Please let the administrator know about this, together with the error message below:<br><br>
              <strong>Error:</strong> {{ getDbUpdateResult.msg }}
            </b-alert>
          </b-col>
        </b-row>
        <b-row class="mt-3" v-if="getDbUpdateResult != null && getDbUpdateResult.queries != null && getDbUpdateResult.queries.length > 0 && getDbUpdateResult.ok == false">
          <b-col offset="1" cols="10">
            <b-row v-bind:key="query.query" v-for="query in getDbUpdateResult.queries">
              <div class="col-sm-12 col-xs-12">
                <span class="glyphicon glyphicon-ok" v-if="query.ok == true"></span>
                <span class="glyphicon glyphicon-remove" v-if="query.ok == false"></span>
                {{ query.query }}
              </div>
            </b-row>
          </b-col>
        </b-row>        
      </form>
    </div>
    <div v-if="getDbVersionInfo != null" slot="modal-footer">
      <b-button type="button" variant="outline-info" v-on:click="updateDb" v-if="getDbVersionInfo.isUpdated == false && getDbIsUpdating == false">
        <b-icon-arrow-repeat></b-icon-arrow-repeat>
        Update
      </b-button>
      <b-button type="button" variant="outline-secondary" disabled v-if="getDbVersionInfo.isUpdated == false && getDbIsUpdating == true">
        <b-icon-hourglass-split></b-icon-hourglass-split>
        Updating
      </b-button>
      <b-button type="button" variant="outline-success" v-on:click="cancel" v-if="getDbVersionInfo.isUpdated == true && getDbIsUpdating == false">
        <b-icon-check></b-icon-check>
        Done
      </b-button>
      <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="cancel" v-if="getDbVersionInfo.isUpdated == false && getDbIsUpdating == false">
        <b-icon-x></b-icon-x>
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue'
import { mapActions, mapGetters } from 'vuex'

export default Vue.extend({
  name: 'DatabaseUpdateModal',
  data() {
    return {
      "errors": [ ],
    }
  },
  methods: {
    ...mapActions('configuration',[
      'queryDbVersionInfo',
      'updateDb'
    ]),
    cancel: function(){
      this.$bvModal.hide('view_database_update_modal')
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getDbVersionInfo',
      'getDbIsUpdating',
      'getDbUpdateResult'
    ]),
    dbVersionText: function() {
      if(this.dbVersionInfo == null){
        return "loading";
      }
      else if(this.dbVersionInfo.dbVersion == this.dbVersionInfo.dbVersionRequired){
        return "v" + this.dbVersionInfo.dbVersion;
      }
      else{
        return "v" + this.dbVersionInfo.dbVersion + " (update available to: v" + this.dbVersionInfo.dbVersionRequired + ")";
      }
    }
  },
  created() {
    this.queryDbVersionInfo()
  },
  mounted() {
    this.$bvModal.show('view_database_update_modal')
  }
})
</script>