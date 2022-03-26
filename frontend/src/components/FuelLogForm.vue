<template>
  <div class="container text-left">
    <div class="row">
      <div class="col-12">
        <warning-box v-if="errors.length > 0" :errors="errors" />
        <form @submit.prevent.self="add">
          <input-engine-hours
            id="engine-hours"
            label="Engine"
            :display-format="getEngineHourFormat"
            v-model="form.engineHours"
            placeholder="0"
            size="small"
          />
          <input-currency
            id="cost"
            label="Cost"
            :currency="getCurrency"
            v-model="form.cost"
            placeholder="0"
            size="small"
          />
          <input-fuel
            id="fuel"
            label="Fuel"
            v-model="form.liters"
            placeholder="0"
            size="small"
          />
          <form-button
            type="submit"
            btn-style="info"
            btn-size="small"
            @click.prevent.self="add"
          >
            Add
          </form-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
import InputEngineHours from "@/components/forms/inputs/InputEngineHours";
import InputCurrency from "@/components/forms/inputs/InputCurrency.vue";
import InputFuel from "@/components/forms/inputs/InputFuel.vue";
import FormButton from "@/components/forms/FormButton.vue";

export default {
  name: "FuelLogForm",
  components: {
    WarningBox,
    InputEngineHours,
    InputCurrency,
    InputFuel,
    FormButton,
  },
  data() {
    return {
      errors: [],
      form: {
        engineHours: null,
        cost: null,
        liters: null,
      },
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getCurrency", "getEngineHourFormat"]),
  },
  methods: {
    add: function () {
      const entry = {
        user_id: this.userInfo.id,
        engine_hours: this.form.engineHours,
        liters: this.form.liters,
        cost: this.form.cost,
      };
      this.addFuelEntry(entry)
        .then(() => {
          this.resetForm();
        })
        .catch((errors) => (this.errors = errors));
    },
    resetForm: function () {
      this.form = {
        engineHours: null,
        cost: null,
        liters: null,
      };
      this.errors = [];
    },
    ...mapActions("boat", ["addFuelEntry"]),
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.queryConfiguration();
  },
};
</script>
