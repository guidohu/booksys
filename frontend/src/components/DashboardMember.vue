<template>
  <div class="row mt-5">
    <div class="col-3" />
    <div class="col-6 text-center">
      <div class="row mt-4">
        <dashboard-button to="/schedule" colWidth="4" offsetWidth="2">
          <div class="bc_icon bc_icon_dashboard bc_icon_calendar" />
          MY CAL
        </dashboard-button>
        <dashboard-button to="/account" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_user" />
          ACCOUNT
        </dashboard-button>
      </div>
      <div class="row mt-4">
        <dashboard-button to="/info" colWidth="4" offsetWidth="2">
          <div class="bc_icon bc_icon_dashboard bc_icon_info" />
          INFO
        </dashboard-button>
        <dashboard-button to="/logout" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_logout" />
          LOGOUT
        </dashboard-button>
      </div>
    </div>
    <div cols="3" />
  </div>
</template>

<script>
import DashboardButton from "./bricks/DashboardButton.vue";

export default {
  name: "DashboardMember",
  components: {
    DashboardButton,
  },
};
</script>
