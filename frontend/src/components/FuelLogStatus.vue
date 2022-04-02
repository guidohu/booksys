<template>
  <div class="container text-left">
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <div v-if="getMyNautiqueEnabled" class="row">
      <label class="col-3 col-form-label">Level</label>
      <div class="col-9">
        <div  class="text-end">
          {{ getMyNautiqueFuelLevel }}% {{ " " }} ( {{ parseInt(getMyNautiqueFuelCapacity * getMyNautiqueFuelLevel / 100) }} L)
        </div>
        <div class="progress" style="height: 3px;">
          <div class="progress-bar bg-info" role="progressbar" :style="progressStatus" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100"></div>
        </div>
      </div>
    </div>
    <div v-if="getMyNautiqueEnabled" class="row">
      <label class="col-3 col-form-label">Consumption</label>
      <div class="col-9">
        <div class="text-end">
          {{ getMyNautiqueFuelLevel }}% {{ " " }} ( {{ parseInt(getMyNautiqueFuelCapacity * getMyNautiqueFuelLevel / 100) }} L/h)
        </div>
        <div class="progress" style="height: 3px;">
          <div class="progress-bar bg-info" role="progressbar" :style="progressStatus" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100"></div>
        </div>
      </div>
    </div>
    <div v-if="getMyNautiqueEnabled" class="row">
      <label class="col-3 col-form-label">Time Left</label>
      <div class="col-9">
        <div class="text-end">
          {{ getMyNautiqueFuelLevel }}% {{ " " }} ( {{ parseInt(getMyNautiqueFuelCapacity * getMyNautiqueFuelLevel / 100) }} hours)
        </div>
        <div class="progress" style="height: 3px;">
          <div class="progress-bar bg-info" role="progressbar" :style="progressStatus" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";

export default {
  name: "FuelLogStatus",
  components: {
    WarningBox,
  },
  data() {
    return {
      errors: [],
      enabled: false,
      fuelLevel: 0,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getCurrency", "getEngineHourFormat", "getMyNautiqueEnabled", "getMyNautiqueBoatId"]),
    ...mapGetters("boat", ["getMyNautiqueFuelLevel", "getMyNautiqueFuelCapacity"]),
    progressStatus: function() {
      return "width: " + this.getMyNautiqueFuelLevel + "%";
    }
  },
  methods: {
    add: function () {
      const entry = {
        user_id: this.userInfo.id,
        engine_hours: this.form.engineHours,
        liters: this.form.liters,
        cost: this.form.cost,
      };
      this.addFuelEntry(entry)
        .then(() => {
          this.resetForm();
        })
        .catch((errors) => (this.errors = errors));
    },
    resetForm: function () {
      this.form = {
        engineHours: null,
        cost: null,
        liters: null,
      };
      this.errors = [];
    },
    ...mapActions("boat", ["addFuelEntry","queryMyNautiqueInfo"]),
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  watch: {
    "getMyNautiqueEnabled": function(enabled) {
      if(enabled){
        this.queryMyNautiqueInfo(this.getMyNautiqueBoatId)
        .then(() => {
          console.log("queried my nautique info");
          // this.enabled = true;
        })
        .catch((error) => {
          this.errors = [error];
        })
      }
      console.log("myNautiqueEnabled value changed to", enabled);
    },
    // "getMyNautiqueFuelLevel": function(value) {
    //   console.log("fuelLevel changed");
    //   this.fuelLevel = value;
    // },
  },
  created() {
    this.queryConfiguration();
    if(this.getMyNautiqueEnabled && this.getMyNautiqueBoatId){
      this.queryMyNautiqueInfo(this.getMyNautiqueBoatId);
    }
  },
};
</script>
