<template>
  <subpage-container title="Book your Session">
    <session-editor-modal
      v-model:visible="showSessionEditorModal"
      :default-values="selectedSession"
    />
    <session-delete-modal
      v-model:visible="showSessionDeleteModal"
      :session="selectedSession"
      @sessionDeletedHandler="sessionDeletedHandler"
    />
    <show-for-desktop>
      <div class="row mx-1">
        <div class="col-8">
          <session-day-card
            v-if="getSessions != null"
            :session-data="getSessions"
            :selected-session="selectedSession"
            :is-mobile="false"
            :timezone="getTimezone"
            @prevDay="prevDay"
            @nextDay="nextDay"
            @selectSessionHandler="selectSlot"
          />
        </div>
        <div class="col-4">
          <session-details-card
            :date="date"
            :session="selectedSession"
            @createSessionHandler="showCreateSession"
            @editSessionHandler="showCreateSession"
            @deleteSessionHandler="showDeleteSession"
          />
          <condition-info-card
            v-if="sunrise != null"
            class="mt-2"
            :sunrise="sunrise"
            :sunset="sunset"
          />
        </div>
      </div>
    </show-for-desktop>
    <show-for-mobile>
      <session-day-card
        class="mx-1 mt-1"
        v-if="getSessions != null"
        :session-data="getSessions"
        :selected-session="selectedSession"
        :is-mobile="true"
        :timezone="getTimezone"
        @prevDay="prevDay"
        @nextDay="nextDay"
        @selectSessionHandler="selectSlot"
      />
      <session-details-card
        class="mx-1 mt-2"
        :date="date"
        :session="selectedSession"
        @createSessionHandler="showCreateSession"
        @editSessionHandler="showCreateSession"
        @deleteSessionHandler="showDeleteSession"
      />
      <condition-info-card
        class="mx-1 mt-2"
        v-if="sunrise != null"
        :sunrise="sunrise"
        :sunset="sunset"
      />
    </show-for-mobile>
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/calendar">
        <i class="bi bi-calendar"></i>
        CALENDAR
      </router-link>
      <router-link
        tag="button"
        class="btn btn-outline-light ms-1"
        to="/dashboard"
      >
        <i class="bi bi-house"></i>
        HOME
      </router-link>
    </template>
  </subpage-container>
</template>

<script setup>
import { defineAsyncComponent, ref, computed, onMounted, watch } from "vue";
import { useStore, mapGetters, mapActions } from "vuex";
import ConditionInfoCard from "booksys/components/ConditionInfoCard.vue";
import SessionDayCard from "booksys/components/SessionDayCard.vue";
import SessionDetailsCard from "booksys/components/SessionDetailsCard.vue";
import Session from "booksys/dataTypes/session";
import * as dayjs from "dayjs";
import * as dayjsCustomParseFormat from "dayjs/plugin/customParseFormat";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";

import difference from "lodash/difference";
import SubpageContainer from "../components/bricks/SubpageContainer.vue";
import ShowForMobile from "../components/bricks/ShowForMobile.vue";
import ShowForDesktop from "../components/bricks/ShowForDesktop.vue";

const SessionEditorModal = defineAsyncComponent(() =>
  import("booksys/components/SessionEditorModal.vue")
);
const SessionDeleteModal = defineAsyncComponent(() =>
  import("booksys/components/SessionDeleteModal.vue")
);

dayjs.extend(dayjsCustomParseFormat);
dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

const store = useStore();

const date = ref(null);
const selectedSession = ref(null);
const showSessionEditorModal = ref(false);
const showSessionDeleteModal = ref(false);

const getTimezone = computed(
  mapGetters("configuration", ["getTimezone"]).getTimezone.bind({ $store: store })
);
const getSessions = computed(
  mapGetters("sessions", ["getSessions"]).getSessions.bind({ $store: store })
);

const sunrise = computed(() => {
  const sessions = getSessions.value;
  if (sessions == null) {
    return null;
  } else {
    return dayjs(sessions.sunrise).format("X");
  }
});

