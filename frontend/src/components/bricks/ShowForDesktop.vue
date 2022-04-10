<template>
  <div v-if="isDesktop">
    <slot></slot>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
export default {
  name: "ShowForDesktop",
  props: ["minWidth"],
  data() {
    return {
      isDesktop: false,
    };
  },
  methods: {
    calculateDesktop: function () {
      if (this.minWidth != null){
        this.isDesktop = !BooksysBrowser.isMobileResponsive(this.minWidth);
      }else{
        this.isDesktop = !BooksysBrowser.isMobileResponsive();
      }
    },
  },
  mounted() {
    window.addEventListener("resize", this.calculateDesktop);
  },
  unmounted() {
    window.removeEventListener("resize", this.calculateDesktop);
  },
  beforeUpdate() {
    this.calculateDesktop();
  },
  created() {
    this.calculateDesktop();
  },
};
</script>
