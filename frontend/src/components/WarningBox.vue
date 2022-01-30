<template>
  <div
    v-if="errors != null && errors.length > 0"
    class="text-left"
  >
    <b-alert
      variant="warning"
      show
      :dismissible="isDismissible"
      @dismissed="dismissedHandler"
    >
      <b>Please correct the following error(s):</b>
      <ul>
        <li
          v-for="err in errors"
          :key="err"
        >
          {{ err }}
        </li>
      </ul>
    </b-alert>
  </div>
  <div v-else>
    <b-alert
      variant="success"
      show
    >
      No issues found.
    </b-alert>
  </div>
</template>

<script>
import { BAlert } from "bootstrap-vue";

export default {
  name: "WarningBox",
  components: {
    BAlert,
  },
  props: ["errors", "dismissible"],
  computed: {
    isDismissible: function () {
      return this.dismissible != null && this.dismissible == "true";
    },
  },
  mounted() {
    console.log("WarningBox just mounted");
    console.log("warningbox:", this.errors);
  },
  methods: {
    dismissedHandler: function () {
      this.$emit("dismissed");
    },
  },
};
</script>
