<template>
  <sectioned-card-module>
    <template v-slot:header>
      <div class="row">
        <div class="col-6">
          <h5>Statistics</h5>
        </div>
      </div>
    </template>
    <template v-slot:body>
      <user-heats-modal v-model:visible="showUserHeatsModal" />
      <div class="row">
        <div class="col-6">
          <div class="row">
            <div class="col-12">Riding Time: {{ heatTimeMinutesYTD }} min </div>
          </div>
          <div class="row">
            <div class="col-12">
              Riding Time (all-time): {{ heatTimeMinutes }} min
            </div>
          </div>
          <div class="row">
            <div class="col-12"> Cost: {{ heatCostYTD }} {{ getCurrency }} </div>
          </div>
          <div class="row">
            <div class="col-12">
              Cost (all-time): {{ heatCost }} {{ getCurrency }}
            </div>
          </div>
        </div>
        <div class="col-6">
          <button type="button" class="btn btn-outline-info btn-sm" @click="showLatestHeats">
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

export default {
  name: "UserStatisticsCard",
  components: {
    UserHeatsModal,
    SectionedCardModule
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
  },
  created() {
    this.queryConfiguration();
  },
};
</script>
