<template>
  <!-- <div class="login-view"> -->
  <b-overlay
    id="overlay-background"
    :show="isLoading"
    spinner-type="border"
    spinner-variant="info"
    rounded="sm"
  >
    <b-container>
      <b-row class="text-center mb-3">
        <b-col cols="12">
          <b-img :src="getLogoFile" fluid :height="100"/>
        </b-col>
      </b-row>
      <login-modal 
        v-if="showLogin"
        :statusMessage="status" 
        :initialUsername="username" 
        @login="handleLogin"
      />
    </b-container>
  </b-overlay>
  <!-- </div> -->
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import LoginModal from '@/components/LoginModal.vue';
import { BooksysBrowser } from '@/libs/browser';
import {
  BContainer,
  BOverlay,
  BImg,
  BRow,
  BCol
} from 'bootstrap-vue';

export default {
  name: 'Login',
  components: {
    LoginModal,
    BContainer,
    BOverlay,
    BImg,
    BRow,
    BCol
  },
  data() {
    return {
      isLoading: true,
      showLogin: false,
      status: null
    }
  },
  computed: {
    isMobile: function () {
      return BooksysBrowser.isMobile()
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    },
    ...mapGetters('login', [
      'username'
    ]),
    ...mapGetters('configuration', [
      'getLogoFile'
    ])
  },
  methods: {
    ...mapActions('login', [
      'getIsLoggedIn',
      'login'
    ]),
    ...mapActions('configuration', [
      'queryLogoFile'
    ]),
    handleLogin: function (username, password) {
      this.login({ 
        username: username, 
        password: password
      })
      .then(() => {
        if(this.$route.query != null && this.$route.query.target != null){
          this.$router.push(this.$route.query.target);
        }else{
          this.$router.push('/dashboard');
        }
      })
      .catch((errors) => {
        if(errors.length > 0){
          this.status = errors[0];
        }
      })
    }
  },
  created () {
    this.isLoading = true;

    this. queryLogoFile()
    .then((uri) => console.log("Loaded logo image", uri))
    .catch((errors) => console.log(errors));

    this.getIsLoggedIn()
    .then((loggedIn) => {
      if(loggedIn == true){
        console.log('User is logged in: redirect to content');
        if(this.$route.query != null && this.$route.query.target != null){
          this.$router.push(this.$route.query.target);
        }else{
          this.$router.push('/dashboard');
        }
      }else{
        this.showLogin = true;
      }
      this.isLoading = false;
    })
    .catch(errors => {
      this.isLoading = false;
      this.showLogin = true;
      console.error("Errors while getIsLoggedIn was called:", errors);
    })
  }
}
</script>
