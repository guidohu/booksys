<template>
  <div class="box">
    <fuel-log-status class="mb-2" />
    <fuel-log-form class="mb-4" />
    <div class="d-grid gap-2 ms-1 me-1 mb-2">
      <button
        class="btn btn-outline-info btn-sm"
        data-bs-toggle="collapse"
        data-bs-target="#fuel-visualization"
        aria-expanded="false"
        aria-controls="collapseExample"
        @click="toggleChart"
      >
        {{ visualizationLabel }}
      </button>
    </div>
    <div class="collapse ms-1 me-1" id="fuel-visualization">
      <div class="card card-body">
        <fuel-log-chart v-if="showChart" />
      </div>
    </div>
    <fuel-log-list />
  </div>
</template>

<script setup>
import { defineAsyncComponent, ref, watch } from "vue";
import FuelLogStatus from "../components/FuelLogStatus.vue";
import FuelLogForm from "../components/FuelLogForm.vue";
import FuelLogList from "../components/FuelLogList.vue";

const FuelLogChart = defineAsyncComponent(
  () => import("booksys/components/FuelLogChart.vue"),
);

const showChart = ref(false);
const visualizationLabel = ref("Show Visualization");

watch(showChart, (newValue) => {
  console.log("showChart changed");
  if (newValue === false) {
    visualizationLabel.value = "Show Visualization";
  } else {
    visualizationLabel.value = "Hide Visualization";
  }
});

function toggleChart() {
  showChart.value = !showChart.value;
}
</script>

<style scoped>
.box {
  display: flex;
  flex-flow: column;
  height: 100%;
  overflow: scroll;
}
</style>
