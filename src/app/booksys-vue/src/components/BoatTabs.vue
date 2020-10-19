<template>
  <b-tabs content-class="mt-3" card>
    <b-tab active>
      <div slot="title">
        <div :class="titleClassEngineHours"></div>
        Engine Hours
      </div>
      <p>I'm the first tab</p>
    </b-tab>
    <b-tab>
      <div slot="title">
        <div :class="titleClassFuel"></div>
        Fuel
      </div>
    </b-tab>
    <b-tab>
      <div slot="title">
        <div :class="titleClassMaintenance"></div>
        Maintenance
      </div>
    </b-tab>            
  </b-tabs>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';

export default Vue.extend({
  name: 'BoatTabs',
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    },
    titleClassEngineHours: function (){
      if(BooksysBrowser.isMobile()){
        return "bc_icon bc_icon_tab bc_icon_tacho align-middle";
      }
      return "bc_icon bc_icon_tab bc_icon_tacho inline align-middle";
    },
    titleClassFuel: function (){
      if(BooksysBrowser.isMobile()){
        return "bc_icon bc_icon_tab bc_icon_fuel align-middle";
      }
      return "bc_icon bc_icon_tab bc_icon_fuel inline align-middle";
    },
    titleClassMaintenance: function (){
      if(BooksysBrowser.isMobile()){
        return "bc_icon bc_icon_tab bc_icon_wrench align-middle";
      }
      return "bc_icon bc_icon_tab bc_icon_wrench inline align-middle";
    }
  }
})
</script>