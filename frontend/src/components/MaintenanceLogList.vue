<template>
  <div>
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <table-module
      v-if="errors.length == 0"
      :rows="items"
      :columns="columns"
      size="small"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import dayjs from "dayjs";
import WarningBox from "booksys/components/WarningBox.vue";
import { formatEngineHour } from "booksys/libs/formatters";
import TableModule from "booksys/components/bricks/TableModule.vue";

const store = useStore();

const items = ref([]);
const columns = ref([]);
const errors = ref([]);

const getMaintenanceLog = computed(() => store.getters["boat/getMaintenanceLog"]);
const getEngineHourFormat = computed(
  () => store.getters["configuration/getEngineHourFormat"]
);

watch(getMaintenanceLog, (newEntries) => {
  setItems(newEntries);
});

watch(getEngineHourFormat, (newFormat, oldFormat) => {
  if (newFormat != oldFormat) {
    setColumns();
  }
});

const queryMaintenanceLog = () => store.dispatch("boat/queryMaintenanceLog");

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
        return dayjs(value * 1000).format("DD.MM.YYYY HH:mm");
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
      key: "user_first_name",
      label: "Driver",
      sortable: true,
    },
    {
      key: "description",
      label: "Maintenance Work",
      sortable: false,
    },
  ];
}

// generate header of table
setColumns();

// query content
queryMaintenanceLog()
  .then(() => {
    setItems(getMaintenanceLog.value);
  })
  .catch((errs) => {
    errors.value = errs;
  });
</script>
