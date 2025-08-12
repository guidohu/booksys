<template>
  <modal-container name="fuelEntrymodal" :visible="visible" :inert="!visible">
    <modal-header
      :closable="true"
      title="Fuel Entry"
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
        <input-text
          id="consumption"
          label="Consumption"
          v-model="form.averageFuelPerHour"
          size="small"
          :disabled="true"
          suffix="ltrs/h"
        />
        <input-engine-hours
          id="engine-hours"
          label="Engine"
          v-model="form.engineHours"
          :displayFormat="getEngineHourFormat"
          size="small"
        />
        <input-fuel
          id="fuel"
          label="Fuel"
          v-model="form.fuel"
          placeholder="0"
          size="small"
        />
        <input-currency
          v-if="!form.isDiscounted"
          id="cost"
          label="Cost"
          :currency="getCurrency"
          v-model="form.cost"
          placeholder="0"
          size="small"
        />
        <input-currency
          v-if="form.isDiscounted"
          id="cost-gross"
          label="Cost (Gross)"
          :currency="getCurrency"
          v-model="form.costGross"
          placeholder="0"
          size="small"
        />
        <input-toggle
          id="discount"
          label="Discount"
          v-model="form.isDiscounted"
          offLabel="No Discount"
          onLabel="Discounted"
        />
        <input-currency
          v-if="form.isDiscounted"
          id="cost-net"
          label="Cost (Net)"
          :currency="getCurrency"
          v-model="form.costNet"
          placeholder="0"
          size="small"
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
import WarningBox from "booksys/components/WarningBox.vue";
import {
  formatCurrency,
  formatFuelConsumption,
  formatFuel,
} from "booksys/libs/formatters";
import dayjs from "dayjs";
import ModalContainer from "booksys/components/bricks/ModalContainer.vue";
import ModalHeader from "booksys/components/bricks/ModalHeader.vue";
import ModalBody from "booksys/components/bricks/ModalBody.vue";
import ModalFooter from "booksys/components/bricks/ModalFooter.vue";
import InputDateTimeLocal from "booksys/components/forms/inputs/InputDateTimeLocal.vue";
import InputText from "booksys/components/forms/inputs/InputText.vue";
import InputEngineHours from "booksys/components/forms/inputs/InputEngineHours.vue";
import InputFuel from "booksys/components/forms/inputs/InputFuel.vue";
import InputCurrency from "booksys/components/forms/inputs/InputCurrency.vue";
import InputToggle from "booksys/components/forms/inputs/InputToggle.vue";

const props = defineProps(["fuelEntry", "visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const errors = ref([]);
const form = ref({ date: null });

const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const getEngineHourFormat = computed(
  () => store.getters["configuration/getEngineHourFormat"],
);

watch(
  () => props.fuelEntry,
  (newValue) => {
    setFormContent(newValue);
  },
);

function setFormContent(entry) {
  if (entry != null) {
    const costGross =
      entry.is_discounted == false ? entry.cost : entry.cost_brutto;
    const costNet = entry.is_discounted == false ? null : entry.cost;
    const cost = entry.is_discounted == false ? entry.cost : entry.cost_brutto;

    form.value = {
      id: entry.id,
      date: dayjs(entry.timestamp * 1000).format("YYYY-MM-DDTHH:mm"),
      isDiscounted: entry.is_discounted,
      cost: formatCurrency(cost, null),
      costGross: formatCurrency(costGross, null),
      costNet: formatCurrency(costNet, null),
      engineHours: entry.engine_hours,
      fuel: formatFuel(entry.liters),
      driver: entry.user_first_name,
      averageFuelPerHour: formatFuelConsumption(entry.avg_liters_per_hour),
    };
  }
}

function close() {
  emit("update:visible", false);
}

function save() {
  if (form.value.isDiscounted && form.value.costNet == null) {
    errors.value = ["Cost (net) cannot be empty if you select a discount."];
    return;
  }

  let cost = 0;
  let costBrutto = null;
  let isDiscounted = false;
  if (form.value.isDiscounted) {
    cost = form.value.costNet;
    costBrutto = form.value.costGross;
    isDiscounted = true;
  } else {
    cost = form.value.cost;
    costBrutto = null;
    isDiscounted = false;
  }

  const updatedEntry = {
    id: form.value.id,
    engine_hours: form.value.engineHours,
    liters: form.value.fuel,
    cost: cost,
    cost_brutto: costBrutto,
    is_discounted: isDiscounted,
  };

  store
    .dispatch("boat/updateFuelEntry", updatedEntry)
    .then(() => close())
    .catch((errs) => (errors.value = errs));
}

function remove() {
  store
    .dispatch("boat/removeFuelEntry", form.value.id)
    .then(() => close())
    .catch((errs) => (errors.value = errs));
}

setFormContent(props.fuelEntry);
store.dispatch("configuration/queryConfiguration");
</script>

<style scoped>
.block-footer {
  display: block;
}
</style>
