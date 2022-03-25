<template>
  <div class="text-begin">
    <warning-box v-if="errors.length > 0" :errors="errors" dismissible="true" />
    <div v-if="errors.length == 0" class="box box-fix-height">
      <div class="row box-fix-content mb-2 mx-1">
        <div class="col-3 col-lg-2 offset-9 offset-lg-10 text-right pe-1">
          <input-select
            v-if="form.years.length > 0"
            id="year"
            v-model="form.selectedYear"
            size="small"
            :options="form.years"
            @changed="yearSelectionChangeHandler"
          />
        </div>
      </div>
      <div class="row box-flex-content text-end">
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ getTotalPayments }}</span> {{ getCurrency }}
            <br />
            Income
            <span v-if="form.selectedYear != 'any'"
              >({{ form.selectedYear }})</span
            >
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ getTotalExpenditures }}</span>
            {{ getCurrency }}
            <br />
            Expenses
            <span v-if="form.selectedYear != 'any'"
              >({{ form.selectedYear }})</span
            >
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ getTotalSessionPayments }}</span>
            {{ getCurrency }}
            <br />
            Sessions Income
            <span v-if="form.selectedYear != 'any'"
              >({{ form.selectedYear }})</span
            >
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ getSessionsBalance }}</span> {{ getCurrency }}
            <br />
            Sessions Credits
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ getBalance }}</span> {{ getCurrency }}
            <br />
            Balance
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ getSessionProfit }}</span> {{ getCurrency }}
            <br />
            Session Profit
            <span v-if="form.selectedYear != 'any'"
              >({{ form.selectedYear }})</span
            >
          </card-module>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import reverse from "lodash/reverse";
import * as dayjs from "dayjs";
import WarningBox from "@/components/WarningBox";
import CardModule from "@/components/bricks/CardModule";
import InputSelect from "@/components/forms/inputs/InputSelect.vue";

export default {
  name: "PaymentDetails",
  components: {
    WarningBox,
    CardModule,
    InputSelect,
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

<style scoped>
.box {
  display: flex;
  flex-flow: column;
  height: 100%;
}

.box-fix-height {
  max-height: 430px;
}

@media (max-width: 992px) {
  .box-fix-height {
    max-height: 90vh;
  }
}

.box-fix-content {
  flex: 0 0 auto;
}

.box-flex-content {
  flex: 1 1 auto;
  overflow: scroll;
}
</style>
