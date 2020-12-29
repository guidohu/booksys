<template>
  <div id="app">
    <b-container v-if="isDesktop">
      <b-row class="display noborder nobackground">
        <div class="vcenter_placeholder">
          <div class="vcenter_content">
            <AlertMessage v-if="backendStatus && backendStatus.ok == false" alertMessage="The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible."/>
            <router-view v-else/>
          </div>
        </div>
      </b-row>
      <footer class="legal-footer">
        Copyright 2013-2020 by Guido Hungerbuehler <a href="https://github.com/guidohu/booksys">Find me on Github</a>
      </footer>
    </b-container>
    <b-container v-if="isMobile" fluid="true">
      <div v-if="isMobile">
        <AlertMessage v-if="backendStatus && backendStatus.ok == false" alertMessage="The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible."/>
        <router-view v-else/>
      </div>
    </b-container>
  </div>
</template>

<script>
import { BooksysBrowser } from '@/libs/browser';
import { BooksysBackend } from '@/libs/backend';
import Backend from '@/api/backend';
import AlertMessage from './components/AlertMessage.vue';
import Vue from 'vue';
import { mapGetters, mapState, mapActions } from 'vuex';

export default Vue.extend({
  name: 'App',
  data: function() {
    return {
      backendAvailable: false, // if the backend is not available
      backendStatus: null,
      backendReachable: true,
      setupDone: false
    }
  },
  components: {
    AlertMessage
  },
  computed: {
    isMobile: function() {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function() {
      return !BooksysBrowser.isMobile()
    },
    ...mapGetters('login', [
      'userInfo'
    ]),
    ...mapState('login', [
      'isLoggedIn'
    ])
  },
  watch: {
    isLoggedIn(newValue){
      if(newValue == false && this.setupDone == true){
        console.log("App.vue: detected not logged in -> forward to /login")
        this.$router.push("/login");
      }else if(newValue == true && this.setupDone == true){
        this.$router.push("/dashboard");
      }
    },
    setupDone(newValue){
      if(newValue == false){
        this.$router.push("/setup");
      }else if(this.isLoggedIn == true){
        this.$router.push("/dashboard");
      }else{
        this.$router.push("/login");
      }
    }
  },
  methods: {
    ...mapActions('login', [
      'getIsLoggedIn'
    ])
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
  created() {
    // short pulse check on the
    // backend to verify whether
    // it is up and configured
    Backend.getStatus()
    .then((status) => {
      console.log(status);
      if(!status.configFile){
        this.setupDone = false;
        this.$router.push("/setup");
      }
    })
    .catch((error) => {
      this.errors = [error];
    })

    if(this.isLoggedIn == false){
      this.$router.push("/login")
    }
    console.log("App.vue: try authentication and user info")
    this.getIsLoggedIn()
  },
  mounted() {
    const resF = (x) => {
      this.backendStatus = x
      if(BooksysBackend.needsSetup(this.backendStatus)){
        this.$router.push('/setup');
      }
    }
    BooksysBackend.getStatus(resF)
  }
})
</script>