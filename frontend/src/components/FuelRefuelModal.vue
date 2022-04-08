<template>
  <modal-container name="refuelModal" :visible="visible">
    <modal-header
      :closable="true"
      title="Refuel"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length > 0" :errors="errors" />
      <form @submit.stop="saveFuel">
        <input-engine-hours
          id="engine-hours"
          label="Engine"
          :display-format="getEngineHourFormat"
          v-model="form.engineHours"
          :description="engineHourDescription"
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
      </form>
    </modal-body>
    <modal-footer>
      <button type="button" class="btn btn-outline-info" @click="saveFuel">
        <i class="bi bi-check" />
        Save
      </button>
      <button
        type="button"
        class="btn btn-outline-danger mb-1"
        @click="$emit('update:visible', false)"
      >
        <i class="bi bi-x" />
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import WarningBox from "@/components/WarningBox";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputEngineHours from "@/components/forms/inputs/InputEngineHours";
import InputCurrency from "@/components/forms/inputs/InputCurrency.vue";
import InputFuel from "@/components/forms/inputs/InputFuel.vue";

export default {
  name: "FuelRefuelModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputEngineHours,
    InputCurrency,
    InputFuel,
  },
  data() {
    return {
      errors: [],
      form: {
        engineHours: 0,
        cost: 0,
        liters: 0,
      },
      engineHourDescription: null,
    };
  },
  props: ["visible"],
  emits: ["saved", "update:visible"],
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("boat", ["getMyNautiqueEngineHours"]),
    ...mapGetters("configuration", ["getCurrency", "getEngineHourFormat"]),
  },
  watch: {
    getMyNautiqueEngineHours: function (newValue) {
      if (this.form.engineHours == null || this.form.engineHours == "") {
        this.form.engineHours = newValue;
        this.engineHourDescription = "prefilled by myNautique";
      }
    },
  },
  methods: {
    saveFuel: function () {
      const entry = {
        user_id: this.userInfo.id,
        engine_hours: this.form.engineHours,
        liters: this.form.liters,
        cost: this.form.cost,
      };
      this.addFuelEntry(entry)
        .then(() => {
          this.resetForm();
          this.close();
        })
        .catch((errors) => (this.errors = errors));
    },
    resetForm: function () {
      this.form.engineHours = "";
      this.form.cost = "";
      this.form.liters = "";
    },
    close: function () {
      this.$emit("update:visible", false);
    },
    ...mapActions("boat", ["addFuelEntry"]),
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.queryConfiguration();

    // set engine hours if known through myNautique
    if (
      this.getMyNautiqueEngineHours != null &&
      this.getMyNautiqueEngineHours > 0
    ) {
      this.form.engineHours = this.getMyNautiqueEngineHours;
      this.engineHourDescription = "prefilled by myNautique";
    }
  },
};
</script>
