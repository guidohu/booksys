<template>
  <card-module :nobody="true">
    <overlay-spinner :active="showOverlay" :full-page="false">
      <warning-box v-if="errors.length > 0" :errors="errors" />
      <table-module :columns="columns" :rows="getLogLines" />
    </overlay-spinner>
  </card-module>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import CardModule from "@/components/bricks/CardModule.vue";
import TableModule from "./bricks/TableModule.vue";
import WarningBox from "@/components/WarningBox";
import OverlaySpinner from "@/components/styling/OverlaySpinner.vue";

export default {
  name: "LogCard",
  components: {
    CardModule,
    WarningBox,
    TableModule,
    OverlaySpinner,
  },
  data() {
    return {
      showOverlay: false,
      errors: [],
      columns: [
        {
          key: "time",
          label: "Time",
        },
        {
          key: "log",
          label: "Message",
        },
      ],
      rows: [],
    };
  },
  computed: {
    ...mapGetters("log", ["getLogLines"]),
  },
  methods: {
    ...mapActions("log", ["queryLogLines"]),
  },
  created() {
    this.showOverlay = true;
    this.queryLogLines()
      .then(() => (this.showOverlay = false))
      .catch((errors) => {
        this.showOverlay = false;
        this.errors = errors;
      });
  },
};
</script>
