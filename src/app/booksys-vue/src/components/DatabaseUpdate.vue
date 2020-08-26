<template>
  <b-modal 
    id="view_database_update_modal"
    title='Database Update'>
    <div v-if="getDbVersionInfo==null" class="text-center">
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
            v{{ getDbVersionInfo.app_version }}
          </b-col>
        </b-row>
        <b-row v-if="getDbVersionInfo != null">
          <b-col offset="1" cols="4" class="font-weight-bold">
            DB Version
          </b-col>
          <b-col cols="6">
            v{{ getDbVersionInfo.db_version}}
          </b-col>
        </b-row>
            <!-- <div class="row row-top-padded">
                <div class="col-sm-1 col-xs-1"></div>
                <div class="col-sm-10 col-xs-10">
                    <div class="alert alert-info" v-if="info.isUpdated == false && info.isUpdating == false && info.ok == true">
                        The database version is outdated and needs to be upgraded. If the upgrade is not done, there might occur data inconsitencies and the app will stop working properly. Thus it is highly recommended to upgrade now.
                    </div>
                    <div class="alert alert-success" v-if="info.isUpdated == false && info.isUpdating == true">
                        Updating database. Please wait, this might take a few minutes.
                    </div>
                    <div class="alert alert-success" v-if="info.isUpdated == true && info.isUpdating == false && info.ok == true">
                        Update successful.
                    </div>
                    <div class="alert alert-danger" v-if="info.isUpdating == false && info.ok == false">
                        Update failed. Please let the administrator know about this, together with the error message below:<br>
                        <strong>Error:</strong> {{ info.message }}
                    </div>
                </div>
                <div class="col-sm-1 col-xs-1"></div>
            </div> -->
            <!-- <div class="row row-top-padded" v-if="info.queries != null && info.queries.length > 0 && info.ok == false">
                <div class="row">
                    <div class="col-sm-1 col-xs-1"></div>
                    <div class="col-sm-10 col-xs-10">
                        <div class="panel panel-default">
                            <div class="panel-body">
                                <div class="row" v-for="query in info.queries">
                                    <div class="col-sm-12 col-xs-12">
                                        <span class="glyphicon glyphicon-ok" v-if="query.ok == true"></span>
                                        <span class="glyphicon glyphicon-remove" v-if="query.ok == false"></span>
                                        {{ query.query }}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="col-sm-1 col-xs-1"></div>
                </div>
            </div> -->
          <!-- <div class="modal-footer">
              <button type="button" class="btn btn-default" v-on:click="update" v-if="info.isUpdated == false && info.isUpdating == false">
                  <span class="glyphicon glyphicon-refresh"></span>
                  Update
              </button>
              <button type="button" class="btn btn-default" disabled v-if="info.isUpdated == false && info.isUpdating == true">
                  <span class="glyphicon glyphicon-refresh"></span>
                  Updating
              </button>
              <button type="button" class="btn btn-default" v-on:click="cancel" v-if="info.isUpdated == true && info.isUpdating == false">
                  <span class="glyphicon glyphicon-ok"></span>
                  Done
              </button>
              <button type="button" class="btn btn-default" v-on:click="cancel" v-if="info.isUpdated == false && info.isUpdating == false">
                  <span class="glyphicon glyphicon-remove"></span>
                  Ignore
              </button>
          </div> -->
      </form>
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
      "info": null,
      "entry": {},
      "errors": [ "pipapo" ],
    }
  },
  methods: {
    ...mapActions('configuration',[
      'queryDbVersionInfo'
    ])
  },
  computed: {
    ...mapGetters('configuration', [
      'getDbVersionInfo'
    ])
  },
  created() {
    this.queryDbVersionInfo()
  },
  mounted() {
    this.$bvModal.show('view_database_update_modal')
  }
})
</script>