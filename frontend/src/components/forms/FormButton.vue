<template>
  <div class="row text-right mt-2">
    <div class="col-9 offset-3">
      <div class="d-grid gap-2">
        <button
          :class="getClass"
          :type="type ? type : 'button'"
          @click="$emit('click', $event)"
        >
          <slot></slot>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps(["type", "btnStyle", "btnSize"]);
defineEmits(["click"]);

const getClass = computed(() => {
  const style =
    props.btnStyle == null
      ? "btn-outline-info"
      : "btn-outline-" + props.btnStyle;
  const size =
    props.btnSize == "small"
      ? "btn-sm"
      : props.btnSize == "large"
        ? "btn-lg"
        : "";
  return "btn " + style + " " + size;
});
</script>
