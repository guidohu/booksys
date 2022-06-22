<template>
  <subpage-container title="Administrator Menu">
    <div class="row mt-5 d-none d-lg-flex text-center">
      <div class="col-1 col-lg-3"></div>
      <div class="col-10 col-lg-6 text-center">
        <div class="row">
          <dashboard-button to="/users" colWidth="4" offsetWidth="2">
            <div class="bc_icon bc_icon_dashboard bc_icon_users" />
            USERS
          </dashboard-button>
          <dashboard-button to="/accounting" colWidth="4">
            <div class="bc_icon bc_icon_dashboard bc_icon_money" />
            FINANCE
          </dashboard-button>
        </div>
        <div class="row mt-4">
          <dashboard-button to="/settings" colWidth="4" offsetWidth="2">
            <div class="bc_icon bc_icon_dashboard bc_icon_wrench" />
            SETTINGS
          </dashboard-button>
          <dashboard-button to="/logs" colWidth="4">
            <div class="bc_icon bc_icon_dashboard bc_icon_statistics" />
            LOGS
          </dashboard-button>
        </div>
      </div>
      <div class="col-1 col-lg-3" />
    </div>
    <div class="row mt-3 d-flex d-lg-none text-center">
      <div class="col-1" />
      <div class="col-5 text-center">
        <div class="row">
          <dashboard-button to="/users" colWidth="12" offsetWidth="0">
            <div class="bc_icon bc_icon_dashboard bc_icon_users" />
            USERS
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/settings" colWidth="12" offsetWidth="0">
            <div class="bc_icon bc_icon_dashboard bc_icon_wrench" />
            SETTINGS
          </dashboard-button>
        </div>
      </div>
      <div class="col-5 text-center">
        <div class="row">
          <dashboard-button to="/accounting" colWidth="12">
            <div class="bc_icon bc_icon_dashboard bc_icon_money" />
            FINANCE
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/logs" colWidth="12">
            <div class="bc_icon bc_icon_dashboard bc_icon_statistics" />
            LOGS
          </dashboard-button>
        </div>
      </div>
    </div>
  </subpage-container>
</template>

<script>
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";
import DashboardButton from "@/components/bricks/DashboardButton.vue";

export default {
  name: "WSAdmin",
  components: {
    SubpageContainer,
    DashboardButton,
  },
};
</script>
