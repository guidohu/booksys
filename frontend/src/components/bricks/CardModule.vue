<template>
  <div class="card">
    <div v-if="!nobody" class="card-body">
      <h5 v-if="title" class="card-title">{{ title }}</h5>
      <div class="card-text">
        <slot></slot>
      </div>
    </div>
    <div v-if="nobody" :class="{ 'card-fill': fill }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: false,
  },
  nobody: {
    type: Boolean,
    required: false,
  },
  fill: {
    type: Boolean,
    required: false,
  },
});
</script>

<style scoped>
/* Opt-in: lets the slot content stretch over the full height of the card
   instead of being sized by its content. Requires the card itself to have
   a height, otherwise this has no effect. */
.card-fill {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
}
</style>
