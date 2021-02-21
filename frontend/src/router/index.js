import Vue from 'vue';
import VueRouter from 'vue-router';
import { loadStoreModules, store } from '@/store';

// Lazy import of all the Views used by the router
const Login     = () => import(/* webpackChunkName: "login" */ '@/views/Login');
const Logout    = () => import(/* webpackChunkName: "logout" */ '@/views/Logout');
const SignUp    = () => import(/* webpackChunkName: "signup" */ '@/views/SignUp');
const Dashboard = () => import(/* webpackChunkName: "dashboard" */ '@/views/Dashboard');
const Account   = () => import(/* webpackChunkName: "account" */ '@/views/Account');
const Info      = () => import(/* webpackChunkName: "info" */ '@/views/Info');
const Schedule  = () => import(/* webpackChunkName: "schedule" */ '@/views/Schedule');
const Today     = () => import(/* webpackChunkName: "today" */ '@/views/Today');
const Calendar  = () => import(/* webpackChunkName: "calendar" */ '@/views/Calendar');
const Boat      = () => import(/* webpackChunkName: "boat" */ '@/views/Boat');
const Ride      = () => import(/* webpackChunkName: "ride" */ '@/views/Ride');
const Watch     = () => import(/* webpackChunkName: "watch" */ '@/views/Watch');
const Admin     = () => import(/* webpackChunkName: "admin" */ '@/views/Admin');
const Users     = () => import(/* webpackChunkName: "users" */ '@/views/admin/Users');
const Accounting = () => import(/* webpackChunkName: "accounting" */ '@/views/admin/Accounting');
const Settings  = () => import(/* webpackChunkName: "settings" */ '@/views/admin/Settings');
const Logs      = () => import(/* webpackChunkName: "logs" */ '@/views/admin/Logs');
const PasswordReset = () => import(/* webpackChunkName: "password-reset" */ '@/views/PasswordReset');
const Setup     = () => import(/* webpackChunkName: "setup" */ '@/views/Setup');

Vue.use(VueRouter);

const loginEnforced = (to, from, next) => {
  if(!store.state.loginStatus.isLoggedIn){
    console.log("Login required to access this route");
    console.log("Navigate from:", from, "To:", to);
    next({ name: "Login", query: { target: to.fullPath }});
    return false;
  }else{
    return true;
  }
}

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/setup',
    name: 'Setup',
    component: Setup
  },
  {
    path: '/login',
    name: 'Login',
    beforeEnter: (to, from, next) => {
      loadStoreModules(['login', 'configuration'], next);
    },
    component: Login
  },
  {
    path: '/logout',
    name: 'Logout',
    beforeEnter: (to, from, next) => {
      loadStoreModules(['login'], next);
    },
    component: Logout
  },
  {
    path: '/signup',
    name: 'SignUp',
    beforeEnter: (to, from, next) => {
      loadStoreModules(['configuration'], next);
    },
    component: SignUp
  },
  {
    path: '/password/reset',
    name: 'PasswordReset',
    beforeEnter: (to, from, next) => {
      loadStoreModules(['configuration'], next);
    },
    component: PasswordReset
  },
  {
    path: '/today',
    name: 'Today',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['configuration', 'sessions', 'user'], next);
      }
    },
    component: Today
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['sessions', 'configuration', 'login'], next);
      }
    },
    component: Dashboard
  },
  {
    path: '/account',
    name: 'Account',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['user', 'configuration', 'login'], next);
      }
    },
    component: Account
  },
  {
    path: '/info',
    name: 'Info',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['configuration'], next);
      }
    },
    component: Info
  },
  {
    path: '/schedule',
    name: 'Schedule',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['configuration', 'user'], next);
      }
    },
    component: Schedule
  },
  {
    path: '/calendar',
    name: 'Calendar',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
       loadStoreModules(['sessions', 'configuration'], next);
      }
    },
    component: Calendar
  },
  {
    path: '/boat',
    name: 'Boat',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['boat', 'configuration', 'login'], next);
      }
    },
    component: Boat
  },
  {
    path: '/ride',
    name: 'Ride',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['sessions', 'configuration', 'stopwatch', 'user'], next);
      }
    },
    component: Ride
  },
  {
    path: '/watch',
    name: 'Watch',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['sessions', 'configuration', 'stopwatch', 'heats'], next);
      }
    },
    component: Watch
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        next();
      }
    }
  },
  {
    path: '/users',
    name: 'Users',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['user', 'configuration'], next);
      }
    },
    component: Users
  },
  {
    path: '/accounting',
    name: 'Accounting',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['configuration', 'accounting', 'user' ], next);
      }
    },
    component: Accounting
  },
  {
    path: '/settings',
    name: 'Settings',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['configuration'], next);
      }
    },
    component: Settings
  },
  {
    path: '/logs',
    name: 'Logs',
    beforeEnter: (to, from, next) => {
      if(loginEnforced(to, from, next)){
        loadStoreModules(['log'], next);
      }
    },
    component: Logs
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
