<template>
  <div>
    <div v-if="isDesktop">
      <subpage-desktop title="Administrator Menu">
        <template v-slot:content>
          <div class="row mt-5">
            <div class="col-3" />
            <div class="col-6 text-center">
              <div class="row mt-4">
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
            <div cols="3" />
          </div>          
        </template>
        <template v-slot:bottom>
          <router-link tag="button" class="btn btn-outline-light" to="/dashboard">
            <i class="bi bi-house"></i>
            HOME
          </router-link>
        </template>
      </subpage-desktop>
    </div>
    <div v-else>
      <navbar-mobile title="Administrator Menu" role="admin-home"/>
      <div class="row">
        <div class="col-1"/>
        <div class="col-5">
          <div class="row mt-3">
            <dashboard-button to="/users" colWidth="4" offsetWidth="2">
              <div class="bc_icon bc_icon_dashboard bc_icon_users" />
              USERS
            </dashboard-button>
          </div>
          <div class="row mt-3">
            <dashboard-button to="/settings" colWidth="4" offsetWidth="2">
              <div class="bc_icon bc_icon_dashboard bc_icon_wrench" />
              SETTINGS
            </dashboard-button>
          </div>
        </div>
        <div class="col-5">
          <div class="row mt-3">
            <dashboard-button to="/accounting" colWidth="4">
              <div class="bc_icon bc_icon_dashboard bc_icon_money" />
              FINANCE
            </dashboard-button>
          </div>
          <div class="row mt-3">
            <dashboard-button to="/logs" colWidth="4">
              <div class="bc_icon bc_icon_dashboard bc_icon_statistics" />
              LOGS
            </dashboard-button>
          </div>
        </div>
      </div>      
    </div>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import SubpageDesktop from "@/components/bricks/SubpageDesktop.vue";
import DashboardButton from "@/components/bricks/DashboardButton.vue";
import NavbarMobile from "@/components/NavbarMobile.vue";

export default {
  name: "WSAdmin",
  components: {
    SubpageDesktop,
    DashboardButton,
    NavbarMobile,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
  },
};
</script>
