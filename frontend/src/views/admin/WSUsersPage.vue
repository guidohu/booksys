<template>
  <div>
    <div v-if="isDesktop" class="display">
      <main-title title-name="User Management" />
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <b-card no-body>
            <UserManagementTabs class="content-max-height" />
          </b-card>
        </b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button variant="outline-light" to="/admin">
          <b-icon-gear />
          ADMIN
        </b-button>
      </div>
    </div>
    <div v-else>
      <NavbarAdminMobile title="User Management" />
      <b-card no-body>
        <UserManagementTabs />
      </b-card>
    </div>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import NavbarAdminMobile from "@/components/NavbarAdminMobile";
import UserManagementTabs from "@/components/UserManagementTabs";
import MainTitle from "@/components/MainTitle";
import { BRow, BCol, BCard, BButton, BIconGear } from "bootstrap-vue";

export default {
  name: "WSUsersPage",
  components: {
    NavbarAdminMobile,
    MainTitle,
    UserManagementTabs,
    BRow,
    BCol,
    BCard,
    BButton,
    BIconGear,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
  },
};
</script>

<style>
.content-max-height {
  min-height: 480px;
  max-height: 480px;
  overflow-y: scroll;
}
</style>
