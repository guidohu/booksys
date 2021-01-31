<template>
  <div id="app">
    <b-container v-if="isDesktop">
      <b-row class="display noborder nobackground">
        <div class="vcenter_placeholder">
          <div class="vcenter_content">
            <alert-message 
              v-if="backendReachable == false" 
              :alertMessage="backendNotReachableAlertMsg"
            />
            <router-view v-else/>
          </div>
        </div>
      </b-row>
      <footer class="legal-footer">
        Copyright 2013-2021 by Guido Hungerbuehler <a href="https://github.com/guidohu/booksys">Find me on Github</a>
      </footer>
    </b-container>
    <b-container v-if="isMobile" fluid="true">
      <alert-message 
        v-if="backendReachable == false" 
        :alertMessage="backendNotReachableAlertMsg"
      />
      <router-view v-else/>
    </b-container>
  </div>
</template>

<script>
import { BooksysBrowser } from '@/libs/browser';
import { getBackendStatus }from '@/api/backend';
import { mapGetters } from 'vuex';
import {
  BContainer,
  BRow
} from 'bootstrap-vue';

// Lazy imports
const AlertMessage = () => import('@/components/AlertMessage');

export default{
  name: 'App',
  data: function() {
    return {
      backendReachable: true,
      backendNotReachableAlertMsg: "The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible. You might try to simply refresh the page if you feel lucky.",
    }
  },
  components: {
    AlertMessage,
    BContainer,
    BRow
  },
  computed: {
    isMobile: function() {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function() {
      return !BooksysBrowser.isMobile()
    },
    ...mapGetters('loginStatus', [
      'isLoggedIn'
    ])
  },
  watch: {
    isLoggedIn: function(newStatus, oldStatus){
      if(newStatus == false && oldStatus == true){
        console.log("Logout detected by App.");
        if(this.$route.path != '/logout'
          && this.$route.path != '/setup'
          && this.$route.path != '/signup'
        ){
          console.log("Redirect to login page.");
          this.$router.push("/login");
        }
      }
    }
  },
  beforeCreate() {
    // set mobile view and style
    if (BooksysBrowser.isMobile()) {
      BooksysBrowser.setViewportMobile()
      BooksysBrowser.setManifest()
      BooksysBrowser.setMetaMobile()
      BooksysBrowser.addMobileCSS()
    }
  },
  mounted() {
    // short pulse check on the
    // backend to verify whether
    // it is up and configured
    getBackendStatus()
    .then((status) => {
      // if we do not have a configFile for the app -> go to setup
      if(!status.configFile || !status.configDb){
        console.log("App Setup Done: no");
        if(this.$route.path !== '/setup'){
          this.$router.push("/setup");
        }
      }else if(!status.dbReachable){
        this.backendReachable = false;
      }else{
        console.log("App Setup Done: yes (no automated forward to setup)");

        // TODO check if user is logged in, if not -> redirect to login page

      }
    })
    .catch((error) => {
      this.errors = [error];
    })
  }
}
</script>
