<template>
  <div class="box">
    <engine-hour-log-form class="box-adjusted-content mb-4" />
    <engine-hour-log-list class="box-fill-content" />
  </div>
</template>

<script setup>
import EngineHourLogForm from "booksys/components/EngineHourLogForm.vue";
import EngineHourLogList from "booksys/components/EngineHourLogList.vue";
</script>

<style scoped>
.box {
  display: flex;
  flex-flow: column;
  flex: 1 1 auto;
  min-height: 0;
  /* Only kicks in when the box is too short for the form plus a usable
     bit of list, e.g. on a phone in landscape */
  overflow-y: auto;
}

.box-adjusted-content {
  flex: 0 0 auto;
}

/* Only the list scrolls, the form above it stays put */
.box-fill-content {
  flex: 1 1 auto;
  min-height: 6rem;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
