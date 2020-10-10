<template>
  <div v-if="errors != null && errors.length > 0">
    <b-alert variant="warning" show>
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="err in errors" :key="err">{{ err }}</li>
      </ul>
    </b-alert>
  </div>
  <div v-else>
    <b-alert variant="success" show>
      No issues found.
    </b-alert>
  </div>
</template>

<script>
import Vue from 'vue'

export default Vue.extend({
  name: "WarningBox",
  props: [ 'errors' ],
  mounted() {
    console.log("WarningBox just mounted");
    console.log("warningbox:", this.errors);
  }
})
</script>