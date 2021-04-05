<template>
  <div>
    <MaintenanceLogForm class="mb-4" />
    <MaintenanceLogList />
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
