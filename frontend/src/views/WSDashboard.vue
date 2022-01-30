<template>
  <div>
    <div
      v-if="
        userInfo != null &&
          userInfo.user_role_name == 'admin' &&
          getDbUpdateStatus != null &&
          getDbUpdateStatus.updateAvailable == true
      "
    >
      <database-update-modal />
    </div>
    <div
      v-if="isDesktop"
      class="display"
    >
      <header
        v-if="userInfo != null"
        class="welcome"
      >
        Welcome,
        <router-link
          to="/account"
          class="header-desktop"
        >
          {{ userInfo.first_name }} {{ userInfo.last_name }}
        </router-link>
      </header>
      <dashboard-admin
        v-if="role && role == 'admin' && getSessions != null"
        :session-data="getSessions"
      />
      <dashboard-member
        v-if="role && role == 'member' && getSessions != null"
        :session-data="getSessions"
      />
      <dashboard-guest
        v-if="role && role == 'guest' && getSessions != null"
        :session-data="getSessions"
      />
    </div>
    <div v-else>
      <dashboard-admin-mobile
        v-if="role && role == 'admin' && getSessions != null"
        :session-data="getSessions"
      />
      <dashboard-member-mobile
        v-if="role && role == 'member' && getSessions != null"
        :session-data="getSessions"
      />
      <dashboard-guest-mobile
        v-if="role && role == 'guest' && getSessions != null"
        :session-data="getSessions"
      />
    </div>
  </div>
</template>

<script>
import { defineAsyncComponent } from 'vue';
import { mapGetters, mapActions } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

const DashboardAdmin = defineAsyncComponent(() => 
  import(/* webpackChunkName: "dashboard-admin" */ "@/components/DashboardAdmin.vue"));
const DashboardMember = defineAsyncComponent(() => 
  import(/* webpackChunkName: "dashboard-member" */ "@/components/DashboardMember.vue"));
const DashboardGuest = defineAsyncComponent(() => 
  import(/* webpackChunkName: "dashboard-guest" */ "@/components/DashboardGuest.vue"));
const DashboardAdminMobile = defineAsyncComponent(() => 
  import(/* webpackChunkName: "dashboard-admin-mobile" */ "@/components/DashboardAdminMobile.vue"));
const DashboardMemberMobile = defineAsyncComponent(() => 
  import(/* webpackChunkName: "dashboard-member-mobile" */ "@/components/DashboardMemberMobile.vue"));
const DashboardGuestMobile = defineAsyncComponent(() => 
  import(/* webpackChunkName: "dashboard-guest-mobile" */ "@/components/DashboardGuestMobile.vue"));
const DatabaseUpdateModal = defineAsyncComponent(() => 
  import(/* webpackChunkName: "database-update" */ "@/components/DatabaseUpdate.vue"));

export default {
  name: "WSDashboard",
  components: {
    DashboardAdmin,
    DashboardMember,
    DashboardGuest,
    DashboardAdminMobile,
    DashboardMemberMobile,
    DashboardGuestMobile,
    DatabaseUpdateModal,
  },
  computed: {
    ...mapGetters("login", ["userInfo", "role"]),
    ...mapGetters("loginStatus", ["isLoggedIn"]),
    ...mapGetters("sessions", ["getSessions"]),
    ...mapGetters("configuration", ["getDbUpdateStatus"]),
    isMobile: function () {
      return BooksysBrowser.isMobile();
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
  },
  watch: {
    role: function () {
      this.dbUpdateCheck();
    },
  },
  methods: {
    ...mapActions("login", ["getUserInfo"]),
    ...mapActions("sessions", ["querySessions"]),
    ...mapActions("configuration", ["queryDbUpdateStatus"]),
    getTimeZone: function () {
      return "Europe/Berlin";
    },
    dbUpdateCheck() {
      if (this.role == "admin") {
        this.queryDbUpdateStatus();
      }
    },
    getSessionInfo() {
      const dateStart = dayjs().tz(this.getTimeZone()).startOf("day");
      const dateEnd = dayjs().tz(this.getTimeZone()).endOf("day");
      console.log("Query sessions from", dateStart, "to", dateEnd);
      this.querySessions({ start: dateStart, end: dateEnd });
    },
  },
  created() {
    // load user info into store
    this.getUserInfo()
      .then(() => {
        console.log("got user info");
        this.dbUpdateCheck();
        this.getSessionInfo();
      })
      .catch((errors) => {
        console.log("Cannot get user info (probably not logged in)");
        if (errors[0] == "login required") {
          this.$router.push("/login");
        }
      });
  },
};
</script>
