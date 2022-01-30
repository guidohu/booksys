<template>
  <b-card
    no-body
    class="text-left"
  >
    <b-card-header>
      <b-row>
        <b-col
          cols="4"
          class="text-right"
        >
          <b-button
            variant="outline-info"
            class="btn-xs"
            @click="prevMonth"
          >
            <b-icon-arrow-left-short />
          </b-button>
        </b-col>
        <b-col
          cols="4"
          class="text-center"
        >
          {{ monthString }}
        </b-col>
        <b-col
          cols="4"
          class="text-left"
        >
          <b-button
            variant="outline-info"
            class="btn-xs"
            @click="nextMonth"
          >
            <b-icon-arrow-right-short />
          </b-button>
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body ref="calendarBody">
      <table>
        <tr>
          <th class="table-title-text">
            Mon
          </th>
          <th class="table-title-text">
            Tue
          </th>
          <th class="table-title-text">
            Wed
          </th>
          <th class="table-title-text">
            Thu
          </th>
          <th class="table-title-text">
            Fri
          </th>
          <th class="table-title-text">
            Sat
          </th>
          <th class="table-title-text">
            Sun
          </th>
        </tr>
        <tr
          v-for="i in Array.from(Array(6).keys())"
          :key="i"
        >
          <td
            v-for="j in Array.from(Array(7).keys())"
            :key="j"
            @click="navigateTo(sessionData[i * 7 + j])"
            @mouseover="mouseOver(sessionData[i * 7 + j])"
          >
            <div :class="getCalendarDayBoxClass(sessionData[i * 7 + j])">
              <BooksysPie
                :session-data="sessionData[i * 7 + j]"
                :properties="properties"
                :pie-id="i * 7 + j"
              />
              <div class="day-number">
                {{ getDay(sessionData[i * 7 + j].window_start) }}
              </div>
            </div>
          </td>
        </tr>
      </table>
    </b-card-body>
  </b-card>
</template>

<script>
import BooksysPie from "./Pie.vue";
import * as dayjs from "dayjs";
import { BooksysBrowser } from "@/libs/browser";
import {
  BCard,
  BCardHeader,
  BCardBody,
  BRow,
  BCol,
  BButton,
  BIconArrowLeftShort,
  BIconArrowRightShort,
} from "bootstrap-vue";

export default {
  name: "SessionMonthCard",
  components: {
    BooksysPie,
    BCard,
    BCardHeader,
    BCardBody,
    BRow,
    BCol,
    BButton,
    BIconArrowLeftShort,
    BIconArrowRightShort,
  },
  props: ["sessionData", "month"],
  data() {
    return {
      properties: {
        containerWidth: 84,
        containerHeight: 64,
        circleX: 45,
        circleY: 34,
        circleRadius: 23,
        animate: false,
        labels: false,
      },
    };
  },
  computed: {
    monthString: function () {
      return dayjs(this.month).format("MMMM YYYY");
    },
  },
  mounted() {
    const totalWidth = this.$refs.calendarBody.clientWidth - 61;
    const cardWidth = totalWidth / 7;
    let cardHeight = cardWidth * 0.75;
    if (this.isMobile()) {
      cardHeight = cardWidth * 1.5;
    }
    const radius = (Math.min(cardWidth, cardHeight) / 2) * 0.7;
    this.properties = {
      containerHeight: cardHeight,
      containerWidth: cardWidth,
      circleX: (cardWidth - 3) / 2,
      circleY: cardHeight / 2,
      circleRadius: radius,
      animate: false,
      labels: false,
    };
  },
  methods: {
    prevMonth: function () {
      this.$emit("prevMonth");
    },
    nextMonth: function () {
      this.$emit("nextMonth");
    },
    getDay: function (isoTime) {
      return dayjs(isoTime).format("DD");
    },
    isMobile: function () {
      return BooksysBrowser.isMobile();
    },
    getCalendarDayBoxClass: function (daySessionData) {
      const boxMonth = dayjs(daySessionData.window_start)
        .startOf("month")
        .format();
      if (this.month == boxMonth) {
        return "calendar-day-box";
      }
      return "different-month";
    },
    navigateTo: function (daySessionData) {
      window.location.href =
        "/today?date=" +
        dayjs(daySessionData.window_start).format("YYYY-MM-DD");
    },
    mouseOver: function (daySessionData) {
      this.$emit("mouseOverHandler", daySessionData);
    },
  },
};
</script>

<style scoped>
.calendar-day-box {
  border-radius: 0.1em;
  border: 0.01em;
  border-color: black;
  background: #eceaea;
  margin-right: 0.01em;
  margin-left: 0.01em;
  margin-top: 0.01em;
  margin-bottom: 0.01em;
  position: relative;
  cursor: pointer;
}

.calendar-day-box:hover {
  background: #c9c9c9;
}

.different-month {
  border-radius: 0.1em;
  border: 0.01em;
  border-color: black;
  background: #fafafa;
  margin-right: 0.01em;
  margin-left: 0.01em;
  margin-top: 0.01em;
  margin-bottom: 0.01em;
  position: relative;
  cursor: pointer;
}

.day-number {
  position: absolute;
  bottom: 0;
  left: 0.1rem;
  font-size: 0.7rem;
}

.table-title-text {
  font-size: 0.7rem;
}

.btn-xs {
  padding: 0.2rem 0.4rem;
  font-size: 0.6rem;
  line-height: 1;
  border-radius: 0.2rem;
}
</style>
