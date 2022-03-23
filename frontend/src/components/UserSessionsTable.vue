<template>
  <table-module size="small" :columns="fields" :rows="items">
    <template #cell(riders)="data">
      <div class="row" v-for="rider in data.cell" :key="rider.id">
        <div class="col-12">{{ rider.first_name }} {{ rider.last_name }}</div>
      </div>
    </template>
    <template #cell(action)="data">
      <button
        type="button"
        class="btn btn-outline-danger btn-sm"
        @click="cancelSession(data.row.id)"
      >
        <i class="bi bi-trash"></i>
      </button>
    </template>
  </table-module>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";
import TableModule from "./bricks/TableModule.vue";

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

export default {
  name: "UserSessionsTable",
  components: {
    TableModule,
  },
  props: ["userSessions", "showCancel"],
  emits: ["cancel"],
  data: function () {
    return {
      fields: [
        {
          key: "start",
          label: "Date",
          formatter: (cell, key, row) => {
            return this.formatTime(row);
          },
        },
        {
          key: "riders",
          label: "Riders",
        },
      ],
      items: this.userSessions,
    };
  },
  computed: {
    ...mapGetters("configuration", ["getConfiguration", "getTimezone"]),
  },
  methods: {
    formatTime: function (item) {
      return (
        dayjs.unix(item.start).tz(this.getTimezone).format("DD.MM.YYYY HH:mm") +
        " - " +
        dayjs.unix(item.end).tz(this.getTimezone).format("HH:mm")
      );
    },
    cancelSession: function (sessionId) {
      console.log("Cancel session with id", sessionId);
      // forward to parent
      this.$emit("cancel", sessionId);
    },
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.queryConfiguration();
  },
  beforeMount() {
    if (this.$props.showCancel == true) {
      this.fields.push({
        key: "action",
        label: "Action",
      });
    }
  },
};
</script>
