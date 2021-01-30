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
import * as dayjs from 'dayjs';
import {
  BCard,
  BCardHeader,
  BCardBody,
  BRow,
  BIconBrightnessAltLowFill,
  BIconBrightnessAltHigh,
  BCol
} from 'bootstrap-vue';

export default {
  name: 'ConditionInfoCard',
  components: {
    BCard,
    BCardHeader,
    BCardBody,
    BRow,
    BCol,
    BIconBrightnessAltLowFill,
    BIconBrightnessAltHigh
  },
  props: [ 
    'sunrise',
    'sunset' ],
  computed: {
    sunriseString: function(){
      if(this.sunrise == null){
        return "n/a";
      }
      return dayjs(this.sunrise).format("HH:mm");
    },
    sunsetString: function(){
      if(this.sunset == null){
        return "n/a";
      }
      return dayjs(this.sunset).format("HH:mm");
    }
  }
}
</script>
