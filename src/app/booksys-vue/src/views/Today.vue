<template>
  <div>
    <SessionEditorModal 
      :defaultValues="slot"
      @sessionCreatedHandler="sessionCreatedHandler"
    />
    <SessionDeleteModal
      :session="slot"
      @sessionDeletedHandler="sessionDeletedHandler"
    />
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
          <SessionDayCard v-if="getSessions != null"
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
                :session="slot"
                @createSessionHandler="showCreateSession"
                @editSessionHandler="showCreateSession"
                @deleteSessionHandler="showDeleteSession"
                @addRidersHandler="addRiders"
              />
            </b-col>
          </b-row>
          <b-row class="mt-1">
            <b-col cols="12">
              <ConditionInfoCard v-if="sunrise!=null"
                :sunrise="sunrise"
                :sunset="sunset"
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
      <SessionDayCard v-if="getSessions != null"
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
        :session="slot"
        @createSessionHandler="showCreateSession"
        @editSessionHandler="showCreateSession"
        @deleteSessionHandler="showDeleteSession"
        @addRidersHandler="addRiders"
      />
      <ConditionInfoCard v-if="sunrise!=null"
        :sunrise="sunrise"
        :sunset="sunset"
      />
    </div>
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
import SessionEditorModal from '@/components/SessionEditorModal';
import SessionDeleteModal from '@/components/SessionDeleteModal';
import moment from 'moment';
import 'moment-timezone';

export default Vue.extend({
  name: 'Today',
  components: {
    NavbarMobile,
    ConditionInfoCard,
    SessionDayCard,
    SessionDetailsCard,
    SessionEditorModal,
    SessionDeleteModal
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
    ]),
    sunrise: function(){
      const sessions = this.getSessions;
      if(sessions == null){
        return null;
      }else{
        return sessions.sunrise;
      }
    },
    sunset: function(){
      const sessions = this.getSessions;
      if(sessions == null){
        return null;
      }else{
        return sessions.sunset;
      }
    }
  },
  methods: {
    ...mapActions('configuration', [
      'queryConfiguration'
    ]),
    ...mapActions('sessions', [
      'querySessions',
      'createSession'
    ]),
    prevDay: function(){
      this.date = moment(this.date).add(-1, 'days');
      this.slot = null;
      this.querySessionsForDate();
    },
    nextDay: function(){
      this.date = moment(this.date).add(1, 'days');
      this.slot = null;
      this.querySessionsForDate();
    },
    querySessionsForDate: function() {
      console.log("the time", this.date);
      const dateStart = moment(this.date).startOf('day').format();
      const dateEnd = moment(this.date).endOf('day').format();
      console.log("dateStart:",dateStart);
      console.log("dateEnd:",dateEnd);

      // query get_booking_day
      this.querySessions({
        start: dateStart,
        end: dateEnd
      });
    },
    selectSlot: function(slot) {
      console.log("selectedSlot", slot);
      console.log(slot.start.format(), slot.end.format());
      this.slot = slot;
    },
    sessionCreatedHandler: function(sessionId) {
      console.log("sessionCreatedHandler: sessionId", sessionId);
      this.slot = null;
    },
    sessionDeletedHandler: function() {
      console.log("sessionDeletedHandler");
      this.slot = null;
    },
    showCreateSession: function(){
      console.log("createSession clicked");
      console.log("for slot:", this.slot);
      console.log("show session editor");
      console.log(this.$bvModal.show('sessionEditorModal'));
    },
    showDeleteSession: function(){
      console.log("showDeleteSession");
      this.$bvModal.show('sessionDeleteModal');
    },
    addRiders: function(){
      console.log("addRiders");
    }
  },
  data() {
    return {
      date: null,
      slot: null
    }
  },
  created() {
    // get day from URL
    const now = moment().format();
    console.log(moment());
    console.log("now:", now);
    this.date = now;
    console.log(this.date);
    // TODO

    // needed to know the timezone
    this.queryConfiguration();

    this.querySessionsForDate();
  }
})
</script>