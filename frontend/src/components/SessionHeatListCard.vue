<template>
  <sectioned-card-module title="Heats">
    <template v-slot:body>
      <warning-box v-if="errors.length > 0" :errors="errors" />
      <heat-entry-modal
        v-model:visible="showHeatEntryModal"
        :heat="selectedHeat"
      />
      <table-module
        v-if="errors.length == 0"
        :columns="columns"
        :rows="getHeatsForSession"
        row-class="clickable"
        size="small"
        @rowclick="rowClick"
      />
    </template>
  </sectioned-card-module>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "../components/WarningBox.vue";
import HeatEntryModal from "../components/HeatEntryModal.vue";
import SectionedCardModule from "./bricks/SectionedCardModule.vue";
import TableModule from "./bricks/TableModule.vue";

const props = defineProps(["sessionId"]);

const store = useStore();

const showHeatEntryModal = ref(false);
const errors = ref([]);
const columns = ref([]);
const rows = ref([{ time: "12" }]);
const selectedHeat = ref(null);

const getHeatsForSession = computed(
  () => store.getters["heats/getHeatsForSession"],
);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);

watch(getHeatsForSession, (newHeats) => {
  // depending whether we have comments or not, we display the column or not
  console.log("SessionHeatListCard, heats changed to", newHeats);
  setColumns();
});

watch(getCurrency, (newCurrency) => {
  // if currency changed -> new formatter
  setColumns();
  console.log(newCurrency);
});

const queryHeatsForSession = (sessionId) =>
  store.dispatch("heats/queryHeatsForSession", sessionId);
const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function formatDuration(durationS) {
  const seconds = durationS % 60;
  const minutes = Math.floor(((durationS - seconds) % 3600) / 60);
  const hours = Math.floor((durationS - seconds * 60 - minutes * 3600) / 3600);

  if (hours > 0) {
    return sprintf("%02d:%02d:%02d", hours, minutes, seconds);
  } else {
    return sprintf("%02d:%02d", minutes, seconds);
  }
}

function formatCost(cost) {
  return sprintf("%.2f", cost) + " " + getCurrency.value;
}


function setColumns() {
  columns.value = [
    {
      key: "first_name",
      label: "Name",
      formatter: (value, key, item) => {
        return item.first_name + " " + item.last_name.substring(0, 1) + ".";
      },
    },
    {
      key: "duration_s",
      label: "Duration",
      formatter: (value) => {
        return formatDuration(value);
      },
    },
    {
      key: "cost",
      label: "Cost",
      formatter: (value) => {
        return formatCost(value);
      },
    },
  ];

  // if any of the entries has a comment, also show the comment column
  const heatsWithComments = getHeatsForSession.value.filter(
    (h) => h.comment != null && h.comment.length > 0,
  );
  if (heatsWithComments.length > 0) {
    columns.value.push({
      key: "comment",
      label: "Comment",
    });
  }
  console.log("Columns:", columns.value);
}

function rowClick(item) {
  console.log("clicked on item", item);
  selectedHeat.value = item;
  showHeatEntryModal.value = true;
}

// get heats
if (props.sessionId != null) {
  queryConfiguration();
  queryHeatsForSession(props.sessionId)
    .then(() => setColumns())
    .catch((errs) => (errors.value = errs));
} else {
  console.error("Cannot query heats for session, as no sessionId is provided");
}
</script>

<style>
tr.clickable {
  cursor: pointer;
}

.scrollable {
  overflow-y: scroll;
}
</style>
