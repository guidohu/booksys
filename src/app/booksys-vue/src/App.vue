<template>
  <b-container id="app" class="fill">
    <b-row v-if="isDesktop" class="display noborder nobackground">
      <div class="vcenter_placeholder">
        <div class="vcenter_content">
          <AlertMessage v-if="backendStatus && backendStatus.ok == false" alertMessage="The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible."/>
          <router-view v-else/>
        </div>
      </div>
    </b-row>
    <div v-if="isMobile">
      <AlertMessage v-if="backendStatus && backendStatus.ok == false" alertMessage="The webpage is currently not working due to the backend not being available. Please let the Administrator know and this will get fixed as soon as possible."/>
      <router-view v-else/>
    </div>
    <footer>
      <a href="https://github.com/guidohu/booksys">Find me on Github</a>
    </footer>
  </b-container>
</template>

<script>
import { BooksysBrowser } from './libs/browser'
import { BooksysBackend } from './libs/backend'
import AlertMessage from './components/AlertMessage.vue'
import Vue from 'vue'
import { mapGetters, mapState } from 'vuex'

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
  methods: {
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
    if(this.isLoggedIn && this.userInfo != null){
      this.$router.push("/dashboard")
    }else{
      this.$router.push("/login")
    }
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
  /* @import "./assets/bootstrap/css/bootstrap.min.css"; */
  @import "./assets/bootstrap/css/bootstrap-theme-bc.css";
  @import "./assets/css/style.css";

  .fill {
    height: 100vh;
    width: 100vw;
  }
</style>
