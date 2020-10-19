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

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/login'
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
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
