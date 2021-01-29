<template>
  <b-form-group
    :id="inputGroupId"
    :label="label"
    :label-for="labelFor"
    label-cols="3"
  >
    <b-input-group size="sm">
      <b-form-input
        :id="id"
        v-model="formValue"
        type="text"
        :state="state"
        :disabled="disabled"
        v-on:blur="propagateValue"
      />
      <b-input-group-append is-text>
        {{ unitText }}
      </b-input-group-append>
    </b-input-group>
  </b-form-group>
</template>

<script>
import { 
  formatEngineHourLabel,
  formatEngineHour,
  isValidEngineHour,
  convertEngineHour
} from '@/libs/formatters';
import {
  BFormGroup,
  BFormInput,
  BInputGroup,
  BInputGroupAppend
} from 'bootstrap-vue';

export default {
  name: "EngineHours",
  props: [ 'label', 'value', 'displayFormat', 'disabled' ],
  components: {
    BFormGroup,
    BFormInput,
    BInputGroup,
    BInputGroupAppend
  },
  data() {
    return {
      formValue: null,
      state: null,
    }
  },
  computed: {
    inputGroupId: function() {
      return 'input-group-' + this.label.toLowerCase();
    },
    labelFor: function() {
      return 'input-' + this.label.toLowerCase();
    },
    id: function() {
      return 'input-' + this.label.toLowerCase();
    },
    unitText: function() {
      return formatEngineHourLabel(this.displayFormat);
    }
  },
  watch: {
    value: function(){
      this.formValue = formatEngineHour(this.value, this.displayFormat);
    },
    formValue: function(){
      if(!isValidEngineHour(this.formValue, this.displayFormat)){
        this.state = false;
      }else{
        this.state = null;
      }
    }
  },
  methods: {
    propagateValue: function(){
      this.$emit('input', convertEngineHour(this.formValue));
    }
  },
  created(){
    this.formValue = formatEngineHour(this.value, this.displayFormat);
  }
}
</script>