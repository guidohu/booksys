<template>
  <div class="card">
    <div v-if="!nobody" class="card-body">
      <h5 v-if="title" class="card-title">{{ title }}</h5>
      <div class="card-text">
        <slot></slot>
      </div>
    </div>
    <div v-if="nobody">
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
});
</script>
