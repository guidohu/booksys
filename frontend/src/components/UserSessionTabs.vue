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
      <div class="tab-pane tab-limited-height active" id="upcoming-sessions">
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

<script setup>
import { computed } from "vue";
import { useStore } from "vuex";
import UserSessionsTable from "booksys/components/UserSessionsTable.vue";

const store = useStore();

const userSchedule = computed(() => store.getters["user/userSchedule"]);

const upcomingSessions = computed(() => {
  if (userSchedule.value.sessions != null) {
    return userSchedule.value.sessions;
  }
  return [];
});

const pastSessions = computed(() => {
  if (userSchedule.value.sessions_old != null) {
    return userSchedule.value.sessions_old;
  }
  return [];
});

const queryUserSchedule = () => store.dispatch("user/queryUserSchedule");
const cancelSession = (sessionId) =>
  store.dispatch("user/cancelSession", sessionId);

function cancelSessionHandler(value) {
  cancelSession(value);
}

queryUserSchedule().catch((error) => {
  console.error(error);
});
</script>

<style scoped>

.tab-limited-height {
  max-height: 400px;
  height: 400px;
  overflow-y: scroll;
}

@media (max-width: 992px) {
  .tab-limited-height {
    height: 100%;
    overflow-y: hide;
  }
}

.nav-link {
  color: #bdbdbd;
}

a {
  color: #bdbdbd;
}
</style>
