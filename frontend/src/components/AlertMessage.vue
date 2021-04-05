<template>
  <b-card>
    <b-alert show variant="danger">
      {{ alertMessage }}
    </b-alert>
  </b-card>
</template>

<script>
import { BCard, BAlert } from "bootstrap-vue";

export default {
  name: "AlertMessage",
  components: {
    BCard,
    BAlert,
  },
  props: ["alertMessage"],
};
</script>
