<template>
  <subpage-container title="Stop Watch">
    <show-for-desktop>
      <div class="row mx-1">
        <div class="col-8">
          <warning-box v-if="errors.length > 0" :errors="errors" />
          <stop-watch-card :session-id="sessionId" />
        </div>
        <div class="col-4">
          <session-heat-list-card
            :session-id="sessionId"
            class="side-bar-heats-component"
          />
          <condition-info-card
            v-if="getSessionConditionInfo != null"
            class="mt-2"
            :sunrise="getSessionConditionInfo.sunrise"
            :sunset="getSessionConditionInfo.sunset"
          />
        </div>
      </div>
    </show-for-desktop>
    <show-for-mobile>
      <warrning-box v-if="errors.length > 0" :errors="errors" />
      <stop-watch-card :session-id="sessionId" />
      <session-heat-list-card :session-id="sessionId" class="mt-2" />
      <condition-info-card
        v-if="getSessionConditionInfo != null"
        class="mt-2"
        :sunrise="getSessionConditionInfo.sunrise"
        :sunset="getSessionConditionInfo.sunset"
      />
    </show-for-mobile>
    <template v-slot:bottom>
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

<script>
import { mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
import StopWatchCard from "@/components/StopWatchCard";
import SessionHeatListCard from "@/components/SessionHeatListCard";
import ConditionInfoCard from "@/components/ConditionInfoCard";
import ShowForMobile from '../components/bricks/ShowForMobile.vue';
import ShowForDesktop from '../components/bricks/ShowForDesktop.vue';
import SubpageContainer from '../components/bricks/SubpageContainer.vue';

export default {
  name: "WSWatch",
  components: {
    WarningBox,
    StopWatchCard,
    SessionHeatListCard,
    ConditionInfoCard,
    ShowForMobile,
    ShowForDesktop,
    SubpageContainer,
  },
  data() {
    return {
      sessionId: null,
      errors: [],
    };
  },
  computed: {
    ...mapGetters("sessions", ["getSessionConditionInfo"]),
  },
  created() {
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = Number(urlParams.get("sessionId"));
    if (sessionId == null || !Number(sessionId) > 0) {
      this.errors = ["No session parameter given in URL"];
    } else {
      this.sessionId = sessionId;
    }
  },
};
</script>

<style>
.side-bar-heats-component {
  max-height: 300px;
}
</style>
