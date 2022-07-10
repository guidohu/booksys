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

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
import { min } from "lodash";
import { sprintf } from "sprintf-js";

export default {
  name: "FuelLogStatus",
  components: {
    WarningBox,
  },
  data() {
    return {
      errors: [],
      fuelLevel: 0,
      timer: null,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", [
      "getCurrency",
      "getEngineHourFormat",
      "getMyNautiqueEnabled",
      "getMyNautiqueBoatId",
    ]),
    ...mapGetters("boat", [
      "getAvgFuelConsumption",
      "getMyNautiqueFuelLevel",
      "getMyNautiqueFuelCapacity",
    ]),
    fuelLevelStyle: function () {
      return "width: " + this.getMyNautiqueFuelLevel + "%";
    },
    fuelConsumptionStyle: function () {
      if (this.getAvgFuelConsumption == null) {
        return "width: 0%";
      }
      return (
        "width: " + min([(this.getAvgFuelConsumption / 35) * 100, 100]) + "%"
      );
    },
    timeTillEmpty: function () {
      if (this.getAvgFuelConsumption == null) {
        return "N/A hours";
      }
      if (this.getMyNautiqueEnabled == true) {
        const hours =
          ((this.getMyNautiqueFuelLevel / 100) *
            this.getMyNautiqueFuelCapacity) /
          this.getAvgFuelConsumption;
        const hoursFloor = parseInt(hours);
        const minutes = parseInt((hours - hoursFloor) * 60);
        return sprintf("up to %d h %02d min left", hoursFloor, minutes);
      }
      return "N/A hours";
    },
  },
  methods: {
    startAutoRefresh: function () {
      this.timer = setInterval(
        this.refresh,
        5000
      );
      // make sure we stop auto refresh again at
      // some point (15min)
      setTimeout(this.stopAutoRefresh, 900000);
    },
    stopAutoRefresh: function () {
      if (this.timer != null) {
        console.log("myNautique auto refresh stopped");
        clearInterval(this.timer);
      }
    },
    refresh: function () {
      if(this.getMyNautiqueEnabled && this.getMyNautiqueBoatId) {
        this.queryMyNautiqueInfo(this.getMyNautiqueBoatId);
      }
    },
    ...mapActions("boat", ["queryMyNautiqueInfo"]),
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.queryConfiguration()
      .then(() => {
        // query myNautique if enabled
        if (this.getMyNautiqueEnabled && this.getMyNautiqueBoatId) {
          console.log("query myNautique");
          this.queryMyNautiqueInfo(this.getMyNautiqueBoatId)
            .then(() => {
              console.log("enable myNautique auto refresh");
              this.startAutoRefresh();
            });
        }
      })
  },
  beforeUnmount() {
    this.stopAutoRefresh();
  }
};
</script>
