<template>
  <b-card no-body class="text-left">
    <b-card-header>
      Statistics
    </b-card-header>
    <b-card-body>
      <b-row>
        <b-col cols="6">
          <b-row>
            <b-col cols="12">
              Riding Time: {{heatTimeMinutesYTD}} min
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              Riding Time (all-time): {{heatTimeMinutes}} min
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              Cost: {{heatCostYTD}} {{currency}}
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              Cost (all-time): {{heatCost}} {{currency}}
            </b-col>
          </b-row>
        </b-col>
        <b-col cols="6">
          <b-button size="sm" type="button" variant="outline-info" v-on:click="showLatestHeats">
            <b-icon-list-ul></b-icon-list-ul>
            Latest Heats
          </b-button>
        </b-col>
      </b-row>
    </b-card-body>
    <UserHeatsModal/>
  </b-card>
</template>

<script>
import Vue from 'vue'
import { mapGetters } from 'vuex'
import UserHeatsModal from './UserHeatsModal'

export default Vue.extend({
  name: 'UserStatisticsCard',
  components: {
    UserHeatsModal
  },
  computed: {
    ...mapGetters('login', [
      'userInfo'
    ]),
    ...mapGetters('configuration', [
      'currency'
    ]),
    ...mapGetters('user', [
      'heatHistory',
      'heatTimeMinutes',
      'heatTimeMinutesYTD',
      'heatCost',
      'heatCostYTD'
    ])
  },
  methods: {
    showLatestHeats: function(){
      this.$bvModal.show('userHeatsModal');
    }
  }
})
</script>