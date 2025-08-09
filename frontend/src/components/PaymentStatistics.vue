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
            <span class="lead">{{ formatNumber(getTotalPayments) }}</span>
            {{ getCurrency }}
            <br />
            Income
            <span v-if="form.selectedYear != 'any'"
              >({{ form.selectedYear }})</span
            >
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ formatNumber(getTotalExpenditures) }}</span>
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
            <span class="lead">{{
              formatNumber(getTotalSessionPayments)
            }}</span>
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
            <span class="lead">{{ formatNumber(getSessionsBalance) }}</span>
            {{ getCurrency }}
            <br />
            Sessions Credits
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ formatNumber(getBalance) }}</span>
            {{ getCurrency }}
            <br />
            Balance
          </card-module>
        </div>
        <div class="col-6 col-md-4">
          <card-module nobody class="mx-1 my-1 pt-3">
            <span class="lead">{{ formatNumber(getSessionProfit) }}</span>
            {{ getCurrency }}
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

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import reverse from "lodash/reverse";
import dayjs from "dayjs";
import WarningBox from "booksys/components/WarningBox.vue";
import CardModule from "booksys/components/bricks/CardModule.vue";
import InputSelect from "booksys/components/forms/inputs/InputSelect.vue";
import { formatNumber as formatNumberLib } from "booksys/libs/formatters.js";

const store = useStore();

const errors = ref([]);
const form = ref({
  years: [],
  selectedYear: "any",
});

const getYears = computed(() => store.getters["accounting/getYears"]);
const getBalance = computed(() => store.getters["accounting/getBalance"]);
const getTotalPayments = computed(
  () => store.getters["accounting/getTotalPayments"],
);
const getTotalExpenditures = computed(
  () => store.getters["accounting/getTotalExpenditures"],
);
const getTotalSessionPayments = computed(
  () => store.getters["accounting/getTotalSessionPayments"],
);
const getSessionsBalance = computed(
  () => store.getters["accounting/getSessionsBalance"],
);
const getSessionProfit = computed(
  () => store.getters["accounting/getSessionProfit"],
);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);

watch(getYears, (newValue) => {
  const availableYears = reverse(newValue);
  console.log(availableYears);
  form.value.years = availableYears.map((v) => {
    return {
      value: v,
      text: v,
    };
  });
});

const queryYears = () => store.dispatch("accounting/queryYears");
const queryStatistics = (year) =>
  store.dispatch("accounting/queryStatistics", year);
const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function yearSelectionChangeHandler() {
  queryStatistics(form.value.selectedYear).catch(
    (errs) => (errors.value = errs),
  );
}

function formatNumber(number) {
  return formatNumberLib(number);
}

queryConfiguration();

queryYears().catch((errs) => (errors.value = errs));

// get current year
const currentYear = dayjs().year();
form.value.selectedYear = currentYear;

queryStatistics(currentYear).catch((errs) => (errors.value = errs));
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
