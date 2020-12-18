import Vue from 'vue';
import VueRouter from 'vue-router';
import Login from '../views/Login';
import Logout from '../views/Logout';
import SignUp from '../views/SignUp';
import Dashboard from '../views/Dashboard';
import Account from '../views/Account';
import Info from '../views/Info';
import Schedule from '../views/Schedule';
import Today from '../views/Today';
import Calendar from '../views/Calendar';
import Boat from '../views/Boat';
import Ride from '../views/Ride';
import Watch from '../views/Watch';
import Admin from '../views/Admin';
import Users from '../views/admin/Users';
import Accounting from '../views/admin/Accounting';
import Settings from '../views/admin/Settings';
import Logs from '../views/admin/Logs';
import PasswordReset from '../views/PasswordReset';
import Setup from '../views/Setup';

Vue.use(VueRouter)

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
    component: SignUp
  },
  {
    path: '/password/reset',
    name: 'PasswordReset',
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
