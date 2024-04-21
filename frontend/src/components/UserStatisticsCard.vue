<template>
  <sectioned-card-module title="Statistics">
    <template v-slot:body>
      <user-heats-modal v-model:visible="showUserHeatsModal" />
      <div class="row">
        <div class="col-6">
          <div class="row">
            <div class="col-12">Riding Time: {{ this.formatHeatTime(heatTimeMinutesYTD) }}</div>
          </div>
          <div class="row">
            <div class="col-12">
              Riding Time (all-time): {{ this.formatHeatTime(heatTimeMinutes) }}
            </div>
          </div>
          <div class="row">
            <div class="col-12">Cost: {{ heatCostYTD }} {{ getCurrency }}</div>
          </div>
          <div class="row">
            <div class="col-12">
              Cost (all-time): {{ this.formatCost(heatCost) }} {{ getCurrency }}
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

<script>
import { mapGetters, mapActions } from "vuex";
import SectionedCardModule from "./bricks/SectionedCardModule.vue";
import UserHeatsModal from "./UserHeatsModal";
import { floor } from "lodash";
import { sprintf } from "sprintf-js";

export default {
  name: "UserStatisticsCard",
  components: {
    UserHeatsModal,
    SectionedCardModule,
  },
  data() {
    return {
      showUserHeatsModal: false,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getCurrency"]),
    ...mapGetters("user", [
      "heatHistory",
      "heatTimeMinutes",
      "heatTimeMinutesYTD",
      "heatCost",
      "heatCostYTD",
    ]),
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    showLatestHeats: function () {
      this.showUserHeatsModal = true;
    },
    // TODO add these formatter to the formatters.js lib.
    formatHeatTime: function(minutes) {
      // get minutes
      let fmtMinutes = 0;
      fmtMinutes = minutes % 60;
      minutes = minutes - fmtMinutes;
      // get hours
      let fmtHours =  (minutes / 60) % 24;
      minutes = minutes - (fmtHours * 60);
      // get days
      let fmtDays = (minutes / 60 / 24);
      if(fmtDays > 0){
        return sprintf("%d days %d hours %d minutes", fmtDays, fmtHours, fmtMinutes);
      }
      if(fmtHours > 0){
        return sprintf("%d hours %d minutes", fmtHours, fmtMinutes);
      }
      return sprintf("%d minutes", fmtMinutes);
    },
    formatCost: function(cost) {
      const c = Number(cost);
      // TODO make locale override configurable.
      return c.toLocaleString('ch-DE');
    }
  },
  created() {
    this.queryConfiguration();
  },
};
</script>
