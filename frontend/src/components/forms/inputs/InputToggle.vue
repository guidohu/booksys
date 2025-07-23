<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div class="d-inline-flex align-items-stretch">
        <div class="form-check form-switch">
          <input
            class="form-check-input"
            type="checkbox"
            role="switch"
            :id="'flex-switch-' + id"
            :checked="formValue"
            @change="update()"
          />
        </div>
        <span class="ps-2 mt-2">
          <label class="form-check-label" :for="'flex-switch-' + id">{{
            formSelectedLabel
          }}</label>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits, ref, computed, watch } from "vue";

const props = defineProps(["id", "label", "offLabel", "onLabel", "modelValue"]);

const emit = defineEmits(["update:modelValue", "change"]);

const formValue = ref(props.modelValue);
const formOffLabel = ref(props.offLabel || "off");
const formOnLabel = ref(props.onLabel || "on");

const update = () => {
  formValue.value = !formValue.value;
  emit("update:modelValue", formValue.value);
  emit("change");
};

watch(
  () => props.modelValue,
  (newValue) => {
    formValue.value = newValue;
  }
);

const formSelectedLabel = computed(() => {
  if (formValue.value == false) {
    return formOffLabel.value;
  }
  return formOnLabel.value;
});
</script>

<style scoped>
.form-switch .form-check-input {
  width: 4em;
  height: 1.9em;
}

.form-check-input:checked {
  background-color: #0dcaf0;
  border-color: rgb(0 0 0 / 25%);
}

.form-check-input:focus {
  box-shadow: 0 0 0 0.25rem rgb(110 110 110 / 25%);
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='-4 -4 8 8'%3e%3ccircle r='3' fill='rgba%280, 0, 0, 0.25%29'/%3e%3c/svg%3e");
  border-color: rgb(0 0 0 / 25%);
}
</style>
