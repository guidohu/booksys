<template>
  <div>
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <engine-hour-entry-modal
      v-model:visible="showEntryHourModal"
      :engine-hour-entry="selectedEngineHourLogEntry"
      :display-format="getEngineHourFormat"
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
import dayjs from "dayjs";
import WarningBox from "booksys/components/WarningBox.vue";
import { formatEngineHour } from "booksys/libs/formatters";
import TableModule from "booksys/components/bricks/TableModule.vue";

const EngineHourEntryModal = defineAsyncComponent(
  () => import("booksys/components/EngineHourEntryModal.vue"),
);

const store = useStore();

const items = ref([]);
const columns = ref([]);
const errors = ref([]);
const selectedEngineHourLogEntry = ref(null);
const showEntryHourModal = ref(false);

const getEngineHourLog = computed(() => store.getters["boat/getEngineHourLog"]);
const getEngineHourFormat = computed(
  () => store.getters["configuration/getEngineHourFormat"],
);

watch(getEngineHourLog, (newEntries) => {
  console.log("getEngineHourLog just changed to", newEntries);
  setItems(newEntries);
});

watch(getEngineHourFormat, (newFormat, oldFormat) => {
  if (newFormat !== oldFormat) {
    setColumns();
  }
});

const queryEngineHourLog = () => store.dispatch("boat/queryEngineHourLog");
const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function setItems(logs) {
  items.value = [];
  logs.slice(0, 200).forEach((l) => {
    items.value.push(l);
  });
}

function setColumns() {
  columns.value = [
    {
      key: "timestamp",
      label: "Date",
      sortable: true,
      formatter: (value) => {
        return dayjs.unix(value).format("DD.MM.YYYY HH:mm");
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
      key: "before_hours",
      label: "Before",
      sortable: true,
      class: "text-end",
      formatter: (value) => {
        return formatEngineHour(value, getEngineHourFormat.value);
      },
    },
    {
      key: "after_hours",
      label: "After",
      sortable: true,
      class: "text-end",
      formatter: (value, key, item) => {
        if (parseFloat(item.after_hours) <= 0.001) {
          return "-";
        }
        return formatEngineHour(value, getEngineHourFormat.value);
      },
    },
    {
      key: "delta_hours",
      label: "Diff",
      sortable: true,
      class: "text-end",
      formatter: (value, key, item) => {
        if (
          parseFloat(item.after_hours) <= 0.001 &&
          parseFloat(value) <= 0.001
        ) {
          return "-";
        }
        return formatEngineHour(value, getEngineHourFormat.value);
      },
    },
  ];
}

function rowClick(item) {
  selectedEngineHourLogEntry.value = item;
  showEntryHourModal.value = true;
}

function rowClass(item) {
  // Type: 1 default session
  // Type: 2 course session
  if (item.type == 1) {
    return "clickable";
  } else {
    return "highlight clickable";
  }
}

queryConfiguration();
setColumns();
queryEngineHourLog()
  .then(() => {
    setItems(getEngineHourLog.value);
  })
  .catch((errs) => {
    errors.value = errs;
  });
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
