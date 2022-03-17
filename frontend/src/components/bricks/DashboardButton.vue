<template>
  <div :class="colClass">
    <router-link tag="button" class="btn btn-light btn-dashboard" :to="to">
      <slot></slot>
    </router-link>
  </div>
</template>

<script>
import { defineComponent } from "vue";

export default defineComponent({
  name: "DashboardButton",
  props: ["to", "colWidth", "offsetWidth"],
  computed: {
    colClass: function () {
      if (this.colWidth == null) {
        return "col-12 text-center";
      }
      if (this.colWidth != null && this.offsetWidth == null) {
        return "col-" + this.colWidth + " text-center";
      }
      return (
        "col-" +
        this.colWidth +
        " " +
        "offset-sm-" +
        this.offsetWidth +
        " text-center"
      );
    },
  },
});
</script>
