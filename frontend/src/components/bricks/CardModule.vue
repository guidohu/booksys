<template>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">{{ title }}</h5>
      <div class="card-text">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CardModule",
  props: ["title"],
};
</script>
