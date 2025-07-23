<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <input
          type="datetime-local"
          class="form-control"
          :disabled="disabled"
          :id="id"
          :value="modelValue"
          @input="changeHandler($event.target.value)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from "vue";

const props = defineProps(["id", "label", "modelValue", "disabled", "size"]);

const emit = defineEmits(["update:modelValue"]);

const changeHandler = (value) => {
  console.log("updateValue to:", value);
  emit("update:modelValue", value);
};

const inputGroupClass = () => {
  if (props.size == null) {
    return "input-group";
  }
  if (props.size == "small") {
    return "input-group input-group-sm";
  }
  if (props.size == "large") {
    return "input-group input-group-lg";
  }
};
</script>
