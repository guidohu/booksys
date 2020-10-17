<template>
  <b-card no-body class="text-left">
    <b-card-header>
      <b-row>
        <b-col cols="4" class="text-right">
          <b-button variant="outline-info" class="btn-xs" v-on:click="prevMonth">
            <b-icon icon="arrow-left-short"></b-icon>
          </b-button>
        </b-col>
        <b-col cols="4" class="text-center">
          {{monthString}}
        </b-col>
        <b-col cols="4" class="text-left">
          <b-button variant="outline-info" class="btn-xs" v-on:click="nextMonth">
            <b-icon icon="arrow-right-short"></b-icon>
          </b-button>
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body>
      <b-row v-for="k in Array.from(Array(6).keys())" :key="k">
        <span v-for="j in Array.from(Array(7).keys())" :key="j">
          <Pie
            :sessionData="sessionData[k*7+j]"
            :properties="properties"
            />
        </span>
      </b-row>
      <!-- <b-row>
        <Pie 
          :sessionData="sessionData" 
          :selectedSession="selectedSession"
          :properties="properties"
          @selectHandler="selectSession"/>
      </b-row> -->
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
        containerWidth: 90,
        containerHeight: 64,
        circleX: 45,
        circleY: 34,
        circleRadius: 23,
        animate: false,
        labels: false
      }
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getTimezone'
    ]),
    monthString: function() {
      return moment(this.month).format("MMMM YYYY");
    }
  },
  methods: {
    prevMonth: function(){
      this.$emit('prevMonth');
    },
    nextMonth: function(){
      this.$emit('nextMonth');
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
  props: ['sessionData', 'month'],
  created() {
    if(this.isMobile != null && this.isMobile == true){
      this.properties = {
        containerHeight: 300,
        containerWidth:  350,
        circleX:         175,
        circleY:         150,
        circleRadius:    90,
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