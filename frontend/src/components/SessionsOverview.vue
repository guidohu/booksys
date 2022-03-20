<template>
  <sectioned-card-module title="Sessions">
    <template v-slot:body>
      <div v-if="sessions == null || sessions.sessions == null" class="row">
        <div class="col-12">No day selected</div>
      </div>
      <div
        v-if="
          sessions != null &&
          sessions.sessions != null &&
          sessions.sessions.length == 0
        "
        class="row"
      >
        <div class="col-12">No sessions for this day.</div>
      </div>
      <div
        v-if="
          sessions != null &&
          sessions.sessions != null &&
          sessions.sessions.length > 0
        "
      >
        <div
          v-for="session in sessions.sessions.slice(0, 2)"
          :key="session.id"
          class="row mt-2 small-text"
        >
          <div class="col-12">
            <div class="row">
              <div class="col-5">
                <i class="bi bi-book" />
                Title
              </div>
              <div class="col-7">
                {{ getTitle(session) }}
              </div>
            </div>
            <div class="row">
              <div class="col-5">
                <i class="bi bi-clock" />
                Time
              </div>
              <div class="col-7">
                {{ getTime(session) }}
              </div>
            </div>
            <div class="row">
              <div class="col-5">
                <i class="bi bi-person" />
                Riders
              </div>
              <div
                v-if="session.rider == null || session.riders.length == 0"
                class="col-7"
              >
                -
              </div>
              <div
                v-if="session.rider != null && session.riders.length > 0"
                class="col-7"
              >
                <div
                  v-for="rider in session.riders.slice(0, 3)"
                  :key="rider.id"
                  class="row"
                >
                  <div class="col-12 text-truncate">
                    {{ rider.name }}
                  </div>
                </div>
                <div v-if="sessions.riders.length > 3">
                  <div class="col-12 text-truncate">...</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div
        v-if="
          sessions != null &&
          sessions.sessions != null &&
          sessions.sessions.length > 2
        "
        class="row small-text"
      >
        <div class="col-12">click for more details...</div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script>
import * as dayjs from "dayjs";
import SectionedCardModule from "./bricks/SectionedCardModule.vue";

export default {
  name: "SessionsOverview",
  components: {
    SectionedCardModule,
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
