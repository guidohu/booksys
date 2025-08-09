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

<script setup>
import { ref, computed, onBeforeMount } from "vue";
import { useStore } from "vuex";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";
import TableModule from "./bricks/TableModule.vue";

dayjs.extend(utc);
dayjs.extend(timezone);

const props = defineProps(["userSessions", "showCancel"]);
const emit = defineEmits(["cancel"]);

const store = useStore();

const fields = ref([
  {
    key: "start",
    label: "Date",
    formatter: (cell, key, row) => {
      return formatTime(row);
    },
  },
  {
    key: "riders",
    label: "Riders",
  },
]);

const items = ref(props.userSessions);

const getTimezone = computed(() => store.getters["configuration/getTimezone"]);

const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function formatTime(item) {
  return (
    dayjs
      .unix(item.start_time)
      .tz(getTimezone.value)
      .format("DD.MM.YYYY HH:mm") +
    " - " +
    dayjs.unix(item.end_time).tz(getTimezone.value).format("HH:mm")
  );
}

function cancelSession(sessionId) {
  console.log("Cancel session with id", sessionId);
  emit("cancel", sessionId);
}

queryConfiguration();

onBeforeMount(() => {
  if (props.showCancel === true) {
    fields.value.push({
      key: "action",
      label: "Action",
    });
  }
});
</script>
