<template>
  <sectioned-card-module title="Statistics">
    <template v-slot:body>
      <user-heats-modal v-model:visible="showUserHeatsModal" />
      <div class="row">
        <div class="col-6">
          <div class="row">
            <div class="col-12">Riding Time: {{ this.formatDuration(heatTimeMinutesYTD) }}</div>
          </div>
          <div class="row">
            <div class="col-12">
              Riding Time (all-time): {{ this.formatDuration(heatTimeMinutes) }}
            </div>
          </div>
          <div class="row">
            <div class="col-12">Cost: {{ this.formatCost(heatCostYTD) }}</div>
          </div>
          <div class="row">
            <div class="col-12">
              Cost (all-time): {{ this.formatCost(heatCost) }}
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
import UserHeatsModal from "./UserHeatsModal.vue";
import { formatDurationString, formatNumber } from "booksys/libs/formatters";

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
    formatDuration: function(minutes) {
      return formatDurationString(minutes*60, false);
    },
    formatCost: function(cost) {
      return formatNumber(cost) + " " + this.getCurrency;
    }
  },
  created() {
    this.queryConfiguration();
  },
};
</script>
