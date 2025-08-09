<template>
  <div class="vld-parent">
    <loading :active="active" :can-cancel="false" :is-full-page="fullPage" />
    <div>
      <slot />
    </div>
  </div>
</template>

<script setup>
import Loading from "vue-loading-overlay";
import "vue-loading-overlay/dist/css/index.css";

defineProps(["active", "fullPage"]);
</script>
