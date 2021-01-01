<template>
  <div>
    <FuelLogForm class="mb-4"/>
    <div class="accordion" role="tablist">
      <b-card no-body class="mb-4">
        <b-card-header header-tag="header" class="p-1" role="tab">
          <b-button size="sm" block v-b-toggle.accordion-1 variant="outline-info" v-on:click="toggleVisualization">{{visualizationLabel}}</b-button>
        </b-card-header>
        <b-collapse id="accordion-1" :visible="showChart" accordion="my-accordion" role="tabpanel">
          <b-card-body>
            <FuelLogChart v-if="showChart"/>
          </b-card-body>
        </b-collapse>
      </b-card>
    </div>
    <FuelLogList/>
  </div>
</template>

<script>
import Vue from 'vue';
import FuelLogForm from '@/components/FuelLogForm';
import FuelLogChart from '@/components/FuelLogChart';
import FuelLogList from '@/components/FuelLogList';
import {
  BCard,
  BCardHeader,
  BButton,
  BCollapse,
  BCardBody
} from 'bootstrap-vue';

export default Vue.extend({
  name: "FuelLogContainer",
  components: {
    FuelLogForm,
    FuelLogChart,
    FuelLogList,
    BCard,
    BCardHeader,
    BButton,
    BCollapse,
    BCardBody
  },
  data() {
    return {
      showChart: false,
      visualizationLabel: "Show Visualization"
    }
  },
  watch: {
    showChart: function(){
      console.log("showChart changed");
      if(this.showChart == false){
        this.visualizationLabel = "Show Visualization";
      }else{
        this.visualizationLabel = "Hide Visualization";
      }
    }
  },
  methods: {
    toggleVisualization: function(){
      this.showChart = !this.showChart;
    }
  }
})
</script>