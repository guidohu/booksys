<template>
  <div>
    <FuelLogForm/>
    <FuelLogChart/>
    <FuelLogList/>
  </div>
</template>

<script>
import Vue from 'vue';
import FuelLogForm from '@/components/FuelLogForm';
import FuelLogChart from '@/components/FuelLogChart';
import FuelLogList from '@/components/FuelLogList';

export default Vue.extend({
  name: "FuelLogContainer",
  components: {
    FuelLogForm,
    FuelLogChart,
    FuelLogList
  }
})
</script>