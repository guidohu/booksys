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
import { mapActions, mapGetters } from 'vuex';
import WarningBox from '@/components/WarningBox';
import {
  BCard,
  BCardBody,
  BTable
} from 'bootstrap-vue';

export default {
  name: "LogCard",
  components: {
    WarningBox,
    BCard,
    BCardBody,
    BTable
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
}
</script>