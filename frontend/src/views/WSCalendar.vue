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

<script setup>
import { ref, computed, onMounted } from "vue";
import { useStore, mapGetters } from "vuex";
import ConditionInfoCard from "booksys/components/ConditionInfoCard.vue";
import SessionMonthCard from "booksys/components/SessionMonthCard.vue";
import SessionsOverview from "booksys/components/SessionsOverview.vue";
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

const store = useStore();

const month = ref(null);
const errors = ref([]);
const sunrise = ref(null);
const sunset = ref(null);
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
  sunrise.value = dayjs(day.sunrise).format("X");
  sunset.value = dayjs(day.sunset).format("X");

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
