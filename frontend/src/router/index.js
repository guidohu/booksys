import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '@/store';

// Lazy import of all the Views used by the router
const Login     = () => import('@/views/Login');
const Logout    = () => import('@/views/Logout');
const SignUp    = () => import('@/views/SignUp');
const Dashboard = () => import('@/views/Dashboard');
const Account   = () => import('@/views/Account');
const Info      = () => import('@/views/Info');
const Schedule  = () => import('@/views/Schedule');
const Today     = () => import('@/views/Today');
const Calendar  = () => import('@/views/Calendar');
const Boat      = () => import('@/views/Boat');
const Ride      = () => import('@/views/Ride');
const Watch     = () => import('@/views/Watch');
const Admin     = () => import('@/views/Admin');
const Users     = () => import('@/views/admin/Users');
const Accounting = () => import('@/views/admin/Accounting');
const Settings  = () => import('@/views/admin/Settings');
const Logs      = () => import('@/views/admin/Logs');
const PasswordReset = () => import('@/views/PasswordReset');
const Setup     = () => import('@/views/Setup');

Vue.use(VueRouter)

const loadStoreModule = (moduleName, next) => {
  if(! store.hasModule(moduleName)){
    import('@/store/modules/' + moduleName).then(module => {
      store.registerModule(moduleName, module.default);
      console.log("Store module registered: " + moduleName);
      next();
    })
  }else{
    next();
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
      loadStoreModule('login', next);
    },
    component: Login
  },
  {
    path: '/logout',
    name: 'Logout',
    component: Logout
  },
  {
    path: '/signup',
    name: 'SignUp',
    beforeEnter: (to, from, next) => {
      loadStoreModule('configuration', next);
    },
    component: SignUp
  },
  {
    path: '/password/reset',
    name: 'PasswordReset',
    beforeEnter: (to, from, next) => {
      loadStoreModule('configuration', next);
    },
    component: PasswordReset
  },
  {
    path: '/today',
    name: 'Today',
    component: Today
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/account',
    name: 'Account',
    component: Account
  },
  {
    path: '/info',
    name: 'Info',
    component: Info
  },
  {
    path: '/schedule',
    name: 'Schedule',
    component: Schedule
  },
  {
    path: '/calendar',
    name: 'Calendar',
    component: Calendar
  },
  {
    path: '/boat',
    name: 'Boat',
    component: Boat
  },
  {
    path: '/ride',
    name: 'Ride',
    component: Ride
  },
  {
    path: '/watch',
    name: 'Watch',
    component: Watch
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin
  },
  {
    path: '/users',
    name: 'Users',
    component: Users
  },
  {
    path: '/accounting',
    name: 'Accounting',
    component: Accounting
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings
  },
  {
    path: '/logs',
    name: 'Logs',
    component: Logs
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
