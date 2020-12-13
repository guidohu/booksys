<template>
  <div id="app">
    <b-container v-if="isDesktop" class="fill">
      <b-row v-if="isDesktop" class="display noborder nobackground">
        <div class="vcenter_placeholder">
          <div class="vcenter_content">
            <AlertMessage v-if="backendStatus && backendStatus.ok == false" alertMessage="The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible."/>
            <router-view v-else/>
          </div>
        </div>
      </b-row>
    </b-container>
    <b-container v-if="isMobile" fluid="true">
      <div v-if="isMobile">
        <AlertMessage v-if="backendStatus && backendStatus.ok == false" alertMessage="The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible."/>
        <router-view v-else/>
      </div>
      <footer class="legal-footer" v-if="isDesktop">
        Copyright 2013-2020 by Guido Hungerbuehler <a href="https://github.com/guidohu/booksys">Find me on Github</a>
      </footer>
    </b-container>
  </div>
</template>

<script>
import { BooksysBrowser } from './libs/browser'
import { BooksysBackend } from './libs/backend'
import AlertMessage from './components/AlertMessage.vue'
import Vue from 'vue'
import { mapGetters, mapState, mapActions } from 'vuex'

export default Vue.extend({
  name: 'App',
  data: function() {
    return {
      backendAvailable: false, // if the backend is not available
      backendStatus: null
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
      if(newValue == false){
        console.log("App.vue: detected not logged in -> forward to /login")
        this.$router.push("/login")
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

<style>
  .fill {
    height: 100vh;
    width: 100vw;
    background-color: transparent;
  }
</style>
