<template>
  <div>
    <maintenance-log-form class="mb-4" />
    <maintenance-log-list />
  </div>
</template>

<script>
import MaintenanceLogForm from "@/components/MaintenanceLogForm";
import MaintenanceLogList from "@/components/MaintenanceLogList";

export default {
  name: "MaintenanceLogContainer",
  components: {
    MaintenanceLogForm,
    MaintenanceLogList,
  },
};
</script>
