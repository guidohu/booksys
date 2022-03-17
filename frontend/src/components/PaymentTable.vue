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

<script>
import { defineAsyncComponent } from "vue";
import { mapGetters, mapActions } from "vuex";
import reverse from "lodash/reverse";
import * as dayjs from "dayjs";
import { formatCurrency } from "@/libs/formatters";
import WarningBox from "@/components/WarningBox";
import TableModule from "./bricks/TableModule.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import OverlaySpinner from "./styling/OverlaySpinner.vue";
import { confirm } from "@/components/bricks/DialogModal";

const IncomeModal = defineAsyncComponent(() =>
  import(/* webpackChunkName: "income-modal" */ "@/components/IncomeModal")
);
const ExpenseModal = defineAsyncComponent(() =>
  import(/* webpackChunkName: "expense-modal" */ "@/components/ExpenseModal")
);

export default {
  name: "PaymentTable",
  components: {
    WarningBox,
    TableModule,
    InputSelect,
    OverlaySpinner,
    IncomeModal,
    ExpenseModal,
  },
  data() {
    return {
      showExpenseModal: false,
      showIncomeModal: false,
      errors: [],
      isLoading: false,
      form: {
        years: [],
        selectedYear: "any",
      },
      items: [],
      selectedItems: [],
      fields: [
        {
          key: "timestamp",
          label: "Date",
          sortable: true,
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
          formatter: (value) => formatCurrency(value, this.getCurrency),
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
      ],
    };
  },
  computed: {
    ...mapGetters("accounting", ["getTransactions", "getYears"]),
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
    ...mapActions("accounting", [
      "queryTransactions",
      "queryYears",
      "deleteTransaction",
    ]),
    ...mapActions("configuration", ["queryConfiguration"]),
    showAddIncome: function () {
      this.showIncomeModal = true;
    },
    showAddExpense: function () {
      this.showExpenseModal = true;
    },
    yearSelectionChangeHandler: function () {
      console.log("year selection has changed to:", this.form.selectedYear);
      this.isLoading = true;

      this.queryTransactions(this.form.selectedYear)
        .then(() => (this.errors = []))
        .catch((errors) => (this.errors = errors))
        .then(() => (this.isLoading = false));
    },
    deleteEntry: function (transaction) {
      console.log("Delete transaction:", transaction);
      confirm({
        title: "Delete Transaction",
        message: "Do you really want to delete this transaction?",
      }).then((value) => {
        // delete transaction
        if (value == true) {
          this.isLoading = true;
          this.deleteTransaction(transaction)
            .then(() => {
              this.errors = [];
            })
            .catch((errors) => {
              this.errors = errors;
            })
            .then(() => (this.isLoading = false));
        }
      });
    },
    dismissedHandler: function () {
      this.errors = [];
    },
  },
  created() {
    this.isLoading = true;
    this.queryConfiguration().catch((errors) => (this.errors = errors));

    this.queryYears().catch((errors) => (this.errors = errors));

    const currentYear = dayjs().year();
    this.form.selectedYear = currentYear;

    this.queryTransactions(currentYear)
      .then(() => (this.isLoading = false))
      .catch((errors) => (this.errors = errors));
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
