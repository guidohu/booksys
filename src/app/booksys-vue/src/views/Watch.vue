<template>
  <div v-if="isDesktop" class="display">
    <b-row>
      <b-col>
        <div class="main-title">
          <h3>
            Stop Watch
          </h3>
        </div>
      </b-col>
    </b-row>
    <b-row class="ml-1 mr-1">
      <b-col cols="8">
        <WarningBox v-if="errors.length > 0" :errors="errors"/>
        <StopWatchCard :sessionId="sessionId"/>
      </b-col>
      <b-col cols="4">
        <SessionHeatListCard :sessionId="sessionId"/>
      </b-col>
    </b-row>
    <div class="bottom mr-2">
      <b-button variant="outline-light" to="/dashboard">
        <b-icon-house></b-icon-house>
        HOME
      </b-button>
    </div>
  </div>
  <div v-else>
    <NavbarMobile title="Stop Watch"/>
    <WarningBox v-if="errors.length > 0" :errors="errors"/>
    <StopWatchCard :sessionId="sessionId"/>
    <SessionHeatListCard :sessionId="sessionId" class="mt-2"/>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarMobile from '@/components/NavbarMobile';
import WarningBox from '@/components/WarningBox';
import StopWatchCard from '@/components/StopWatchCard';
import SessionHeatListCard from '@/components/SessionHeatListCard';

export default Vue.extend({
  name: "Watch",
  components: {
    NavbarMobile,
    WarningBox,
    StopWatchCard,
    SessionHeatListCard
  },
  data() {
    return {
      sessionId: null,
      errors: []
    }
  },
  computed: {
    isDesktop: function() {
      return !BooksysBrowser.isMobile();
    }
  },
  created() {
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = Number(urlParams.get('sessionId'));
    if(sessionId == null || !Number(sessionId) > 0){
      this.errors = [ "No session parameter given in URL" ];
    }else{
      this.sessionId = sessionId;
    }
  }
})
</script>