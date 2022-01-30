<template>
  <div class="row mt-5">
    <div class="col-3" />
    <div class="col-6 text-center">
      <div class="row mt-4">
        <dashboard-button to="/today" colWidth="4">
          <booksys-pie :session-data="sessionData" :properties="properties" />
          TODAY
        </dashboard-button>
        <dashboard-button to="/boat" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_boat" />
          BOAT
        </dashboard-button>
        <dashboard-button to="/ride" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_wakeboard" />
          RIDE
        </dashboard-button>
      </div>
      <div class="row mt-4">
        <dashboard-button to="/calendar" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_calendar_add" />
          BOOK
        </dashboard-button>
        <dashboard-button to="/schedule" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_calendar" />
          MY CAL
        </dashboard-button>
        <dashboard-button to="/admin" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_setting" />
          ADMIN
        </dashboard-button>
      </div>
      <div class="row mt-4">
        <dashboard-button to="/account" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_user" />
          ACCOUNT
        </dashboard-button>
        <dashboard-button to="/info" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_info" />
          INFO
        </dashboard-button>
        <dashboard-button to="/logout" colWidth="4">
          <div class="bc_icon bc_icon_dashboard bc_icon_logout" />
          LOGOUT
        </dashboard-button>
      </div>
    </div>
    <div class="col-3" />
  </div>
</template>

<script>
import BooksysPie from "./Pie.vue";
import DashboardButton from "./bricks/DashboardButton.vue";

export default {
  name: "DashboardAdmin",
  components: {
    BooksysPie,
    DashboardButton,
  },
  props: ["sessionData"],
  data() {
    return {
      properties: {
        containerWidth: 80,
        containerHeight: 85,
        animate: false,
      },
    };
  },
};
</script>
