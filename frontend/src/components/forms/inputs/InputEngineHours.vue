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
          @focus="$event.target.select()"
          @blur="changeHandler($event.target.value)"
          :placeholder="formPlaceholder"
        />
        <span class="input-group-text">{{ unitText }}</span>
      </div>
      <div v-if="description" :id="id + '-description'" class="form-text">
        {{ description }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits, ref, computed, watch, onMounted } from "vue";
import {
  formatEngineHourLabel,
  formatEngineHour,
  isValidEngineHour,
  convertEngineHour,
} from "booksys/libs/formatters";

const props = defineProps([
  "id",
  "label",
  "modelValue",
  "description",
  "disabled",
  "displayFormat",
  "placeholder",
  "size",
]);

const emit = defineEmits(["update:modelValue"]);

const formValue = ref(null);
const formPlaceholder = ref(null);
const state = ref(null);

const changeHandler = (value) => {
  if (formValue.value != null) {
    emit("update:modelValue", convertEngineHour(value));
  } else {
    emit("update:modelValue", null);
  }
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

const unitText = computed(() => {
  return formatEngineHourLabel(props.displayFormat);
});

watch(
  () => props.modelValue,
  (newValue) => {
    formValue.value = formatEngineHour(newValue, props.displayFormat);
  },
);

watch(
  () => formValue.value,
  (newValue) => {
    if (newValue == null || newValue == "") {
      state.value = null;
    } else if (!isValidEngineHour(newValue, props.displayFormat)) {
      state.value = false;
    } else {
      state.value = null;
    }
  },
);

onMounted(() => {
  formValue.value = formatEngineHour(props.modelValue, props.displayFormat);

  if (props.placeholder != null) {
    formPlaceholder.value = formatEngineHour(
      props.placeholder,
      props.displayFormat,
    );
  }
});
</script>
