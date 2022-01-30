<template>
  <div
    v-if="isDesktop"
    class="display"
  >
    <main-title title-name="Stop Watch" />
    <b-row class="ml-1 mr-1">
      <b-col cols="8">
        <WarningBox
          v-if="errors.length > 0"
          :errors="errors"
        />
        <StopWatchCard :session-id="sessionId" />
      </b-col>
      <b-col cols="4">
        <SessionHeatListCard
          :session-id="sessionId"
          class="side-bar-heats-component"
        />
        <ConditionInfoCard
          v-if="getSessionConditionInfo != null"
          class="mt-2"
          :sunrise="getSessionConditionInfo.sunrise"
          :sunset="getSessionConditionInfo.sunset"
        />
      </b-col>
    </b-row>
    <div class="bottom mr-2">
      <b-button
        variant="outline-light"
        to="/dashboard"
      >
        <b-icon-house />
        HOME
      </b-button>
    </div>
  </div>
  <div v-else>
    <NavbarMobile title="Stop Watch" />
    <WarningBox
      v-if="errors.length > 0"
      :errors="errors"
    />
    <StopWatchCard :session-id="sessionId" />
    <SessionHeatListCard
      :session-id="sessionId"
      class="mt-2"
    />
    <ConditionInfoCard
      v-if="getSessionConditionInfo != null"
      class="mt-2"
      :sunrise="getSessionConditionInfo.sunrise"
      :sunset="getSessionConditionInfo.sunset"
    />
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import WarningBox from "@/components/WarningBox";
import StopWatchCard from "@/components/StopWatchCard";
import SessionHeatListCard from "@/components/SessionHeatListCard";
import ConditionInfoCard from "@/components/ConditionInfoCard";
import MainTitle from "@/components/MainTitle";
import { BRow, BCol, BButton, BIconHouse } from "bootstrap-vue";

export default {
  name: "WSWatch",
  components: {
    NavbarMobile,
    MainTitle,
    WarningBox,
    StopWatchCard,
    SessionHeatListCard,
    ConditionInfoCard,
    BRow,
    BCol,
    BButton,
    BIconHouse,
  },
  data() {
    return {
      sessionId: null,
      errors: [],
    };
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
    ...mapGetters("sessions", ["getSessionConditionInfo"]),
  },
  created() {
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = Number(urlParams.get("sessionId"));
    if (sessionId == null || !Number(sessionId) > 0) {
      this.errors = ["No session parameter given in URL"];
    } else {
      this.sessionId = sessionId;
    }
  },
};
</script>

<style>
.side-bar-heats-component {
  max-height: 300px;
  overflow-y: scroll;
}
</style>
