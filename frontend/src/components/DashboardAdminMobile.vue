<template>
  <div>
    <navbar-mobile title="Dashboard" role="admin-home" />
    <div class="bk-page bk-bento">
      <!-- Today leads: it is the only destination whose content changes
           during the day, so it gets the full width and says what is behind
           it rather than just naming itself. -->
      <dashboard-tile
        to="/today"
        label="Today"
        :hint="todayHint"
        icon="bi-water"
        wide
      />

      <dashboard-tile to="/calendar" label="Book" icon="bi-calendar-plus" />
      <dashboard-tile to="/ride" label="Ride" icon="bi-stopwatch" />
      <dashboard-tile
        to="/schedule"
        label="My calendar"
        icon="bi-calendar-week"
      />
      <dashboard-tile to="/boat" label="Boat" icon="bi-speedometer2" />
      <dashboard-tile to="/account" label="Account" icon="bi-person" />
      <dashboard-tile to="/admin" label="Admin" icon="bi-sliders" />

      <dashboard-tile
        to="/info"
        label="Info"
        hint="Location and map"
        icon="bi-geo-alt"
        wide
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import NavbarMobile from "booksys/components/NavbarMobile.vue";
import DashboardTile from "booksys/components/bricks/DashboardTile.vue";

const props = defineProps({
  sessionData: {
    type: Object,
    required: true,
  },
});

// The sessions pie used to sit here, but at the size a tile allows it was a
// dark disc with unreadable slivers, and it pulled the whole raphael chart
// bundle into the first screen after login. The count says the same thing in
// a glance and costs nothing.
const todayHint = computed(() => {
  const sessions = props.sessionData && props.sessionData.sessions;
  if (!Array.isArray(sessions)) {
    return "Sessions on the water";
  }
  if (sessions.length === 0) {
    return "Nothing booked yet";
  }
  return sessions.length === 1
    ? "1 session today"
    : `${sessions.length} sessions today`;
});
</script>
