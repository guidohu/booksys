<template>
  <div class="card m-5 p-5">
    <div class="alert alert-danger" role="alert">
      {{ alertMessage }}
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  alertMessage: {
    type: String,
    required: true,
  },
});
</script>
