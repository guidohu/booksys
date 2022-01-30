<template>
  <b-container float class="text-left">
    <b-row>
      <b-col cols="12">
        <WarningBox v-if="errors.length > 0" :errors="errors" />
        <b-form @submit="add">
          <engine-hours
            v-model="form.engineHours"
            label="Engine Hrs"
            :display-format="getEngineHourFormat"
          />
          <b-form-group
            id="input-group-description"
            label="Description"
            label-for="input-description"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-textarea
                id="input-description"
                v-model="form.description"
                type="text"
                placeholder="Add your maintenance notes here..."
                rows="3"
                max-rows="3"
              />
            </b-input-group>
          </b-form-group>
          <b-row class="text-right">
            <b-col cols="9" offset="3">
              <b-button block variant="outline-info" size="sm" @click="add">
                Add
              </b-button>
            </b-col>
          </b-row>
        </b-form>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
import EngineHours from "@/components/forms/inputs/EngineHours";
import {
  BContainer,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BInputGroup,
  BFormTextarea,
  BButton,
} from "bootstrap-vue";

export default {
  name: "MaintenanceLogForm",
  components: {
    WarningBox,
    EngineHours,
    BContainer,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BInputGroup,
    BFormTextarea,
    BButton,
  },
  data() {
    return {
      errors: [],
      form: {
        engineHours: null,
        description: null,
      },
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getEngineHourFormat"]),
  },
  methods: {
    add: function () {
      const entry = {
        user_id: this.userInfo.id,
        engine_hours: this.form.engineHours,
        description: this.form.description,
      };
      this.addMaintenanceEntry(entry)
        .then(() => {
          this.resetForm();
        })
        .catch((errors) => (this.errors = errors));
    },
    resetForm: function () {
      this.form = {
        engineHours: null,
        description: null,
      };
      this.errors = [];
    },
    ...mapActions("boat", ["addMaintenanceEntry"]),
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.queryConfiguration();
  },
};
</script>
