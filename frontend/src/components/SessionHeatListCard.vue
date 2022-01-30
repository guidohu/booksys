<template>
  <b-card no-body class="text-left">
    <b-card-body>
      <HeatEntryModal
        v-model:visible="showHeatEntryModal"
        :heat="selectedHeat"
      />
      <b-table
        v-if="errors.length == 0"
        striped
        hover
        small
        :fields="columns"
        :items="getHeatsForSession"
        tbody-tr-class="clickable"
        @row-clicked="rowClick"
      />
      <WarningBox v-if="errors.length > 0" :errors="errors" />
    </b-card-body>
  </b-card>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { sprintf } from "sprintf-js";
import WarningBox from "@/components/WarningBox";
import HeatEntryModal from "@/components/HeatEntryModal";
import { BCard, BCardBody, BTable } from "bootstrap-vue";

export default {
  name: "SessionHeatListCard",
  components: {
    WarningBox,
    HeatEntryModal,
    BCard,
    BCardBody,
    BTable,
  },
  props: ["sessionId"],
  data() {
    return {
      showHeatEntryModal: false,
      errors: [],
      columns: [],
      rows: [{ time: "12" }],
      selectedHeat: null,
    };
  },
  computed: {
    ...mapGetters("heats", ["getHeatsForSession"]),
    ...mapGetters("configuration", ["getCurrency"]),
  },
  watch: {
    getHeatsForSession: function (newHeats) {
      // depending whether we have comments or not, we display the column or not
      console.log("SessionHeatListCard, heats changed to", newHeats);
      this.setColumns();
    },
    getCurrency: function (newCurrency) {
      // if currency changed -> new formatter
      this.setColumns();
      console.log(newCurrency);
    },
  },
  methods: {
    ...mapActions("heats", ["queryHeatsForSession"]),
    ...mapActions("configuration", ["queryConfiguration"]),
    formatDuration: function (durationS) {
      const seconds = durationS % 60;
      const minutes = Math.floor(((durationS - seconds) % 3600) / 60);
      const hours = Math.floor(
        (durationS - seconds * 60 - minutes * 3600) / 3600
      );

      if (hours > 0) {
        return sprintf("%02d:%02d:%02d", hours, minutes, seconds);
      } else {
        return sprintf("%02d:%02d", minutes, seconds);
      }
    },
    setColumns: function () {
      this.columns = [
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
            return this.formatDuration(value);
          },
        },
        {
          key: "cost",
          label: "Cost",
          formatter: (value) => value + " " + this.getCurrency,
        },
      ];

      // if any of the entries has a comment, also show the comment column
      const heatsWithComments = this.getHeatsForSession.filter(
        (h) => h.comment != null && h.comment.length > 0
      );
      if (heatsWithComments.length > 0) {
        this.columns.push({
          key: "comment",
          label: "Comment",
        });
      }
    },
    rowClick: function (item) {
      console.log("clicked on item", item);
      this.selectedHeat = item;
      this.showHeatEntryModal = true;
    },
  },
  created() {
    // get heats
    if (this.sessionId != null) {
      this.queryConfiguration();
      this.queryHeatsForSession(this.sessionId)
        .then(() => this.setColumns())
        .catch((errors) => (this.errors = errors));
    } else {
      console.error(
        "Cannot query heats for session, as no sessionId is provided"
      );
    }
  },
};
</script>

<style>
tr.clickable {
  cursor: pointer;
}
</style>
