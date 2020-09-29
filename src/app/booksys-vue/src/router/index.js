import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../views/Login'
import Logout from '../views/Logout'
import SignUp from '../views/SignUp'
import Dashboard from '../views/Dashboard'
import Account from '../views/Account'
import Info from '../views/Info'
import Schedule from '../views/Schedule'

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
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
