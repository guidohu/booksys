<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <input
          type="text"
          class="form-control"
          :disabled="disabled"
          :id="id"
          :value="modelValue"
          :placeholder="placeholder"
          @focus="$event.target.select()"
          @input="changeHandler($event.target.value)"
        />
        <span class="input-group-text">{{ currency }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "InputCurrency",
  props: [
    "id",
    "label",
    "modelValue",
    "disabled",
    "size",
    "currency",
    "placeholder",
  ],
  emits: ["update:modelValue"],
  methods: {
    changeHandler(value) {
      console.log("updateValue to:", value);
      this.$emit("update:modelValue", value);
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
