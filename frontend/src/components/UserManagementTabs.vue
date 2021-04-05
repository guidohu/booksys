<template>
  <b-tabs content-class="mt-3">
    <b-tab active>
      <div slot="title">
        <p class="h6 mt-2 mb-2">
          <b-icon-person />
          Users
        </p>
      </div>
      <div :class="tabClass">
        <UserTable />
      </div>
    </b-tab>
    <b-tab>
      <div slot="title">
        <p class="h6 mt-2 mb-2">
          <b-icon-people />
          User Groups
        </p>
      </div>
      <div :class="tabClass">
        <UserGroupTable />
      </div>
    </b-tab>
  </b-tabs>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import UserTable from "@/components/UserTable";
import UserGroupTable from "@/components/UserGroupTable";
import { BTabs, BTab, BIconPerson, BIconPeople } from "bootstrap-vue";

export default {
  name: "UserManagementTabs",
  components: {
    UserTable,
    UserGroupTable,
    BTabs,
    BTab,
    BIconPerson,
    BIconPeople,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
    tabClass: function () {
      if (!BooksysBrowser.isMobile()) {
        return "user-management-tab-limited-height";
      } else {
        return "";
      }
    },
  },
};
</script>

<style scoped>
.user-management-tab-limited-height {
  max-height: 400px;
  overflow: scroll;
}
</style>
