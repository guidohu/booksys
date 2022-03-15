<template>
  <modal-container
    name="expense-modal"
    :visible="visible"
  >
    <modal-header
      :closable="true"
      title="Add New Expense"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length" :errors="errors" />
      <form @submit.prevent="save">
        <input-select
          id="type-select"
          label="Type"
          v-model="form.type"
          :options="expenseTypes"
          :description="typeDescription"
          size="small"
        />
        <input-date-time-local
          id="date"
          label="Date"
          v-model="form.date"
          size="small"
        />
        <input-select
          v-if="form.type != null"
          id="user-select"
          :label="userLabel"
          v-model="form.user"
          :options="users"
          :description="userDescription"
          size="small"
        />
        <input-currency
          v-if="form.type != null"
          id="amount"
          label="Cost"
          size="small"
          v-model="form.amount"
          placeholder="0.00"
          :currency="getCurrency"
        />
        <input-fuel
          v-if="form.type != null && form.type == 0"
          id="fuel-liters"
          label="Fuel"
          size="small"
          placeholder="0.00"
          v-model="form.fuelLiters"
        />
        <input-engine-hours
          v-if="form.type != null && form.type == 0"
          id="engine-hours"
          label="Engine Hours"
          size="small"
          placeholder="0"
          v-model="form.engineHours"
          :display-format="getEngineHourFormat"
        />
        <input-text-multiline
          v-if="form.type != null && form.type != 0"
          id="description"
          label="Description"
          rows="2"
          v-model="form.description"
        />
        <div class="alert alert-warning" v-if="form.type == 0 && getFuelPaymentType == 'billed'">
          Entry will not be visible in the payment section as it is not a direct expense. According to your settings fuel consumption is 'billed' and thus only 'fuel bill' has an impact on the balance.
        </div>
      </form>
    </modal-body>
    <modal-footer>
      <button
        type="button"
        class="btn btn-outline-info mr-1"
        @click="add"
      >
        <i class="bi bi-check"/>
        Add
      </button>
      <button
        type="button"
        class="btn btn-outline-danger"
        @click="close"
      >
        <i class="bi bi-x"/>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import * as dayjs from "dayjs";
import orderBy from "lodash/orderBy";
import ModalContainer from './bricks/ModalContainer.vue';
import ModalHeader from './bricks/ModalHeader.vue';
import ModalBody from './bricks/ModalBody.vue';
import ModalFooter from './bricks/ModalFooter.vue';
import InputSelect from './forms/inputs/InputSelect.vue';
import InputTextMultiline from './forms/inputs/InputTextMultiline.vue';
import InputCurrency from './forms/inputs/InputCurrency.vue';
import InputDateTimeLocal from './forms/inputs/InputDateTimeLocal.vue';
import InputFuel from './forms/inputs/InputFuel.vue';
import InputEngineHours from './forms/inputs/InputEngineHours.vue';

