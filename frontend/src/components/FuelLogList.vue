<template>
  <div>
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <fuel-entry-modal
      v-model:visible="showFuelEntryModal"
      :fuel-entry="selectedFuelEntry"
    />
    <table-module
      v-if="errors.length == 0"
      :rows="items"
      :columns="columns"
      :rowClassFunction="rowClass"
      size="small"
      @rowclick="rowClick"
    />
  </div>
</template>

<script setup>
import { defineAsyncComponent, ref, computed, watch } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import { BooksysBrowser } from "booksys/libs/browser";
import remove from "lodash/remove";
import {
  formatEngineHour,
  formatCurrency,
  formatFuel,
  formatFuelConsumption,
} from "booksys/libs/formatters";
import dayjs from "dayjs";
import TableModule from "booksys/components/bricks/TableModule.vue";

const FuelEntryModal = defineAsyncComponent(() =>
  import("booksys/components/FuelEntryModal.vue")
);

const store = useStore();

const errors = ref([]);
const items = ref([]);
const columns = ref([]);
const selectedFuelEntry = ref(null);
const showFuelEntryModal = ref(false);

const getFuelLog = computed(() => store.getters["boat/getFuelLog"]);
const getEngineHourFormat = computed(() => store.getters["configuration/getEngineHourFormat"]);

const queryFuelLog = () => store.dispatch("boat/queryFuelLog");

function setItems(logs) {
  items.value = [];

  if (logs == null) {
    return;
  }

  logs.slice(0, 200).forEach((l) => {
    items.value.push(l);
  });
}

function setColumns() {
  var cols = [
    {
      key: "timestamp",
      label: "Date",
      sortable: true,
      formatter: (value) => {
        return dayjs(value * 1000).format("DD.MM.YYYY HH:mm");
      },
    },
    {
      key: "user_first_name",
      label: "Driver",
      sortable: true,
      formatter: (value, key, item) => {
        return item.user_first_name;
      },
    },
    {
      key: "engine_hours",
      label: "EngineHrs",
      sortable: true,
      class: "text-end",
      formatter: (value) => {
        return formatEngineHour(value, getEngineHourFormat.value);
      },
    },
    {
      key: "liters",
      label: "Fuel",
      sortable: true,
      class: "text-end",
      formatter: (value) => formatFuel(value),
    },
    {
      key: "avg_liters_per_hour",
      label: "L/hr",
      sortable: true,
      class: "text-end",
      formatter: (value) => formatFuelConsumption(value),
    },
    {
      key: "cost",
      label: "Cost",
      sortable: true,
      class: "text-end",
      formatter: (value, key, item) => {
        return getFuelCost(item);
      },
    },
  ];
  if (BooksysBrowser.isMobileResponsive()) {
    remove(cols, function (n, idx) {
      return idx == 4;
    });
  }
  columns.value = cols;
}

function getFuelCost(entry) {
  // returns either net or gross values for the cost
  if (entry.cost != null) {
    return formatCurrency(Number(entry.cost), null);
  } else {
    return formatCurrency(Number(entry.cost_brutto), null);
  }
}

function rowClick(item) {
  selectedFuelEntry.value = item;
  showFuelEntryModal.value = true;
}

function rowClass(item) {
  // we highlight entries that have a deduction
  if (item.is_discounted == true) {
    return "highlight clickable";
  } else {
    return "clickable";
  }
}

watch(getFuelLog, (newValues) => {
  setItems(newValues);
});

watch(getEngineHourFormat, (newFormat, oldFormat) => {
  if (newFormat != oldFormat) {
    setColumns();
  }
});

queryFuelLog().catch((errs) => (errors.value = errs));

setColumns();
</script>

<style>
tr.highlight {
  color: #ffffff;
  background-color: rgb(91, 192, 222);
}

tr.clickable {
  cursor: pointer;
}
</style>
