<template>
  <b-form-group
    :id="inputGroupId"
    :label="label"
    :label-for="labelFor"
    label-cols="3"
  >
    <b-input-group :size="formSize">
      <b-form-input
        :id="id"
        v-model="formValue"
        type="text"
        :state="state"
        :disabled="disabled"
        v-on:blur="propagateValue"
        :placeholder="formPlaceholder"
      />
      <b-input-group-append is-text>
        {{ unitText }}
      </b-input-group-append>
    </b-input-group>
  </b-form-group>
</template>

<script>
import {
  formatEngineHourLabel,
  formatEngineHour,
  isValidEngineHour,
  convertEngineHour,
} from "@/libs/formatters";
import {
  BFormGroup,
  BFormInput,
  BInputGroup,
  BInputGroupAppend,
} from "bootstrap-vue";

export default {
  name: "EngineHours",
  props: ["label", "value", "displayFormat", "disabled", "size", "placeholder"],
  components: {
    BFormGroup,
    BFormInput,
    BInputGroup,
    BInputGroupAppend,
  },
  data() {
    return {
      formValue: null,
      formSize: "sm",
      formPlaceholder: null,
      state: null,
    };
  },
  computed: {
    inputGroupId: function () {
      return "input-group-" + this.label.toLowerCase().replace(" ", "-");
    },
    labelFor: function () {
      return "input-" + this.label.toLowerCase().replace(" ", "-");
    },
    id: function () {
      return "input-" + this.label.toLowerCase().replace(" ", "-");
    },
    unitText: function () {
      return formatEngineHourLabel(this.displayFormat);
    },
  },
  watch: {
    value: function () {
      this.formValue = formatEngineHour(this.value, this.displayFormat);
    },
    formValue: function () {
      if (this.formValue == null || this.formValue == "") {
        this.state = null;
      } else if (!isValidEngineHour(this.formValue, this.displayFormat)) {
        this.state = false;
      } else {
        this.state = null;
      }
    },
  },
  methods: {
    propagateValue: function () {
      if (this.formValue != null) {
        this.$emit("input", convertEngineHour(this.formValue));
      } else {
        this.$emit("input", null);
      }
    },
  },
  created() {
    this.formValue = formatEngineHour(this.value, this.displayFormat);
    if (this.size != null) {
      this.formSize = this.size;
    }
    if (this.placeholder != null) {
      this.formPlaceholder = formatEngineHour(
        this.placeholder,
        this.displayFormat
      );
    }
  },
};
</script>
