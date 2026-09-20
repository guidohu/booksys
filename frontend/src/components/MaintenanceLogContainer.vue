<template>
  <div class="box">
    <maintenance-log-form class="box-adjusted-content mb-4" />
    <maintenance-log-list class="box-fill-content" />
  </div>
</template>

<script setup>
import MaintenanceLogForm from "booksys/components/MaintenanceLogForm.vue";
import MaintenanceLogList from "booksys/components/MaintenanceLogList.vue";
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