export default {
  name: "ExpenseModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputSelect,
    InputCurrency,
    InputTextMultiline,
    InputDateTimeLocal,
    InputFuel,
    InputEngineHours,
  },
  props: ["visible"],
  data() {
    return {
      isLoading: false,
      errors: [],
      expenseTypes: [],
      users: [],
      userLabel: "User",
      userDescription: "",
      typeDescription: "",
      form: {
        type: null,
        user: null,
        date: null,
      },
    };
  },
  computed: {
    ...mapGetters("configuration", ["getCurrency", "getEngineHourFormat", "getFuelPaymentType"]),
    ...mapGetters("user", ["userList"]),
    ...mapGetters("accounting", ["getExpenseTypes"]),
  },
  watch: {
    getExpenseTypes: function (newValues) {
      this.buildTypeSelect(newValues);
    },
    userList: function (newValues) {
      // console.log("TODO, add userList:", newValues);
      this.buildUserSelect(newValues);
    },
    "form.type": function (newValue, oldValue) {
      console.log(newValue, oldValue);
      if (newValue != oldValue) {
        // if the type is either
        if ([0].includes(Number(newValue))) {
          // 0: fuel
          this.userLabel = "Driver";
          this.userDescription = "Driver that fueled the boat";
          this.typeDescription = "Add a fuel entry.";
        } else if ([9].includes(Number(newValue))) {
          // 9: fuel bill
          this.userLabel = "Payer";
          this.userDescription = "Pays the bill";
          this.typeDescription = "Fuel bill";
        } else if ([3].includes(Number(newValue))) {
          // 3: invest
          this.userLabel = "Payee";
          this.userDescription = "Payee or internal reference.";
          this.typeDescription = "";
        } else if ([1].includes(Number(newValue))) {
          // 1: maintenance
          this.userLabel = "Payee";
          this.userDescription = "Payee or internal reference.";
          this.typeDescription = "Expenses related to maintenance";
        } else if ([2].includes(Number(newValue))) {
          // 2: material
          this.userLabel = "Payee";
          this.userDescription = "Payee or internal reference.";
          this.typeDescription = "Expenses related to aquired parts.";
        } else if ([6, 4].includes(Number(newValue))) {
          // 6: membership fee
          // 4: session
          this.userLabel = "Account";
          this.userDescription = "Refund to this account.";
          this.typeDescription = "Refund of payments.";
        } else if ([5].includes(Number(newValue))) {
          // 5: other
          this.userLabel = "Payee";
          this.userDescription = "Payee or internal reference";
          this.typeDescription =
            "everything that does not fit another category";
        } else if ([7].includes(Number(newValue))) {
          // 7: salary
          this.userLabel = "Driver";
          this.userDescription = "Driver that we pay a compensation.";
          this.typeDescription = "Compensation payments";
        } else if ([8].includes(Number(newValue))) {
          // 8: owners refund
          this.userLabel = "Payee";
          this.userDescription = "Owner that gets a refund.";
          this.typeDescription = "";
        } else {
          // Default case
          this.userLabel = "Payee";
          this.userDescription = "";
          this.typeDescription = "";
        }
      }
    },
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    ...mapActions("user", ["queryUserList"]),
    ...mapActions("accounting", [
      "queryTransactions",
      "queryStatistics",
      "queryExpenseTypes",
      "addExpense",
    ]),
    ...mapActions("boat", ["addFuelEntry"]),
    buildTypeSelect: function (types) {
      this.expenseTypes = types.map((t) => {
        return {
          value: t.id,
          text: t.name,
        };
      });
      this.expenseTypes = orderBy(this.expenseTypes, ["text"], ["asc"]);
      this.expenseTypes.unshift({ value: null, text: "Please select" });
    },
    buildUserSelect: function (users) {
      this.users = users.map((u) => {
        return {
          value: u.id,
          text: u.firstName + " " + u.lastName,
          lastName: u.lastName,
          firstName: u.firstName,
        };
      });
      this.users = orderBy(this.users, ["text"], ["asc"]);
      this.users.unshift({ value: null, text: "Please select" });
    },
    clearForm: function () {
      this.form = {
        amount: null,
        engineHours: null,
        fuelLiters: null,
        type: null,
        date: dayjs().format("YYYY-MM-DD"),
        user: null,
        description: null,
      };
      this.isLoading = false;
    },
    addDefault: function () {
      this.isLoading = true;
      const expense = {
        amount: Number(this.form.amount),
        typeId: Number(this.form.type),
        date: this.form.date,
        userId: Number(this.form.user),
        comment: this.form.description,
      };

      this.addExpense(expense)
        .then(() => {
          this.errors = [];
          this.close();
        })
        .catch((errors) => {
          this.errors = errors;
          this.isLoading = false;
        });
    },
    addFuel: function () {
      this.isLoading = true;
      const fuelEntry = {
        user_id: Number(this.form.user),
        engine_hours: Number(this.form.engineHours),
        liters: Number(this.form.fuelLiters),
        cost: Number(this.form.amount),
        date: this.form.date,
      };

      this.addFuelEntry(fuelEntry)
        .then(() => {
          this.errors = [];

          // refresh the expenses in statistics and transaction views
          // - this is a special case as we add a fuel entry from within accounting
          // - this could also be handled by the boat module in vuex.
          this.queryStatistics();
          this.queryTransactions()
            .catch((errors) => {
              console.error(
                "Updating transactions logged the following error(s):",
                errors
              );
              this.close();
            })
            .then(() => this.close());
        })
        .catch((errors) => {
          this.errors = errors;
          this.isLoading = false;
        });
    },
    add: function () {
      if (this.form.type == null) {
        return;
      }

      // fuel case
      if (Number(this.form.type) == 0) {
        this.addFuel();
      } else {
        this.addDefault();
      }
    },
    save: function () {
      this.add();
    },
    close: function () {
      this.clearForm();
      this.$emit("update:visible", false);
    },
  },
  created() {
    this.queryConfiguration();

    this.queryUserList()
      .then(() => this.buildUserSelect(this.userList))
      .catch((errors) => this.errors.push(...errors));

    this.queryExpenseTypes()
      .then(() => this.buildTypeSelect(this.getExpenseTypes))
      .catch((errors) => this.errors.push(...errors));

    this.form.date = dayjs().format("YYYY-MM-DD");
  },
};
</script>
