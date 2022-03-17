<template>
  <div v-if="isMobile">
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: "ShowFoMobile",
  props: ["maxWidth"],
  data() {
    return {
      screenWidth: 0,
      isMobile: false,
      defaultMaxWidth: 992,
    };
  },
  methods: {
    calculateMobile: function () {
      let screenWidth = window.innerWidth;
      this.screenWidth = screenWidth;
      if (this.maxWidth != null && this.screenWidth > this.maxWidth) {
        this.isMobile = false;
        return;
      }
      if (this.maxWidth == null && this.screenWidth > this.defaultMaxWidth) {
        this.isMobile = false;
        return;
      }
      this.isMobile = true;
      return;
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
