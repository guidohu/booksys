<template>
  <div class="boat-tabs">
    <ul class="nav nav-tabs">
      <li class="nav-item">
        <a
          class="nav-link active"
          data-bs-toggle="tab"
          href="#engine-hours"
          @click="selectTab('engine-hours')"
        >
          <div
            class="bc_icon bc_icon_tab bc_icon_tacho align-middle nav-tab-icon"
          />
          Engine Hours
        </a>
      </li>
      <li class="nav-item">
        <a
          class="nav-link"
          data-bs-toggle="tab"
          href="#fuel"
          @click="selectTab('fuel')"
        >
          <div
            class="bc_icon bc_icon_tab bc_icon_fuel align-middle nav-tab-icon"
          />
          Fuel
        </a>
      </li>
      <li class="nav-item">
        <a
          class="nav-link"
          data-bs-toggle="tab"
          href="#maintenance"
          @click="selectTab('maintenance')"
        >
          <div
            class="bc_icon bc_icon_tab bc_icon_wrench align-middle nav-tab-icon"
          />
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

<script setup>
import EngineHourLogContainer from "booksys/components/EngineHourLogContainer.vue";
import FuelLogContainer from "booksys/components/FuelLogContainer.vue";
import MaintenanceLogContainer from "booksys/components/MaintenanceLogContainer.vue";
import { Tab } from "bootstrap";
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const selectedTab = ref(null);

const selectTab = (tabName) => {
  router.push("/boat/" + tabName);
};

selectedTab.value = route.params.tab;

onMounted(() => {
  if (
    selectedTab.value == "engine-hours" ||
    selectedTab.value == "fuel" ||
    selectedTab.value == "maintenance"
  ) {
    // show specific tab
    let tabTrigger = document.querySelector(
      "[href='#" + selectedTab.value + "']",
    );
    let tab = new Tab(tabTrigger);
    tab.show();
  }
});
</script>

<style scoped>
.nav-link {
  color: #bdbdbd;
  /* The icon is a block element on mobile, where it sits above the label,
     so center both of them within the tab */
  text-align: center;
}

a {
  color: #bdbdbd;
}

/* The tab strip keeps its own height, the active pane takes the rest */
.boat-tabs {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
}

.nav-tabs {
  flex: 0 0 auto;
}

.tab-content {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
}

.tab-class {
  flex: 1 1 auto;
  min-height: 0;
}

/* Bootstrap sets `display: block` on the active pane, this selector is
   specific enough to turn it into a flex container so that the content
   inside can stretch over the pane instead of overflowing it */
.tab-content > .tab-class.active {
  display: flex;
  flex-direction: column;
}

.nav-tab-icon {
  display: inline;
}

@media (max-width: 992px) {
  .nav-tab-icon {
    display: block;
  }
}
</style>
