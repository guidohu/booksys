<template>
  <b-card no-body>
    <HeatCommentModal :defaultComment="comment" @commentChangeHandler="changeComment"/>
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
          <b-form-select v-model="selectedRiderId" :options="selectableRiders" :disabled="getIsRunning"/>
        </b-form-group>
      </b-col>
    </b-row>
    <b-row class="ml-1 mr-1 mb-2">
      <b-col cols="12">
        <b-button-group v-if="selectedRiderId != null">
          <b-button size="lg" variant="info" v-if="!getIsRunning">
            <b-icon icon="pencil-square"/>
            Back
          </b-button>
          <b-button size="lg" variant="info" v-if="!getIsRunning" v-on:click="startTakingTime">
            <b-icon icon="pencil-square"/>
            Start
          </b-button>
          <b-button size="lg" variant="info" v-if="getIsRunning && !getIsPaused" v-on:click="pauseTakingTime">
            <b-icon icon="pencil-square"/>
            Pause
          </b-button>
          <b-button size="lg" variant="info" v-if="getIsPaused" v-on:click="resumeTakingTime">
            <b-icon icon="pencil-square"/>
            Resume
          </b-button>
          <b-button size="lg" variant="info" v-if="getIsRunning" v-on:click="finishTakingTime">
            <b-icon icon="pencil-square"/>
            Finish
          </b-button>
          <b-dropdown right variant="info" v-if="selectedRiderId != null || getIsRunning">
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
import HeatCommentModal from '@/components/HeatCommentModal';

export default Vue.extend({
  name: 'StopWatchCard',
  components: {
    HeatCommentModal
  },
  data() {
    return {
      selectedRiderId: null,
      selectableRiders: [],
      comment: null,
      displayTime: "00:00"
    }
  },
  computed: {
    ...mapGetters('stopwatch', [
      'getDisplayTime',
      'getIsPaused',
      'getIsRunning',
    ])
  },
  watch: {
    selectedRiderId: function(newRiderId) {
      if(newRiderId != null){
        this.setUserId = newRiderId;
      }
    }
  },
  methods: {
    showHeatCommentModal: function(){
      this.$bvModal.show('heatCommentModal');
    },
    changeComment: function(newComment){
      this.comment = newComment;
    },
    ...mapActions('stopwatch', [
      'startTakingTime',
      'pauseTakingTime',
      'resumeTakingTime',
      'finishTakingTime',
      'setSessionId',
      'setUserId'
    ])
  },
  created() {
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = Number(urlParams.get('sessionId'));
    if(sessionId == null){
      this.errors = [ "No session selected" ];
    }else{
      this.setSessionId(sessionId);
    }
  }
})
</script>

<style scoped>
  .stopwatch-text {
    font-size: 4rem
  }
</style>