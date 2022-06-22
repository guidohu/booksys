<template>
  <div v-if="isMobile">
    <slot></slot>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
export default {
  name: "ShowFoMobile",
  props: ["maxWidth"],
  data() {
    return {
      isMobile: false,
    };
  },
  methods: {
    calculateMobile: function () {
      if (this.maxWidth != null){
        this.isMobile = BooksysBrowser.isMobileResponsive(this.maxWidth);
      }else{
        this.isMobile = BooksysBrowser.isMobileResponsive();
      }
    },
  },
  mounted() {
    window.addEventListener("resize", this.calculateMobile);
  },
  unmounted() {
    window.removeEventListener("resize", this.calculateMobile);
  },
  beforeUpdate() {
    this.calculateMobile();
  },
  created() {
    this.calculateMobile();
  },
};
</script>
