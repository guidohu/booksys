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
        >
          <option v-for="o in options" :key="o.value" :value="o.value">
            {{ o.text }}
          </option>
        </select>
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
    >
      <option v-for="o in options" :key="o.value" :value="o.value">
        {{ o.text }}
      </option>
    </select>
  </div>
</template>

<script>
export default {
  name: "InputSelect",
  props: ["id", "label", "modelValue", "options", "disabled", "size"],
  emits: ["update:modelValue"],
  data() {
    return {
      selectedValue: 0
    }
  },
  watch: {
    modelValue: function(newValue) {
      this.selectedValue = newValue;
    }
  },
  methods: {
    changeHandler(event) {
      const value = event.target.value;
      console.log("selected value", value);
      this.$emit("update:modelValue", value);
    },
    isSelected: function(option) {
      return option.value == this.modelValue;
    },
    inputGroupClass() {
      if(this.size == null){
        return "input-group";
      }
      if(this.size == "small"){
        return "input-group input-group-sm"
      }
      if(this.size == "large"){
        return "input-group input-group-lg"
      }
    },
    formSelectClass() {
      if(this.size == null){
        return "form-select";
      }
      if(this.size == "small"){
        return "form-select form-select-sm"
      }
      if(this.size == "large"){
        return "form-select form-select-lg"
      }
    }
  },
  created() {
    this.selectedValue = this.modelValue;
  }
};
</script>
