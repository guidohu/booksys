<template>
  <div class="home">
    <LoginModal msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
// @ is an alias to /src
import LoginModal from '@/components/LoginModal.vue'

export default {
  name: 'Login',
  components: {
    LoginModal
  }
}
</script>
