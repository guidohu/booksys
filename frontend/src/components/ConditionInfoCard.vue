<template>
  <sectioned-card-module>
    <template v-slot:header>
      <div class="row">
        <div class="col-8">
          <h5>Condition Info</h5>
        </div>
      </div>
    </template>
    <template v-slot:body>
      <div class="row">
        <div class="col-5">
          <i class="bi bi-sunrise"></i>
          Sunrise
        </div>
        <div class="col-7">
          {{ sunriseString }}
        </div>
      </div>
      <div class="row">
        <div class="col-5">
          <i class="bi bi-sunset-fill"></i>
          Sunset
        </div>
        <div class="col-7">
          {{ sunsetString }}
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script>
import * as dayjs from "dayjs";
import SectionedCardModule from "@/components/bricks/SectionedCardModule.vue";

export default {
  name: "ConditionInfoCard",
  components: {
    SectionedCardModule,
  },
  props: ["sunrise", "sunset"],
  computed: {
    sunriseString: function () {
      if (this.sunrise == null) {
        return "n/a";
      }
      return dayjs(this.sunrise).format("HH:mm");
    },
    sunsetString: function () {
      if (this.sunset == null) {
        return "n/a";
      }
      return dayjs(this.sunset).format("HH:mm");
    },
  },
};
</script>
