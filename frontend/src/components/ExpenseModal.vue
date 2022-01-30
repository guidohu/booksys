<template>
  <b-modal
    id="expenseModal"
    title="Add New Expense"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <b-row v-if="errors.length > 0">
      <b-col cols="1" class="d-none d-sm-block" />
      <b-col cols="12" sm="10">
        <WarningBox :errors="errors" dismissible="true" />
      </b-col>
      <b-col cols="1" class="d-none d-sm-block" />
    </b-row>
    <b-form @submit="save">
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block" />
        <b-col cols="12" sm="10">
          <b-form-group
            id="type-select-group"
            label="Type"
            label-for="type-select"
            :description="typeDescription"
            label-cols="3"
          >
            <b-form-select
              id="type-select"
              v-model="form.type"
              :options="expenseTypes"
            />
          </b-form-group>
          <b-form-group
            id="date-group"
            label="Date"
            label-for="date-select"
            description=""
            label-cols="3"
          >
            <b-form-input id="date-select" v-model="form.date" type="date" />
          </b-form-group>
          <b-form-group
            v-if="form.type != null"
            id="user-select-group"
            :label="userLabel"
            label-for="user-select"
            :description="userDescription"
            label-cols="3"
          >
            <b-form-select
              id="user-select"
              v-model="form.user"
              :options="users"
            />
          </b-form-group>
          <b-form-group
            v-if="form.type != null"
            id="amount"
            label="Cost"
            label-for="amount-input"
            description=""
            label-cols="3"
          >
            <b-input-group>
              <b-form-input
                id="amount-input"
                v-model="form.amount"
                type="text"
                placeholder="0.00"
              />
              <b-input-group-append is-text>
                {{ getCurrency }}
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            v-if="form.type != null && form.type == 0"
            id="fuel-liters"
            label="Fuel"
            label-for="fuel-liters-input"
            description=""
            label-cols="3"
          >
            <b-input-group>
              <b-form-input
                id="fuel-liters-input"
                v-model="form.fuelLiters"
                type="text"
                placeholder="0.00"
              />
              <b-input-group-append is-text> ltrs </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <engine-hours
            v-if="form.type != null && form.type == 0"
            v-model="form.engineHours"
            label="Engine Hours"
            :display-format="getEngineHourFormat"
            size="md"
            placeholder="0"
          />
          <b-form-group
            v-if="form.type != null && form.type != 0"
            id="description"
            label="Description"
            label-for="description-input"
            description=""
            label-cols="3"
          >
            <b-form-textarea
              id="description-input"
              v-model="form.description"
              type="text"
              rows="2"
              placeholder=""
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <b-overlay
        id="overlay-background"
        :show="isLoading"
        spinner-type="border"
        spinner-variant="info"
        rounded="sm"
      >
        <div class="text-right d-inline">
          <b-button
            v-if="form.type != null"
            class="ml-4"
            type="button"
            variant="outline-info"
            @click="add"
          >
            <b-icon-check />
            Add
          </b-button>
          <b-button
            class="ml-1"
            type="button"
            variant="outline-danger"
            @click="close"
          >
            <b-icon-x />
            Cancel
          </b-button>
        </div>
      </b-overlay>
    </div>
  </b-modal>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import EngineHours from "@/components/forms/inputs/EngineHours";
import * as dayjs from "dayjs";
import orderBy from "lodash/orderBy";
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormSelect,
  BFormInput,
  BInputGroup,
  BInputGroupAppend,
  BFormTextarea,
  BOverlay,
  BButton,
  BIconCheck,
  BIconX,
} from "bootstrap-vue";

export default {
  name: "ExpenseModal",
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormSelect,
    BFormInput,
    BInputGroup,
    BInputGroupAppend,
    BFormTextarea,
    BOverlay,
    BButton,
    BIconCheck,
    BIconX,
    EngineHours,
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
    ...mapGetters("configuration", ["getCurrency", "getEngineHourFormat"]),
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
