<template>
  <b-form-group
    id="input-group-fuel-payment-type-selector"
    label="Fuel Payment*"
    label-for="input-fuel-payment-type-selector"
    label-cols="3"
    description="select between direct payment and credited payment where you add payments for fuel in the payments section separately"
  >
    <b-form-select
      v-model="formSelected"
      :options="options"
    />
  </b-form-group>
</template>

<script>
import { BFormGroup, BFormSelect } from "bootstrap-vue";

export default {
  name: "FuelPaymentTypeSelector",
  components: {
    BFormGroup,
    BFormSelect,
  },
  props: ["selected"],
  data() {
    return {
      formSelected: "instant",
      options: [
        { value: "instant", text: "pay directly" },
        { value: "billed", text: "pay by bill" },
      ],
    };
  },
  watch: {
    formSelected: function () {
      this.$emit("update:selected", this.formSelected);
    },
    selected: function () {
      this.setFormSelected(this.selected);
    },
  },
  created() {
    this.setFormSelected(this.selected);
  },
  methods: {
    setFormSelected: function (format) {
      if (format == null) {
        console.log("no default payment type specified");
        this.formSelected = "instant";
      } else if (format == "instant" || format == "billed") {
        this.formSelected = format;
      } else {
        this.formSelected = "instant";
      }
    },
  },
};
</script>
