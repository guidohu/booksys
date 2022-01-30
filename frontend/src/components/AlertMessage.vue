<template>
  <div class="card">
    <div
      class="alert alert-danger"
      role="alert"
    >
      {{ alertMessage }}
    </div>
  </div>
</template>

<script>
export default {
  name: "AlertMessage",
  props: ["alertMessage"],
};
</script>
