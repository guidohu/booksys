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

<script>
import { defineAsyncComponent } from "vue";
import FuelLogStatus from "@/components/FuelLogStatus";
import FuelLogForm from "@/components/FuelLogForm";
import FuelLogList from "@/components/FuelLogList";

const FuelLogChart = defineAsyncComponent(() =>
  import(/* webpackChunkName: "fuel-log-chart" */ "@/components/FuelLogChart")
);

export default {
  name: "FuelLogContainer",
  components: {
    FuelLogStatus,
    FuelLogForm,
    FuelLogChart,
    FuelLogList,
  },
  data() {
    return {
      showChart: false,
      visualizationLabel: "Show Visualization",
    };
  },
  watch: {
    showChart: function () {
      console.log("showChart changed");
      if (this.showChart == false) {
        this.visualizationLabel = "Show Visualization";
      } else {
        this.visualizationLabel = "Hide Visualization";
      }
    },
  },
  methods: {
    toggleChart: function () {
      this.showChart = !this.showChart;
    },
  },
};
</script>

<style scoped>
.box {
  display: flex;
  flex-flow: column;
  height: 100%;
  overflow: scroll;
}
</style>
