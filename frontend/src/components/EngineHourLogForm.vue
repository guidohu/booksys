<template>
  <div class="container text-left">
    <div class="row">
      <div class="col-12">
        <warning-box v-if="errors.length > 0" :errors="errors" />
        <form @submit.prevent="add">
          <input-text
            id="driver"
            label="Driver"
            v-model="form.driverName"
            size="small"
            :disabled="true"
            type="text"
          />
          <input-engine-hours
            id="before-hours"
            label="Before"
            :display-format="getEngineHourFormat"
            :disabled="disableBefore"
            v-model="form.beforeHours"
            :description="beforeDescription"
            placeholder="0"
            size="small"
          />
          <input-engine-hours
            v-if="showAfter"
            id="after-hours"
            label="After"
            :display-format="getEngineHourFormat"
            :disabled="!showAfter"
            v-model="form.afterHours"
            :description="afterDescription"
            placeholder="0"
            size="small"
          />
          <input-toggle
            id="type"
            label="Type"
            off-label="Private"
            on-label="Course"
            v-model="form.type"
          />
          <form-button
            type="submit"
            btn-style="info"
            btn-size="small"
            @click.prevent="add"
          >
            {{ showAfter ? "Check Out" : "Check In" }}
          </form-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import InputEngineHours from "@/components/forms/inputs/InputEngineHours.vue";
import InputText from "@/components/forms/inputs/InputText.vue";
import InputToggle from "@/components/forms/inputs/InputToggle.vue";
import FormButton from "@/components/forms/FormButton.vue";
import WarningBox from "@/components/WarningBox.vue";

export default {
  name: "EngineHourLogForm",
  components: {
    InputEngineHours,
    InputText,
    InputToggle,
    FormButton,
    WarningBox,
  },
  data() {
    return {
      form: {
        driverId: null,
        driverName: "",
        beforeHours: null,
        afterHours: null,
        type: false,
      },
      beforeDescription: null,
      afterDescription: null,
      disableBefore: false,
      showAfter: false,
      errors: [],
    };
  },
  computed: {
    ...mapGetters("boat", ["getEngineHourLogLatest", "getMyNautiqueEngineHours"]),
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getEngineHourFormat", "getMyNautiqueEnabled", "getMyNautiqueBoatId"]),
  },
  watch: {
    userInfo: function () {
      this.setDriver();
    },
    getEngineHourLogLatest: function (newData) {
      console.log("engineHourLatest just changed", newData);
      if (newData == null) {
        // no engineHourLatest so far
        this.form.beforeHours = null;
        this.form.beforeDescription = null;
        this.prefillBefore();
        this.form.afterHours = null;
      } else if (newData.before_hours != null && newData.after_hours != null) {
        // the latest engine entry is complete
        // (has before and after)        
        this.form.beforeHours = null;
        this.form.type = false;
        this.prefillBefore();
        this.form.afterHours = null;
      } else {
        this.form.beforeHours = newData.before_hours;
        this.form.afterHours = null;
        this.prefillAfter();
      }

      this.setDriver();
      this.setDisableBefore();
      this.setShowAfter();
    },
    getMyNautiqueBoatId: function(boatId) {
      this.queryMyNautiqueInfo(boatId);
    },
    getMyNautiqueEngineHours: function (newEngineHours) {
      if(this.disableBefore == false) {
        this.prefillBefore();
      } else {
        this.prefillAfter();
      }
    },
  },
  methods: {
    ...mapActions("boat", ["queryEngineHourLogLatest", "addEngineHours", "queryMyNautiqueInfo"]),
    ...mapActions("configuration", ["queryConfiguration"]),
    prefillBefore: function() {
      if(this.getMyNautiqueEnabled && this.getMyNautiqueEngineHours != null){
        this.form.beforeHours = this.getMyNautiqueEngineHours;
        this.beforeDescription = "prefilled by myNautique";
        this.afterDescription = null;
      }
    },
    prefillAfter: function() {
      if(this.getMyNautiqueEnabled && this.getMyNautiqueEngineHours != null){
        this.form.afterHours = this.getMyNautiqueEngineHours;
        this.beforeDescription = null;
        this.afterDescription = "prefilled by myNautique";
      }
    },
    setDisableBefore: function () {
      if (
        this.getEngineHourLogLatest != null &&
        this.getEngineHourLogLatest.after_hours != null
      ) {
        this.disableBefore = false;
        return;
      } else if (this.getEngineHourLogLatest == null) {
        this.disableBefore = false;
        return;
      }
      this.disableBefore = true;
    },
    setShowAfter: function () {
      if (
        this.getEngineHourLogLatest != null &&
        this.getEngineHourLogLatest.after_hours == null
      ) {
        this.showAfter = true;
        return;
      }
      this.showAfter = false;
    },
    setDriver: function () {
      if (
        this.getEngineHourLogLatest != null &&
        this.getEngineHourLogLatest.after_hours == null
      ) {
        this.form.driverName =
          this.getEngineHourLogLatest.user_first_name +
          " " +
          this.getEngineHourLogLatest.user_last_name;
        this.form.driverId = this.getEngineHourLogLatest.user_id;
      } else if (this.userInfo != null) {
        this.form.driverName =
          this.userInfo.first_name + " " + this.userInfo.last_name;
        this.form.driverId = this.userInfo.id;
      } else {
        this.form.driverName = "unknown";
        this.form.driverId = null;
      }
    },
    add: function () {
      // get the type
      // 0: private
      // 1: course
      const type = this.form.type ? 1 : 0;

      const data = {
        user_id: this.form.driverId,
        engine_hours_before: this.form.beforeHours,
        engine_hours_after: this.form.afterHours,
        type: type,
      };
      this.addEngineHours(data)
        .then(() => {
          this.errors = [];
        })
        .catch((errors) => (this.errors = errors));
    },
  },
  created() {
    this.queryConfiguration();
    this.setDriver();
    this.queryEngineHourLogLatest();
    if(this.getMyNautiqueBoatId){
      this.queryMyNautiqueInfo(this.getMyNautiqueBoatId);
    }
    this.setDisableBefore();
    this.setShowAfter();
  },
};
</script>
