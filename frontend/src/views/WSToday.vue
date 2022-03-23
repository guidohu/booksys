<template>
  <subpage-container title="Book your Session">
    <session-editor-modal
      v-model:visible="showSessionEditorModal"
      :default-values="selectedSession"
    />
    <session-delete-modal
      v-model:visible="showSessionDeleteModal"
      :session="selectedSession"
      @sessionDeletedHandler="sessionDeletedHandler"
    />
    <show-for-desktop>
      <div class="row mx-1">
        <div class="col-8">
          <session-day-card
            v-if="getSessions != null"
            :session-data="getSessions"
            :selected-session="selectedSession"
            :is-mobile="false"
            :timezone="getTimezone"
            @prevDay="prevDay"
            @nextDay="nextDay"
            @selectSessionHandler="selectSlot"
          />
        </div>
        <div class="col-4">
          <session-details-card
            :date="date"
            :session="selectedSession"
            @createSessionHandler="showCreateSession"
            @editSessionHandler="showCreateSession"
            @deleteSessionHandler="showDeleteSession"
          />
          <condition-info-card
            v-if="sunrise != null"
            class="mt-2"
            :sunrise="sunrise"
            :sunset="sunset"
          />
        </div>
      </div>
    </show-for-desktop>
    <show-for-mobile>
      <session-day-card
        class="mx-1 mt-1"
        v-if="getSessions != null"
        :session-data="getSessions"
        :selected-session="selectedSession"
        :is-mobile="true"
        :timezone="getTimezone"
        @prevDay="prevDay"
        @nextDay="nextDay"
        @selectSessionHandler="selectSlot"
      />
      <session-details-card
        class="mx-1 mt-2"
        :date="date"
        :session="selectedSession"
        @createSessionHandler="showCreateSession"
        @editSessionHandler="showCreateSession"
        @deleteSessionHandler="showDeleteSession"
      />
      <condition-info-card
        class="mx-1 mt-2"
        v-if="sunrise != null"
        :sunrise="sunrise"
        :sunset="sunset"
      />
    </show-for-mobile>
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/calendar">
        <i class="bi bi-calendar"></i>
        CALENDAR
      </router-link>
      <router-link
        tag="button"
        class="btn btn-outline-light ms-1"
        to="/dashboard"
      >
        <i class="bi bi-house"></i>
        HOME
      </router-link>
    </template>
  </subpage-container>
</template>

<script>
import { defineAsyncComponent } from "vue";
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import ConditionInfoCard from "@/components/ConditionInfoCard";
import SessionDayCard from "@/components/SessionDayCard";
import SessionDetailsCard from "@/components/SessionDetailsCard";
import Session from "@/dataTypes/session";
import * as dayjs from "dayjs";
import * as dayjsCustomParseFormat from "dayjs/plugin/customParseFormat";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";

import difference from "lodash/difference";
import SubpageContainer from "../components/bricks/SubpageContainer.vue";
import ShowForMobile from "../components/bricks/ShowForMobile.vue";
import ShowForDesktop from "../components/bricks/ShowForDesktop.vue";

const SessionEditorModal = defineAsyncComponent(() =>
  import("@/components/SessionEditorModal")
);
const SessionDeleteModal = defineAsyncComponent(() =>
  import("@/components/SessionDeleteModal")
);

dayjs.extend(dayjsCustomParseFormat);
dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

