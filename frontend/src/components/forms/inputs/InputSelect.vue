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
          :size="selectSize"
          :multiple="selectMode=='multiple'"
        >
          <option v-for="o in options" :key="o.value" :value="o.value">
            {{ o.text }}
          </option>
        </select>
      </div>
      <div v-if="description" :id="id+'-description'" class="form-text">
        {{ description }}
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
      :size="selectSize"
      :multiple="selectMode=='multiple'"
    >
      <option v-for="o in options" :key="o.value" :value="o.value">
        {{ o.text }}
      </option>
    </select>
    <div v-if="description" :id="id+'-description'" class="form-text">
      {{ description }}
    </div>
  </div>
</template>

<script>
export default {
  name: "InputSelect",
  props: ["id", "label", "modelValue", "options", "disabled", "size", "selectMode", "selectSize", "description"],
  emits: ["update:modelValue", "changed"],
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
      console.log(event);
      if (this.selectMode != 'multiple') {
        const value = event.target.value;
        console.log("selected value", value);
        this.$emit("update:modelValue", value);
      } else {
        let array = [];
        const options = event.target.selectedOptions;
        for (let i = 0; i < options.length; i++) { 
          array.push(options[i].value);
        }
        console.log("selected values", array);
        this.$emit("update:modelValue", array);
      }
      this.$emit("changed");
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
