<template>
  <modal-container
    name="heatEntryModal"
    title="Heat Entry"
    :visible="visible"
    :inert="!visible"
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
        <input-text id="rider" label="Rider" v-model="form.rider" disabled />
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
        <div class="col-4 text-start">
          <button
            type="button"
            class="btn btn-outline-danger"
            @click.stop="remove"
          >
            <i class="bi bi-trash"></i>
            Delete
          </button>
        </div>
        <div class="col-8 text-end">
          <button
            type="submit"
            class="btn btn-outline-info me-2"
            @click.prevent.self="save"
          >
            <i class="bi bi-check"></i>
            Save
          </button>
          <button
            type="button"
            class="btn btn-outline-danger"
            @click.prevent.self="close"
          >
            <i class="bi bi-x"></i>
            Cancel
          </button>
        </div>
      </div>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import dayjs from "dayjs";
import { sprintf } from "sprintf-js";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "booksys/components/bricks/ModalContainer.vue";
import ModalHeader from "booksys/components/bricks/ModalHeader.vue";
import ModalBody from "booksys/components/bricks/ModalBody.vue";
import ModalFooter from "booksys/components/bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputTextMultiline from "./forms/inputs/InputTextMultiline.vue";
import InputDateTimeLocal from "./forms/inputs/InputDateTimeLocal.vue";

const props = defineProps(["heat", "visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const errors = ref([]);
const form = ref({ date: null });

const getCurrency = computed(() => store.getters["configuration/getCurrency"]);

watch(
  () => props.heat,
  (newHeat) => {
    setFormDefaults(newHeat);
  },
);

function setFormDefaults(heatData) {
  form.value = {
    id: heatData.heat_id,
    date: dayjs(heatData.timestamp * 1000).format("YYYY-MM-DDTHH:mm"),
    userId: heatData.user_id,
    rider: `${heatData.first_name} ${heatData.last_name}`,
    fare: sprintf("%.2f", heatData.price_per_min),
    duration: formatDuration(heatData.duration_s),
    cost: sprintf("%.2f", heatData.cost),
    comment: heatData.comment,
  };
}

function formatDuration(durationS) {
  const seconds = durationS % 60;
  const minutes = Math.floor((durationS - seconds) / 60);
  return sprintf("%02d:%02d", minutes, seconds);
}

function durationChangeHandler() {
  console.log("duration changed to:", form.value.duration);
  if (form.value.duration == null || !form.value.duration.match(/^\d+:\d+$/)) {
    form.value.cost = "...";
    return;
  } else {
    const durationParts = form.value.duration.split(":");
    const seconds = Number(durationParts[1]);
    const minutes = Number(durationParts[0]);
    const durationSeconds = seconds + 60 * minutes;
    form.value.cost =
      Math.round(((durationSeconds * form.value.fare) / 60) * 100) / 100;
  }
}

function close() {
  emit("update:visible", false);
}

function remove() {
  store
    .dispatch("heats/removeHeat", form.value.id)
    .then(() => close())
    .catch((errs) => (errors.value = errs));
}

function save() {
  if (form.value.duration == null) {
    errors.value = ["No duration provided"];
    return;
  }
  const durations = form.value.duration.split(":");
  if (durations.length !== 2) {
    errors.value = ["Duration does not have a valid format such as 23:15."];
    return;
  }
  const seconds = Number(durations[1]);
  const minutes = Number(durations[0]);
  if (isNaN(seconds) || isNaN(minutes)) {
    errors.value = ["The duration must use numbers such as 23:15."];
    return;
  }

  const durationSeconds = seconds + 60 * minutes;
  const heatUpdate = {
    id: form.value.id,
    userId: form.value.userId,
    duration: durationSeconds,
    comment: form.value.comment,
  };
  store
    .dispatch("heats/updateHeat", heatUpdate)
    .then(() => {
      errors.value = [];
      close();
    })
    .catch((errs) => (errors.value = errs));
}

if (props.heat != null) {
  setFormDefaults(props.heat);
}
</script>

<style scoped>
.block-footer {
  display: block;
}
</style>
