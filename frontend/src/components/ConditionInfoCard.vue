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

<script>
import * as dayjs from "dayjs";
import SectionedCardModule from "@/components/bricks/SectionedCardModule.vue";

export default {
  name: "ConditionInfoCard",
  components: {
    SectionedCardModule,
  },
  data() {
    return {
      sunriseString: "n/a",
      sunsetString: "n/a",
    }
  },
  props: ["sunrise", "sunset"],
  watch: {
    sunrise: function(newVal, oldVal) {
      console.debug("ConditionInfoCard: sunrise changed to", newVal)
      this.setSunrise(newVal);
    },
    sunset: function(newVal, oldVal) {
      console.debug("ConditionInfoCard: sunset changed to", newVal)
      this.setSunset(newVal);
    },
  },
  methods: {
    setSunrise: function(time) {
      if (time === null) {
        this.sunriseString = "n/a";
        return;
      }
      this.sunriseString = dayjs.unix(time).format("HH:mm");
      console.debug("ConditionInfoCard: new sunriseString", this.sunriseString);
    },
    setSunset: function(time) {
      if (time === null) {
        this.sunsetString = "n/a";
        return;
      }
      this.sunsetString = dayjs.unix(time).format("HH:mm");
      console.debug("ConditionInfoCard: new sunsetString", this.sunsetString);
    }
  },
  created() {
    console.debug("ConditionInfoCard: sunrise", this.sunrise, "sunset", this.sunset);
    this.setSunrise(this.sunrise);
    this.setSunset(this.sunset);
  }
  // computed: {
  //   sunriseString: function () {
  //     if (this.sunrise == null) {
  //       return "n/a";
  //     }
  //     return dayjs(this.sunrise).format("HH:mm");
  //   },
  //   sunsetString: function () {
  //     if (this.sunset == null) {
  //       return "n/a";
  //     }
  //     return dayjs(this.sunset).format("HH:mm");
  //   },
  // },
};
</script>
