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

<script>
import { defineAsyncComponent } from "vue";
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
// import FuelEntryModal from "@/components/FuelEntryModal";
import { BooksysBrowser } from "@/libs/browser";
import remove from "lodash/remove";
import {
  formatEngineHour,
  formatCurrency,
  formatFuel,
  formatFuelConsumption,
} from "@/libs/formatters";
import * as dayjs from "dayjs";
import TableModule from "@/components/bricks/TableModule.vue";

const FuelEntryModal = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "fuel-entry-modal" */ "@/components/FuelEntryModal"
  )
);

export default {
  name: "FuelLogList",
  components: {
    WarningBox,
    FuelEntryModal,
    TableModule,
  },
  data() {
    return {
      errors: [],
      items: [],
      columns: [],
      selectedFuelEntry: null,
      showFuelEntryModal: false,
    };
  },
  computed: {
    ...mapGetters("boat", ["getFuelLog"]),
    ...mapGetters("configuration", ["getEngineHourFormat"]),
  },
  methods: {
    ...mapActions("boat", ["queryFuelLog"]),
    setItems: function (logs) {
      this.items = [];

      if (logs == null) {
        return;
      }

      logs.slice(0, 200).forEach((l) => {
        this.items.push(l);
      });
    },
    setColumns: function () {
      var columns = [
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
            return formatEngineHour(value, this.getEngineHourFormat);
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
            return this.getFuelCost(item);
          },
        },
      ];
      if (BooksysBrowser.isMobile()) {
        remove(columns, function (n, idx) {
          return idx == 4;
        });
      }
      this.columns = columns;
    },
    getFuelCost: function (entry) {
      // returns either net or gross values for the cost
      if (entry.cost != null) {
        return formatCurrency(Number(entry.cost), null);
      } else {
        return formatCurrency(Number(entry.cost_brutto), null);
      }
    },
    rowClick: function (item) {
      this.selectedFuelEntry = item;
      this.showFuelEntryModal = true;
    },
    rowClass: function (item) {
      // we highlight entries that have a deduction
      if (item.cost_brutto != null) {
        return "highlight clickable";
      } else {
        return "clickable";
      }
    },
  },
  watch: {
    getFuelLog: function (newValues) {
      this.setItems(newValues);
    },
    getEngineHourFormat: function (newFormat, oldFormat) {
      if (newFormat != oldFormat) {
        this.setColumns();
      }
    },
  },
  created() {
    this.queryFuelLog().catch((errors) => (this.errors = errors));

    this.setColumns();
  },
};
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
