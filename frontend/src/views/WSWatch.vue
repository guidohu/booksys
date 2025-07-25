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

<script setup>
import { ref, computed, onMounted } from "vue";
import { useStore, mapGetters, mapActions } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import StopWatchCard from "booksys/components/StopWatchCard.vue";
import SessionHeatListCard from "booksys/components/SessionHeatListCard.vue";
import ConditionInfoCard from "booksys/components/ConditionInfoCard.vue";
import ShowForMobile from "../components/bricks/ShowForMobile.vue";
import ShowForDesktop from "../components/bricks/ShowForDesktop.vue";
import SubpageContainer from "../components/bricks/SubpageContainer.vue";

const store = useStore();

const sessionId = ref(null);
const errors = ref([]);

const getSessionConditionInfo = computed(
  mapGetters("sessions", ["getSessionConditionInfo"]).getSessionConditionInfo.bind({
    $store: store,
  })
);

const querySessionMetadata = mapActions("sessions", [
  "querySessionMetadata",
]).querySessionMetadata.bind({ $store: store });

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const id = Number(urlParams.get("sessionId"));
  if (id == null || !Number(id) > 0) {
    errors.value = ["No session parameter given in URL"];
  } else {
    sessionId.value = id;
  }

  querySessionMetadata(sessionId.value)
    .then(() => console.log("queried session metadata"))
    .catch((err) => (errors.value = err));
});
</script>

<style>
.side-bar-heats-component {
  max-height: 300px;
}
</style>