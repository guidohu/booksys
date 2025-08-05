<template>
  <div :id="pieElementId" class="text-center full-width" />
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue";
import BooksysPie from "../libs/pie";
import dayjs from "dayjs";

const props = defineProps(["sessionData", "selectedSession", "properties", "pieId"]);
const emit = defineEmits(["selectHandler"]);

const pieSessions = ref([]);
const pieElementId = ref("someId");
const pie = ref(null);

const date = computed(() => {
  if (props.sessionData != null) {
    return dayjs(props.sessionData.window_start).format("YYYY-MM-DD");
  } else {
    return "unknown";
  }
});

watch(
  () => props.sessionData,
  () => {
    // repaint the pie upon any change
    repaint();
  }
);

watch(
  () => props.selectedSession,
  (newSelectedSession) => {
    if (newSelectedSession == null) {
      pie.value.resetSelection();
      return;
    }

    if (newSelectedSession.id != null) {
      for (let i = 0; i < pieSessions.value.length; i++) {
        if (pieSessions.value[i].id == newSelectedSession.id) {
          pie.value.selectSector(i);
          break;
        }
      }
    }
  }
);

watch(
  () => props.properties,
  () => {
    repaint();
  }
);

onMounted(() => {
  let el = document.getElementById(pieElementId.value);
  const newPie = new BooksysPie();
  pie.value = newPie;
  pieSessions.value = newPie.drawPie(
    el,
    props.sessionData,
    selectHandler,
    props.properties
  );
});

if (props.pieId != null) {
  pieElementId.value = "pie" + props.pieId;
}

function selectHandler(selectedId) {
  emit("selectHandler", pieSessions.value[selectedId]);
}

function repaint() {
  let el = document.getElementById(pieElementId.value);
  el.innerHTML = "";
  pieSessions.value = pie.value.drawPie(
    el,
    props.sessionData,
    selectHandler,
    props.properties
  );
}
</script>

<style scoped>
.full-width {
  width: 100%;
}
</style>
