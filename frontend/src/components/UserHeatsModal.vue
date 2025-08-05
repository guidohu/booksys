<template>
  <modal-container name="user-heats-modal" :visible="visible">
    <modal-header
      :closable="true"
      title="Heat History"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <div class="heatTableHeight">
        <table-module size="small" :columns="fields" :rows="heatHistory" />
      </div>
    </modal-body>
    <modal-footer>
      <button type="submit" class="btn btn-outline-info" @click="close">
        <i class="bi bi-check"></i>
        OK
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { computed } from "vue";
import { useStore } from "vuex";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import TableModule from "./bricks/TableModule.vue";
import { formatCost as formatCostFormatter } from "booksys/libs/formatters";

defineProps(["visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const heatHistory = computed(() => store.getters["user/heatHistory"]);

const fields = computed(() => [
  {
    key: "date",
    label: "Date",
  },
  {
    key: "cost",
    label: "Cost",
    formatter: (value) => formatCost(value),
  },
  {
    key: "duration",
    label: "Duration",
    formatter: (value) => value,
  },
]);

const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function close() {
  emit("update:visible", false);
}

function formatCost(value) {
  return formatCostFormatter(value, getCurrency.value);
}

queryConfiguration();
</script>

<style>
.heatTableHeight {
  overflow: auto;
  overflow-x: scroll;
  max-height: 300px;
}
</style>
