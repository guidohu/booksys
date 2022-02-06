<template>
  <div class="row mb-2">
    <label :for="id" class="col-3 col-form-label">{{ label }}</label>
    <div class="col-9">
      <div class="form-check form-switch align-middle">
        <input
          class="form-check-input"
          type="checkbox"
          role="switch"
          id="flexSwitchCheckDefault"
          :checked="formValue"
          @change="update()"
        />
        <span class="align-middle mt-4">
          <label
            class="form-check-label ms-3 mt-2"
            for="flexSwitchCheckDefault"
            >{{ formSelectedLabel }}</label
          >
        </span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "InputToggle",
  data() {
    return {
      formOffLabel: "off",
      formOnLabel: "on",
      formValue: false,
    };
  },
  props: ["id", "label", "offLabel", "onLabel", "modelValue"],
  emits: ["update:modelValue"],
  methods: {
    update() {
      this.formValue = !this.formValue;
      this.$emit("update:modelValue", this.formValue);
    },
  },
  watch: {
    modelValue: function (newValue) {
      this.formValue = newValue;
    },
  },
  computed: {
    formSelectedLabel: function () {
      if (this.formValue == false) {
        return this.formOffLabel;
      }
      return this.formOnLabel;
    },
  },
  created() {
    this.formOffLabel = this.offLabel;
    this.formOnLabel = this.onLabel;
    this.formValue = this.modelValue;
  },
};
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
