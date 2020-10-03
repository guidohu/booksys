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
      <Pie 
        :sessionData="sessionData" 
        :properties="properties"
        @selectHandler="selectSession"/>
    </b-card-body>
  </b-card>
</template>

<script>
import Vue from 'vue';
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
    dateString: function() {
      console.log("sessionData", this.sessionData);
      return moment(this.sessionData.window_start, "X").format("dddd DD.MM.YYYY");
    }
  },
  methods: {
    prevDay: function(){
      console.log('prevDay');
      this.$emit('prevDay');
    },
    nextDay: function(){
      console.log('nextDay');
      this.$emit('nextDay');
    },
    selectSession: function(slot) {
      this.$emit('selectSessionHandler', slot);
    }
  },
  components: {
    Pie
  },
  props: ['isMobile', 'sessionData', 'timezone'],
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