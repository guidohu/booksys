<template>
  <b-card no-body>
    <HeatCommentModal :defaultComment="comment" @commentChangeHandler="changeComment"/>
    <b-row v-if="errors.length > 0">
      <b-col cols="12">
        <WarningBox :errors="errors"/>
      </b-col>
    </b-row>
    <b-row class="mt-4 ml-1 mr-1 mb-4">
      <b-col cols="12">
        <b-card bg-variant="light">
          <b-row class="stopwatch-text text-center">
            <b-col cols="12">
              {{getDisplayTime}}
            </b-col>
          </b-row>
        </b-card>
      </b-col>
    </b-row>
    <b-row v-if="comment!=null" class="ml-1 mr-1">
      <b-col cols="12">
        <b-form-group>
          <b-input-group>
            <b-form-input v-model="comment"/>
            <template #append>
              <b-button v-on:click="showHeatCommentModal">
                <b-icon icon="pencil-square"/>
                <!-- Edit -->
              </b-button>
            </template>
          </b-input-group>
        </b-form-group>
      </b-col>
    </b-row>
    <b-row class="ml-1 mr-1">
      <b-col cols="12">
        <b-form-group>
          <b-form-select v-model="selectedRiderId" :options="selectableRiders" :disabled="getIsRunning || sessionId==null"/>
        </b-form-group>
      </b-col>
    </b-row>
    <b-row class="ml-1 mr-1 mb-2 text-center">
      <b-col cols="12">
        <b-button-group v-if="selectedRiderId != null">
          <b-button size="lg" variant="outline-info" v-if="!getIsRunning">
            <b-icon icon="arrow-left"/>
            Back
          </b-button>
          <b-button size="lg" variant="outline-info" v-if="!getIsRunning" v-on:click="startTakingTime">
            <b-icon icon="play"/>
            Start
          </b-button>
          <b-button size="lg" variant="outline-info" v-if="getIsRunning && !getIsPaused" v-on:click="pauseTakingTime">
            <b-icon icon="pause"/>
            Pause
          </b-button>
          <b-button size="lg" variant="outline-info" v-if="getIsPaused" v-on:click="resumeTakingTime">
            <b-icon icon="eject" rotate="90"/>
            Resume
          </b-button>
          <b-button size="lg" variant="outline-info" v-if="getIsRunning" v-on:click="finish">
            <b-icon icon="check-square"/>
            Finish
          </b-button>
          <b-dropdown right variant="outline-info" v-if="selectedRiderId != null || getIsRunning">
            <b-dropdown-item v-if="selectedRiderId != null" v-on:click="showHeatCommentModal">Add Comment</b-dropdown-item>
            <b-dropdown-divider v-if="getIsRunning && selectedRiderId != null"/>
            <b-dropdown-item v-if="getIsRunning">Cancel Heat</b-dropdown-item>
          </b-dropdown>
        </b-button-group>
      </b-col>
    </b-row>
  </b-card>
</template>

<script>
import Vue from 'vue';
import { mapGetters, mapActions } from 'vuex';
import WarningBox from '@/components/WarningBox';
import HeatCommentModal from '@/components/HeatCommentModal';

export default Vue.extend({
  name: 'StopWatchCard',
  components: {
    WarningBox,
    HeatCommentModal
  },
  props: [ 'sessionId' ],
  data() {
    return {
      selectedRiderId: null,
      selectableRiders: [],
      comment: null,
      displayTime: "00:00",
      errors: []
    }
  },
  computed: {
    ...mapGetters('stopwatch', [
      'getDisplayTime',
      'getIsPaused',
      'getIsRunning',
      'getUserId'
    ]),
    ...mapGetters('sessions', [
      'getSession'
    ])
  },
  watch: {
    getUserId: function(selectedUserId) {
      if(selectedUserId != this.selectedRiderId){
        this.selectedRiderId = selectedUserId;
      }
    },
    selectedRiderId: function(newRiderId) {
      if(newRiderId != null){
        console.log("setUserId to", newRiderId);
        this.setUserId(newRiderId);
      }
    },
    getSession: function(newSession) {
      console.log("session data changed to", newSession);
      if(newSession != null){
        this.setSelectableRiders(newSession.riders);
      }
    }
  },
  methods: {
    setSelectableRiders: function(riders){
      console.log("set selectable riders", riders);
      if(riders != null && riders.length > 0){
        this.selectableRiders = [{
          value: null,
          text: 'Please select a rider'
        }];
        const ridersToAdd = riders.map(r => {
          return {
            value: r.id,
            text: r.first_name + " " + r.last_name
          }
        })
        console.log("ridersToAdd", ridersToAdd);
        this.selectableRiders.push(...ridersToAdd);
        console.log("new selectableRiders", this.selectableRiders);
      }
    },
    showHeatCommentModal: function(){
      this.$bvModal.show('heatCommentModal');
    },
    changeComment: function(newComment){
      this.comment = newComment;
    },
    finish: function(){
      this.finishTakingTime()
      .catch(errors => this.errors = errors);
    },
    ...mapActions('stopwatch', [
      'startTakingTime',
      'pauseTakingTime',
      'resumeTakingTime',
      'finishTakingTime',
      'setSessionId',
      'setUserId'
    ]),
    ...mapActions('sessions', [
      'querySession'
    ])
  },
  created() {
    if(this.sessionId == null){
      this.errors = [ "No session selected" ];
    }else{
      this.setSessionId(this.sessionId);
      this.querySession(this.sessionId)
      .then(() => console.log("queried session"))
      .catch((errors) => this.errors = errors);
    }
  }
})
</script>

<style scoped>
  .stopwatch-text {
    font-size: 4rem
  }
</style>