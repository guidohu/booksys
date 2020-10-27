<template>
  <div v-if="isDesktop" class="display">
    <b-row>
      <b-col>
        <div class="main-title">
          <h3>
            Stop Watch
          </h3>
        </div>
      </b-col>
    </b-row>
    <b-row class="ml-1 mr-1">
      <b-col cols="8">
        <StopWatchCard/>
      </b-col>
      <b-col cols="4">
      </b-col>
    </b-row>
  </div>
  <div v-else>
    <NavbarMobile title="Stop Watch"/>
    <StopWatchCard/>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarMobile from '@/components/NavbarMobile';
import StopWatchCard from '@/components/StopWatchCard';


export default Vue.extend({
  name: "Watch",
  components: {
    NavbarMobile,
    StopWatchCard
  },
  computed: {
    isDesktop: function() {
      return !BooksysBrowser.isMobile();
    }
  }
})
</script>