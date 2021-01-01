<template>

<div v-if="isDesktop" class="display">
  <main-title title-name="Account"/>
  <b-row class="ml-1 mr-1">
    <b-col cols="12">
      <UserProfileCard/>
    </b-col>
  </b-row>
  <b-row class="ml-1 mr-1 mt-1">
    <b-col cols="12">
      <UserStatisticsCard/>
    </b-col>
  </b-row>
  <b-row class="ml-1 mr-1 mt-1">
    <b-col cols="12">
      <UserBalanceCard/>
    </b-col>
  </b-row>
  <div class="bottom mr-2">
    <b-button variant="outline-light" to="/dashboard">
      <b-icon-house/>
      HOME
    </b-button>
  </div>
</div>
<div v-else>
  <NavbarMobile title="Account"/>
  <UserProfileCard/>
  <UserStatisticsCard class="mt-2"/>
  <UserBalanceCard class="mt-2"/>
</div>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import { BooksysBrowser } from '@/libs/browser';
import NavbarMobile from '@/components/NavbarMobile';
import UserProfileCard from '@/components/UserProfileCard';
import UserStatisticsCard from '@/components/UserStatisticsCard';
import UserBalanceCard from '@/components/UserBalanceCard';
import MainTitle from '@/components/MainTitle';
import {
  BRow,
  BCol,
  BButton,
  BIconHouse
} from 'bootstrap-vue';

export default Vue.extend({
  name: 'Account',
  components: {
    NavbarMobile,
    MainTitle,
    UserProfileCard,
    UserStatisticsCard,
    UserBalanceCard,
    BRow,
    BCol,
    BButton,
    BIconHouse
  },
  computed: {
    ...mapGetters('configuration', [
      'getLocationAddress',
      'getLocationMap'
    ]),
    ...mapGetters('login', [
      'userInfo'
    ]),
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    },
    mapHeight: function() {
      if(this.isMobile){
        return 350
      }else{
        return 300
      }
    },
    mapWidth: function() {
      if(this.isMobile){
        return 340
      }else{
        return 600
      }
    }
  },
  methods: {
    ...mapActions('configuration', [
      'queryConfiguration'
    ]),
    ...mapActions('user', [
      'queryHeatHistory'
    ])
  },
  created () {
    console.log("Info.vue: Try to get all information required")
    this.queryConfiguration();
    this.queryHeatHistory();
  }
})
</script>

<style>
  div.main-title {
    font-family: Arial, Helvetica, sans-serif;
    color: #a1a1a1;
    border-bottom: 1px solid #eee;
    margin: 0px 0px 10px;
    padding-top: 10px;
    padding-bottom: 0px;
    padding-left: 15px;
  }

  div.main-color {
    color: #a1a1a1;
  }
</style>