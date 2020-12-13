<template>
  <div>
    <WarningBox v-if="errors.length > 0" :errors="errors"/>
    <EngineHourEntryModal
      :engineHourEntry="selectedEngineHourLogEntry"
    />
    <b-table v-if="errors.length == 0" 
      hover 
      small
      :items="items" 
      :fields="columns" 
      @row-clicked="rowClick" 
      :tbody-tr-class="rowClass"
      class="text-left"
    />
  </div>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import moment from "moment-timezone";
import WarningBox from '@/components/WarningBox';
import EngineHourEntryModal from '@/components/EngineHourEntryModal';

export default Vue.extend({
  name: "EngineHourLogList",
  components: {
    WarningBox,
    EngineHourEntryModal
  },
  data() {
    return {
      items: [],
      columns: [],
      errors: [],
      selectedEngineHourLogEntry: null,
    };
  },
  computed: {
    ...mapGetters('boat',[
      'getEngineHourLog'
    ])
  },
  watch: {
    getEngineHourLog: function(newEntries){
      console.log("getEngineHourLog just changed to", newEntries);
      this.setItems(newEntries);
    }
  },
  methods: {
    ...mapActions('boat',[
      'queryEngineHourLog'
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
          key: 'time',
          label: 'Date',
          sortable: true,
          formatter: (value) => { return moment(value, "X").format("DD.MM.YYYY HH:mm"); }
        },
        {
          key: 'user_first_name',
          label: 'Driver',
          sortable: true,
          formatter: (value, key, item) => { return item.user_first_name }
        },
        {
          key: 'before_hours',
          label: 'Before',
          sortable: true
        },
        {
          key: 'after_hours',
          label: 'After',
          sortable: true
        },
        {
          key: 'delta_hours',
          label: 'Diff',
          sortable: true
        }
      ]
    },
    rowClick: function(item){
      this.selectedEngineHourLogEntry = item;
      this.$bvModal.show('engineHourEntryModal');
    },
    rowClass: function(item) {
      if(item.type == 0){
        return "clickable";
      }else{
        return "highlight clickable";
      }
    }
  },
  created() {
    // generate header of table
    this.setColumns();

    // query content
    this.queryEngineHourLog()
    .then( () => {
      this.setItems(this.getEngineHourLog);
    })
    .catch( (errors) => {
      this.errors = errors;
    });
  }
})
</script>

<style>
  tr.highlight {
    color: #ffffff;
    background-color: rgb(91,192,222);
  }

  tr.clickable {
    cursor: pointer;
  }
</style>