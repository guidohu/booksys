<template>
  <div>
    <canvas
      id="fuelChart"
      width="400"
      height="200"
    />
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import Chart from "chart.js/dist/Chart";
import "chartjs-plugin-colorschemes/src/plugins/plugin.colorschemes";
import { ClassicBlue7 } from "chartjs-plugin-colorschemes/src/colorschemes/colorschemes.tableau";
import groupBy from "lodash/groupBy";
import max from "lodash/max";
import sum from "lodash/sum";
import * as dayjs from "dayjs";
import * as dayjsDayOfYear from "dayjs/plugin/dayOfYear";

dayjs.extend(dayjsDayOfYear);

export default {
  name: "FuelLogChart",
  data() {
    return {
      fuelChartObject: null,
      datasets: [],
      options: {
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
      },
    };
  },
  computed: {
    ...mapGetters("boat", ["getFuelLog"]),
  },
  watch: {
    getFuelLog: function (newValues) {
      this.setDatasets(newValues);
      this.drawChart();
    },
  },
  mounted() {
    this.setDatasets(this.getFuelLog);
    this.drawChart();
  },
  methods: {
    setDatasets: function (fuelLog) {
      if (fuelLog == null || fuelLog == 0) {
        this.datasets = [];
        return;
      }

      // split by year
      const fuelLogByYear = groupBy(fuelLog, (x) =>
        Number(dayjs.unix(x.timestamp).format("YYYY"))
      );

      // get all the years with the latest ones first
      const years = Object.keys(fuelLogByYear).sort().reverse();
      const displayYears = years.slice(0, 3); // only display the last 3 years

      // build the datasets per year
      const datasets = [];
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
          hidden: i != 0, // hide the stats for previous years by default
        };

        let lastEngineHours = 0;
        const dayKeys = Object.keys(fuelValuesPerDay).sort((a, b) => a - b);
        for (let j = 0; j < dayKeys.length; j++) {
          const dayNumber = dayKeys[j];
          const fuelEntries = fuelValuesPerDay[dayNumber];

          // calculate the engine hours of that day and the liters fueled in total
          const engineHours = max(
            fuelEntries.map((entry) => Number(entry.engine_hours))
          );
          const fuelConsumption = sum(
            fuelEntries.map((entry) => Number(entry.liters))
          );

          // first value, we cannot calculate an average
          if (lastEngineHours == 0) {
            lastEngineHours = engineHours;
            continue;
          }
          // check if engine hours decreased
          if (lastEngineHours > engineHours) {
            lastEngineHours = engineHours;
            continue;
          }

          // calculate engine hour difference
          const engineHoursDiff = engineHours - lastEngineHours;
          if (engineHoursDiff == 0) {
            console.log("FuelLogChart: division by 0, skip value");
            lastEngineHours = engineHours;
            continue;
          }

          // calculate fuel consumption per engine hour
          const fuelConsumptionPerHour = fuelConsumption / engineHoursDiff;

          // add the entry to the graph dataset
          dataset.data.push({
            x: dayjs(
              dayjs.unix(fuelEntries[0].timestamp).set("year", 1970)
            ).toDate(),
            y: fuelConsumptionPerHour,
          });

          // update the lastEngineHour reference
          lastEngineHours = engineHours;
        }

        // add the entire dataset to the collection of datasets
        datasets.push(dataset);
      }

      this.datasets = datasets;
    },
    drawChart: function () {
      var ctx = document.getElementById("fuelChart").getContext("2d");
      const options = this.options;
      const datasets = this.datasets;
      this.fuelChartObject = new Chart(ctx, {
        type: "scatter",
        data: {
          datasets: datasets,
        },
        options: options,
      });
    },
  },
};
</script>
