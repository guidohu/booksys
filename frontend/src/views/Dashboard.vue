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
    <div v-if="isDesktop" class="display">
      <header v-if="userInfo != null" class="welcome">
        Welcome,
        <b-link to="/account" class="header-desktop">
          {{ userInfo.first_name }} {{ userInfo.last_name }}
        </b-link>
      </header>
      <DashboardAdmin
        v-if="role && role == 'admin' && getSessions != null"
        v-bind:sessionData="getSessions"
      />
      <DashboardMember
        v-if="role && role == 'member' && getSessions != null"
        v-bind:sessionData="getSessions"
      />
      <DashboardGuest
        v-if="role && role == 'guest' && getSessions != null"
        v-bind:sessionData="getSessions"
      />
    </div>
    <div v-else>
      <DashboardAdminMobile
        v-if="role && role == 'admin' && getSessions != null"
        v-bind:sessionData="getSessions"
      />
      <DashboardMemberMobile
        v-if="role && role == 'member' && getSessions != null"
        v-bind:sessionData="getSessions"
      />
      <DashboardGuestMobile
        v-if="role && role == 'guest' && getSessions != null"
        v-bind:sessionData="getSessions"
      />
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";
import { BLink } from "bootstrap-vue";

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

const DashboardAdmin = () =>
  import(
    /* webpackChunkName: "dashboard-admin" */ "@/components/DashboardAdmin"
  );
const DashboardMember = () =>
  import(
    /* webpackChunkName: "dashboard-member" */ "@/components/DashboardMember"
  );
const DashboardGuest = () =>
  import(
    /* webpackChunkName: "dashboard-guest" */ "@/components/DashboardGuest"
  );
const DashboardAdminMobile = () =>
  import(
    /* webpackChunkName: "dashboard-admin-mobile" */ "@/components/DashboardAdminMobile"
  );
const DashboardMemberMobile = () =>
  import(
    /* webpackChunkName: "dashboard-member-mobile" */ "@/components/DashboardMemberMobile"
  );
const DashboardGuestMobile = () =>
  import(
    /* webpackChunkName: "dashboard-guest-mobile" */ "@/components/DashboardGuestMobile"
  );
const DatabaseUpdateModal = () =>
  import(
    /* webpackChunkName: "database-update" */ "@/components/DatabaseUpdate"
  );

export default {
  name: "Dashboard",
  components: {
    DashboardAdmin,
    DashboardMember,
    DashboardGuest,
    DashboardAdminMobile,
    DashboardMemberMobile,
    DashboardGuestMobile,
    DatabaseUpdateModal,
    BLink,
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
