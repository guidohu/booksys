// TODO
// - add users for mobile
// - fix timezone consistency in pie, edit, create, ...
<template>
  <div>
    <SessionEditorModal 
      :defaultValues="selectedSession"
    />
    <SessionDeleteModal
      :session="selectedSession"
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
            :selectedSession="selectedSession"
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
                :session="selectedSession"
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
        :selectedSession="selectedSession"
        :isMobile="isMobile"
        :timezone="getTimezone"
        @prevDay="prevDay"
        @nextDay="nextDay"
        @selectSessionHandler="selectSlot"
      />
      <SessionDetailsCard
        class="mb-1"
        :date="date"
        :session="selectedSession"
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
import Session from '@/dataTypes/session';
import moment from 'moment';
import 'moment-timezone';
import _ from 'lodash';

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
      this.selectedSession = null;
      this.querySessionsForDate();
    },
    nextDay: function(){
      this.date = moment(this.date).add(1, 'days');
      this.selectedSession = null;
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
    selectSlot: function(selectedSession) {
      console.log("selectedSlot", selectedSession);
      console.log(moment(selectedSession.start).format(), moment(selectedSession.end).format());
      const sessionWithSelectedId = this.getSessions.sessions.find(s => selectedSession.id == s.id);
      if(sessionWithSelectedId == null){
        this.selectedSession = new Session(
          selectedSession.id,
          null,
          null,
          selectedSession.start,
          selectedSession.end
        )
      }else{
        this.selectedSession = sessionWithSelectedId;
      }
    },
    sessionDeletedHandler: function() {
      console.log("sessionDeletedHandler");
      this.selectedSession = null;
    },
    showCreateSession: function(){
      console.log("createSession clicked");
      console.log("for selectedSession:", this.selectedSession);
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
  watch: {
    getSessions: function(newInfo, oldInfo){
      // in case we get an update affecting the sessions
      // we will update our selected session too
      if(this.selectedSession != null 
        && this.selectedSession.id != null
        && newInfo.sessions.map(s => s.id).includes(this.selectedSession.id)
      ){
        this.selectedSession = newInfo.sessions.filter(s => s.id == this.selectedSession.id)[0];
      }

      // in case a new session has been created, we select it
      if(newInfo.sessions != null
        && (oldInfo.sessions == null || newInfo.sessions.length > oldInfo.sessions.length)
      ){
        // find the session that is new
        const newIds = newInfo.sessions.map(s => s.id);
        const oldIds = oldInfo.sessions.map(s => s.id);
        const difference = _.difference(newIds, oldIds);
        if(difference.length == 1){
          this.selectedSession = newInfo.sessions.find(s => s.id == difference[0])
        }else{
          console.error("old and new session info differs by more than one session");
        }
      }else if(newInfo.sessions.length < oldInfo.sessions.length){
        // a session has been deleted -> reset selected session
        this.selectedSession = null;
      }
    }
  },
  data() {
    return {
      date: null,
      selectedSession: null
    }
  },
  created() {
    // needed to know the timezone
    this.queryConfiguration();
    // TODO block till we have the config
    // then when promise resolved, call
    // a function that sets the date

    // get day from URL
    const urlParams = new URLSearchParams(window.location.search);
    const urlDate = urlParams.get('date');
    if(urlDate != null){
      const urlDateParsed =  moment(urlDate, 'YYYY-MM-DD').tz(this.getTimezone).startOf('day').format();
      this.date = urlDateParsed;
    }else{
      this.date = moment().tz(this.getTimezone).startOf('day').format();
    }    

    this.querySessionsForDate();
  }
})
</script>