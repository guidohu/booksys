<template>
  <div v-if="label" class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <select
          :id="id"
          :class="formSelectClass()"
          :disabled="disabled"
          @click.stop
          @change="changeHandler($event)"
          v-model="selectedValue"
          :size="selectSize"
          :multiple="selectMode == 'multiple'"
        >
          <option v-for="o in options" :key="o.value" :value="o.value">
            {{ o.text }}
          </option>
        </select>
      </div>
      <div v-if="description" :id="id + '-description'" class="form-text">
        {{ description }}
      </div>
    </div>
  </div>
  <div v-else>
    <select
      :id="id"
      :class="formSelectClass()"
      :disabled="disabled"
      @click.stop
      @change="changeHandler($event)"
      v-model="selectedValue"
      :size="selectSize"
      :multiple="selectMode == 'multiple'"
    >
      <option v-for="o in options" :key="o.value" :value="o.value">
        {{ o.text }}
      </option>
    </select>
    <div v-if="description" :id="id + '-description'" class="form-text">
      {{ description }}
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits, ref, watch } from "vue";

const props = defineProps([
  "id",
  "label",
  "modelValue",
  "options",
  "disabled",
  "size",
  "selectMode",
  "selectSize",
  "description",
]);

const emit = defineEmits(["update:modelValue", "changed"]);

const selectedValue = ref(props.modelValue);

watch(
  () => props.modelValue,
  (newValue) => {
    selectedValue.value = newValue;
  }
);

const changeHandler = (event) => {
  console.log(event);
  if (props.selectMode != "multiple") {
    const value = event.target.value;
    console.log("selected value", value);
    emit("update:modelValue", value);
  } else {
    let array = [];
    const options = event.target.selectedOptions;
    for (let i = 0; i < options.length; i++) {
      array.push(options[i].value);
    }
    console.log("selected values", array);
    emit("update:modelValue", array);
  }
  emit("changed");
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

const formSelectClass = () => {
  if (props.size == null) {
    return "form-select";
  }
  if (props.size == "small") {
    return "form-select form-select-sm";
  }
  if (props.size == "large") {
    return "form-select form-select-lg";
  }
};
</script>
