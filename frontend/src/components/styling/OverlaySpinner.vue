<template>
  <div class="vld-parent">
    <loading :active="active" :can-cancel="false" :is-full-page="fullPage" />
    <div>
      <slot />
    </div>
  </div>
</template>

<script>
import { defineComponent } from "vue";
import Loading from "vue-loading-overlay";
import "vue-loading-overlay/dist/vue-loading.css";

export default defineComponent({
  name: "OverlaySpinner",
  components: {
    Loading,
  },
  props: ["active", "fullPage"],
});
</script>
