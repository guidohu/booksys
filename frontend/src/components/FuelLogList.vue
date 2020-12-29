<template>
  <div>
    <div v-if="errors.length > 0">
      <WarningBox :errors="errors"/>
    </div>
    <div v-else>
      <FuelEntryModal
        :fuelEntry="selectedFuelEntry"
      />
      <b-table
        hover
        small
        :items="items"
        :fields="columns"
        @row-clicked="rowClick"
        :tbody-tr-class="rowClass"
        class="text-left"
      />
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import WarningBox from '@/components/WarningBox';
import FuelEntryModal from '@/components/FuelEntryModal';
import moment from "moment-timezone";

export default Vue.extend({
  name: "FuelLogList",
  components: {
    WarningBox,
    FuelEntryModal
  },
  data() {
    return {
      errors: [],
      items: [],
      columns: [],
      selectedFuelEntry: null
    }
  },
  computed: {
    ...mapGetters('boat', [
      'getFuelLog'
    ])
  },
  methods: {
    ...mapActions('boat', [
      'queryFuelLog'
    ]),
    setItems: function(logs){
      console.log("set Items to:", logs);
      this.items = [];

      if(logs == null){
        return;
      }

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
          key: 'user_first_name',
          label: 'Driver',
          sortable: true,
          formatter: (value, key, item) => { return item.user_first_name }
        },
        {
          key: 'engine_hours',
          label: 'EngineHrs',
          sortable: true
        },
        {
          key: 'liters',
          label: 'Fuel',
          sortable: true
        },
        {
          key: 'cost',
          label: 'Cost',
          sortable: true,
          formatter: (value, key, item) => { return this.getFuelCost(item) }
        }
      ]
    },
    getFuelCost: function(entry){
      // returns either net or gross values for the cost
      if(entry.cost != null){
        return Number(entry.cost);
      }else{
        return Number(entry.cost_brutto);
      }
    },
    rowClick: function(item, index, event){
      console.log("rowClick item", item, "index", index, "event", event);
      this.selectedFuelEntry = item;
      this.$bvModal.show('fuelEntryModal');
    },
    rowClass: function(item) {
      // we highlight entries that have a deduction
      if(item.cost_brutto != null){
        return "highlight clickable"
      }else{
        return "clickable"
      }
    }
  },
  watch: {
    getFuelLog: function(newValues) {
      console.log("fuel log values changed to:", newValues);
      this.setItems(newValues);
    }
  },
  created() {
    this.queryFuelLog()
    .catch((errors) => this.errors = errors);

    this.setColumns();
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