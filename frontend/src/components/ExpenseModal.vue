<template>
  <modal-container name="expense-modal" :visible="visible">
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
          v-if="form.type != null && form.type == 1"
          id="fuel-liters"
          label="Fuel"
          size="small"
          placeholder="0.00"
          v-model="form.fuelLiters"
        />
        <input-engine-hours
          v-if="form.type != null && form.type == 1"
          id="engine-hours"
          label="Engine"
          size="small"
          placeholder="0"
          v-model="form.engineHours"
          :display-format="getEngineHourFormat"
        />
        <input-text-multiline
          v-if="form.type != null && form.type != 1"
          id="description"
          label="Description"
          rows="2"
          v-model="form.description"
        />
        <div
          class="alert alert-warning"
          v-if="form.type == 1 && getFuelPaymentType == 'billed'"
        >
          Entry will not be visible in the payment section as it is not a direct
          expense. According to your settings fuel consumption is 'billed' and
          thus only 'fuel bill' has an impact on the balance.
        </div>
      </form>
    </modal-body>
    <modal-footer>
      <button type="button" class="btn btn-outline-info mr-1" @click="add">
        <i class="bi bi-check" />
        Add
      </button>
      <button type="button" class="btn btn-outline-danger" @click="close">
        <i class="bi bi-x" />
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import dayjs from "dayjs";
import orderBy from "lodash/orderBy";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import InputTextMultiline from "./forms/inputs/InputTextMultiline.vue";
import InputCurrency from "./forms/inputs/InputCurrency.vue";
import InputDateTimeLocal from "./forms/inputs/InputDateTimeLocal.vue";
import InputFuel from "./forms/inputs/InputFuel.vue";
import InputEngineHours from "./forms/inputs/InputEngineHours.vue";

const props = defineProps(["visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const isLoading = ref(false);
const errors = ref([]);
const expenseTypes = ref([]);
const users = ref([]);
const userLabel = ref("User");
const userDescription = ref("");
const typeDescription = ref("");
const form = ref({
  type: null,
  user: null,
  date: null,
});

const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const getEngineHourFormat = computed(() => store.getters["configuration/getEngineHourFormat"]);
const getFuelPaymentType = computed(() => store.getters["configuration/getFuelPaymentType"]);
const userList = computed(() => store.getters["user/userList"]);
const getExpenseTypes = computed(() => store.getters["accounting/getExpenseTypes"]);

watch(getExpenseTypes, (newValues) => {
  buildTypeSelect(newValues);
});

watch(userList, (newValues) => {
  buildUserSelect(newValues);
});

watch(() => form.value.type, (newValue, oldValue) => {
  if (newValue !== oldValue) {
    const typeActions = {
      1: { userLabel: "Driver", userDescription: "Driver that fueled the boat", typeDescription: "Add a fuel entry." },
      10: { userLabel: "Payer", userDescription: "Pays the bill", typeDescription: "Fuel bill" },
      4: { userLabel: "Payee", userDescription: "Payee or internal reference.", typeDescription: "" },
      2: { userLabel: "Payee", userDescription: "Payee or internal reference.", typeDescription: "Expenses related to maintenance" },
      3: { userLabel: "Payee / Payer", userDescription: "Payee or internal reference.", typeDescription: "Expenses related to aquired parts." },
      5: { userLabel: "Account", userDescription: "Refund to this account.", typeDescription: "Refund of payments." },
      7: { userLabel: "Account", userDescription: "Refund to this account.", typeDescription: "Refund of payments." },
      6: { userLabel: "Payee", userDescription: "Payee or internal reference", typeDescription: "everything that does not fit another category" },
      8: { userLabel: "Driver", userDescription: "Driver that we pay a compensation.", typeDescription: "Compensation payments" },
      9: { userLabel: "Payee", userDescription: "Owner that gets a refund.", typeDescription: "" },
    };
    const action = typeActions[Number(newValue)] || { userLabel: "Payee", userDescription: "", typeDescription: "" };
    userLabel.value = action.userLabel;
    userDescription.value = action.userDescription;
    typeDescription.value = action.typeDescription;
  }
});

const queryConfiguration = () => store.dispatch("configuration/queryConfiguration");
const queryUserList = () => store.dispatch("user/queryUserList");
const queryTransactions = () => store.dispatch("accounting/queryTransactions");
const queryStatistics = () => store.dispatch("accounting/queryStatistics");
const queryExpenseTypes = () => store.dispatch("accounting/queryExpenseTypes");
const addExpense = (expense) => store.dispatch("accounting/addExpense", expense);
const addFuelEntry = (fuelEntry) => store.dispatch("boat/addFuelEntry", fuelEntry);

function buildTypeSelect(types) {
  let mappedTypes = types.map((t) => ({ value: t.id, text: t.name }));
  mappedTypes = orderBy(mappedTypes, ["text"], ["asc"]);
  mappedTypes.unshift({ value: null, text: "Please select" });
  expenseTypes.value = mappedTypes;
}

function buildUserSelect(usersData) {
  let mappedUsers = usersData.map((u) => ({
    value: u.id,
    text: `${u.firstName} ${u.lastName}`,
    lastName: u.lastName,
    firstName: u.firstName,
  }));
  mappedUsers = orderBy(mappedUsers, ["text"], ["asc"]);
  mappedUsers.unshift({ value: null, text: "Please select" });
  users.value = mappedUsers;
}

function clearForm() {
  form.value = {
    amount: null,
    engineHours: null,
    fuelLiters: null,
    type: null,
    date: dayjs().format("YYYY-MM-DD"),
    user: null,
    description: null,
  };
  isLoading.value = false;
}

function addDefault() {
  isLoading.value = true;
  const expense = {
    amount: Number(form.value.amount),
    typeId: Number(form.value.type),
    date: form.value.date,
    userId: Number(form.value.user),
    comment: form.value.description,
  };

  addExpense(expense)
    .then(() => {
      errors.value = [];
      close();
    })
    .catch((errs) => {
      errors.value = errs;
      isLoading.value = false;
    });
}

function addFuel() {
  isLoading.value = true;
  const fuelEntry = {
    user_id: Number(form.value.user),
    engine_hours: Number(form.value.engineHours),
    liters: Number(form.value.fuelLiters),
    cost: Number(form.value.amount),
    date: form.value.date,
  };

  addFuelEntry(fuelEntry)
    .then(() => {
      errors.value = [];
      queryStatistics();
      queryTransactions()
        .catch((errs) => {
          console.error("Updating transactions logged the following error(s):", errs);
        })
        .finally(() => close());
    })
    .catch((errs) => {
      errors.value = errs;
      isLoading.value = false;
    });
}

function add() {
  if (form.value.type == null) {
    return;
  }
  if (Number(form.value.type) === 1) {
    addFuel();
  } else {
    addDefault();
  }
}

function save() {
  add();
}

function close() {
  clearForm();
  emit("update:visible", false);
}

queryConfiguration();
queryUserList()
  .then(() => buildUserSelect(userList.value))
  .catch((errs) => errors.value.push(...errs));
queryExpenseTypes()
  .then(() => buildTypeSelect(getExpenseTypes.value))
  .catch((errs) => errors.value.push(...errs));

form.value.date = dayjs().format("YYYY-MM-DD");
</script>
