<template>
  <div v-if="isDesktop">
    <slot></slot>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, onBeforeUpdate } from "vue";
import { BooksysBrowser } from "booksys/libs/browser";

const props = defineProps({
  minWidth: {
    type: String,
    required: false,
  },
});

const isDesktop = ref(false);

const calculateDesktop = () => {
  if (props.minWidth != null) {
    isDesktop.value = !BooksysBrowser.isMobileResponsive(props.minWidth);
  } else {
    isDesktop.value = !BooksysBrowser.isMobileResponsive();
  }
};

onMounted(() => {
  window.addEventListener("resize", calculateDesktop);
  calculateDesktop();
});

onUnmounted(() => {
  window.removeEventListener("resize", calculateDesktop);
});

onBeforeUpdate(() => {
  calculateDesktop();
});
</script>
