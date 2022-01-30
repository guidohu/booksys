<template>
  <b-tabs content-class="mt-3">
    <b-tab active>
      <div slot="title">
        <div :class="titleClassEngineHours" />
        Engine Hours
      </div>
      <div :class="tabClass">
        <EngineHourLogContainer />
      </div>
    </b-tab>
    <b-tab>
      <div slot="title">
        <div :class="titleClassFuel" />
        Fuel
      </div>
      <div :class="tabClass">
        <FuelLogContainer />
      </div>
    </b-tab>
    <b-tab>
      <div slot="title">
        <div :class="titleClassMaintenance" />
        Maintenance
      </div>
      <div :class="tabClass">
        <MaintenanceLogContainer />
      </div>
    </b-tab>
  </b-tabs>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import EngineHourLogContainer from "@/components/EngineHourLogContainer";
import FuelLogContainer from "@/components/FuelLogContainer";
import MaintenanceLogContainer from "@/components/MaintenanceLogContainer";
import { BTabs, BTab } from "bootstrap-vue";

export default {
  name: "BoatTabs",
  components: {
    EngineHourLogContainer,
    FuelLogContainer,
    MaintenanceLogContainer,
    BTabs,
    BTab,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
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
  overflow: scroll;
}
</style>
