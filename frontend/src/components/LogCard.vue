<template>
  <b-card
    no-body
    class="text-left"
  >
    <b-card-body>
      <b-overlay
        id="overlay-background"
        :show="showOverlay"
        spinner-type="border"
        spinner-variant="info"
        rounded="sm"
      >
        <b-table
          v-if="errors.length == 0"
          striped
          hover
          small
          :fields="columns"
          :items="getLogLines"
          empty-text="no log entries to display"
        />
        <WarningBox
          v-if="errors.length > 0"
          :errors="errors"
        />
      </b-overlay>
    </b-card-body>
  </b-card>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
import { BCard, BCardBody, BTable, BOverlay } from "bootstrap-vue";

export default {
  name: "LogCard",
  components: {
    WarningBox,
    BCard,
    BCardBody,
    BTable,
    BOverlay,
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
