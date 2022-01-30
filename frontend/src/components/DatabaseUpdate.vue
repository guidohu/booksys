<template>
  <modal-container name="dbUpdateModal" :visible="visible">    
    <modal-header title="Database Update" closable="true" @close="cancel">
      <h5 class="modal-title">Database Update</h5>
      <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
    </modal-header>
    <modal-body>
      <div
        v-if="getDbVersionInfo == null"
        class="text-center"
      >
        <div class="spinner-border text-primary" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div v-else>
        <form>
          <div class="row" v-if="errors.length">
            <div class="col-10 offset-sm-1">
              <div class="alert alert-warning">
                <b>Please correct the following error(s):</b>
                <ul>
                  <li
                    v-for="error in errors"
                    :key="error"
                  >
                    {{ error }}
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div class="row" v-if="getDbVersionInfo != null">
            <div class="col-4 offset-sm-1 font-weight-bold">
              App Version
            </div>
            <div class="col-6">
              v{{ getDbVersionInfo.appVersion }}
            </div>
          </div>
          <div class="row" v-if="getDbVersionInfo != null">
            <div class="col-4 offset-sm-1 font-weight-bold">
              DB Version
            </div>
            <div class="col-6">
              {{ dbVersionText }}
            </div>
          </div>
          <div class="row mt-3">
            <div class="col-10 offset-sm-1 font-weight-bold">
              <div class="alert alert-info"
                v-if="
                  getDbVersionInfo != null &&
                    getDbVersionInfo.isUpdated == false &&
                    getDbUpdateResult == null
                "
              >
                The database version is outdated and needs to be upgraded. If the
                upgrade is not done, there might occur data inconsisstencies and
                the app will stop working properly. Thus it is highly recommended
                to upgrade now.
              </div>
              <div class="alert alert-success"
                v-if="
                  getDbVersionInfo.isUpdated == false && getDbIsUpdating == true
                "
              >
                Updating database. Please wait, this might take a few minutes.
              </div>
              <div class="alert alert-success"
                v-if="
                  getDbVersionInfo.isUpdated == true &&
                    getDbIsUpdating == false &&
                    (getDbUpdateResult == null || getDbUpdateResult.success == true)
                "
              >
                Update successful.
              </div>
              <div class="alert alert-danger"
                v-if="
                  getDbIsUpdating == false &&
                    getDbUpdateResult != null &&
                    getDbUpdateResult.success == false
                "
                variant="danger"
                show
              >
                Update failed. Please let the administrator know about this,
                together with the error message below:<br><br>
                <strong>Error:</strong> {{ getDbUpdateResult.msg }}
              </div>
            </div>
          </div>
          <div
            v-if="
              getDbUpdateResult != null &&
                getDbUpdateResult.queries != null &&
                getDbUpdateResult.queries.length > 0 &&
                getDbUpdateResult.ok == false
            "
            class="row mt-3"
          >
            <div class="col-10 offset-sm-1">
              <div class="row"
                v-for="query in getDbUpdateResult.queries"
                :key="query.query"
              >
                <div class="col-sm-12 col-xs-12">
                  <span
                    v-if="query.ok == true"
                    class="glyphicon glyphicon-ok"
                  />
                  <span
                    v-if="query.ok == false"
                    class="glyphicon glyphicon-remove"
                  />
                  {{ query.query }}
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>
    </modal-body>
    <modal-footer>
      <button
        v-if="getDbVersionInfo != null && getDbVersionInfo.isUpdated == false && getDbIsUpdating == false"
        type="button"
        class="btn btn-outline-info"
        @click="updateDb"
      >
        <i class="bi bi-arrow-clockwise"></i>
        Update
      </button>
      <button
        v-if="getDbVersionInfo != null && getDbVersionInfo.isUpdated == false && getDbIsUpdating == true"
        type="button"
        class="btn btn-outline-secondary"
        disabled
      >
        <i class="bi bi-hourglass-split"></i>
        Updating
      </button>
      <button
        v-if="getDbVersionInfo != null && getDbVersionInfo.isUpdated == true && getDbIsUpdating == false"
        type="button"
        class="btn btn-outline-success"
        @click="cancel"
      >
        <i class="bi bi-check2"></i>
        Done
      </button>
      <button
        v-if="getDbVersionInfo != null && getDbVersionInfo.isUpdated == false && getDbIsUpdating == false"
        class="btn btn-outline-danger ml-1"
        type="button"
        @click="cancel"
      >
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import ModalContainer from '@/components/bricks/ModalContainer.vue';
import ModalHeader from '@/components/bricks/ModalHeader.vue';
import ModalBody from '@/components/bricks/ModalBody.vue';
import ModalFooter from '@/components/bricks/ModalFooter.vue';

export default {
  name: "DatabaseUpdateModal",
  components: {
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter
  },
  data() {
    return {
      errors: [],
      dbVersionText: "n/a",
      visible: false,
      modal: null,
    };
  },
  watch: {
    'getDbVersionInfo.dbVersion': function(newValue) {
      console.info("getDbVersionInfo.dbVersion value changed to", newValue);
      this.dbVersionText = this.getVersionText(this.getDbVersionInfo);
    } 
  },
  methods: {
    ...mapActions("configuration", ["queryDbVersionInfo", "updateDb"]),
    cancel: function () {
      this.visible = false;
    },
    getVersionText: function (versionInfo) {
      console.log()
      if (versionInfo == null && this.getDbVersionInfo == null) {
        return "loading";
      }

      if (versionInfo == null) {
        versionInfo = this.getDbVersionInfo;
      }

      if (versionInfo.dbVersion == versionInfo.dbVersionRequired) {
        return "v" + versionInfo.dbVersion;
      } else {
        return (
          "v" +
          versionInfo.dbVersion +
          " (update available to: v" +
          versionInfo.dbVersionRequired +
          ")"
        );
      }
    },
  },
  computed: {
    ...mapGetters("configuration", [
      "getDbVersionInfo",
      "getDbIsUpdating",
      "getDbUpdateResult",
    ]),
  },
  created() {
    this.queryDbVersionInfo()
      .then(() => {
        this.dbVersionText = this.getVersionText();
      })
      .catch((errors) => {
        this.errors = errors
      });
  },
  mounted() {
    this.visible = true;
    this.$store.subscribe((mutation, state) => {
      console.log("Type:",mutation.type)
      console.log("Payload",mutation.payload)
      console.log("state",state)
    });
  },
};
</script>