export default {
  name: "WSToday",
  components: {
    ConditionInfoCard,
    SessionDayCard,
    SessionDetailsCard,
    SessionEditorModal,
    SessionDeleteModal,
    ShowForMobile,
    ShowForDesktop,
    SubpageContainer,
  },
  computed: {
    isMobile: function () {
      return BooksysBrowser.isMobile();
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
    ...mapGetters("configuration", ["getTimezone"]),
    ...mapGetters("sessions", ["getSessions"]),
    sunrise: function () {
      const sessions = this.getSessions;
      if (sessions == null) {
        return null;
      } else {
        return sessions.sunrise;
      }
    },
    sunset: function () {
      const sessions = this.getSessions;
      if (sessions == null) {
        return null;
      } else {
        return sessions.sunset;
      }
    },
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    ...mapActions("sessions", ["querySessions", "createSession"]),
    prevDay: function () {
      this.date = dayjs(this.date).add(-1, "days");
      this.selectedSession = null;
      this.querySessionsForDate();
    },
    nextDay: function () {
      this.date = dayjs(this.date).add(1, "days");
      this.selectedSession = null;
      this.querySessionsForDate();
    },
    querySessionsForDate: function () {
      const dateStart = dayjs(this.date).startOf("day").format();
      const dateEnd = dayjs(this.date).endOf("day").format();

      // query get_booking_day
      this.querySessions({
        start: dateStart,
        end: dateEnd,
      });
    },
    selectSlot: function (selectedSession) {
      console.log("selectedSlot", selectedSession);
      console.log(
        dayjs(selectedSession.start).format(),
        dayjs(selectedSession.end).format()
      );
      const sessionWithSelectedId = this.getSessions.sessions.find(
        (s) => selectedSession.id == s.id
      );
      if (sessionWithSelectedId == null) {
        this.selectedSession = new Session(
          selectedSession.id,
          null,
          null,
          selectedSession.start,
          selectedSession.end
        );
      } else {
        this.selectedSession = sessionWithSelectedId;
      }
    },
    sessionDeletedHandler: function () {
      console.log("sessionDeletedHandler");
      this.selectedSession = null;
    },
    showCreateSession: function () {
      console.log("createSession clicked");
      console.log("for selectedSession:", this.selectedSession);
      console.log("show session editor");
      this.showSessionEditorModal = true;
    },
    showDeleteSession: function () {
      console.log("showDeleteSession");
      this.showSessionDeleteModal = true;
    },
  },
  watch: {
    getSessions: function (newInfo, oldInfo) {
      if (newInfo == null || newInfo.sessions == null) {
        this.selectedSession = null;
        return;
      }

      // in case we get an update affecting the sessions
      // we will update our selected session too
      if (
        this.selectedSession != null &&
        this.selectedSession.id != null &&
        newInfo.sessions.map((s) => s.id).includes(this.selectedSession.id)
      ) {
        this.selectedSession = newInfo.sessions.filter(
          (s) => s.id == this.selectedSession.id
        )[0];
      }

      // in case a new session has been created, we select it
      if (
        newInfo.sessions != null &&
        (oldInfo == null ||
          oldInfo.sessions == null ||
          newInfo.sessions.length > oldInfo.sessions.length)
      ) {
        // find the session that is new
        const newIds = newInfo.sessions.map((s) => s.id);
        const oldIds =
          oldInfo != null && oldInfo.sessions != null
            ? oldInfo.sessions.map((s) => s.id)
            : [];
        const diff = difference(newIds, oldIds);
        if (diff.length == 1) {
          this.selectedSession = newInfo.sessions.find((s) => s.id == diff[0]);
        } else if (diff.length == 0) {
          // nothing changed, happens during the initial phase when
          // both have no sessions
          this.selectedSession = null;
        } else {
          console.error(
            "old and new session info differs by more than one session"
          );
          console.error(oldInfo, newInfo);
          console.error("difference:", diff);
        }
      } else if (newInfo.sessions.length < oldInfo.sessions.length) {
        // a session has been deleted -> reset selected session
        this.selectedSession = null;
      }
    },
  },
  data() {
    return {
      date: null,
      selectedSession: null,
      showSessionEditorModal: false,
      showSessionDeleteModal: false,
    };
  },
  created() {
    // needed to know the timezone
    this.queryConfiguration();
    // TODO block till we have the config
    // then when promise resolved, call
    // a function that sets the date

    // get day from URL
    const urlParams = new URLSearchParams(window.location.search);
    const urlDate = urlParams.get("date");
    if (urlDate != null) {
      const urlDateParsed = dayjs(urlDate, "YYYY-MM-DD")
        .tz(this.getTimezone)
        .startOf("day")
        .format();
      this.date = urlDateParsed;
    } else {
      this.date = dayjs().tz(this.getTimezone).startOf("day").format();
    }

    this.querySessionsForDate();
  },
};
</script>
