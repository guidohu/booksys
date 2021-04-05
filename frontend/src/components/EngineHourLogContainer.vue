<template>
  <div>
    <EngineHourLogForm class="mb-4" />
    <EngineHourLogList />
  </div>
</template>

<script>
import EngineHourLogForm from "@/components/EngineHourLogForm";
import EngineHourLogList from "@/components/EngineHourLogList";

export default {
  name: "EngineHourLogContainer",
  components: {
    EngineHourLogForm,
    EngineHourLogList,
  },
};
</script>
