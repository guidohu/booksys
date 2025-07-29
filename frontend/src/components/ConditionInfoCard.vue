<template>
  <sectioned-card-module title="Condition Info">
    <template v-slot:body>
      <div class="row">
        <div class="col-4">
          <i class="bi bi-sunrise"></i>
          Sunrise
        </div>
        <div class="col-8">
          {{ sunriseString }}
        </div>
      </div>
      <div class="row">
        <div class="col-4">
          <i class="bi bi-sunset-fill"></i>
          Sunset
        </div>
        <div class="col-8">
          {{ sunsetString }}
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script setup>
import dayjs from "dayjs";
import SectionedCardModule from "booksys/components/bricks/SectionedCardModule.vue";
import { ref, watch } from "vue";

const props = defineProps({
  sunrise: {
    type: Number,
    required: true,
  },
  sunset: {
    type: Number,
    required: true,
  },
});

const sunriseString = ref("n/a");
const sunsetString = ref("n/a");

const setSunrise = (time) => {
  if (time === null) {
    sunriseString.value = "n/a";
    return;
  }
  sunriseString.value = dayjs.unix(time).format("HH:mm");
};

const setSunset = (time) => {
  if (time === null) {
    sunsetString.value = "n/a";
    return;
  }
  sunsetString.value = dayjs.unix(time).format("HH:mm");
};

watch(() => props.sunrise, (newVal) => {
  console.debug("ConditionInfoCard: sunrise changed to", newVal);
  setSunrise(newVal);
});

watch(() => props.sunset, (newVal) => {
  console.debug("ConditionInfoCard: sunset changed to", newVal);
  setSunset(newVal);
});

console.debug("ConditionInfoCard: sunrise", props.sunrise, "sunset", props.sunset);
setSunrise(props.sunrise);
setSunset(props.sunset);
</script>
