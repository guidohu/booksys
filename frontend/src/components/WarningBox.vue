<template>
  <div v-if="errors != null && errors.length > 0">
    <div class="alert alert-warning alert-dismissible fade show" role="alert">
      <button
        v-if="dismissible"
        type="button"
        class="btn-close"
        data-bs-dismiss="alert"
        aria-label="Close"
      ></button>
      <b v-if="!hideLabel">Please correct the following error(s):</b>
      <div v-if="errors.length > 1">
        <ul>
          <li v-for="err in errors" :key="err">
            {{ err }}
          </li>
        </ul>
      </div>
      <div v-else>
        {{ errors[0] }}
      </div>
    </div>
  </div>
  <div v-else>
    <div class="alert alert-success show">No issues found.</div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps(["errors", "dismissible", "hideLabel"]);

const isDismissible = computed(() => {
  return props.dismissible != null && props.dismissible == "true";
});
</script>
