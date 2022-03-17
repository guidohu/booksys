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

<script>
export default {
  name: "InputPassword",
  props: [
    "id",
    "label",
    "modelValue",
    "type",
    "disabled",
    "size",
    "suffix",
    "autocomplete",
  ],
  emits: ["update:modelValue", "input"],
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
