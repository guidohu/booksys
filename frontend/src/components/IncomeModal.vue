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
        <input-text-multiline
          v-if="form.type != null && form.type != 4"
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

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import * as dayjs from "dayjs";
import orderBy from "lodash/orderBy";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import InputTextMultiline from "./forms/inputs/InputTextMultiline.vue";
import InputCurrency from "./forms/inputs/InputCurrency.vue";
import InputDateTimeLocal from "./forms/inputs/InputDateTimeLocal.vue";

export default {
  name: "IncomeModal",
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
  },
  props: ["visible"],
  data() {
    return {
      errors: [],
      incomeTypes: [],
      users: [],
      userLabel: "User",
      userDescription: "",
      form: {
        type: null,
        user: null,
        date: null,
        comment: null,
      },
    };
  },
  computed: {
    ...mapGetters("configuration", ["getCurrency"]),
    ...mapGetters("user", ["userList"]),
    ...mapGetters("accounting", ["getIncomeTypes"]),
  },
  watch: {
    getIncomeTypes: function (newValues) {
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
        if ([6, 4].includes(Number(newValue))) {
          // 6: membership fee
          // 4: session
          this.userLabel = "Account";
          this.userDescription = "User that paid for sessions or membership.";
        } else if ([3, 5].includes(Number(newValue))) {
          // 3: invest
          // 5: other
          this.userLabel = "User";
          this.userDescription = "Payer or internal reference";
        } else if ([7].includes(Number(newValue))) {
          // 7: salary
          this.userLabel = "Driver";
          this.userDescription = "Driver that got paid for a session.";
        } else {
          // Default case
          this.userLabel = "User";
          this.userDescription = "";
        }
      }
    },
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    ...mapActions("user", ["queryUserList"]),
    ...mapActions("accounting", ["queryIncomeTypes", "addIncome"]),
    buildTypeSelect: function (types) {
      this.incomeTypes = types.map((t) => {
        return {
          value: t.id,
          text: t.name,
        };
      });
      this.incomeTypes = orderBy(this.incomeTypes, ["text"], ["asc"]);
      this.incomeTypes.unshift({ value: null, text: "Please select" });
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
        type: null,
        date: dayjs().format("YYYY-MM-DDTHH:mm"),
        user: null,
        comment: null,
      };
    },
    add: function () {
      const income = {
        amount: Number(this.form.amount),
        typeId: Number(this.form.type),
        date: this.form.date,
        userId: Number(this.form.user),
        comment: this.form.comment,
      };

      this.addIncome(income)
        .then(() => {
          this.errors = [];
          this.close();
        })
        .catch((errors) => (this.errors = errors));
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

    this.queryIncomeTypes()
      .then(() => this.buildTypeSelect(this.getIncomeTypes))
      .catch((errors) => this.errors.push(...errors));

    this.form.date = dayjs().format("YYYY-MM-DD");
  },
};
</script>
