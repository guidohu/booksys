<template>
  <div v-if="errors != null && errors.length > 0" class="text-left">
    <div class="alert alert-warning alert-dismissible fade show" role="alert">
      <button
        v-if="dismissible"
        type="button"
        class="btn-close"
        data-bs-dismiss="alert"
        aria-label="Close"
      ></button>
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="err in errors" :key="err">
          {{ err }}
        </li>
      </ul>
    </div>
  </div>
  <div v-else>
    <div class="alert alert-success show">No issues found.</div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps(["errors", "dismissible"]);

const isDismissible = computed(() => {
  return props.dismissible != null && props.dismissible == "true";
});
</script>
