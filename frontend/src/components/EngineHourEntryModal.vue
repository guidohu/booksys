<template>
  <modal-container name="engineHourentryModal" :visible="visible">
    <modal-header
      :closable="true"
      title="Engine Hour Entry"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <div class="row" v-if="errors.length">
        <div class="col-1 d-none d-sm-block" />
        <div class="col-12 col-sm-10">
          <warning-box :errors="errors" />
        </div>
      </div>
      <form @submit.prevent.self="save">
        <input-date-time-local
          id="date"
          label="Date"
          v-model="form.date"
          size="small"
          :disabled="true"
        />
        <input-text
          id="driver"
          label="Driver"
          v-model="form.driver"
          size="small"
          :disabled="true"
        />
        <input-engine-hours
          id="before"
          label="Before"
          v-model="form.beforeHours"
          :displayFormat="displayFormat"
          size="small"
          :disabled="true"
        />
        <input-engine-hours
          id="after"
          label="After"
          v-model="form.afterHours"
          :displayFormat="displayFormat"
          size="small"
          :disabled="true"
        />
        <input-engine-hours
          id="delta"
          label="Difference"
          v-model="form.deltaHours"
          :displayFormat="displayFormat"
          size="small"
          :disabled="true"
        />
        <input-toggle
          id="type"
          label="Type"
          v-model="form.type"
          offLabel="Private"
          onLabel="Course"
        />
      </form>
    </modal-body>
    <modal-footer>
      <button type="button" class="btn btn-outline-info" @click="save">
        <i class="bi bi-check"></i>
        Save
      </button>
      <button type="button" class="btn btn-outline-danger" @click="close">
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import ModalContainer from "@/components/bricks/ModalContainer.vue";
import ModalHeader from "@/components/bricks/ModalHeader.vue";
import ModalBody from "@/components/bricks/ModalBody.vue";
import ModalFooter from "@/components/bricks/ModalFooter.vue";
import * as dayjs from "dayjs";
import InputDateTimeLocal from "@/components/forms/inputs/InputDateTimeLocal.vue";
import InputText from "@/components/forms/inputs/InputText.vue";
import InputEngineHours from "@/components/forms/inputs/InputEngineHours.vue";
import InputToggle from "@/components/forms/inputs/InputToggle.vue";

export default {
  name: "EngineHourEntryModal",
  components: {
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputDateTimeLocal,
    InputText,
    InputEngineHours,
    InputToggle,
    WarningBox,
  },
  props: ["engineHourEntry", "visible", "displayFormat"],
  data() {
    return {
      errors: [],
      form: {
        date: null,
      },
    };
  },
  computed: {
    toggleWidth: function () {
      return 70;
    },
  },
  watch: {
    engineHourEntry: function (newValue) {
      this.setFormContent(newValue);
    },
  },
  methods: {
    setFormContent: function (entry) {
      if (entry != null) {
        this.form = {
          id: entry.id,
          date: dayjs(entry.time * 1000).format("YYYY-MM-DDTHH:mm"),
          driver: entry.user_first_name,
          beforeHours: entry.before_hours,
          afterHours: entry.after_hours,
          deltaHours: entry.delta_hours,
          type: entry.type == 0 ? false : true,
        };
      }
    },
    toggleType: function () {
      this.form.type = !this.form.type;
    },
    close: function () {
      this.$emit("update:visible", false);
    },
    save: function () {
      const update = {
        id: this.form.id,
        type: this.form.type == true ? 1 : 0,
      };
      this.updateEngineHours(update)
        .then(() => this.close())
        .catch((errors) => (this.errors = errors));
    },
    ...mapActions("boat", ["updateEngineHours"]),
  },
  created() {
    this.setFormContent(this.engineHourEntry);
  },
};
</script>
