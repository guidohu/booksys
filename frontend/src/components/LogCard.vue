<template>
  <card-module :nobody="true">
    <overlay-spinner :active="showOverlay" :full-page="false">
      <warning-box v-if="errors.length > 0" :errors="errors" />
      <table-module :columns="columns" :rows="getLogLines" />
    </overlay-spinner>
  </card-module>
</template>

<script setup>
import { ref, computed } from "vue";
import { useStore } from "vuex";
import CardModule from "booksys/components/bricks/CardModule.vue";
import TableModule from "./bricks/TableModule.vue";
import WarningBox from "booksys/components/WarningBox.vue";
import OverlaySpinner from "booksys/components/styling/OverlaySpinner.vue";

const store = useStore();

const showOverlay = ref(false);
const errors = ref([]);
const columns = ref([
  {
    key: "time",
    label: "Time",
  },
  {
    key: "log",
    label: "Message",
  },
]);

const getLogLines = computed(() => store.getters["log/getLogLines"]);

const queryLogLines = () => store.dispatch("log/queryLogLines");

showOverlay.value = true;
queryLogLines()
  .then(() => (showOverlay.value = false))
  .catch((errs) => {
    showOverlay.value = false;
    errors.value = errs;
  });
</script>
