import { createRouter, createWebHistory } from "vue-router";
import { loadStoreModules, store } from "store";

// Lazy import of all the Views used by the router
const Login = () => import("booksys/views/WSLogin.vue");
const Logout = () =>
  import("booksys/views/WSLogout.vue");
const SignUp = () =>
  import("booksys/views/WSSignUp.vue");
const Dashboard = () =>
  import("booksys/views/WSDashboard.vue");
const Account = () =>
  import("booksys/views/WSAccount.vue");
const Info = () => import("booksys/views/WSInfo.vue");
const Schedule = () =>
  import("booksys/views/WSSchedule.vue");
const Today = () => import("booksys/views/WSToday.vue");
const Calendar = () =>
  import("booksys/views/WSCalendar.vue");
const Boat = () => import("booksys/views/WSBoat.vue");
const Ride = () => import("booksys/views/WSRide.vue");
const Watch = () => import("booksys/views/WSWatch.vue");
const Admin = () => import("booksys/views/WSAdmin.vue");
const Users = () =>
  import("booksys/views/admin/WSUsers.vue");
const Accounting = () =>
  import("booksys/views/admin/WSAccounting.vue");
const Settings = () =>
  import("booksys/views/admin/WSSettings.vue");
const Logs = () =>
  import( "booksys/views/admin/WSLogs.vue");
const PasswordReset = () =>
  import( "booksys/views/WSPasswordReset.vue");
const Setup = () =>
  import( "booksys/views/WSSetupPage.vue");

const loginEnforced = (to, from, next) => {
  if (!store.state.loginStatus.isLoggedIn) {
    console.log("Login required to access this route");
    console.log("Navigate from:", from, "To:", to);
    next({ name: "Login", query: { target: to.fullPath } });
    return false;
  } else {
    console.log("Login required to access:", to.fullPath, "User is logged in.");
    return true;
  }
};

const routes = [
  {
    path: "/",
    redirect: "/login",
  },
  {
    path: "/setup",
    name: "Setup",
    component: Setup,
  },
  {
    path: "/setup/database",
    name: "Setup Database",
    component: Setup,
  },
  {
    path: "/login",
    name: "Login",
    beforeEnter: (to, from, next) => {
      loadStoreModules(["login", "configuration"], next);
    },
    component: Login,
  },
  {
    path: "/logout",
    name: "Logout",
    beforeEnter: (to, from, next) => {
      loadStoreModules(["login"], next);
    },
    component: Logout,
  },
  {
    path: "/signup",
    name: "SignUp",
    beforeEnter: (to, from, next) => {
      loadStoreModules(["configuration"], next);
    },
    component: SignUp,
  },
  {
    path: "/password/reset",
    name: "PasswordReset",
    beforeEnter: (to, from, next) => {
      loadStoreModules(["configuration"], next);
    },
    component: PasswordReset,
  },
  {
    path: "/today",
    name: "Today",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["configuration", "sessions", "user"], next);
      }
    },
    component: Today,
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["sessions", "configuration", "login"], next);
      }
    },
    component: Dashboard,
  },
  {
    path: "/account",
    name: "Account",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["user", "configuration", "login"], next);
      }
    },
    component: Account,
  },
  {
    path: "/info",
    name: "Info",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["configuration"], next);
      }
    },
    component: Info,
  },
  {
    path: "/schedule",
    name: "Schedule",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["configuration", "user"], next);
      }
    },
    component: Schedule,
  },
  {
    path: "/calendar",
    name: "Calendar",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["sessions", "configuration"], next);
      }
    },
    component: Calendar,
  },
  {
    path: "/boat/:tab?",
    name: "Boat",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["boat", "configuration", "login"], next);
      }
    },
    component: Boat,
  },
  {
    path: "/ride",
    name: "Ride",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(
          ["sessions", "configuration", "stopwatch", "user"],
          next
        );
      }
    },
    component: Ride,
  },
  {
    path: "/watch",
    name: "Watch",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(
          ["sessions", "configuration", "stopwatch", "heats"],
          next
        );
      }
    },
    component: Watch,
  },
  {
    path: "/admin",
    name: "Admin",
    component: Admin,
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        next();
      }
    },
  },
  {
    path: "/users",
    name: "Users",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["user", "configuration"], next);
      }
    },
    component: Users,
  },
  {
    path: "/accounting",
    name: "Accounting",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["configuration", "accounting", "user", "boat"], next);
      }
    },
    component: Accounting,
  },
  {
    path: "/settings",
    name: "Settings",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["configuration"], next);
      }
    },
    component: Settings,
  },
  {
    path: "/logs",
    name: "Logs",
    beforeEnter: (to, from, next) => {
      if (loginEnforced(to, from, next)) {
        loadStoreModules(["log"], next);
      }
    },
    component: Logs,
  },
];

const router = createRouter({
  history: createWebHistory(),
  base: import.meta.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  document.title = to.name;
  next();
});

export default router;
