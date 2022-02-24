<template>
  <div v-if="isDesktop">
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: "ShowForDesktop",
  props: ["minWidth"],
  data() {
    return {
      screenWidth: 0,
      isDesktop: false,
      defaultMinWidth: 992,
    }
  },
  methods: {
    calculateDesktop: function() {
      const screenWidth = window.innerWidth;
      this.screenWidth = screenWidth;

      if(this.minWidth != null && this.screenWidth > this.minWidth) {
        this.isDesktop = true;
        return;
      }
      if(this.minWidth == null && this.screenWidth > this.defaultMinWidth) {
        this.isDesktop = true;
        return;
      }
      this.isDesktop = false;
      return;
    }
  },
  mounted() {
    window.addEventListener('resize', this.calculateDesktop);
  },
  unmounted() {
    window.removeEventListener('resize', this.calculateDesktop);
  },
  beforeUpdate() {
    this.calculateDesktop();
  },
  created() {
    this.calculateDesktop();
  }
}
</script>