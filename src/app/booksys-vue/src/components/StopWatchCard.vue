<template>
  <b-card no-body>
    <HeatCommentModal :defaultComment="comment" @commentChangeHandler="changeComment"/>
    <b-row class="mt-4 ml-1 mr-1 mb-4">
      <b-col cols="12">
        <b-card bg-variant="light">
          <b-row class="stopwatch-text text-center">
            <b-col cols="12">
              00:00
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
          <b-form-select v-model="selectedRiderId" :options="selectableRiders"/>
        </b-form-group>
      </b-col>
    </b-row>
    <b-row class="ml-1 mr-1 mb-2">
      <b-col cols="12">
        <b-button-group>
          <b-button size="lg" variant="info">
            <b-icon icon="pencil-square"/>
            Back
          </b-button>
          <b-button size="lg" variant="info">
            <b-icon icon="pencil-square"/>
            Start
          </b-button>
          <b-button size="lg" variant="info">
            <b-icon icon="pencil-square"/>
            Pause
          </b-button>
          <b-button size="lg" variant="info">
            <b-icon icon="pencil-square"/>
            Resume
          </b-button>
          <b-button size="lg" variant="info">
            <b-icon icon="pencil-square"/>
            Finish
          </b-button>
          <b-dropdown right variant="info">
            <b-dropdown-item v-on:click="showHeatCommentModal">Add Comment</b-dropdown-item>
            <b-dropdown-divider/>
            <b-dropdown-item>Cancel Heat</b-dropdown-item>
          </b-dropdown>
        </b-button-group>
      </b-col>
    </b-row>
  </b-card>
</template>

<script>
import Vue from 'vue';
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
      comment: null
    }
  },
  methods: {
    showHeatCommentModal: function(){
      this.$bvModal.show('heatCommentModal');
    },
    changeComment: function(newComment){
      this.comment = newComment;
    }
  }
})
</script>

<style scoped>
  .stopwatch-text {
    font-size: 4rem
  }
</style>