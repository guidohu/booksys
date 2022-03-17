<template>
  <div>
    <navbar-mobile title="Dashboard" role="member-home" />

    <div class="row">
      <div class="col-3" />
      <div class="col-6 text-center">
        <div class="row mt-3">
          <dashboard-button to="/schedule">
            <div class="bc_icon bc_icon_dashboard bc_icon_calendar" />
            MY CAL
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/account">
            <div class="bc_icon bc_icon_dashboard bc_icon_user" />
            ACCOUNT
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/info">
            <div class="bc_icon bc_icon_dashboard bc_icon_info" />
            INFO
          </dashboard-button>
        </div>
      </div>
      <div class="col-3" />
    </div>
  </div>
</template>

<script>
import NavbarMobile from "@/components/NavbarMobile.vue";
import DashboardButton from "@/components/bricks/DashboardButton.vue";

export default {
  name: "DashboardMemberMobile",
  components: {
    NavbarMobile,
    DashboardButton,
  },
};
</script>
