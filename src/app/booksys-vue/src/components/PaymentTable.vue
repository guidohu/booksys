<template>
  <div class="text-left">
    <WarningBox v-if="errors.length > 0" :errors="errors" dismissible="true" @dismissed="dismissedHandler"/>
    <div v-else>
      <b-row v-if="form.years.length>0" class="text-right">
        <b-col offset="6" cols="6" class="d-sm-none mr-1">
          <b-form-group
            id="year"
            label="Year"
            label-for="year-select"
            description=""
            label-cols="3"
          >
            <b-form-select
              id="year-select" 
              v-model="form.selectedYear" 
              :options="form.years" 
              @change="yearSelectionChangeHandler($event)"
            />
          </b-form-group>
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="8">
          <b-button v-on:click="showAddIncome" size="sm" variant="outline-info" class="mr-1 mb-2"><b-icon icon="plus"/>Income</b-button>
          <b-button v-on:click="showAddExpenditure" size="sm" variant="outline-info" class="mb-2"><b-icon icon="minus"/>Expenditure</b-button>
        </b-col>
        <b-col cols="4" class="d-none d-sm-block">
          <b-form-group
            id="year"
            label="Year"
            label-for="year-select"
            description=""
            label-cols="3"
          >
            <b-form-select
              id="year-select" 
              v-model="form.selectedYear" 
              :options="form.years" 
              @change="yearSelectionChangeHandler($event)"
            />
          </b-form-group>
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="12">
          <b-table 
            striped
            hover
            response
            small
            sticky-header
            borderless
            sort-by="date"
            :items="getTransactions"
            :fields="fields"
            :selectable="false"
            show-empty
            empty-text="No records to show"
          >
            <template #cell(action)="data">
              <div class="text-center">
                <b-button size="sm" style="font-size: 0.8em;" variant="light">
                  <b-icon v-on:click="deleteTransaction(data.item)" icon="trash" variant="danger"/>
                </b-button>
              </div>
            </template>
          </b-table>
        </b-col>
      </b-row>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapGetters, mapActions } from 'vuex';
import _ from 'lodash';
import moment from 'moment';
import WarningBox from '@/components/WarningBox';

export default Vue.extend({
  name: "PaymentTable",
  components: {
    WarningBox
  },
  data() {
    return {
      errors: [],
      form: {
        years: [],
        selectedYear: "any"
      },
      items: [],
      selectedItems: [],
      fields: [
        {
          key: "timestamp",
          label: "Date",
          sortable: true
        },
        {
          key: "type_name",
          label: "Type",
          sortable: true
        },
        {
          key: "name",
          label: "Name",
          sortable: true,
          formatter: (value, key, item) => item.fn + " " + item.ln
        },
        {
          key: "amount",
          label: "Amount",
          sortable: true,
          formatter: (value) => value + " " + this.getCurrency,
          tdClass: "text-right",
          thClass: "text-right"
        },
        {
          key: "comment",
          label: "Comment",
          sortable: false
        },
        {
          key: "action",
          label: "Action",
          sortable: false
        }
      ]
    }
  },
  computed: {
    ...mapGetters('accounting', [
      'getTransactions',
      'getYears'
    ]),
    ...mapGetters('configuration', [
      'getCurrency'
    ])
  },
  watch: {
    getYears: function(newValue){
      const availableYears = _.reverse(newValue);
      console.log(availableYears);
      this.form.years = availableYears.map(v => {
        return {
          value: v,
          text: v
        }
      });
    }
  },
  methods: {
    ...mapActions('accounting', [
      'queryTransactions',
      'queryYears'
    ]),
    ...mapActions('configuration', [
      'queryConfiguration'
    ]),
    showAddIncome: function() {
      console.log("TODO showAddIncome");
    },
    showAddExpenditure: function() {
      console.log("TODO showAddExpenditure");
    },
    yearSelectionChangeHandler: function(selection){
      this.queryTransactions(selection)
      .catch((errors) => this.errors = errors);
    },
    deleteTransaction: function(transaction){
      console.log("TODO delete transaction:", transaction);
    }
  },
  created(){
    this.queryConfiguration()
    .catch((errors) => this.errors = errors);

    this.queryYears()
    .catch((errors) => this.errors = errors);

    const currentYear = moment().year();
    this.form.selectedYear = currentYear;

    this.queryTransactions(currentYear)
    .catch((errors) => this.errors = errors);
  }
})
</script>