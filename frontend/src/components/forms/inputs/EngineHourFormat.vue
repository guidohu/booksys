<template>
  <b-form-group
    id="input-group-engine-hour-format"
    label="Engine Hour Format*"
    label-for="input-engine-hour-format"
    label-cols="3"
    description="select the format your engine hour display uses"
  >
    <b-form-select
      v-model="formSelected"
      :options="options"
    />
  </b-form-group>
</template>

<script>
import { BFormGroup, BFormSelect } from "bootstrap-vue";

export default {
  name: "EngineHourFormat",
  components: {
    BFormGroup,
    BFormSelect,
  },
  props: ["selected"],
  data() {
    return {
      formSelected: "hh.h",
      options: [
        { value: "hh.h", text: "hh.h  - such as 9.7" },
        { value: "hh:mm", text: "hh:mm - such as 9:42" },
      ],
    };
  },
  watch: {
    formSelected: function () {
      this.$emit("update:selected", this.formSelected);
    },
    selected: function () {
      this.setFormSelected(this.selected);
    },
  },
  created() {
    this.setFormSelected(this.selected);
  },
  methods: {
    setFormSelected: function (format) {
      if (format == null) {
        console.log("no default engine hour format specified");
        this.formSelected = "hh.h";
      } else if (format == "hh.h" || format == "hh:mm") {
        this.formSelected = format;
      } else {
        this.formSelected = "hh.h";
      }
    },
  },
};
</script>
