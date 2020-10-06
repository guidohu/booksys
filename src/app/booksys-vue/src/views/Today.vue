<template>
  <div v-if="isDesktop" class="display">
    <b-row>
      <b-col>
        <div class="main-title">
          <h3>
            Book Your Session
          </h3>
        </div>
      </b-col>
    </b-row>
    <b-row class="ml-1 mr-1">
      <b-col cols="8">
        <SessionDayCard
          :sessionData="getSessions"
          :isMobile="isMobile"
          :timezone="getTimezone"
          @prevDay="prevDay"
          @nextDay="nextDay"
          @selectSessionHandler="selectSlot"
        />
      </b-col>
      <b-col cols="4">
        <b-row>
          <b-col cols="12">
            <SessionDetailsCard            
              :date="date"
              :sessionTime="slot"
              @createSessionHandler="createSession"
              @addRidersHandler="addRiders"
            />
          </b-col>
        </b-row>
        <b-row class="mt-1">
          <b-col cols="12">
            <ConditionInfoCard
              :sunrise="getSessions.sunrise"
              :sunset="getSessions.sunset"
            />
          </b-col>
        </b-row>
      </b-col>
    </b-row>
    <div class="bottom mr-2">
      <b-button class="mr-1" variant="outline-light" to="/calendar">
        <b-icon-calendar3></b-icon-calendar3>
        CALENDAR
      </b-button>
      <b-button variant="outline-light" to="/dashboard">
        <b-icon-house></b-icon-house>
        HOME
      </b-button>
    </div>
  </div>
  <div v-else>
    <NavbarMobile title="Book Your Session"/>
    <SessionDayCard
      class="mb-1"
      :sessionData="getSessions"
      :isMobile="isMobile"
      :timezone="getTimezone"
      @prevDay="prevDay"
      @nextDay="nextDay"
      @selectSessionHandler="selectSlot"
    />
    <SessionDetailsCard
      class="mb-1"
      :date="date"
      :sessionTime="slot"
      @createSessionHandler="createSession"
      @addRidersHandler="addRiders"
    />
    <ConditionInfoCard
      :sunrise="getSessions.sunrise"
      :sunset="getSessions.sunset"
    />
  </div>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import { BooksysBrowser } from '@/libs/browser';
import NavbarMobile from '@/components/NavbarMobile';
import ConditionInfoCard from '@/components/ConditionInfoCard';
import SessionDayCard from '@/components/SessionDayCard';
import SessionDetailsCard from '@/components/SessionDetailsCard';
import moment from 'moment';
import 'moment-timezone';

export default Vue.extend({
  name: 'Today',
  components: {
    NavbarMobile,
    ConditionInfoCard,
    SessionDayCard,
    SessionDetailsCard
  },
  computed: {
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    },
    ...mapGetters('configuration', [
      'getTimezone'
    ]),
    ...mapGetters('sessions', [
      'getSessions'
    ])
  },
  methods: {
    ...mapActions('configuration', [
      'queryConfiguration'
    ]),
    ...mapActions('sessions', [
      'querySessions'
    ]),
    prevDay: function(){
      this.date.add(-1, 'days');
      this.slot = null;
      this.querySessionsForDate();
    },
    nextDay: function(){
      this.date.add(1, 'days');
      this.slot = null;
      this.querySessionsForDate();
    },
    querySessionsForDate: function() {
      console.log("the time", this.date);
      const dateStart = this.date.tz(this.getTimezone).startOf('day');
      const dateEnd = this.date.tz(this.getTimezone).endOf('day');
      console.log("dateStart:",dateStart);
      console.log("dateEnd:",dateEnd);
      console.log("startWindow", dateStart.format("X"));
      console.log("endWindow", dateEnd.format("X"));

      // query get_booking_day
      this.querySessions({
        start: dateStart,
        end: dateEnd
      });
    },
    selectSlot: function(slot) {
      console.log("selectedSlot", slot);
      this.slot = slot;
    },
    createSession: function(){
      console.log("createSession clicked");
    },
    addRiders: function(){
      console.log("addRiders");
    }
  },
  data() {
    return {
      date: moment(),
      slot: null
    }
  },
  created() {
    // get day from URL
    // TODO

    // needed to know the timezone
    this.queryConfiguration();

    this.querySessionsForDate();
  }
})
</script>