<template>
  <div>
    <b-table 
      sticky-header 
      striped 
      hover 
      responsive
      :items="items"
      :fields="fields"
    >
      <template v-slot:cell(riders)="data">
          <b-row v-for="rider in data.item.riders" :key="rider.id">
            <b-col>
              {{ rider.first_name}} {{ rider.last_name}}
            </b-col>
          </b-row>
      </template>
      <template v-slot:cell(action)="data">
        <b-button type="button" size="sm" variant="outline-danger" v-on:click="cancelSession(data.item.id)"><b-icon-x/>Cancel Session</b-button>
      </template>
    </b-table>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import moment from 'moment-timezone/moment-timezone';
import {
  BTable,
  BRow,
  BCol,
  BButton,
  BIconX
} from 'bootstrap-vue';

export default Vue.extend({
  name: 'UserSessionsTable',
  components: {
    BTable,
    BRow,
    BCol,
    BButton,
    BIconX
  },
  props: [
    'userSessions',
    'showCancel'
  ],
  data: function(){
    return {
      fields: [
        {
          key: "start",
          label: "Date",
          formatter: (value, key, item) => { return this.formatTime(item) }
        },
        {
          key: "riders",
          label: "Riders",
        }
      ],
      items: this.userSessions
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getConfiguration',
      'getTimezone'
    ])
  },
  methods: {
    formatTime: function(item) {
      return moment(item.start, "X").tz(this.getTimezone).format("DD.MM.YYYY HH:mm")
      + " - " + moment(item.end, "X").tz(this.getTimezone).format("HH:mm");
    },
    cancelSession: function(sessionId) {
      console.log("Cancel session with id", sessionId);
      // forward to parent
      this.$emit('cancel', sessionId);
    },
    ...mapActions('configuration', [
      'queryConfiguration'
    ])
  },
  created() {
    this.queryConfiguration();
  },
  beforeMount() {
    if(this.$props.showCancel == true){
      this.fields.push({
        key: "action",
        label: "Action"
      });
    }
  }
})
</script>