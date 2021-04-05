<template>
  <div>
    <WarningBox v-if="errors.length > 0" :errors="errors" dismissible="true" />
    <b-row v-if="form.years.length > 0" class="text-right">
      <b-col cols="12">
        <b-form-group
          id="year"
          label="Year"
          label-for="year-select"
          description=""
          label-cols="8"
        >
          <b-form-select
            id="year-select"
            v-model="form.selectedYear"
            :options="form.years"
            @change="yearSelectionChangeHandler($event)"
          />
        </b-form-group>
      </b-col>
    </b-row>
    <b-row class="text-right">
      <b-col cols="6" md="4">
        <b-card border-variant="info" class="m-1">
          <span class="lead">{{ getTotalPayments }}</span> {{ getCurrency }}
          <br />
          Income
          <span v-if="form.selectedYear != 'any'"
            >({{ form.selectedYear }})</span
          >
        </b-card>
      </b-col>
      <b-col cols="6" md="4">
        <b-card border-variant="info" class="m-1">
          <span class="lead">{{ getTotalExpenditures }}</span> {{ getCurrency }}
          <br />
          Expenses
          <span v-if="form.selectedYear != 'any'"
            >({{ form.selectedYear }})</span
          >
        </b-card>
      </b-col>
      <b-col cols="6" md="4">
        <b-card border-variant="info" class="m-1">
          <span class="lead">{{ getTotalSessionPayments }}</span>
          {{ getCurrency }}
          <br />
          Sessions Income
          <span v-if="form.selectedYear != 'any'"
            >({{ form.selectedYear }})</span
          >
        </b-card>
      </b-col>
      <b-col cols="6" md="4">
        <b-card border-variant="info" class="m-1">
          <span class="lead">{{ getSessionsBalance }}</span> {{ getCurrency }}
          <br />
          Sessions Credits
        </b-card>
      </b-col>
      <b-col cols="6" md="4">
        <b-card border-variant="info" class="m-1">
          <span class="lead">{{ getBalance }}</span> {{ getCurrency }}
          <br />
          Balance
        </b-card>
      </b-col>
      <b-col cols="6" md="4">
        <b-card border-variant="info" class="m-1">
          <span class="lead">{{ getSessionProfit }}</span> {{ getCurrency }}
          <br />
          Session Profit
          <span v-if="form.selectedYear != 'any'"
            >({{ form.selectedYear }})</span
          >
        </b-card>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import reverse from "lodash/reverse";
import * as dayjs from "dayjs";
import WarningBox from "@/components/WarningBox";
import { BRow, BCol, BFormGroup, BFormSelect, BCard } from "bootstrap-vue";

export default {
  name: "PaymentDetails",
  components: {
    WarningBox,
    BRow,
    BCol,
    BFormGroup,
    BFormSelect,
    BCard,
  },
  data() {
    return {
      errors: [],
      form: {
        years: [],
        selectedYear: "any",
      },
    };
  },
  computed: {
    ...mapGetters("accounting", [
      "getYears",
      "getBalance",
      "getTotalPayments",
      "getTotalExpenditures",
      "getTotalSessionPayments",
      "getSessionsBalance",
      "getSessionProfit",
    ]),
    ...mapGetters("configuration", ["getCurrency"]),
  },
  watch: {
    getYears: function (newValue) {
      const availableYears = reverse(newValue);
      console.log(availableYears);
      this.form.years = availableYears.map((v) => {
        return {
          value: v,
          text: v,
        };
      });
    },
  },
  methods: {
    ...mapActions("accounting", ["queryYears", "queryStatistics"]),
    ...mapActions("configuration", ["queryConfiguration"]),
    yearSelectionChangeHandler: function (selection) {
      this.queryStatistics(selection).catch((errors) => (this.errors = errors));
    },
  },
  created() {
    this.queryConfiguration();

    this.queryYears().catch((errors) => (this.errors = errors));

    // get current year
    const currentYear = dayjs().year();
    this.form.selectedYear = currentYear;

    this.queryStatistics(currentYear).catch((errors) => (this.errors = errors));
  },
};
</script>
