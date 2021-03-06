<template>
  <div class="text-left">
    <WarningBox
      v-if="errors.length > 0"
      :errors="errors"
      dismissible="true"
      @dismissed="dismissedHandler"
    />
    <div v-else>
      <IncomeModal :visible.sync="showIncomeModal" />
      <ExpenseModal :visible.sync="showExpenseModal" />
      <b-row v-if="form.years.length > 0" class="text-right">
        <b-col offset="6" cols="6" class="d-sm-none mr-1">
          <b-form-group
            id="year"
            label="Year"
            label-for="year-select"
            description=""
            label-cols="3"
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
      <b-row>
        <b-col cols="8">
          <b-button
            v-on:click="showAddIncome"
            size="sm"
            variant="outline-info"
            class="mr-1 mb-2"
            ><b-icon-plus />Income</b-button
          >
          <b-button
            v-on:click="showAddExpense"
            size="sm"
            variant="outline-info"
            class="mb-2"
            ><b-icon-dash />Expense</b-button
          >
        </b-col>
        <b-col cols="4" class="d-none d-sm-block">
          <b-form-group
            id="year"
            label="Year"
            label-for="year-select"
            description=""
            label-cols="3"
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
      <b-row>
        <b-col cols="12">
          <b-overlay
            id="overlay-background"
            :show="isLoading"
            spinner-type="border"
            spinner-variant="info"
            rounded="sm"
          >
            <b-table
              striped
              hover
              small
              sticky-header
              borderless
              sort-by="date"
              :items="getTransactions"
              :fields="fields"
              :selectable="false"
              table-class="payment-table"
              tbody-class="payment-table"
              show-empty
              empty-text="No records to show"
            >
              <template #cell(action)="data">
                <div class="text-center">
                  <b-button size="sm" style="font-size: 0.8em" variant="light">
                    <b-icon-trash
                      v-on:click="deleteEntry(data.item)"
                      variant="danger"
                    />
                  </b-button>
                </div>
              </template>
            </b-table>
          </b-overlay>
        </b-col>
      </b-row>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import { mapGetters, mapActions } from "vuex";
import reverse from "lodash/reverse";
import * as dayjs from "dayjs";
import { formatCurrency } from "@/libs/formatters";
import WarningBox from "@/components/WarningBox";
import {
  BRow,
  BCol,
  BFormGroup,
  BFormSelect,
  BButton,
  BIconPlus,
  BIconDash,
  BIconTrash,
  BOverlay,
  BTable,
  ModalPlugin,
} from "bootstrap-vue";

const IncomeModal = () =>
  import(/* webpackChunkName: "income-modal" */ "@/components/IncomeModal");
const ExpenseModal = () =>
  import(/* webpackChunkName: "expense-modal" */ "@/components/ExpenseModal");

Vue.use(ModalPlugin);

export default {
  name: "PaymentTable",
  components: {
    WarningBox,
    IncomeModal,
    ExpenseModal,
    BRow,
    BCol,
    BFormGroup,
    BFormSelect,
    BButton,
    BIconPlus,
    BIconDash,
    BIconTrash,
    BOverlay,
    BTable,
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
    yearSelectionChangeHandler: function (selection) {
      this.isLoading = true;

      this.queryTransactions(selection)
        .then(() => (this.errors = []))
        .catch((errors) => (this.errors = errors))
        .then(() => (this.isLoading = false));
    },
    deleteEntry: function (transaction) {
      const message = "Do you really want to delete this transaction?";
      this.$bvModal
        .msgBoxConfirm(message, {
          title: "Delete Transaction",
          size: "sm",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "Delete",
          cancelTitle: "Cancel",
          footerClass: "p-2",
          hideHeaderClose: false,
          centered: true,
        })
        .then((value) => {
          // delete transaction
          if (value == true) {
            this.isLoading = true;
            this.deleteTransaction(transaction)
              .then(() => {
                this.errors = [];
              })
              .catch((errors) => (this.errors = errors));
          }
        })
        .catch((err) => {
          this.errors = [err];
        })
        .then(() => (this.isLoading = false));
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
@media (min-width: 450px) {
  .b-table-sticky-header {
    min-height: 330px;
    max-height: 330px;
    height: 330px;
  }
}

@media (max-width: 450px) {
  .b-table-sticky-header {
    min-height: 540px;
    max-height: 540px;
    height: 540px;
  }
}
</style>
