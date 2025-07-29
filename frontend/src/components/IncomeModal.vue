<template>
  <modal-container name="income-modal" :visible="visible">
    <modal-header
      :closable="true"
      title="Add New Income"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length" :errors="errors" />
      <form @submit.prevent="save">
        <input-select
          id="type-select"
          label="Type"
          v-model="form.type"
          :options="incomeTypes"
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
          label="Amount"
          size="small"
          v-model="form.amount"
          placeholder="0.00"
          :currency="getCurrency"
        />
        <!-- Income type 5 (session payment) does not need a comment
        it will be added automatically. -->
        <input-text-multiline
          v-if="form.type != null && form.type != 5"
          id="description"
          label="Description"
          rows="2"
          v-model="form.comment"
        />
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

const props = defineProps(["visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const errors = ref([]);
const incomeTypes = ref([]);
const users = ref([]);
const userLabel = ref("User");
const userDescription = ref("");
const form = ref({
  type: null,
  user: null,
  date: null,
  comment: null,
});

const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const userList = computed(() => store.getters["user/userList"]);
const getIncomeTypes = computed(() => store.getters["accounting/getIncomeTypes"]);

watch(getIncomeTypes, (newValues) => {
  buildTypeSelect(newValues);
});

watch(userList, (newValues) => {
  buildUserSelect(newValues);
});

watch(() => form.value.type, (newValue, oldValue) => {
  if (newValue !== oldValue) {
    const typeActions = {
      6: { userLabel: "Account", userDescription: "User that paid for sessions or membership." },
      4: { userLabel: "Account", userDescription: "User that paid for sessions or membership." },
      3: { userLabel: "User", userDescription: "Payer or internal reference" },
      5: { userLabel: "User", userDescription: "Payer or internal reference" },
      7: { userLabel: "Driver", userDescription: "Driver that got paid for a session." },
    };
    const action = typeActions[Number(newValue)] || { userLabel: "User", userDescription: "" };
    userLabel.value = action.userLabel;
    userDescription.value = action.userDescription;
  }
});

const queryConfiguration = () => store.dispatch("configuration/queryConfiguration");
const queryUserList = () => store.dispatch("user/queryUserList");
const queryIncomeTypes = () => store.dispatch("accounting/queryIncomeTypes");
const addIncome = (income) => store.dispatch("accounting/addIncome", income);

function buildTypeSelect(types) {
  let mappedTypes = types.map((t) => ({ value: t.id, text: t.name }));
  mappedTypes = orderBy(mappedTypes, ["text"], ["asc"]);
  mappedTypes.unshift({ value: null, text: "Please select" });
  incomeTypes.value = mappedTypes;
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
    type: null,
    date: dayjs().format("YYYY-MM-DDTHH:mm"),
    user: null,
    comment: null,
  };
}

function add() {
  const income = {
    amount: Number(form.value.amount),
    typeId: Number(form.value.type),
    date: form.value.date,
    userId: Number(form.value.user),
    comment: form.value.comment,
  };

  addIncome(income)
    .then(() => {
      errors.value = [];
      close();
    })
    .catch((errs) => (errors.value = errs));
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
queryIncomeTypes()
  .then(() => buildTypeSelect(getIncomeTypes.value))
  .catch((errs) => errors.value.push(...errs));

form.value.date = dayjs().format("YYYY-MM-DD");
</script>