const sunset = computed(() => {
  const sessions = getSessions.value;
  if (sessions == null) {
    return null;
  } else {
    return dayjs(sessions.sunset).format("X");
  }
});

const queryConfiguration = mapActions("configuration", [
  "queryConfiguration",
]).queryConfiguration.bind({ $store: store });
const querySessions = mapActions("sessions", ["querySessions"]).querySessions.bind({
  $store: store,
});

const prevDay = () => {
  date.value = dayjs(date.value).add(-1, "days");
  selectedSession.value = null;
  querySessionsForDate();
};

const nextDay = () => {
  date.value = dayjs(date.value).add(1, "days");
  selectedSession.value = null;
  querySessionsForDate();
};

const querySessionsForDate = () => {
  const dateStart = dayjs(date.value).startOf("day").format();
  const dateEnd = dayjs(date.value).endOf("day").format();

  // query get_booking_day
  querySessions({
    start: dateStart,
    end: dateEnd,
  });
};

const selectSlot = (session) => {
  console.log("selectedSlot", session);
  console.log(dayjs(session.start).format(), dayjs(session.end).format());
  const sessionWithSelectedId = getSessions.value.sessions.find(
    (s) => session.id == s.id
  );
  if (sessionWithSelectedId == null) {
    selectedSession.value = new Session(
      session.id,
      null,
      null,
      session.start,
      session.end
    );
  } else {
    selectedSession.value = sessionWithSelectedId;
  }
};

const sessionDeletedHandler = () => {
  console.log("sessionDeletedHandler");
  selectedSession.value = null;
};

const showCreateSession = () => {
  console.log("createSession clicked");
  console.log("for selectedSession:", selectedSession.value);
  console.log("show session editor");
  showSessionEditorModal.value = true;
};

const showDeleteSession = () => {
  console.log("showDeleteSession");
  showSessionDeleteModal.value = true;
};

watch(getSessions, (newInfo, oldInfo) => {
  if (newInfo == null || newInfo.sessions == null) {
    selectedSession.value = null;
    return;
  }

  // in case we get an update affecting the sessions
  // we will update our selected session too
  if (
    selectedSession.value != null &&
    selectedSession.value.id != null &&
    newInfo.sessions.map((s) => s.id).includes(selectedSession.value.id)
  ) {
    selectedSession.value = newInfo.sessions.filter(
      (s) => s.id == selectedSession.value.id
    )[0];
  }

  // in case a new session has been created, we select it
  if (
    newInfo.sessions != null &&
    (oldInfo == null ||
      oldInfo.sessions == null ||
      newInfo.sessions.length > oldInfo.sessions.length)
  ) {
    // find the session that is new
    const newIds = newInfo.sessions.map((s) => s.id);
    const oldIds =
      oldInfo != null && oldInfo.sessions != null
        ? oldInfo.sessions.map((s) => s.id)
        : [];
    const diff = difference(newIds, oldIds);
    if (diff.length == 1) {
      selectedSession.value = newInfo.sessions.find((s) => s.id == diff[0]);
    } else if (diff.length == 0) {
      // nothing changed, happens during the initial phase when
      // both have no sessions
      selectedSession.value = null;
    } else {
      console.error(
        "old and new session info differs by more than one session"
      );
      console.error(oldInfo, newInfo);
      console.error("difference:", diff);
    }
  } else if (newInfo.sessions.length < oldInfo.sessions.length) {
    // a session has been deleted -> reset selected session
    selectedSession.value = null;
  }
});

onMounted(() => {
  // needed to know the timezone
  queryConfiguration();
  // TODO block till we have the config
  // then when promise resolved, call
  // a function that sets the date

  // get day from URL
  const urlParams = new URLSearchParams(window.location.search);
  const urlDate = urlParams.get("date");
  if (urlDate != null) {
    const urlDateParsed = dayjs(urlDate, "YYYY-MM-DD")
      .tz(getTimezone.value)
      .startOf("day")
      .format();
    date.value = urlDateParsed;
  } else {
    date.value = dayjs().tz(getTimezone.value).startOf("day").format();
  }

  querySessionsForDate();
});
</script>
