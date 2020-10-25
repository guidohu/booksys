<template>
  <div>
    <MaintenanceLogForm class="mb-4"/>
    <MaintenanceLogList/>
  </div>
</template>

<script>
import Vue from 'vue';
import MaintenanceLogForm from '@/components/MaintenanceLogForm';
import MaintenanceLogList from '@/components/MaintenanceLogList';

export default Vue.extend({
  name: "MaintenanceLogContainer",
  components: {
    MaintenanceLogForm,
    MaintenanceLogList
  }
})
</script>