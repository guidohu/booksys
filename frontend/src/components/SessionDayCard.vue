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
            v-if="!disableDayBrowsing"
            variant="outline-info"
            class="btn-xs"
            @click="prevDay"
          >
            <b-icon-arrow-left-short />
          </b-button>
        </b-col>
        <b-col
          cols="4"
          class="text-center"
        >
          {{ dateString }}
        </b-col>
        <b-col
          cols="4"
          class="text-left"
        >
          <b-button
            v-if="!disableDayBrowsing"
            variant="outline-info"
            class="btn-xs"
            @click="nextDay"
          >
            <b-icon-arrow-right-short />
          </b-button>
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body>
      <b-row>
        <BooksysPie
          :session-data="sessionData"
          :selected-session="selectedSession"
          :properties="properties"
          @selectHandler="selectSession"
        />
      </b-row>
      <b-row
        v-if="
          isToday &&
            selectedSession != null &&
            selectedSession.id != null &&
            selectedSession.riders.length > 0
        "
        class="text-center"
      >
        <b-col
          cols="12"
          class="text-center"
        >
          <b-button
            type="button"
            variant="outline-success"
            @click="navigateSessionStart"
          >
            Start Session
          </b-button>
        </b-col>
      </b-row>
    </b-card-body>
  </b-card>
</template>

<script>
import { mapGetters } from "vuex";
import BooksysPie from "./Pie.vue";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

import {
  BCard,
  BCardHeader,
  BRow,
  BCol,
  BButton,
  BIconArrowLeftShort,
  BIconArrowRightShort,
  BCardBody,
} from "bootstrap-vue";

export default {
  name: "SessionDayCard",
  data() {
    return {
      properties: {
        containerWidth: 700,
        containerHeight: 350,
        circleX: 300,
        circleY: 170,
        circleRadius: 100,
        animate: true,
        labels: true,
      },
    };
  },
  computed: {
    ...mapGetters("configuration", ["getTimezone"]),
    dateString: function () {
      console.log("sessionData", this.sessionData);
      return dayjs(this.sessionData.window_start).format("dddd DD.MM.YYYY");
    },
    isToday: function () {
      console.log(this.getTimezone);
      const today = dayjs().tz(this.getTimezone).startOf("day").format();
      const sessionDay = dayjs(this.sessionData.window_start)
        .tz(this.getTimezone)
        .startOf("day")
        .format();
      if (today == sessionDay) {
        return true;
      }
      return false;
    },
  },
  created() {
    if (this.isMobile != null && this.isMobile == true) {
      this.properties = {
        containerHeight: 300,
        containerWidth: 350,
        circleX: 175,
        circleY: 150,
        circleRadius: 90,
        animation: false,
        labels: true,
      };
    }
    if (this.timezone != null) {
      this.properties.timezone = this.timezone;
    } else {
      this.properties.timezone = "UTC";
    }
  },
  methods: {
    prevDay: function () {
      this.$emit("prevDay");
    },
    nextDay: function () {
      this.$emit("nextDay");
    },
    selectSession: function (slot) {
      this.$emit("selectSessionHandler", slot);
    },
    navigateSessionStart: function () {
      window.location.href = "/watch?sessionId=" + this.selectedSession.id;
    },
  },
  components: {
    BooksysPie,
    BCard,
    BCardHeader,
    BRow,
    BCol,
    BButton,
    BIconArrowLeftShort,
    BIconArrowRightShort,
    BCardBody,
  },
  props: [
    "isMobile",
    "sessionData",
    "selectedSession",
    "timezone",
    "disableDayBrowsing",
  ],
};
</script>

<style scoped>
.btn-xs {
  padding: 0.2rem 0.4rem;
  font-size: 0.6rem;
  line-height: 1;
  border-radius: 0.2rem;
}
</style>
