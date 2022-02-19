<template>
  <div class="card">
    <div v-if="$slots.header" class="card-header">
      <slot name="header"></slot>
    </div>
    <div v-if="$slots.body" class="card-body">
      <slot name="body"></slot>
    </div>
    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "SectionedCardModule",
  // methods: {
  //   hasHeader: function() {
  //     return this.$slots.header() == true;
  //   },
  //   hasBody: function() {
  //     return (this.$slots.body == true);
  //   },
  //   hasFooter: function() {
  //     return (this.$slots.footer == true);
  //   }
  // }
};
</script>
