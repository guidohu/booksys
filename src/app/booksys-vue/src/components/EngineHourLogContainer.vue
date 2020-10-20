<template>
  <div>
    <EngineHourLogForm class="mb-4"/>
    <EngineHourLogList/>
  </div>
</template>

<script>
import Vue from 'vue';
import EngineHourLogForm from '@/components/EngineHourLogForm';
import EngineHourLogList from '@/components/EngineHourLogList';

export default Vue.extend({
  name: "EngineHourLogContainer",
  components: {
    EngineHourLogForm,
    EngineHourLogList
  }
})
</script>