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
    <!-- Mobile: the same bento grid the dashboard uses, so /admin does not
         drop back to the old square buttons halfway through a journey. -->
    <div class="bk-bento d-lg-none">
      <dashboard-tile to="/users" label="Users" icon="bi-people" />
      <dashboard-tile to="/accounting" label="Finance" icon="bi-cash-stack" />
      <dashboard-tile to="/settings" label="Settings" icon="bi-gear" />
      <dashboard-tile to="/logs" label="Logs" icon="bi-journal-text" />
    </div>
  </subpage-container>
</template>

<script setup>
import SubpageContainer from "booksys/components/bricks/SubpageContainer.vue";
import DashboardButton from "booksys/components/bricks/DashboardButton.vue";
import DashboardTile from "booksys/components/bricks/DashboardTile.vue";
</script>
