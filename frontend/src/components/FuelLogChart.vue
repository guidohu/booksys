<template>
  <div>
    <canvas id="fuelChart" width="400" height="200" />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue";
import { useStore } from "vuex";
import Chart from "chart.js/dist/Chart";
import "chartjs-plugin-colorschemes/src/plugins/plugin.colorschemes";
import { ClassicBlue7 } from "chartjs-plugin-colorschemes/src/colorschemes/colorschemes.tableau";
import groupBy from "lodash/groupBy";
import max from "lodash/max";
import sum from "lodash/sum";
import dayjs from "dayjs";
import dayOfYear from "dayjs/plugin/dayOfYear";

dayjs.extend(dayOfYear);

const store = useStore();

const fuelChartObject = ref(null);
const datasets = ref([]);
const options = ref({
  title: {
    display: true,
    text: "Average Fuel Consumption",
    fontStyle: "none",
  },
  legend: {
    position: "bottom",
  },
  tooltips: {
    enabled: true,
    callbacks: {
      title: function () {
        return "Average Fuel Consumption";
      },
      label: function (tooltipItem, data) {
        let time =
          data.datasets[tooltipItem.datasetIndex].data[
            tooltipItem.index
          ].x.format("DD. MMM");
        let liters =
          data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index]
            .y;
        liters = Math.round(liters * 10) / 10;
        return time + " - " + liters + " L";
      },
      afterLabel: function () {
        return "";
      },
    },
  },
  scales: {
    xAxes: [
      {
        type: "time",
        distribution: "series",
        time: {
          unit: "day",
        },
      },
    ],
    yAxis: [
      {
        ticks: {
          min: 0,
          max: 40,
          suggestedMin: 0,
          suggestedMax: 40,
        },
        scaleLabel: {
          display: true,
          labelString: "liter/hour",
        },
      },
    ],
  },
  maintainAspectRatio: false,
  responsive: true,
  plugins: {
    colorschemes: {
      scheme: ClassicBlue7,
    },
  },
});

const getFuelLog = computed(() => store.getters["boat/getFuelLog"]);

watch(getFuelLog, (newValues) => {
  setDatasets(newValues);
  drawChart();
});

onMounted(() => {
  setDatasets(getFuelLog.value);
  drawChart();
});

function setDatasets(fuelLog) {
  if (fuelLog == null || fuelLog.length === 0) {
    datasets.value = [];
    return;
  }

  const fuelLogByYear = groupBy(fuelLog, (x) =>
    Number(dayjs.unix(x.timestamp).format("YYYY"))
  );

  const years = Object.keys(fuelLogByYear).sort().reverse();
  const displayYears = years.slice(0, 3);

  const newDatasets = [];
  for (let i = 0; i < displayYears.length; i++) {
    const year = displayYears[i];
    const fuelValues = fuelLogByYear[year];
    const fuelValuesPerDay = groupBy(fuelValues, (x) =>
      dayjs.unix(x.timestamp).dayOfYear()
    );

    const dataset = {
      label: year,
      data: [],
      showLine: true,
      fill: false,
      hidden: i !== 0,
    };

    let lastEngineHours = 0;
    const dayKeys = Object.keys(fuelValuesPerDay).sort((a, b) => a - b);
    for (let j = 0; j < dayKeys.length; j++) {
      const dayNumber = dayKeys[j];
      const fuelEntries = fuelValuesPerDay[dayNumber];

      const engineHours = max(
        fuelEntries.map((entry) => Number(entry.engine_hours))
      );
      const fuelConsumption = sum(
        fuelEntries.map((entry) => Number(entry.liters))
      );

      if (lastEngineHours === 0) {
        lastEngineHours = engineHours;
        continue;
      }
      if (lastEngineHours > engineHours) {
        lastEngineHours = engineHours;
        continue;
      }

      const engineHoursDiff = engineHours - lastEngineHours;
      if (engineHoursDiff === 0) {
        console.log("FuelLogChart: division by 0, skip value");
        lastEngineHours = engineHours;
        continue;
      }

      const fuelConsumptionPerHour = fuelConsumption / engineHoursDiff;

      dataset.data.push({
        x: dayjs(
          dayjs.unix(fuelEntries[0].timestamp).set("year", 1970)
        ).toDate(),
        y: fuelConsumptionPerHour,
      });

      lastEngineHours = engineHours;
    }

    newDatasets.push(dataset);
  }

  datasets.value = newDatasets;
}

function drawChart() {
  if (fuelChartObject.value) {
    fuelChartObject.value.destroy();
  }
  var ctx = document.getElementById("fuelChart").getContext("2d");
  fuelChartObject.value = new Chart(ctx, {
    type: "scatter",
    data: {
      datasets: datasets.value,
    },
    options: options.value,
  });
}
</script>
