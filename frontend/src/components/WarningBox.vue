<template>
  <div v-if="errors != null && errors.length > 0" class="text-left">
    <div class="alert alert-warning alert-dismissible fade show" role="alert">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="err in errors" :key="err">
          {{ err }}
        </li>
      </ul>
    </div>
  </div>
  <div v-else>
    <div class="alert alert-success show"> No issues found. </div>
  </div>
</template>

<script>

export default {
  name: "WarningBox",

  props: ["errors", "dismissible"],
  computed: {
    isDismissible: function () {
      return this.dismissible != null && this.dismissible == "true";
    },
  },
  mounted() {
    console.log("warningbox:", this.errors);
  },
  // methods: {
  //   dismissedHandler: function () {
  //     this.$emit("dismissed");
  //   },
  // },
};
</script>
