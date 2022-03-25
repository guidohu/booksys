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

<script>
export default {
  name: "InputDateTimeLocal",
  props: ["id", "label", "modelValue", "disabled", "size"],
  data() {
    return {
      formattedDate: "",
    };
  },
  emits: ["update:modelValue"],
  methods: {
    changeHandler: function (value) {
      console.log("updateValue to:", value);
      this.$emit("update:modelValue", value);
    },
    inputGroupClass: function () {
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
