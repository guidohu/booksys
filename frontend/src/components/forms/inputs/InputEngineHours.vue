<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <input
          type="text"
          class="form-control"
          :id="id"
          :state="state"
          :disabled="disabled"
          v-model="formValue"
          @blur="changeHandler($event.target.value)"
          :placeholder="formPlaceholder"
        />
        <span class="input-group-text">{{ unitText }}</span>
      </div>
    </div>
  </div>
</template>

<script>
import {
  formatEngineHourLabel,
  formatEngineHour,
  isValidEngineHour,
  convertEngineHour,
} from "@/libs/formatters";

export default {
  name: "InputEngineHours",
  props: [
    "id",
    "label",
    "modelValue",
    "disabled",
    "displayFormat",
    "placeholder",
    "size",
  ],
  data() {
    return {
      formValue: null,
      formPlaceholder: null,
      state: null,
    };
  },
  emits: ["update:modelValue"],
  methods: {
    changeHandler(value) {
      if (this.formValue != null) {
        this.$emit("update:modelValue", convertEngineHour(value));
      } else {
        this.$emit("update:modelValue", null);
      }
    },
    inputGroupClass() {
      if(this.size == null){
        return "input-group";
      }
      if(this.size == "small"){
        return "input-group input-group-sm"
      }
      if(this.size == "large"){
        return "input-group input-group-lg"
      }
    }
  },
  computed: {
    unitText: function () {
      return formatEngineHourLabel(this.displayFormat);
    },
  },
  watch: {
    modelValue: function() {
      this.formValue = formatEngineHour(this.modelValue, this.displayFormat);
    },
    formValue: function() {
      if (this.formValue == null || this.formValue == "") {
        this.state = null;
      } else if (!isValidEngineHour(this.formValue, this.displayFormat)) {
        this.state = false;
      } else {
        this.state = null;
      }
    }
  },
  created() {
    this.formValue = formatEngineHour(this.modelValue, this.displayFormat);

    if(this.placeholder != null){
      this.formPlaceholder = formatEngineHour(
        this.placeholder,
        this.displayFormat
      );
    }
  }
};
</script>
