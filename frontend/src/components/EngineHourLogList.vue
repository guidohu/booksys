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

<script>
import { defineAsyncComponent } from "vue";
import { mapActions, mapGetters } from "vuex";
import dayjs from "dayjs";
import WarningBox from "booksys/components/WarningBox.vue";
import { formatEngineHour } from "booksys/libs/formatters";
import TableModule from "booksys/components/bricks/TableModule.vue";

const EngineHourEntryModal = defineAsyncComponent(() =>
  import(
    "booksys/components/EngineHourEntryModal.vue"
  )
);

export default {
  name: "EngineHourLogList",
  components: {
    WarningBox,
    TableModule,
    EngineHourEntryModal,
  },
  data() {
    return {
      items: [],
      columns: [],
      errors: [],
      selectedEngineHourLogEntry: null,
      showEntryHourModal: false,
    };
  },
  computed: {
    ...mapGetters("boat", ["getEngineHourLog"]),
    ...mapGetters("configuration", ["getEngineHourFormat"]),
  },
  watch: {
    getEngineHourLog: function (newEntries) {
      console.log("getEngineHourLog just changed to", newEntries);
      this.setItems(newEntries);
    },
    getEngineHourFormat: function (newFormat, oldFormat) {
      if (newFormat != oldFormat) {
        this.setColumns();
      }
    },
  },
  methods: {
    ...mapActions("boat", ["queryEngineHourLog"]),
    ...mapActions("configuration", ["queryConfiguration"]),
    setItems: function (logs) {
      this.items = [];
      logs.slice(0, 200).forEach((l) => {
        this.items.push(l);
      });
    },
    setColumns: function () {
      this.columns = [
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
            return formatEngineHour(value, this.getEngineHourFormat);
          },
        },
        {
          key: "after_hours",
          label: "After",
          sortable: true,
          class: "text-end",
          formatter: (value, key, item) => {
            if (parseFloat(item.after_hours) <= 0.001) {
              return "-"
            }
            return formatEngineHour(value, this.getEngineHourFormat);
          },
        },
        {
          key: "delta_hours",
          label: "Diff",
          sortable: true,
          class: "text-end",
          formatter: (value, key, item) => {
            if (parseFloat(item.after_hours) <= 0.001 && parseFloat(value) <= 0.001) {
              return '-'
            }
            return formatEngineHour(value, this.getEngineHourFormat);
          },
        },
      ];
    },
    rowClick: function (item) {
      this.selectedEngineHourLogEntry = item;
      this.showEntryHourModal = true;
    },
    rowClass: function (item) {
      // Type: 1 default session
      // Type: 2 course session
      if (item.type == 1) {
        return "clickable";
      } else {
        return "highlight clickable";
      }
    },
  },
  created() {
    this.queryConfiguration();

    // generate header of table
    this.setColumns();

    // query content
    this.queryEngineHourLog()
      .then(() => {
        this.setItems(this.getEngineHourLog);
      })
      .catch((errors) => {
        this.errors = errors;
      });
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
