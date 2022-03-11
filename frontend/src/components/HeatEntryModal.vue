<template>
  <modal-container
    name="heatEntryModal"
    title="Heat Entry"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <modal-header
      :closable="true"
      title="Heat Entry"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length > 0" :errors="errors" />
      <form @submit.stop="save">
        <input-date-time-local
          id="date"
          label="Date"
          v-model="form.date"
          disabled
        />
        <input-text
          id="rider"
          label="Rider"
          v-model="form.rider"
          disabled
        />
        <input-text
          id="fare"
          label="Fare"
          v-model="form.fare"
          :suffix="getCurrency + '/min'"
          disabled
        />
        <input-text
          id="duration"
          label="Duration"
          v-model="form.duration"
          @input="durationChangeHandler"
          suffix="mm:ss"
        />
        <input-text
          id="cost"
          label="Cost"
          v-model="form.cost"
          disabled
          :suffix="getCurrency"
        />
        <input-text-multiline
          id="comment"
          label="Comment"
          v-model="form.comment"
          rows="3"
        />
      </form>
    </modal-body>
    <modal-footer class="block-footer">
      <div class="row">
        <div class="col-6 text-start">
          <button type="button" class="btn btn-outline-danger" @click.stop="remove">
            <i class="bi bi-trash"></i>
            Delete
          </button>
        </div>
        <div class="col-6 text-end">
          <button type="submit" class="btn btn-outline-info me-2" @click.prevent.self="save">
            <i class="bi bi-check"></i>
            Save
          </button>
          <button type="button" class="btn btn-outline-danger" @click.prevent.self="close">
            <i class="bi bi-x"></i>
            Cancel
          </button>
        </div>
      </div>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import * as dayjs from "dayjs";
import { sprintf } from "sprintf-js";
import WarningBox from "@/components/WarningBox";
import ModalContainer from "@/components/bricks/ModalContainer.vue";
import ModalHeader from "@/components/bricks/ModalHeader.vue";
import ModalBody from "@/components/bricks/ModalBody.vue";
import ModalFooter from "@/components/bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputTextMultiline from "./forms/inputs/InputTextMultiline.vue"
import InputDateTimeLocal from './forms/inputs/InputDateTimeLocal.vue';

export default {
  name: "HeatEntryModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputText,
    InputTextMultiline,
    InputDateTimeLocal,
  },
  props: ["heat", "visible"],
  data() {
    return {
      errors: [],
      form: {
        date: null,
      },
    };
  },
  computed: {
    ...mapGetters("configuration", ["getCurrency"]),
  },
  watch: {
    heat: function (newHeat) {
      this.setFormDefaults(newHeat);
    },
  },
  methods: {
    setFormDefaults: function (heatData) {
      this.form = {
        id: heatData.heat_id,
        date: dayjs(heatData.timestamp * 1000).format("YYYY-MM-DDTHH:mm"),
        userId: heatData.user_id,
        rider: heatData.first_name + " " + heatData.last_name,
        fare: sprintf("%.2f", heatData.price_per_min),
        duration: this.formatDuration(heatData.duration_s),
        cost: heatData.cost,
        comment: heatData.comment,
      };
    },
    formatDuration: function (durationS) {
      const seconds = durationS % 60;
      const minutes = Math.floor((durationS - seconds) / 60);
      return sprintf("%02d:%02d", minutes, seconds);
    },
    durationChangeHandler: function () {
      console.log("duration changed to:", this.form.duration);
      if (
        this.form.duration == null ||
        !this.form.duration.match(/^\d+:\d+$/)
      ) {
        this.form.cost = "...";
        return;
      } else {
        const durationParts = this.form.duration.split(":");
        const seconds = Number(durationParts[1]);
        const minutes = Number(durationParts[0]);
        const durationSeconds = seconds + 60 * minutes;
        this.form.cost =
          Math.round(((durationSeconds * this.form.fare) / 60) * 100) / 100;
        return;
      }
    },
    close: function () {
      this.$emit("update:visible", false);
    },
    remove: function () {
      this.removeHeat(this.form.id)
        .then(() => this.close())
        .catch((errors) => (this.errors = errors));
    },
    save: function () {
      if (this.form.duration == null) {
        this.errors = ["No duration provided"];
        return;
      }
      const durations = this.form.duration.split(":");
      if (durations.length != 2) {
        this.errors = ["Duration does not have a valid format such as 23:15."];
        return;
      }
      const seconds = Number(durations[1]);
      const minutes = Number(durations[0]);
      if (isNaN(seconds) || isNaN(minutes)) {
        this.errors = ["The duration must use numbers such as 23:15."];
        return;
      }

      const durationSeconds = seconds + 60 * minutes;
      const heatUpdate = {
        id: this.form.id,
        userId: this.form.userId,
        duration: durationSeconds,
        comment: this.form.comment,
      };
      this.updateHeat(heatUpdate)
        .then(() => {
          this.errors = [];
          this.close();
        })
        .catch((errors) => (this.errors = errors));
    },
    ...mapActions("heats", ["removeHeat", "updateHeat"]),
  },
  created() {
    if (this.heat != null) {
      this.setFormDefaults(this.heat);
    }
  },
};
</script>

<style scoped>
  .block-footer {
    display: block;
  }
</style>