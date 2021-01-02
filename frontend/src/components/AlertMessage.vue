<template>
  <div class="alert alert-danger">
    {{ alertMessage }}
  </div>
</template>

<script>
export default {
  name: 'AlertMessage',
  components: {
  },
  props: ['alertMessage']
}
</script>
