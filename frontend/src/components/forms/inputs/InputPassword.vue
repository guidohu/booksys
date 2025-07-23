<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <input
          type="password"
          class="form-control"
          :disabled="disabled"
          :id="id"
          :value="modelValue"
          :autocomplete="autocomplete"
          @input="changeHandler($event.target.value)"
        />
        <span v-if="suffix != null" class="input-group-text">{{ suffix }}</span>
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
  "type",
  "disabled",
  "size",
  "suffix",
  "autocomplete",
]);

const emit = defineEmits(["update:modelValue", "input"]);

const changeHandler = (value) => {
  emit("update:modelValue", value);
  emit("input");
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
