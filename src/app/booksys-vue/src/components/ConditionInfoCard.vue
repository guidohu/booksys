<template>
  <b-card no-body class="text-left">
    <b-card-header>
      Condition Info
    </b-card-header>
    <b-card-body>
      <b-row>
        <b-col cols="5">
          <b-icon-brightness-alt-high></b-icon-brightness-alt-high>
          Sunrise
        </b-col>
        <b-col cols="7">
          {{ sunriseString }}
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="5">
          <b-icon-brightness-alt-low-fill></b-icon-brightness-alt-low-fill>
          Sunset
        </b-col>
        <b-col cols="7">
          {{ sunsetString }}
        </b-col>
      </b-row>
    </b-card-body>
  </b-card>
</template>

<script>
import Vue from 'vue';
import moment from 'moment-timezone';

export default Vue.extend({
  name: 'ConditionInfoCard',
  props: [ 
    'sunrise',
    'sunset' ],
  computed: {
    sunriseString: function(){
      return moment(this.sunrise, "X").format("HH:mm");
    },
    sunsetString: function(){
      return moment(this.sunset, "X").format("HH:mm");
    }
  }
})
</script>
