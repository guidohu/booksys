<template>
  <div>
    <ul class="nav nav-tabs" id="user-management-tab" role="tablist">
      <li class="nav-item" role="presentation">
        <a
          class="nav-link active"
          data-bs-toggle="tab"
          href="#users"
          id="user-tab"
          data-bs-target="#users"
          role="tab"
          aria-controls="users"
          aria-selected="true"
        >
          <i class="bi bi-person"></i>
          Users
        </a>
      </li>
      <li class="nav-item" role="presentation">
        <a
          class="nav-link"
          data-bs-toggle="tab"
          href="#user-groups"
          id="user-groups-tab"
          data-bs-target="#user-groups"
          role="tab"
          aria-controls="user-groups"
          aria-selected="false"
        >
          <i class="bi bi-people"></i>
          User Groups
        </a>
      </li>
    </ul>

    <!-- Responsive tab content to handle different
         height requirements -->
    <div class="tab-content">
      <div
        class="tab-pane active"
        id="users"
        role="tabpanel"
        aria-labelledby="users-tab"
      >
        <show-for-desktop>
          <div class="tab-limited-height">
            <user-table />
          </div>
        </show-for-desktop>
        <show-for-mobile>
          <user-table />
        </show-for-mobile>
      </div>
      <div
        class="tab-pane"
        id="user-groups"
        role="tabpanel"
        aria-labelledby="user-groups-tab"
      >
        <show-for-desktop>
          <div class="tab-limited-height">
            <user-group-table />
          </div>
        </show-for-desktop>
        <show-for-mobile>
          <user-group-table />
        </show-for-mobile>
      </div>
    </div>
  </div>
</template>

<script>
import UserTable from "@/components/UserTable";
import UserGroupTable from "@/components/UserGroupTable";
import ShowForMobile from "@/components/bricks/ShowForMobile.vue";
import ShowForDesktop from "@/components/bricks/ShowForDesktop.vue";

export default {
  name: "UserManagementTabs",
  components: {
    ShowForMobile,
    ShowForDesktop,
    UserTable,
    UserGroupTable,
  },
};
</script>

<style scoped>
.tab-limited-height {
  max-height: 420px;
  height: 420px;
  overflow-y: scroll;
}

.nav-link {
  color: #bdbdbd;
}

a {
  color: #bdbdbd;
}
</style>
