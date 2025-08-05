<template>
  <sectioned-card-module>
    <template v-slot:header>
      <div class="row">
        <div class="col-4 text-end">
          <button
            v-if="!disableDayBrowsing"
            type="button"
            class="btn btn-outline-info btn-xs"
            @click="prevDay"
          >
            <i class="bi bi-arrow-left-short" />
          </button>
        </div>
        <div class="col-4 text-center">
          {{ dateString }}
        </div>
        <div class="col-4 text-start">
          <button
            v-if="!disableDayBrowsing"
            type="button"
            class="btn btn-outline-info btn-xs"
            @click="nextDay"
          >
            <i class="bi bi-arrow-right-short" />
          </button>
        </div>
      </div>
    </template>
    <template v-slot:body>
      <div class="row">
        <booksys-pie
          :session-data="sessionData"
          :selected-session="selectedSession"
          :properties="properties"
          @selectHandler="selectSession"
        />
      </div>
      <div
        v-if="
          isToday &&
          selectedSession != null &&
          selectedSession.id != null &&
          selectedSession.riders.length > 0
        "
        class="row text-center"
      >
        <div class="col-12 text-center">
          <button
            type="button"
            class="btn btn-outline-success"
            @click="navigateSessionStart"
          >
            Start Session
          </button>
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script setup>
import { computed } from "vue";
import { useStore } from "vuex";
import BooksysPie from "./BooksysPie.vue";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";
import SectionedCardModule from "booksys/components/bricks/SectionedCardModule.vue";

dayjs.extend(utc);
dayjs.extend(timezone);

const props = defineProps([
  "isMobile",
  "sessionData",
  "selectedSession",
  "timezone",
  "disableDayBrowsing",
]);

const emit = defineEmits(["prevDay", "nextDay", "selectSessionHandler"]);

const store = useStore();

const getTimezone = computed(() => store.getters["configuration/getTimezone"]);

const properties = computed(() => {
  const baseProps = props.isMobile
    ? {
        containerHeight: 300,
        containerWidth: 350,
        circleX: 175,
        circleY: 150,
        circleRadius: 90,
        animation: false,
        labels: true,
      }
    : {
        containerWidth: 700,
        containerHeight: 350,
        circleX: 300,
        circleY: 170,
        circleRadius: 100,
        animate: true,
        labels: true,
      };
  baseProps.timezone = props.timezone || "UTC";
  return baseProps;
});

const dateString = computed(() => {
  console.log("sessionData", props.sessionData);
  return dayjs(props.sessionData.window_start).format("dddd DD.MM.YYYY");
});

const isToday = computed(() => {
  console.log(getTimezone.value);
  const today = dayjs().tz(getTimezone.value).startOf("day").format();
  const sessionDay = dayjs(props.sessionData.window_start)
    .tz(getTimezone.value)
    .startOf("day")
    .format();
  return today === sessionDay;
});

function prevDay() {
  emit("prevDay");
}

function nextDay() {
  emit("nextDay");
}

function selectSession(slot) {
  emit("selectSessionHandler", slot);
}

function navigateSessionStart() {
  window.location.href = "/watch?sessionId=" + props.selectedSession.id;
}
</script>

<style scoped>
.btn-xs {
  padding: 0.2rem 0.4rem;
  font-size: 0.6rem;
  line-height: 1;
  border-radius: 0.2rem;
}
</style>
