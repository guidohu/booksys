<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div :class="inputGroupClass()">
        <input
          type="file"
          class="form-control"
          :disabled="disabled"
          :id="id"
          :accept="accept"
          :placeholder="placeholder"
          @change="changeHandler($event)"
        />
        <span v-if="suffix != null" class="input-group-text">{{ suffix }}</span>
      </div>
      <div v-if="description" :id="id + '-description'" class="form-text">
        {{ description }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "InputFile",
  props: [
    "id",
    "label",
    "modelValue",
    "disabled",
    "size",
    "suffix",
    "accept",
    "placeholder",
    "description",
  ],
  emits: ["update:modelValue"],
  methods: {
    changeHandler(event) {
      this.$emit("update:modelValue", event.target.files[0]);
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
