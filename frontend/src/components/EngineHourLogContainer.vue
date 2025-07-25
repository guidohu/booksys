<template>
  <div class="box">
    <engine-hour-log-form class="box-adjusted-content mb-4" />
    <engine-hour-log-list class="box-fill-content" />
  </div>
</template>

<script>
import EngineHourLogForm from "booksys/components/EngineHourLogForm.vue";
import EngineHourLogList from "booksys/components/EngineHourLogList.vue";

export default {
  name: "EngineHourLogContainer",
  components: {
    EngineHourLogForm,
    EngineHourLogList,
  },
};
</script>

<style scoped>
.box {
  display: flex;
  flex-flow: column;
  height: 100%;
}

.box-adjusted-content {
  flex: 0 0 auto;
}

.box-fill-content {
  flex: 1 1 auto;
  overflow: scroll;
}
</style>
