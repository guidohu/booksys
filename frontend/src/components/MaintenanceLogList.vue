<template>
  <div>
    <WarningBox v-if="errors.length > 0" :errors="errors"/>
    <b-table v-if="errors.length == 0" 
      hover 
      small
      :items="items" 
      :fields="columns"
      class="text-left"
    />
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import moment from "moment-timezone";
import WarningBox from '@/components/WarningBox';
import { BTable } from 'bootstrap-vue';

export default {
  name: "MaintenanceLogList",
  components: {
    WarningBox,
    BTable
  },
  data() {
    return {
      items: [],
      columns: [],
      errors: []
    };
  },
  computed: {
    ...mapGetters('boat',[
      'getMaintenanceLog'
    ])
  },
  watch: {
    getMaintenanceLog: function(newEntries){
      this.setItems(newEntries);
    }
  },
  methods: {
    ...mapActions('boat',[
      'queryMaintenanceLog'
    ]),
    setItems: function(logs){
      this.items = [];
      logs.slice(0,200).forEach( l => {
        this.items.push(l);
      })
    },
    setColumns: function(){
      this.columns = [
        {
          key: 'timestamp',
          label: 'Date',
          sortable: true,
          formatter: (value) => { return moment(value, "X").format("DD.MM.YYYY HH:mm"); }
        },
        {
          key: 'engine_hours',
          label: 'EngineHrs',
          sortable: true
        },
        {
          key: 'user_first_name',
          label: 'Driver',
          sortable: true
        },
        {
          key: 'description',
          label: 'Maintenance Work',
          sortable: false
        }
      ]
    }
  },
  created() {
    // generate header of table
    this.setColumns();

    // query content
    this.queryMaintenanceLog()
    .then( () => {
      this.setItems(this.getMaintenanceLog);
    })
    .catch( (errors) => {
      this.errors = errors;
    });
  }
}
</script>