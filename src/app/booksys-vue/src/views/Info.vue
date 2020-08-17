<template>
<div v-if="isDesktop" class="display">
  <b-card no-body>
    <b-card-header>
      <h3>Location</h3>
    </b-card-header>
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
    <b-card-footer class="text-right">
      <!-- <div class="bottom"> -->
        <b-button variant="outline-dark" to="/dashboard">
          <b-icon-house></b-icon-house>
          HOME
        </b-button>
      <!-- </div> -->
    </b-card-footer>
  </b-card>
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
import Vue from 'vue'
import { mapActions, mapGetters } from 'vuex'
import { BooksysBrowser } from '@/libs/browser'
import NavbarMobile from '@/components/NavbarMobile.vue'

export default Vue.extend({
  name: 'Info',
  components: {
    NavbarMobile
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
})
</script>