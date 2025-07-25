<template>
  <modal-container name="fuelEntrymodal" :visible="visible">
    <modal-header
      :closable="true"
      title="Fuel Entry"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <div class="row" v-if="errors.length">
        <div class="col-1 d-none d-sm-block" />
        <div class="col-12 col-sm-10">
          <warning-box :errors="errors" />
        </div>
      </div>
      <form @submit.prevent.self="save">
        <input-date-time-local
          id="date"
          label="Date"
          v-model="form.date"
          size="small"
          :disabled="true"
        />
        <input-text
          id="driver"
          label="Driver"
          v-model="form.driver"
          size="small"
          :disabled="true"
        />
        <input-text
          id="consumption"
          label="Consumption"
          v-model="form.averageFuelPerHour"
          size="small"
          :disabled="true"
          suffix="ltrs/h"
        />
        <input-engine-hours
          id="engine-hours"
          label="Engine"
          v-model="form.engineHours"
          :displayFormat="getEngineHourFormat"
          size="small"
        />
        <input-fuel
          id="fuel"
          label="Fuel"
          v-model="form.fuel"
          placeholder="0"
          size="small"
        />
        <input-currency
          v-if="!form.isDiscounted"
          id="cost"
          label="Cost"
          :currency="getCurrency"
          v-model="form.cost"
          placeholder="0"
          size="small"
        />
        <input-currency
          v-if="form.isDiscounted"
          id="cost-gross"
          label="Cost (Gross)"
          :currency="getCurrency"
          v-model="form.costGross"
          placeholder="0"
          size="small"
        />
        <input-toggle
          id="discount"
          label="Discount"
          v-model="form.isDiscounted"
          offLabel="No Discount"
          onLabel="Discounted"
        />
        <input-currency
          v-if="form.isDiscounted"
          id="cost-net"
          label="Cost (Net)"
          :currency="getCurrency"
          v-model="form.costNet"
          placeholder="0"
          size="small"
        />
      </form>
    </modal-body>
    <modal-footer>
      <button
        type="submit"
        class="btn btn-outline-info"
        @click.prevent.self="save"
      >
        <i class="bi bi-check"></i>
        Save
      </button>
      <button
        type="button"
        class="btn btn-outline-danger"
        @click.prevent.self="close"
      >
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import {
  formatCurrency,
  formatFuelConsumption,
  formatFuel,
} from "booksys/libs/formatters";
import dayjs from "dayjs";
import ModalContainer from "booksys/components/bricks/ModalContainer.vue";
import ModalHeader from "booksys/components/bricks/ModalHeader.vue";
import ModalBody from "booksys/components/bricks/ModalBody.vue";
import ModalFooter from "booksys/components/bricks/ModalFooter.vue";
import InputDateTimeLocal from "booksys/components/forms/inputs/InputDateTimeLocal.vue";
import InputText from "booksys/components/forms/inputs/InputText.vue";
import InputEngineHours from "booksys/components/forms/inputs/InputEngineHours.vue";
import InputFuel from "booksys/components/forms/inputs/InputFuel.vue";
import InputCurrency from "booksys/components/forms/inputs/InputCurrency.vue";
import InputToggle from "booksys/components/forms/inputs/InputToggle.vue";

export default {
  name: "FuelEntryModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputDateTimeLocal,
    InputText,
    InputEngineHours,
    InputFuel,
    InputCurrency,
    InputToggle,
  },
  props: ["fuelEntry", "visible"],
  data() {
    return {
      errors: [],
      form: {
        date: null,
      },
    };
  },
  computed: {
    toggleWidth: function () {
      return 100;
    },
    ...mapGetters("configuration", ["getCurrency", "getEngineHourFormat"]),
  },
  watch: {
    fuelEntry: function (newValue) {
      this.setFormContent(newValue);
    },
  },
  methods: {
    setFormContent: function (entry) {
      if (entry != null) {
        // apply some logic depending on the backend data to get the right values for
        // net/gross and cost
        const costGross = entry.is_discounted == false ? entry.cost : entry.cost_brutto;
        const costNet = entry.is_discounted == false ? null : entry.cost;
        const cost = entry.is_discounted == false ? entry.cost : entry.cost_brutto;

        this.form = {
          id: entry.id,
          date: dayjs(entry.timestamp * 1000).format("YYYY-MM-DDTHH:mm"),
          isDiscounted: entry.is_discounted,
          cost: formatCurrency(cost, null),
          costGross: formatCurrency(costGross, null),
          costNet: formatCurrency(costNet, null),
          engineHours: entry.engine_hours,
          fuel: formatFuel(entry.liters),
          driver: entry.user_first_name,
          averageFuelPerHour: formatFuelConsumption(entry.avg_liters_per_hour),
        };
      }
    },
    toggleDiscount: function () {
      this.form.isDiscounted = !this.form.isDiscounted;
    },
    costChange: function () {
      this.form.costGross = this.form.cost;
    },
    costGrossChange: function () {
      this.form.cost = this.form.costGross;
    },
    close: function () {
      this.$emit("update:visible", false);
    },
    save: function () {
      if (this.form.isDiscounted && this.form.costNet == null) {
        this.errors = ["Cost (net) cannot be empty if you select a discount."];
        return;
      }

      let cost = 0;
      let costBrutto = null;
      let isDiscounted = false;
      if (this.form.isDiscounted) {
        cost = this.form.costNet;
        costBrutto = this.form.costGross;
        isDiscounted = true;
      } else {
        cost = this.form.cost;
        costBrutto = null;
        isDiscounted = false
      }

      const updatedEntry = {
        id: this.form.id,
        engine_hours: this.form.engineHours,
        liters: this.form.fuel,
        cost: cost,
        cost_brutto: costBrutto,
        is_discounted: isDiscounted,
      };

      this.updateFuelEntry(updatedEntry)
        .then(() => this.close())
        .catch((errors) => (this.errors = errors));
    },
    ...mapActions("boat", ["updateFuelEntry"]),
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.setFormContent(this.fuelEntry);
    this.queryConfiguration();
  },
};
</script>
