<template>
  <div>
    <FuelLogForm class="mb-4" />
    <div class="accordion" role="tablist">
      <b-card no-body class="mb-4">
        <b-card-header header-tag="header" class="p-1" role="tab">
          <b-button
            v-b-toggle.accordion-1
            size="sm"
            block
            variant="outline-info"
          >
            {{ visualizationLabel }}
          </b-button>
        </b-card-header>
        <b-collapse
          id="accordion-1"
          accordion="my-accordion"
          role="tabpanel"
          @shown="displayChart"
          @hidden="hideChart"
        >
          <b-card-body>
            <FuelLogChart v-if="showChart" />
          </b-card-body>
        </b-collapse>
      </b-card>
    </div>
    <FuelLogList />
  </div>
</template>

<script>
import FuelLogForm from "@/components/FuelLogForm";
import FuelLogList from "@/components/FuelLogList";
import {
  BCard,
  BCardHeader,
  BButton,
  BCollapse,
  BCardBody,
  VBToggle,
} from "bootstrap-vue";

const FuelLogChart = () =>
  import(/* webpackChunkName: "fuel-log-chart" */ "@/components/FuelLogChart");

export default {
  name: "FuelLogContainer",
  components: {
    FuelLogForm,
    FuelLogChart,
    FuelLogList,
    BCard,
    BCardHeader,
    BButton,
    BCollapse,
    BCardBody,
  },
  directives: {
    "b-toggle": VBToggle,
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
    displayChart: function () {
      this.showChart = true;
    },
    hideChart: function () {
      this.showChart = false;
    },
  },
};
</script>
