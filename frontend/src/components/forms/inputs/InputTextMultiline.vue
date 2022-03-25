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

<script>
export default {
  name: "InputText",
  props: [
    "id",
    "label",
    "modelValue",
    "disabled",
    "size",
    "placeholder",
    "rows",
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
