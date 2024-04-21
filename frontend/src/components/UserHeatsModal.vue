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

<script>
import { sprintf } from "sprintf-js";
import { mapGetters, mapActions } from "vuex";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import TableModule from "./bricks/TableModule.vue";

export default {
  name: "UserHeatsModal",
  components: {
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    TableModule,
  },
  props: ["visible"],
  data: function () {
    return {
      fields: [
        {
          key: "date",
          label: "Date",
        },
        {
          key: "cost",
          label: "Cost",
          formatter: (value) => {
            const v = parseFloat(value);
            if(isNaN(v)) {
              return "NaN"
            }
            return this.formatCost(v);
          },
        },
        {
          key: "duration",
          label: "Duration",
          formatter: (value) => {
            return value;
          },
        },
      ],
    };
  },
  computed: {
    ...mapGetters("configuration", ["getConfiguration", "getCurrency"]),
    ...mapGetters("user", ["heatHistory"]),
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    close: function () {
      this.$emit("update:visible", false);
    },
    formatCost: function (value) {
      return sprintf("%.02f %s", value, this.getCurrency);
    },
  },
  created() {
    this.queryConfiguration();
  },
};
</script>

<style>
.heatTableHeight {
  overflow: auto;
  overflow-x: scroll;
  max-height: 300px;
}
</style>
