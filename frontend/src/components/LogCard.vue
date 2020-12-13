<template>
  <b-card no-body class="text-left">
    <b-card-body>
      <b-table
        v-if="errors.length==0"
        striped
        hover
        small
        :fields="columns"
        :items="getLogLines"
      />
      <WarningBox v-if="errors.length>0" :errors="errors"/>
    </b-card-body>
  </b-card>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import WarningBox from '@/components/WarningBox';

export default Vue.extend({
  name: "LogCard",
  components: {
    WarningBox
  },
  data() {
    return {
      errors: [],
      columns: [
        {
          key: "time",
          label: "Time"
        },
        {
          key: "log",
          label: "Message"
        }
      ],
      rows: []
    }
  },
  computed: {
    ...mapGetters('log', [
      'getLogLines'
    ])
  },
  methods: {
    ...mapActions('log', [
      'queryLogLines'
    ])
  },
  created() {
    this.queryLogLines()
    .catch((errors) => this.errors = errors);
  }
})
</script>