<template>
  <div class="text-left">
    <warning-box
      v-if="errors.length > 0"
      :errors="errors"
      dismissible="true"
      @dismissed="dismissedHandler"
    />
    <div v-if="errors.length == 0" class="box box-fix-height">
      <income-modal v-model:visible="showIncomeModal" />
      <expense-modal v-model:visible="showExpenseModal" />
      <div class="row box-fix-content mb-2 mx-1">
        <div class="col-6 text-left ps-1">
          <button
            class="btn btn-outline-info btn-sm me-1"
            @click="showAddIncome"
          >
            <i class="bi bi-plus" />
            Income
          </button>
          <button class="btn btn-outline-info btn-sm" @click="showAddExpense">
            <i class="bi bi-dash" />
            Expense
          </button>
        </div>
        <div
          class="col-3 col-sm-3 col-lg-2 offset-3 offset-sm-3 offset-lg-4 text-right pe-1"
        >
          <input-select
            id="year"
            v-model="form.selectedYear"
            size="small"
            :options="form.years"
            @changed="yearSelectionChangeHandler"
          />
        </div>
      </div>
      <div class="row box-flex-content">
        <div class="col-12">
          <overlay-spinner :active="isLoading">
            <table-module
              :columns="fields"
              :rows="getTransactions"
              size="small"
            >
              <template #cell(action)="item">
                <div class="text-center">
                  <button class="btn btn-light btn-sm" style="font-size: 0.8em">
                    <i
                      class="bi bi-trash text-alert"
                      @click="deleteEntry(item.row)"
                    />
                  </button>
                </div>
              </template>
            </table-module>
          </overlay-spinner>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineAsyncComponent, ref, computed, watch } from "vue";
import { useStore } from "vuex";
import reverse from "lodash/reverse";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import { formatCurrency } from "booksys/libs/formatters";
import WarningBox from "booksys/components/WarningBox.vue";
import TableModule from "./bricks/TableModule.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import OverlaySpinner from "./styling/OverlaySpinner.vue";
import { confirm } from "booksys/components/bricks/DialogModal.js";

const IncomeModal = defineAsyncComponent(
  () => import("booksys/components/IncomeModal.vue"),
);
const ExpenseModal = defineAsyncComponent(
  () => import("booksys/components/ExpenseModal.vue"),
);

dayjs.extend(utc);

const store = useStore();

const showExpenseModal = ref(false);
const showIncomeModal = ref(false);
const errors = ref([]);
const isLoading = ref(false);
const form = ref({
  years: [],
  selectedYear: "any",
});

const getTransactions = computed(
  () => store.getters["accounting/getTransactions"],
);
const getYears = computed(() => store.getters["accounting/getYears"]);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);

const fields = computed(() => [
  {
    key: "timestamp",
    label: "Date",
    sortable: true,
    formatter: (value, key, item) => dayjs.utc(value).format("DD.MM.YYYY"),
  },
  {
    key: "type_name",
    label: "Type",
    sortable: true,
  },
  {
    key: "name",
    label: "Name",
    sortable: true,
    formatter: (value, key, item) => item.fn + " " + item.ln,
  },
  {
    key: "amount",
    label: "Amount",
    sortable: true,
    formatter: (value) => formatCurrency(value, getCurrency.value),
    tdClass: "text-right",
    thClass: "text-right",
  },
  {
    key: "comment",
    label: "Comment",
    sortable: false,
  },
  {
    key: "action",
    label: "Action",
    sortable: false,
  },
]);

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

const queryTransactions = (year) =>
  store.dispatch("accounting/queryTransactions", year);
const queryYears = () => store.dispatch("accounting/queryYears");
const deleteTransaction = (transaction) =>
  store.dispatch("accounting/deleteTransaction", transaction);
const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function showAddIncome() {
  showIncomeModal.value = true;
}

function showAddExpense() {
  showExpenseModal.value = true;
}

function yearSelectionChangeHandler() {
  console.log("year selection has changed to:", form.value.selectedYear);
  isLoading.value = true;

  queryTransactions(form.value.selectedYear)
    .then(() => (errors.value = []))
    .catch((errs) => (errors.value = errs))
    .finally(() => (isLoading.value = false));
}

function deleteEntry(transaction) {
  console.log("Delete transaction:", transaction);
  confirm({
    title: "Delete Transaction",
    message: "Do you really want to delete this transaction?",
  }).then((value) => {
    if (value == true) {
      isLoading.value = true;
      deleteTransaction(transaction)
        .then(() => {
          errors.value = [];
        })
        .catch((errs) => {
          errors.value = errs;
        })
        .finally(() => (isLoading.value = false));
    }
  });
}

function dismissedHandler() {
  errors.value = [];
}

isLoading.value = true;
queryConfiguration().catch((errs) => (errors.value = errs));

queryYears().catch((errs) => (errors.value = errs));

const currentYear = dayjs().year();
form.value.selectedYear = currentYear;

queryTransactions(currentYear)
  .catch((errs) => (errors.value = errs))
  .finally(() => (isLoading.value = false));
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
