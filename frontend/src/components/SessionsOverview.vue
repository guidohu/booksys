<template>
  <b-card
    no-body
    class="text-left"
  >
    <b-card-header>
      <b-row>
        <b-col cols="12">
          Sessions
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body>
      <div v-if="sessions != null && sessions.sessions != null">
        <b-row v-if="sessions.sessions.length == 0">
          <b-col cols="12">
            No sessions for this day.
          </b-col>
        </b-row>
        <b-row
          v-for="session in sessions.sessions.slice(0, 2)"
          :key="session.id"
          class="mt-2 small-text"
        >
          <b-col cols="12">
            <b-row>
              <b-col cols="5">
                <b-icon-book />
                Title
              </b-col>
              <b-col cols="7">
                {{ getTitle(session) }}
              </b-col>
            </b-row>
            <b-row>
              <b-col cols="5">
                <b-icon-clock />
                Time
              </b-col>
              <b-col cols="7">
                {{ getTime(session) }}
              </b-col>
            </b-row>
            <b-row>
              <b-col cols="5">
                <b-icon-person />
                Riders
              </b-col>
              <b-col
                v-if="session.riders != null && session.riders.length > 0"
                cols="7"
              >
                <b-row
                  v-for="rider in session.riders.slice(0, 3)"
                  :key="rider.id"
                >
                  <b-col
                    cols="12"
                    class="text-truncate"
                  >
                    {{ rider.name }}
                  </b-col>
                </b-row>
                <b-row v-if="session.riders.length > 3">
                  <b-col
                    cols="12"
                    class="text-truncate"
                  >
                    ...
                  </b-col>
                </b-row>
              </b-col>
              <b-col
                v-else
                cols="7"
              >
                -
              </b-col>
            </b-row>
          </b-col>
        </b-row>
        <b-row
          v-if="sessions.sessions.length > 2"
          class="mt-4 small-text"
        >
          <b-col cols="12">
            click for more details...
          </b-col>
        </b-row>
      </div>
      <div v-if="sessions == null || sessions.sessions == null">
        <b-row>
          <b-col cols="12">
            No day selected
          </b-col>
        </b-row>
      </div>
    </b-card-body>
  </b-card>
</template>

<script>
import * as dayjs from "dayjs";
import {
  BCard,
  BCardHeader,
  BCardBody,
  BRow,
  BCol,
  BIconBook,
  BIconClock,
  BIconPerson,
} from "bootstrap-vue";

export default {
  name: "SessionsOverview",
  components: {
    BCard,
    BCardHeader,
    BCardBody,
    BRow,
    BCol,
    BIconBook,
    BIconClock,
    BIconPerson,
  },
  props: ["sessions"],
  computed: {
    date: function () {
      if (this.sessions != null && this.sessions.window_start != null) {
        return dayjs(this.sessions.window_start).format("ddd DD.MM.YYYY");
      }
      return "no date selected";
    },
  },
  methods: {
    getTitle: function (session) {
      if (session.title) {
        return session.title;
      } else {
        return "-";
      }
    },
    getTime: function (session) {
      const startStr = dayjs(session.start).format("HH:mm");
      const endStr = dayjs(session.end).format("HH:mm");
      return startStr + " - " + endStr;
    },
  },
};
</script>

<style scoped>
.small-text {
  font-size: 0.7rem;
}
</style>
