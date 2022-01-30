<template>
  <div>
    <navbar-mobile :title="Home" role="admin-home"/>

    <div class="row">
      <div class="col-1"/>
      <div class="col-5">
        <div class="row mt-3">
          <dashboard-button to="/boat">
            <div id="boat" class="bc_icon bc_icon_dashboard bc_icon_boat" />
            BOAT
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/today">
            <BooksysPie
              :session-data="sessionData"
              :properties="properties"
            />
            TODAY
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/admin">
            <div class="bc_icon bc_icon_dashboard bc_icon_setting" />
            ADMIN
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/account">
            <div class="bc_icon bc_icon_dashboard bc_icon_user" />
            ACCOUNT
          </dashboard-button>
        </div>
      </div>
      <div class="col-5">
        <div class="row mt-3">
          <dashboard-button to="/ride">
            <div class="bc_icon bc_icon_dashboard bc_icon_wakeboard" />
            RIDE
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/calendar">
            <div class="bc_icon bc_icon_dashboard bc_icon_calendar_add" />
            BOOK
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/schedule">
            <div class="bc_icon bc_icon_dashboard bc_icon_calendar" />
            MY CAL
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/info">
            <div class="bc_icon bc_icon_dashboard bc_icon_info" />
            INFO
          </dashboard-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import BooksysPie from "./Pie.vue";
import NavbarMobile from "@/components/NavbarMobile.vue";
import DashboardButton from '@/components/bricks/DashboardButton.vue';

export default {
  name: "DashboardAdminMobile",
  components: {
    BooksysPie,
    NavbarMobile,
    DashboardButton
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
