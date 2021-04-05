<template>
  <b-container float class="text-left">
    <b-row>
      <b-col cols="12">
        <WarningBox v-if="errors.length > 0" :errors="errors" />
        <b-form @submit="add">
          <b-form-group
            id="input-group-driver"
            label="Driver"
            label-for="input-driver"
            label-cols="3"
          >
            <b-form-input
              id="input-driver"
              v-model="form.driverName"
              type="text"
              size="sm"
              disabled
            />
          </b-form-group>
          <engine-hours
            label="Before"
            v-model="form.beforeHours"
            :display-format="getEngineHourFormat"
            :disabled="disableBefore"
          />
          <engine-hours
            v-if="showAfter"
            label="After"
            v-model="form.afterHours"
            :display-format="getEngineHourFormat"
          />
          <b-form-group
            id="input-group-type"
            label="Type"
            label-for="input-type"
            label-cols="3"
          >
            <toggle-button
              id="input-type"
              :value="false"
              @change="toggleType"
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{ checked: 'Course', unchecked: 'Private' }"
            />
          </b-form-group>
          <b-row class="text-right">
            <b-col cols="9" offset="3">
              <b-button
                v-if="!disableBefore"
                block
                variant="outline-info"
                size="sm"
                v-on:click="add"
                >Start</b-button
              >
              <b-button
                v-if="disableBefore"
                block
                variant="outline-info"
                size="sm"
                v-on:click="add"
                >Finish</b-button
              >
            </b-col>
          </b-row>
        </b-form>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { ToggleButton } from "vue-js-toggle-button";
import WarningBox from "@/components/WarningBox";
import EngineHours from "@/components/forms/inputs/EngineHours";
import {
  BContainer,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BButton,
} from "bootstrap-vue";

export default {
  name: "EngineHourLogForm",
  components: {
    ToggleButton,
    WarningBox,
    EngineHours,
    BContainer,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
    BButton,
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
      engineHourLabel: null,
      disableBefore: false,
      showAfter: false,
      errors: [],
    };
  },
  computed: {
    ...mapGetters("boat", ["getEngineHourLogLatest"]),
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getEngineHourFormat"]),
    toggleWidth: function () {
      return 70;
    },
  },
  watch: {
    userInfo: function () {
      console.log("userInfo just changed");
      this.setDriver();
    },
    getEngineHourLogLatest: function (newData) {
      console.log("engineHourLatest just changed", newData);
      if (newData == null) {
        this.form.beforeHours = null;
      } else if (newData.before_hours != null && newData.after_hours != null) {
        this.form.beforeHours = null;
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
    toggleType: function () {
      this.form.type = !this.form.type;
    },
    add: function () {
      const type = this.form.type ? 1 : 0;

      const data = {
        user_id: this.form.driverId,
        engine_hours_before: this.form.beforeHours,
        engine_hours_after: this.form.afterHours,
        type: type,
      };
      this.addEngineHours(data)
        .then(() => (this.errors = []))
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
