<template>
  <div>
    <ul class="nav nav-tabs">
      <li class="nav-item">
        <a class="nav-link active" data-bs-toggle="tab" href="#engine-hours">
          <div :class="titleClassEngineHours" />
          Engine Hours
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-bs-toggle="tab" href="#fuel">
          <div :class="titleClassFuel" />
          Fuel
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-bs-toggle="tab" href="#maintenance">
          <div :class="titleClassMaintenance" />
          Maintenance
        </a>
      </li>
    </ul>

    <div class="tab-content">
      <div :class="'tab-pane ' + tabClass + ' active'" id="engine-hours">
        <engine-hour-log-container />
      </div>
      <div :class="'tab-pane ' + tabClass" id="fuel">
        <fuel-log-container />
      </div>
      <div :class="'tab-pane ' + tabClass" id="maintenance">
        <maintenance-log-container />
      </div>
    </div>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import EngineHourLogContainer from "@/components/EngineHourLogContainer";
import FuelLogContainer from "@/components/FuelLogContainer";
import MaintenanceLogContainer from "@/components/MaintenanceLogContainer";

export default {
  name: "BoatTabs",
  components: {
    EngineHourLogContainer,
    FuelLogContainer,
    MaintenanceLogContainer,
  },
  computed: {
    titleClassEngineHours: function () {
      if (BooksysBrowser.isMobile()) {
        return "bc_icon bc_icon_tab bc_icon_tacho align-middle";
      }
      return "bc_icon bc_icon_tab bc_icon_tacho inline align-middle";
    },
    titleClassFuel: function () {
      if (BooksysBrowser.isMobile()) {
        return "bc_icon bc_icon_tab bc_icon_fuel align-middle";
      }
      return "bc_icon bc_icon_tab bc_icon_fuel inline align-middle";
    },
    titleClassMaintenance: function () {
      if (BooksysBrowser.isMobile()) {
        return "bc_icon bc_icon_tab bc_icon_wrench align-middle";
      }
      return "bc_icon bc_icon_tab bc_icon_wrench inline align-middle";
    },
    tabClass: function () {
      if (!BooksysBrowser.isMobile()) {
        return "boat-tab-limited-height";
      } else {
        return "";
      }
    },
  },
};
</script>

<style scoped>
.boat-tab-limited-height {
  max-height: 400px;
  height: 400px;
}

.nav-link {
  color: #bdbdbd;
}

a {
  color: #bdbdbd;
}
</style>
