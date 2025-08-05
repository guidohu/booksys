<template>
  <div>
    <navbar-mobile title="Dashboard" role="admin-home" />
    <div class="row navbar-padding">
      <div class="col-1" />
      <div class="col-5">
        <div class="row mt-3">
          <dashboard-button to="/boat">
            <div id="boat" class="bc_icon bc_icon_dashboard bc_icon_boat" />
            BOAT
          </dashboard-button>
        </div>
        <div class="row mt-3">
          <dashboard-button to="/today">
            <BooksysPie :session-data="sessionData" :properties="properties" />
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

<script setup>
import BooksysPie from "./BooksysPie.vue";
import NavbarMobile from "booksys/components/NavbarMobile.vue";
import DashboardButton from "booksys/components/bricks/DashboardButton.vue";
import { ref } from "vue";

const props = defineProps({
  sessionData: {
    type: Object,
    required: true,
  },
});

const properties = ref({
  containerWidth: 80,
  containerHeight: 85,
  animate: false,
});
</script>

<style scoped>
.navbar-padding {
  padding-top: 60px;
}
</style>
