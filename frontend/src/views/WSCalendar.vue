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
            :sunrise="sunriseUnix"
            :sunset="sunsetUnix"
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

<script setup>
import { ref, computed, onMounted } from "vue";
import { useStore, mapGetters } from "vuex";
import ConditionInfoCard from "booksys/components/ConditionInfoCard.vue";
import SessionMonthCard from "booksys/components/SessionMonthCard.vue";
import SessionsOverview from "booksys/components/SessionsOverview.vue";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";
import customParseFormat from "dayjs/plugin/customParseFormat";
import ShowForMobile from "../components/bricks/ShowForMobile.vue";
import ShowForDesktop from "../components/bricks/ShowForDesktop.vue";
import SubpageContainer from "../components/bricks/SubpageContainer.vue";

dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(customParseFormat);

const store = useStore();

const month = ref(null);
const errors = ref([]);
const sunriseUnix = ref(0);
const sunsetUnix = ref(0);
const sessionsOverview = ref(null);

const getTimezone = computed(
  mapGetters("configuration", ["getTimezone"]).getTimezone.bind({ $store: store })
);
const getSessionsCalendar = computed(
  mapGetters("sessions", ["getSessionsCalendar"]).getSessionsCalendar.bind({
    $store: store,
  })
);

const prevMonth = () => {
  month.value = dayjs(month.value).add(-1, "month").format();
  querySessionsForMonth();
  console.log("month changed to", month.value);
};

const nextMonth = () => {
  month.value = dayjs(month.value).add(1, "month").format();
  querySessionsForMonth();
  console.log("month changed to", month.value);
};

const querySessionsForMonth = () => {
  // query get_booking_day
  store
    .dispatch("sessions/querySessionsCalendar", month.value)
    .then(() => {
      console.log("Calendar ready");
    })
    .catch((err) => (errors.value = err));
};

const mouseOverDayHandler = (day) => {
  const sunrise = parseInt(dayjs(day.sunrise).format("X"));
  if (Number.isNaN(sunrise)) {
    sunriseUnix.value = 0;
  } else {
    sunriseUnix.value = sunrise;
  }
  const sunset = parseInt(dayjs(day.sunset).format("X"));
  if (Number.isNaN(sunset)) {
    sunsetUnix.value = 0;
  } else {
    sunsetUnix.value = sunset;
  }
  sessionsOverview.value = day;
};

onMounted(() => {
  // needed to know the timezone
  store.dispatch("configuration/queryConfiguration");
  // TODO block till we have the config
  // then when promise resolved, call
  // a function that sets the date

  // get day from URL
  const urlParams = new URLSearchParams(window.location.search);
  const urlDate = urlParams.get("date");
  if (urlDate != null) {
    const urlDateParsed = dayjs(urlDate, "YYYY-MM-DD")
      .tz(getTimezone.value)
      .startOf("month")
      .format();
    month.value = urlDateParsed;
  } else {
    month.value = dayjs().tz(getTimezone.value).startOf("month").format();
    console.log("calculated month to be:", month.value);
  }

  querySessionsForMonth();
});
</script>
