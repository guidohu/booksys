<template>
  <div class="alert alert-danger">
    {{ alertMessage }}
  </div>
</template>

<script>
import Vue from 'vue'

export default Vue.extend({
  name: 'AlertMessage',
  components: {
  },
  props: ['alertMessage']
})
</script>
