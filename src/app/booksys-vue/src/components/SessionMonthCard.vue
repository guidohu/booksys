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
    <b-card-body ref="calendarBody">
      <table>
        <tr>
          <th class="table-title-text">Monday</th>
          <th class="table-title-text">Tuesday</th>
          <th class="table-title-text">Wednesday</th>
          <th class="table-title-text">Thursday</th>
          <th class="table-title-text">Friday</th>
          <th class="table-title-text">Saturday</th>
          <th class="table-title-text">Sunday</th>
        </tr>
        <tr v-for="i in Array.from(Array(6).keys())" :key="i">
          <td v-for="j in Array.from(Array(7).keys())" :key="j">
            <div class="calendar-day-box">
              <Pie
                :sessionData="sessionData[i*7+j]"
                :properties="properties"
                :key="i*7+j"
              />
              <div class="day-number">
                {{ getDay(sessionData[i*7+j].window_start) }}
              </div>
            </div>
          </td>
        </tr> 
      </table>
    </b-card-body>
  </b-card>
</template>

<script>
import Vue from 'vue';
import { mapGetters } from 'vuex';
import Pie from "./Pie.vue";
import moment from 'moment-timezone';
import { BooksysBrowser } from '@/libs/browser';

export default Vue.extend({
  name: 'SessionMonthCard',
  data() {
    return {
      properties: {
        containerWidth: 84,
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
    },
    getDay: function(isoTime) {
      return moment(isoTime).format("DD");
    },
    isMobile: function() {
      return BooksysBrowser.isMobile();
    }
  },
  components: {
    Pie
  },
  props: ['sessionData', 'month'],
  mounted() {
    // if(this.isMobile()){
      console.log("set mobile properties, mounted:",this.$refs.calendarBody.clientWidth);
      const totalWidth = this.$refs.calendarBody.clientWidth - 61;
      const cardWidth = (totalWidth / 7);
      const cardHeight = (cardWidth * .75);
      const radius     = Math.min(cardWidth, cardHeight) / 2 *.7
      this.properties = {
          containerHeight: cardHeight,
          containerWidth:  cardWidth,
          circleX:         (cardWidth - 3) / 2,
          circleY:         cardHeight / 2,
          circleRadius:    radius,
          animation:       false,
          labels:          false,
      }
    // }
  }
})
</script>

<style scoped>
  .calendar-day-box {
    border-radius: .1em;
    border: 0.01em;
    border-color: black;
    background: #eceaea;
    margin-right: 0.01em;
    margin-left: 0.01em;
    margin-top: 0.01em;
    margin-bottom: 0.01em;
    position: relative;
  }

  .calendar-day-box:hover {
    background: #c9c9c9;
  } 

  .day-number { 
    position: absolute;
    bottom: 0;
    left: 0.1rem;
    font-size: 0.5rem;
  }

  .table-title-text {
    font-size: 0.3rem;
  }

  .btn-xs {
    padding: .2rem .4rem;
    font-size: .6rem;
    line-height: 1;
    border-radius: .2rem;
  }
</style>