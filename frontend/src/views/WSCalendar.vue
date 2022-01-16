<template>
  <div>
    <div v-if="isDesktop" class="display">
      <main-title title-name="Calendar" />
      <b-row class="ml-1 mr-1">
        <b-col cols="8">
          <SessionMonthCard
            v-if="getSessionsCalendar != null"
            :sessionData="getSessionsCalendar"
            :month="month"
            @prevMonth="prevMonth"
            @nextMonth="nextMonth"
            @mouseOverHandler="mouseOverDayHandler"
          />
        </b-col>
        <b-col cols="4">
          <b-row>
            <b-col cols="12">
              <SessionsOverview :sessions="sessionsOverview" />
            </b-col>
          </b-row>
          <b-row class="mt-2">
            <b-col cols="12">
              <ConditionInfoCard :sunrise="sunrise" :sunset="sunset" />
            </b-col>
          </b-row>
        </b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button variant="outline-light" to="/dashboard">
          <b-icon-house />
          HOME
        </b-button>
      </div>
    </div>
    <div v-else>
      <NavbarMobile title="Calendar" />
      <SessionMonthCard
        v-if="getSessionsCalendar != null"
        :sessionData="getSessionsCalendar"
        :month="month"
        @prevMonth="prevMonth"
        @nextMonth="nextMonth"
      />
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import ConditionInfoCard from "@/components/ConditionInfoCard";
import SessionMonthCard from "@/components/SessionMonthCard";
import SessionsOverview from "@/components/SessionsOverview";
import MainTitle from "@/components/MainTitle";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";
import * as dayjsCustomParseFormat from "dayjs/plugin/customParseFormat";

import { BRow, BCol, BButton, BIconHouse } from "bootstrap-vue";

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);
dayjs.extend(dayjsCustomParseFormat);

export default {
  name: "WSCalendar",
  components: {
    NavbarMobile,
    MainTitle,
    ConditionInfoCard,
    SessionMonthCard,
    SessionsOverview,
    BRow,
    BCol,
    BButton,
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
    ...mapGetters("sessions", ["getSessionsCalendar"]),
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    ...mapActions("sessions", ["querySessionsCalendar"]),
    prevMonth: function () {
      this.month = dayjs(this.month).add(-1, "month").format();
      this.querySessionsForMonth();
    },
    nextMonth: function () {
      this.month = dayjs(this.month).add(1, "month").format();
      this.querySessionsForMonth();
    },
    querySessionsForMonth: function () {
      // query get_booking_day
      this.querySessionsCalendar(this.month)
        .then(() => {
          console.log("Calendar ready");
        })
        .catch((errors) => (this.errors = errors));
    },
    mouseOverDayHandler: function (day) {
      this.sunrise = day.sunrise;
      this.sunset = day.sunset;

      this.sessionsOverview = day;
    },
  },
  data() {
    return {
      month: null,
      errors: [],
      sunrise: null,
      sunset: null,
      sessionsOverview: null,
      // selectedSession: null
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
        .startOf("month")
        .format();
      this.month = urlDateParsed;
    } else {
      this.month = dayjs().tz(this.getTimezone).startOf("month").format();
      console.log("calculated month to be:", this.month);
    }

    this.querySessionsForMonth();
  },
};
</script>
