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

<script>
export default {
  name: "CardModule",
  props: ["title", "nobody"],
};
</script>
