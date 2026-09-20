<template>
  <div class="text-begin">
    <warning-box v-if="errors.length > 0" :errors="errors" dismissible="true" />
    <div v-if="errors.length == 0" class="box box-fix-height">
      <div class="box-fix-content stats-header">
        <span class="stats-title">Overview</span>
        <div class="stats-filter">
          <label for="year" class="stats-filter-label">
            <i class="bi bi-calendar3"></i>
            Year
          </label>
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
      <div class="box-flex-content stats-plane">
        <div class="row g-2">
          <div class="col-12 col-md-6">
            <stat-tile
              hero
              tone-value
              label="Balance"
              icon="bi-bank"
              :value="formatNumber(getBalance)"
              :unit="getCurrency"
              :tone="signTone(getBalance)"
              context="All time"
            />
          </div>
          <div class="col-12 col-md-6">
            <stat-tile
              label="Sessions Credits"
              icon="bi-wallet2"
              tone="neutral"
              :value="formatNumber(getSessionsBalance)"
              :unit="getCurrency"
              context="Outstanding, all time"
            />
          </div>
          <div class="col-6 col-md-3">
            <stat-tile
              label="Income"
              icon="bi-arrow-down-left-circle"
              tone="positive"
              :value="formatNumber(getTotalPayments)"
              :unit="getCurrency"
              :context="periodLabel"
            />
          </div>
          <div class="col-6 col-md-3">
            <stat-tile
              label="Expenses"
              icon="bi-arrow-up-right-circle"
              tone="negative"
              :value="formatNumber(getTotalExpenditures)"
              :unit="getCurrency"
              :context="periodLabel"
            />
          </div>
          <div class="col-6 col-md-3">
            <stat-tile
              label="Sessions Income"
              icon="bi-cash-coin"
              tone="positive"
              :value="formatNumber(getTotalSessionPayments)"
              :unit="getCurrency"
              :context="periodLabel"
            />
          </div>
          <div class="col-6 col-md-3">
            <stat-tile
              tone-value
              label="Session Profit"
              :icon="
                isNegative(getSessionProfit)
                  ? 'bi-graph-down-arrow'
                  : 'bi-graph-up-arrow'
              "
              :value="formatNumber(getSessionProfit)"
              :unit="getCurrency"
              :tone="signTone(getSessionProfit)"
              :context="periodLabel"
            />
          </div>
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
import StatTile from "booksys/components/bricks/StatTile.vue";
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

// The period the year-bound figures cover, shown on those tiles.
const periodLabel = computed(() =>
  form.value.selectedYear == "any"
    ? "All time"
    : String(form.value.selectedYear),
);

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

function isNegative(number) {
  return Number(number) < 0;
}

function signTone(number) {
  return isNegative(number) ? "negative" : "positive";
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
  overflow-x: hidden;
  overflow-y: auto;
}

/* Filter row, sits above the tiles */
.stats-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  padding: 0 0.25rem 0.6rem;
}

.stats-title {
  color: #6c757d;
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.stats-filter {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.stats-filter-label {
  margin-bottom: 0;
  color: #6c757d;
  font-size: 0.75rem;
  white-space: nowrap;
}

/* Faint plane so the white tiles read as cards sitting on the dashboard */
.stats-plane {
  padding: 0.6rem;
  background-color: #f5f5f4;
  border: 1px solid rgba(11, 11, 11, 0.07);
  border-radius: 0.6rem;
}
</style>
