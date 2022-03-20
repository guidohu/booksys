<template>
  <div class="card">
    <div v-if="$slots.header" class="card-header">
      <slot name="header"></slot>
    </div>
    <div v-if="$slots.body" :class="getBodyClass()">
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
  props: [ "scrollable" ],
  methods: {
    getBodyClass: function() {
      let style = "card-body margin padding";
      if(this.scrollable) {
        style += " scrollable";
      }
      return style;
    }
  }
};
</script>

<style scoped>
.scrollable {
  overflow-y: scroll;
}

.margin {
  margin: 1rem;
}

.padding {
  padding: 0rem;
}
</style>
