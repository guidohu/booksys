<template>
  <sectioned-card-module title="Statistics">
    <template v-slot:body>
      <user-heats-modal v-model:visible="showUserHeatsModal" />
      <div class="row">
        <div class="col-6">
          <div class="row">
            <div class="col-12">Riding Time: {{ formatDuration(heatTimeMinutesYTD) }}</div>
          </div>
          <div class="row">
            <div class="col-12">
              Riding Time (all-time): {{ formatDuration(heatTimeMinutes) }}
            </div>
          </div>
          <div class="row">
            <div class="col-12">Cost: {{ formatCost(heatCostYTD) }}</div>
          </div>
          <div class="row">
            <div class="col-12">
              Cost (all-time): {{ formatCost(heatCost) }}
            </div>
          </div>
        </div>
        <div class="col-6">
          <button
            type="button"
            class="btn btn-outline-info btn-sm"
            @click="showLatestHeats"
          >
            <i class="bi bi-list-ul"></i>
            Latest Heats
          </button>
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script setup>
import { ref, computed } from "vue";
import { useStore } from "vuex";
import SectionedCardModule from "./bricks/SectionedCardModule.vue";
import UserHeatsModal from "./UserHeatsModal.vue";
import { formatDurationString, formatNumber } from "booksys/libs/formatters";

const store = useStore();

const showUserHeatsModal = ref(false);

const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const heatTimeMinutes = computed(() => store.getters["user/heatTimeMinutes"]);
const heatTimeMinutesYTD = computed(() => store.getters["user/heatTimeMinutesYTD"]);
const heatCost = computed(() => store.getters["user/heatCost"]);
const heatCostYTD = computed(() => store.getters["user/heatCostYTD"]);

const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function showLatestHeats() {
  showUserHeatsModal.value = true;
}

function formatDuration(minutes) {
  return formatDurationString(minutes * 60, false);
}

function formatCost(cost) {
  return formatNumber(cost) + " " + getCurrency.value;
}

queryConfiguration();
</script>
