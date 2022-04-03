<template>
  <div class="container text-left">
    <div class="row">
      <div class="col-12">
        <warning-box v-if="errors.length > 0" :errors="errors" />
        <form @submit.prevent="add">
          <input-engine-hours
            id="engine-hours"
            label="Engine"
            :display-format="getEngineHourFormat"
            v-model="form.engineHours"
            placeholder="0"
            size="small"
          />
          <input-text-multiline
            id="description"
            label="Description"
            v-model="form.description"
            rows="3"
            placeholder="Add your notes here..."
            size="small"
          />
          <form-button
            type="submit"
            btn-style="info"
            btn-size="small"
            @click.prevent="add"
          >
            Add
          </form-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
import InputEngineHours from "@/components/forms/inputs/InputEngineHours";
import InputTextMultiline from "@/components/forms/inputs/InputTextMultiline";
import FormButton from "@/components/forms/FormButton.vue";

export default {
  name: "MaintenanceLogForm",
  components: {
    WarningBox,
    InputEngineHours,
    InputTextMultiline,
    FormButton,
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
