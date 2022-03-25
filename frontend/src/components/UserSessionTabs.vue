<template>
  <div>
    <ul class="nav nav-tabs">
      <li class="nav-item">
        <a
          class="nav-link active"
          data-bs-toggle="tab"
          href="#upcoming-sessions"
        >
          <i class="bi bi-calendar"></i>
          Upcoming Sessions
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-bs-toggle="tab" href="#past-sessions">
          <i class="bi bi-calendar-check"></i>
          Past Sessions
        </a>
      </li>
    </ul>

    <div class="tab-content">
      <div :class="'tab-pane ' + tabClass + ' active'" id="upcoming-sessions">
        <user-sessions-table
          v-if="upcomingSessions.length > 0"
          :user-sessions="upcomingSessions"
          :show-cancel="true"
          @cancel="cancelSessionHandler"
        />
        <div
          v-if="upcomingSessions.length == 0"
          class="alert alert-info text-center m-2"
          role="alert"
        >
          You do not have any upcoming sessions scheduled.
        </div>
      </div>
      <div :class="'tab-pane ' + tabClass" id="past-sessions">
        <user-sessions-table
          v-if="pastSessions.length > 0"
          :user-sessions="pastSessions"
        />
        <div
          v-if="pastSessions.length == 0"
          class="alert alert-info text-center m-2"
          role="alert"
        >
          You did not take part in a session yet.
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import UserSessionsTable from "@/components/UserSessionsTable.vue";

export default {
  name: "UserSessionTabs",
  components: {
    UserSessionsTable,
  },
  computed: {
    ...mapGetters("user", ["userSchedule"]),
    upcomingSessions: function () {
      if (this.userSchedule.sessions != null) {
        return this.userSchedule.sessions;
      }
      return [];
    },
    pastSessions: function () {
      if (this.userSchedule.sessions_old != null) {
        return this.userSchedule.sessions_old;
      }
      return [];
    },
    tabClass: function () {
      if (!BooksysBrowser.isMobile()) {
        return "tab-limited-height";
      } else {
        return "";
      }
    },
  },
  methods: {
    ...mapActions("user", ["queryUserSchedule", "cancelSession"]),
    cancelSessionHandler(value) {
      console.log("Cancel session with ID:", value, "in parent");
      this.cancelSession(value);
    },
  },
  created() {
    this.queryUserSchedule()
      .then(() => {
        console.log("Retrieved userSchedule");
        console.debug(this.userSchedule);
      })
      .catch((error) => {
        console.error(error);
      });
  },
};
</script>

<style scoped>
.tab-limited-height {
  max-height: 400px;
  height: 400px;
  overflow-y: scroll;
}

.nav-link {
  color: #bdbdbd;
}

a {
  color: #bdbdbd;
}
</style>
