<template>
  <div>
    <SessionEditorModal
      :defaultValues="selectedSession"
      :visible.sync="showSessionEditorModal"
    />
    <SessionDeleteModal
      :session="selectedSession"
      @sessionDeletedHandler="sessionDeletedHandler"
      :visible.sync="showSessionDeleteModal"
    />
    <div v-if="isDesktop" class="display">
      <main-title title-name="Select Session" />
      <b-row class="ml-1 mr-1">
        <b-col cols="8">
          <SessionDayCard
            v-if="getSessions != null"
            :sessionData="getSessions"
            :selectedSession="selectedSession"
            :isMobile="isMobile"
            :timezone="getTimezone"
            disableDayBrowsing="true"
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
              />
            </b-col>
          </b-row>
          <b-row class="mt-1">
            <b-col cols="12">
              <ConditionInfoCard
                v-if="sunrise != null"
                :sunrise="sunrise"
                :sunset="sunset"
              />
            </b-col>
          </b-row>
        </b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button class="mr-1" variant="outline-light" to="/calendar">
          <b-icon-calendar3 />
          CALENDAR
        </b-button>
        <b-button variant="outline-light" to="/dashboard">
          <b-icon-house />
          HOME
        </b-button>
      </div>
    </div>
    <div v-else>
      <NavbarMobile title="Select Session" />
      <SessionDayCard
        v-if="getSessions != null"
        class="mb-1"
        :sessionData="getSessions"
        :selectedSession="selectedSession"
        :isMobile="isMobile"
        :timezone="getTimezone"
        disableDayBrowsing="true"
        @selectSessionHandler="selectSlot"
      />
      <SessionDetailsCard
        class="mb-1"
        :date="date"
        :session="selectedSession"
        @createSessionHandler="showCreateSession"
        @editSessionHandler="showCreateSession"
        @deleteSessionHandler="showDeleteSession"
      />
      <ConditionInfoCard
        v-if="sunrise != null"
        :sunrise="sunrise"
        :sunset="sunset"
      />
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import ConditionInfoCard from "@/components/ConditionInfoCard";
import SessionDayCard from "@/components/SessionDayCard";
import SessionDetailsCard from "@/components/SessionDetailsCard";
import MainTitle from "@/components/MainTitle";
import Session from "@/dataTypes/session";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";
import difference from "lodash/difference";
import { BRow, BCol, BButton, BIconCalendar3, BIconHouse } from "bootstrap-vue";

const SessionEditorModal = () =>
  import(
    /* webpackChunkName: "session-editor-modal" */ "@/components/SessionEditorModal"
  );
const SessionDeleteModal = () =>
  import(
    /* webpackChunkName: "session-delete-modal" */ "@/components/SessionDeleteModal"
  );

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

export default {
  name: "Ride",
  components: {
    NavbarMobile,
    MainTitle,
    ConditionInfoCard,
    SessionDayCard,
    SessionDetailsCard,
    SessionEditorModal,
    SessionDeleteModal,
    BRow,
    BCol,
    BButton,
    BIconCalendar3,
    BIconHouse,
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
    ...mapGetters("stopwatch", ["getIsRunning", "getSessionId"]),
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
      this.selectedSession = null;
    },
    showCreateSession: function () {
      this.showSessionEditorModal = true;
    },
    showDeleteSession: function () {
      this.showSessionDeleteModal = true;
    },
  },
  watch: {
    getSessions: function (newInfo, oldInfo) {
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
          oldInfo == null || oldInfo.session == null
            ? []
            : oldInfo.sessions.map((s) => s.id);
        const diff = difference(newIds, oldIds);
        if (diff.length == 1) {
          this.selectedSession = newInfo.sessions.find((s) => s.id == diff[0]);
        } else {
          console.error(
            "old and new session info differs by more than one session"
          );
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

    // check if the watch is running currently -> then go to watch
    if (this.getIsRunning == true) {
      this.$router.push("/watch?sessionId=" + this.getSessionId);
    }
  },
};
</script>
