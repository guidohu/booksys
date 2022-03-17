<template>
  <subpage-container title="Calendar">
    <show-for-desktop>
      <div class="row mx-1">
        <div class="col-8">
          <session-month-card
            v-if="getSessionsCalendar != null"
            :session-data="getSessionsCalendar"
            :month="month"
            @prevMonth="prevMonth"
            @nextMonth="nextMonth"
            @mouseOverHandler="mouseOverDayHandler"
          />
        </div>
        <div class="col-4">
          <sessions-overview :sessions="sessionsOverview" />
          <condition-info-card
            class="mt-2"
            :sunrise="sunrise"
            :sunset="sunset"
          />
        </div>
      </div>
    </show-for-desktop>
    <show-for-mobile>
      <session-month-card
        class="mx-1 mt-1"
        v-if="getSessionsCalendar != null"
        :session-data="getSessionsCalendar"
        :month="month"
        @prevMonth="prevMonth"
        @nextMonth="nextMonth"
      />
    </show-for-mobile>
  </subpage-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import ConditionInfoCard from "@/components/ConditionInfoCard";
import SessionMonthCard from "@/components/SessionMonthCard";
import SessionsOverview from "@/components/SessionsOverview";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";
import * as dayjsCustomParseFormat from "dayjs/plugin/customParseFormat";
import ShowForMobile from "../components/bricks/ShowForMobile.vue";
import ShowForDesktop from "../components/bricks/ShowForDesktop.vue";
import SubpageContainer from "../components/bricks/SubpageContainer.vue";

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);
dayjs.extend(dayjsCustomParseFormat);

export default {
  name: "WSCalendar",
  components: {
    ConditionInfoCard,
    SessionMonthCard,
    SessionsOverview,
    ShowForMobile,
    ShowForDesktop,
    SubpageContainer,
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
      console.log("month changed to", this.month);
    },
    nextMonth: function () {
      this.month = dayjs(this.month).add(1, "month").format();
      this.querySessionsForMonth();
      console.log("month changed to", this.month);
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
