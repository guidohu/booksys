<template>
  <b-card no-body class="text-left">
    <b-card-header>
      <b-row>
        <b-col cols="4" class="text-right">
          <b-button variant="outline-info" class="btn-xs" v-on:click="prevDay">
            <b-icon icon="arrow-left-short"></b-icon>
          </b-button>
        </b-col>
        <b-col cols="4" class="text-center">
          {{dateString}}
        </b-col>
        <b-col cols="4" class="text-left">
          <b-button variant="outline-info" class="btn-xs" v-on:click="nextDay">
            <b-icon icon="arrow-right-short"></b-icon>
          </b-button>
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body>
      <b-row>
        <Pie 
          :sessionData="sessionData" 
          :selectedSession="selectedSession"
          :properties="properties"
          @selectHandler="selectSession"/>
      </b-row>
      <b-row v-if="isToday && selectedSession != null && selectedSession.id != null" class="text-center">
        <b-col cols="12" class="text-center">
          <b-button v-on:click="navigateSessionStart" type="button" variant="outline-success">Start Session</b-button>
        </b-col>
      </b-row>
    </b-card-body>
  </b-card>
</template>

<script>
import Vue from 'vue';
import { mapGetters } from 'vuex';
import Pie from "./Pie.vue";
import moment from 'moment-timezone';

export default Vue.extend({
  name: 'SessionDayCard',
  data() {
    return {
      properties: {
        containerWidth: 700,
        containerHeight: 350,
        circleX: 300,
        circleY: 170,
        circleRadius: 100,
        animate: true,
        labels: true
      }
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getTimezone'
    ]),
    dateString: function() {
      console.log("sessionData", this.sessionData);
      return moment(this.sessionData.window_start, "X").format("dddd DD.MM.YYYY");
    },
    isToday: function() {
      console.log(this.getTimezone);
      const today = moment().tz(this.getTimezone).startOf('day').format();
      const sessionDay = moment(this.sessionData.window_start, "X").tz(this.getTimezone).startOf('day').format();
      if(today == sessionDay){
        return true;
      }
      return false;
    }
  },
  methods: {
    prevDay: function(){
      this.$emit('prevDay');
    },
    nextDay: function(){
      this.$emit('nextDay');
    },
    selectSession: function(slot) {
      this.$emit('selectSessionHandler', slot);
    },
    navigateSessionStart: function() {
      // TODO adjust correct link
      console.log("TODO adjust correct link");
      window.location.href = '/watch.html?sessionId='+this.selectedSession.id;
    }
  },
  components: {
    Pie
  },
  props: ['isMobile', 'sessionData', 'selectedSession', 'timezone'],
  created() {
    if(this.isMobile != null && this.isMobile == true){
      this.properties = {
        containerHeight: 300,
        containerWidth:  350,
        circleX:         175,
        circleY:         150,
        circleRadius:    100,
        animation:       false,
        labels:          true,
      }
    }
    if(this.timezone != null){
      this.properties.timezone = this.timezone;
    }else{
      this.properties.timezone = "UTC";
    }
  }
})
</script>

<style scoped>

  .btn-xs {
    padding: .2rem .4rem;
    font-size: .6rem;
    line-height: 1;
    border-radius: .2rem;
  }
</style>