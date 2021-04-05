<template>
  <b-modal
    id="userHeatsModal"
    title="Heat History"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <b-row class="text-left">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <b-table
          class="heatTableHeight"
          sticky-header
          small
          striped
          hover
          responsive
          :items="heatHistory"
          :fields="fields"
        ></b-table>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"> </b-col>
    </b-row>
    <div slot="modal-footer">
      <b-button type="button" variant="outline-info" v-on:click="close">
        <b-icon-check />
        OK
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import { BModal, BRow, BCol, BTable, BButton, BIconCheck } from "bootstrap-vue";

export default {
  name: "UserHeatsModal",
  components: {
    BModal,
    BRow,
    BCol,
    BTable,
    BButton,
    BIconCheck,
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
            return this.formatCost(value);
          },
        },
        {
          key: "duration",
          label: "Duration",
          formatter: (value) => {
            return value + " min";
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
      console.log(value);
      console.log(this.$store);
      return value + " " + this.getCurrency;
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
