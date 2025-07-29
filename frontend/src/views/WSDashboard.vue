<template>
  <div>
    <show-for-desktop>
      <div class="display">
        <header v-if="userInfo != null" class="welcome">
          Welcome,
          <router-link to="/account" class="header-desktop">
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
    </show-for-desktop>
    <show-for-mobile>
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
    </show-for-mobile>
  </div>
</template>

<script setup>
import { defineAsyncComponent, computed, onMounted } from "vue";
import { useStore, mapGetters } from "vuex";
import { useRouter } from "vue-router";
import ShowForMobile from "booksys/components/bricks/ShowForMobile.vue";
import ShowForDesktop from "booksys/components/bricks/ShowForDesktop.vue";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);

const DashboardAdmin = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "dashboard-admin" */ "booksys/components/DashboardAdmin.vue"
  )
);
const DashboardMember = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "dashboard-member" */ "booksys/components/DashboardMember.vue"
  )
);
const DashboardGuest = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "dashboard-guest" */ "booksys/components/DashboardGuest.vue"
  )
);
const DashboardAdminMobile = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "dashboard-admin-mobile" */ "booksys/components/DashboardAdminMobile.vue"
  )
);
const DashboardMemberMobile = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "dashboard-member-mobile" */ "booksys/components/DashboardMemberMobile.vue"
  )
);
const DashboardGuestMobile = defineAsyncComponent(() =>
  import(
    /* webpackChunkName: "dashboard-guest-mobile" */ "booksys/components/DashboardGuestMobile.vue"
  )
);

const store = useStore();
const router = useRouter();

const userInfo = computed(mapGetters("login", ["userInfo"]).userInfo.bind({ $store: store }));
const role = computed(mapGetters("login", ["role"]).role.bind({ $store: store }));
const isLoggedIn = computed(
  mapGetters("loginStatus", ["isLoggedIn"]).isLoggedIn.bind({ $store: store })
);
const getSessions = computed(
  mapGetters("sessions", ["getSessions"]).getSessions.bind({ $store: store })
);

const getTimeZone = () => {
  return "Europe/Zurich";
};

const getSessionInfo = () => {
  const dateStart = dayjs().tz(getTimeZone()).startOf("day");
  const dateEnd = dayjs().tz(getTimeZone()).endOf("day");
  console.log("Query sessions from", dateStart, "to", dateEnd);
  store.dispatch("sessions/querySessions", { start: dateStart, end: dateEnd });
};

onMounted(() => {
  // load user info into store
  store
    .dispatch("login/getUserInfo")
    .then(() => {
      console.log("Get session info.");
      getSessionInfo();
    })
    .catch((errors) => {
      console.log("Cannot get user info (probably not logged in)");
      if (errors[0] == "login required") {
        router.push("/login");
      }
    });
});
</script>