<template>
  <div>
    <div v-if="isDesktop" class="display">
      <b-row>
        <b-col>
          <div class="main-title">
            <h3>
              Boat
            </h3>
          </div>
        </b-col>
      </b-row>
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <b-card no-body>
            <BoatTabs/>
          </b-card>
        </b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button variant="outline-light" to="/dashboard">
          <b-icon-house></b-icon-house>
          HOME
        </b-button>
      </div>
    </div>
    <div v-else>
      <NavbarMobile title="Boat"/>
      <b-card no-body>
        <BoatTabs/>  
      </b-card> 
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarMobile from '@/components/NavbarMobile';
import BoatTabs from '@/components/BoatTabs';

export default Vue.extend({
  name: 'Boat',
  components: {
    NavbarMobile,
    BoatTabs
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  }
})
</script>