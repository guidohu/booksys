<template>
  <div class="vld-parent" :class="{ 'spinner-fill': fill }">
    <loading :active="active" :can-cancel="false" :is-full-page="fullPage" />
    <div :class="{ 'spinner-fill': fill }">
      <slot />
    </div>
  </div>
</template>

<script setup>
import Loading from "vue-loading-overlay";
import "vue-loading-overlay/dist/css/index.css";

defineProps({
  // Handed straight to vue-loading-overlay. Deliberately left untyped so that
  // an unset `fullPage` stays undefined and the library applies its own
  // default, instead of being cast to false by Vue's boolean handling.
  active: null,
  fullPage: null,
  fill: {
    type: Boolean,
    required: false,
    default: false,
  },
});
</script>

<style scoped>
/* Opt-in: lets the slot content stretch over the full height of the
   surrounding container instead of being sized by its content. Requires the
   parent to be a flex column with a height, otherwise this has no effect. */
.spinner-fill {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
}
</style>
