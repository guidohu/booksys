<template>
  <div :class="colClass">
    <router-link tag="button" class="btn btn-light btn-dashboard" :to="to">
      <slot></slot>
    </router-link>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  to: {
    type: String,
    required: true,
  },
  colWidth: {
    type: String,
    required: false,
  },
  offsetWidth: {
    type: String,
    required: false,
  },
});

const colClass = computed(() => {
  if (props.colWidth == null) {
    return "col-12 text-center";
  }
  if (props.colWidth != null && props.offsetWidth == null) {
    return "col-" + props.colWidth + " text-center";
  }
  return (
    "col-" +
    props.colWidth +
    " " +
    "offset-sm-" +
    props.offsetWidth +
    " text-center"
  );
});
</script>
