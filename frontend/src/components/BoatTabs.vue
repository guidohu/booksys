<template>
  <div>
    <ul class="nav nav-tabs">
      <li class="nav-item">
        <a class="nav-link active" data-bs-toggle="tab" href="#engine-hours">
          <div class="bc_icon bc_icon_tab bc_icon_tacho align-middle nav-tab-icon" />
          Engine Hours
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-bs-toggle="tab" href="#fuel">
          <div class="bc_icon bc_icon_tab bc_icon_fuel align-middle nav-tab-icon" />
          Fuel
        </a>
      </li>
      <li class="nav-item">
        <a class="nav-link" data-bs-toggle="tab" href="#maintenance">
          <div class="bc_icon bc_icon_tab bc_icon_wrench align-middle nav-tab-icon" />
          Maintenance
        </a>
      </li>
    </ul>

    <div class="tab-content">
      <div class="tab-pane tab-class active" id="engine-hours">
        <engine-hour-log-container />
      </div>
      <div class="tab-pane tab-class" id="fuel">
        <fuel-log-container />
      </div>
      <div class="tab-pane tab-class" id="maintenance">
        <maintenance-log-container />
      </div>
    </div>
  </div>
</template>

<script>
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
};
</script>

<style scoped>
.nav-link {
  color: #bdbdbd;
}

a {
  color: #bdbdbd;
}

.tab-class {

}

.nav-tab-icon {
  display: inline;
}

@media (max-width: 992px) {
  .tab-class {
    max-height: 400px;
    height: 400px;
  }

  .nav-tab-icon {
    display: block;
  }
}
</style>
