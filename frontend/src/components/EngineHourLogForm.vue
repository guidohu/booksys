<template>
  <div class="container text-left">
    <div class="row">
      <div class="col-12">
        <warning-box v-if="errors.length > 0" :errors="errors" />
        <form @submit.prevent.self="add">
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
            @click.prevent.self="add"
          >
            {{ showAfter ? "Finish" : "Start" }}
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
      disableBefore: false,
      showAfter: false,
      errors: [],
    };
  },
  computed: {
    ...mapGetters("boat", ["getEngineHourLogLatest"]),
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getEngineHourFormat"]),
  },
  watch: {
    userInfo: function () {
      this.setDriver();
    },
    getEngineHourLogLatest: function (newData) {
      console.log("engineHourLatest just changed", newData);
      if (newData == null) {
        this.form.beforeHours = null;
      } else if (newData.before_hours != null && newData.after_hours != null) {
        this.form.beforeHours = null;
        this.form.type = false;
      } else {
        this.form.beforeHours = newData.before_hours;
      }
      this.form.afterHours = null;

      this.setDriver();
      this.setDisableBefore();
      this.setShowAfter();
    },
  },
  methods: {
    ...mapActions("boat", ["queryEngineHourLogLatest", "addEngineHours"]),
    ...mapActions("configuration", ["queryConfiguration"]),
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
    this.setDisableBefore();
    this.setShowAfter();
  },
};
</script>
