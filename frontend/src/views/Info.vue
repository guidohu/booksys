<template>

<div v-if="isDesktop" class="display">
  <main-title title-name="Location"/>
  <b-row>
    <b-col cols="12" class="mt-5 text-center">
      <div cols="12" v-if="this.getLocationAddress != null" v-html="this.getLocationAddress" class="main-color"/>
      <div v-else>no address set by the site-owner</div>
    </b-col>
  </b-row>
  <b-row>
    <b-col cols="12" class="text-center mt-3">
      <iframe v-if="this.getLocationMap != null" :src="this.getLocationMap" frameborder="0" style="{border:0;}" :width="mapWidth" :height="this.mapHeight"/>
      <div v-else>no map to display for this adress</div>
    </b-col>
  </b-row>
  <div class="bottom">
    <b-button variant="outline-light" to="/dashboard">
      <b-icon-house/>
      HOME
    </b-button>
  </div>
</div>
<div v-else>
  <NavbarMobile title="Location"/>
  <b-card no-body>
    <b-card-body>
      <b-card-text>
        <b-row class="text-center">
          <b-col cols="12" class="text-center">
            <div cols="12" v-if="this.getLocationAddress != null" v-html="this.getLocationAddress"/>
            <div v-else>no address is configured</div>
          </b-col>
        </b-row>
        <b-row>
          <b-col cols="12" class="text-center mt-3">
            <iframe v-if="this.getLocationMap != null" :src="this.getLocationMap" frameborder="0" style="{border:0;}" :width="mapWidth" :height="this.mapHeight"/>
            <div v-else>no map is configured</div>
          </b-col>
        </b-row>
      </b-card-text>
    </b-card-body>
  </b-card>
</div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { BooksysBrowser } from '@/libs/browser';
import NavbarMobile from '@/components/NavbarMobile';
import MainTitle from '@/components/MainTitle';
import {
  BRow,
  BCol,
  BButton,
  BIconHouse,
  BCard,
  BCardBody,
  BCardText
} from 'bootstrap-vue';

export default {
  name: 'Info',
  components: {
    MainTitle,
    NavbarMobile,
    BRow,
    BCol,
    BButton,
    BIconHouse,
    BCard,
    BCardBody,
    BCardText
  },
  computed: {
    ...mapGetters('configuration', [
      'getLocationAddress',
      'getLocationMap',
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
    ])
  },
  created () {
    console.log("Info.vue: Try to get all information required")
    this.queryConfiguration()
  }
}
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