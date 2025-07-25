<template>
  <div>
    <maintenance-log-form class="mb-4" />
    <maintenance-log-list />
  </div>
</template>

<script>
import MaintenanceLogForm from "booksys/components/MaintenanceLogForm.vue";
import MaintenanceLogList from "booksys/components/MaintenanceLogList.vue";

export default {
  name: "MaintenanceLogContainer",
  components: {
    MaintenanceLogForm,
    MaintenanceLogList,
  },
};
</script>
