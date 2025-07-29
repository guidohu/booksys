<template>
  <modal-container name="refuelModal" :visible="visible">
    <modal-header
      :closable="true"
      title="Refuel"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length > 0" :errors="errors" />
      <form @submit.stop="saveFuel">
        <input-engine-hours
          id="engine-hours"
          label="Engine"
          :display-format="getEngineHourFormat"
          v-model="form.engineHours"
          :description="engineHourDescription"
          placeholder="0"
          size="small"
        />
        <input-currency
          id="cost"
          label="Cost"
          :currency="getCurrency"
          v-model="form.cost"
          placeholder="0"
          size="small"
        />
        <input-fuel
          id="fuel"
          label="Fuel"
          v-model="form.liters"
          placeholder="0"
          size="small"
        />
      </form>
    </modal-body>
    <modal-footer>
      <button type="button" class="btn btn-outline-info" @click="saveFuel">
        <i class="bi bi-check" />
        Save
      </button>
      <button
        type="button"
        class="btn btn-outline-danger mb-1"
        @click="$emit('update:visible', false)"
      >
        <i class="bi bi-x" />
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputEngineHours from "booksys/components/forms/inputs/InputEngineHours.vue";
import InputCurrency from "./forms/inputs/InputCurrency.vue";
import InputFuel from "./forms/inputs/InputFuel.vue";

const props = defineProps(["visible"]);
const emit = defineEmits(["saved", "update:visible"]);

const store = useStore();

const errors = ref([]);
const form = ref({
  engineHours: 0,
  cost: 0,
  liters: 0,
});
const engineHourDescription = ref(null);

const userInfo = computed(() => store.getters["login/userInfo"]);
const getMyNautiqueEngineHours = computed(() => store.getters["boat/getMyNautiqueEngineHours"]);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const getEngineHourFormat = computed(() => store.getters["configuration/getEngineHourFormat"]);

watch(getMyNautiqueEngineHours, (newValue) => {
  if (form.value.engineHours == null || form.value.engineHours === "") {
    form.value.engineHours = newValue;
    engineHourDescription.value = "prefilled by myNautique";
  }
});

function saveFuel() {
  const engineHours = isNaN(form.value.engineHours) || form.value.engineHours === "" ? null : form.value.engineHours;
  const liters = isNaN(form.value.liters) || form.value.liters === "" ? null : form.value.liters;
  const cost = isNaN(form.value.cost) || form.value.cost === "" ? null : form.value.cost;
  const entry = {
    user_id: parseInt(userInfo.value.id),
    engine_hours: engineHours,
    liters: liters,
    cost: cost,
  };
  store.dispatch("boat/addFuelEntry", entry)
    .then(() => {
      resetForm();
      close();
    })
    .catch((errs) => (errors.value = errs));
}

function resetForm() {
  form.value.engineHours = "";
  form.value.cost = "";
  form.value.liters = "";
}

function close() {
  emit("update:visible", false);
}

store.dispatch("configuration/queryConfiguration");

if (getMyNautiqueEngineHours.value != null && getMyNautiqueEngineHours.value > 0) {
  form.value.engineHours = getMyNautiqueEngineHours.value;
  engineHourDescription.value = "prefilled by myNautique";
}
</script>
