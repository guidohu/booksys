<template>
  <div>
    <navbar-mobile title="Dashboard" role="guest-home"/>

    <div class="row">
      <div class="col-3"/>
      <div class="col-6 text-center">
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
      <div class="col-3"/>
    </div>
  </div>
</template>

<script>
import NavbarMobile from "@/components/NavbarMobile.vue";
import DashboardButton from '@/components/bricks/DashboardButton.vue';

export default {
  name: "DashboardGuestMobile",
  components: {
    NavbarMobile,
    DashboardButton
  },
};
</script>
