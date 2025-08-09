<template>
  <div class="container text-left">
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <div v-if="getMyNautiqueEnabled" class="row">
      <label class="col-3 col-form-label">Level</label>
      <div class="col-9">
        <div class="text-end">
          {{ getMyNautiqueFuelLevel }}% {{ " " }} (
          {{
            parseInt((getMyNautiqueFuelCapacity * getMyNautiqueFuelLevel) / 100)
          }}
          L)
        </div>
        <div class="progress" style="height: 3px">
          <div
            class="progress-bar bg-info"
            role="progressbar"
            :style="fuelLevelStyle"
            aria-valuenow="25"
            aria-valuemin="0"
            aria-valuemax="100"
          ></div>
        </div>
      </div>
    </div>
    <div class="row">
      <label class="col-3 col-form-label">Consumption</label>
      <div class="col-9">
        <div class="text-end">{{ getAvgFuelConsumption }} {{ " " }} L/hour</div>
        <div class="progress" style="height: 3px">
          <div
            class="progress-bar bg-info"
            role="progressbar"
            :style="fuelConsumptionStyle"
            aria-valuenow="25"
            aria-valuemin="0"
            aria-valuemax="100"
          ></div>
        </div>
      </div>
    </div>
    <div v-if="getMyNautiqueEnabled" class="row">
      <label class="col-3 col-form-label">Time Left</label>
      <div class="col-9">
        <div class="text-end">
          {{ timeTillEmpty }}
        </div>
        <div class="progress" style="height: 3px">
          <div
            class="progress-bar bg-info"
            role="progressbar"
            :style="fuelLevelStyle"
            aria-valuenow="25"
            aria-valuemin="0"
            aria-valuemax="100"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onBeforeUnmount } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import { min } from "lodash";
import { sprintf } from "sprintf-js";

const store = useStore();

const errors = ref([]);
const timer = ref(null);

const getMyNautiqueEnabled = computed(
  () => store.getters["configuration/getMyNautiqueEnabled"],
);
const getMyNautiqueBoatId = computed(
  () => store.getters["configuration/getMyNautiqueBoatId"],
);
const getAvgFuelConsumption = computed(
  () => store.getters["boat/getAvgFuelConsumption"],
);
const getMyNautiqueFuelLevel = computed(
  () => store.getters["boat/getMyNautiqueFuelLevel"],
);
const getMyNautiqueFuelCapacity = computed(
  () => store.getters["boat/getMyNautiqueFuelCapacity"],
);

const fuelLevelStyle = computed(
  () => `width: ${getMyNautiqueFuelLevel.value}%`,
);

const fuelConsumptionStyle = computed(() => {
  if (getAvgFuelConsumption.value == null) {
    return "width: 0%";
  }
  return `width: ${min([(getAvgFuelConsumption.value / 35) * 100, 100])}%`;
});

const timeTillEmpty = computed(() => {
  if (getAvgFuelConsumption.value == null || !getMyNautiqueEnabled.value) {
    return "N/A hours";
  }
  const hours =
    ((getMyNautiqueFuelLevel.value / 100) * getMyNautiqueFuelCapacity.value) /
    getAvgFuelConsumption.value;
  const hoursFloor = parseInt(hours);
  const minutes = parseInt((hours - hoursFloor) * 60);
  return sprintf("up to %d h %02d min left", hoursFloor, minutes);
});

function startAutoRefresh() {
  timer.value = setInterval(refresh, 5000);
  setTimeout(stopAutoRefresh, 900000);
}

function stopAutoRefresh() {
  if (timer.value != null) {
    console.log("myNautique auto refresh stopped");
    clearInterval(timer.value);
  }
}

function refresh() {
  if (getMyNautiqueEnabled.value && getMyNautiqueBoatId.value) {
    store.dispatch("boat/queryMyNautiqueInfo", getMyNautiqueBoatId.value);
  }
}

store
  .dispatch("configuration/queryConfiguration")
  .then(() => {
    if (getMyNautiqueEnabled.value) {
      console.log("query myNautique with boat ID", getMyNautiqueBoatId.value);
      store
        .dispatch("boat/queryMyNautiqueInfo", getMyNautiqueBoatId.value)
        .then(() => {
          console.log("enable myNautique auto refresh");
          startAutoRefresh();
        })
        .catch((errs) => (errors.value = errs));
    }
  })
  .catch((errs) => (errors.value = errs));

onBeforeUnmount(() => {
  stopAutoRefresh();
});
</script>
