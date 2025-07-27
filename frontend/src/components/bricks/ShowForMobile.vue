<template>
  <div v-if="isMobile">
    <slot></slot>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, onBeforeUpdate } from "vue";
import { BooksysBrowser } from "booksys/libs/browser";

const props = defineProps({
  maxWidth: {
    type: String,
    required: false,
  },
});

const isMobile = ref(false);

const calculateMobile = () => {
  if (props.maxWidth != null) {
    isMobile.value = BooksysBrowser.isMobileResponsive(props.maxWidth);
  } else {
    isMobile.value = BooksysBrowser.isMobileResponsive();
  }
};

onMounted(() => {
  window.addEventListener("resize", calculateMobile);
  calculateMobile();
});

onUnmounted(() => {
  window.removeEventListener("resize", calculateMobile);
});

onBeforeUpdate(() => {
  calculateMobile();
});
</script>
