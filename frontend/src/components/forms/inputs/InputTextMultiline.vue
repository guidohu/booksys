<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <textarea
          :id="id"
          class="form-control"
          :rows="rows"
          :placeholder="placeholder"
          :value="modelValue"
          :disabled="disabled"
          @input="changeHandler($event.target.value)"
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from "vue";

const props = defineProps([
  "id",
  "label",
  "modelValue",
  "disabled",
  "size",
  "placeholder",
  "rows",
]);

const emit = defineEmits(["update:modelValue"]);

const changeHandler = (value) => {
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
