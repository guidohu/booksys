<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <input
          :type="type"
          class="form-control"
          :disabled="disabled"
          :id="id"
          :value="modelValue"
          :placeholder="placeholder"
          @input="changeHandler($event.target.value)"
        />
        <span v-if="suffix != null" class="input-group-text">{{ suffix }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "InputText",
  props: ["id", "label", "modelValue", "type", "disabled", "size", "suffix", "placeholder"],
  emits: ["update:modelValue"],
  methods: {
    changeHandler(value) {
      this.$emit("update:modelValue", value);
      this.$emit("input");
    },
    inputGroupClass() {
      if (this.size == null) {
        return "input-group";
      }
      if (this.size == "small") {
        return "input-group input-group-sm";
      }
      if (this.size == "large") {
        return "input-group input-group-lg";
      }
    },
  },
};
</script>
