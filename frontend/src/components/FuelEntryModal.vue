<template>
  <b-modal
    id="fuelEntryModal"
    title="Fuel Entry"
  >
    <b-row v-if="errors.length">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <WarningBox :errors="errors"/>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-form @submit="save">
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <b-form-group
            id="date"
            label="Date"
            label-for="date-input"
            description=""
          >
            <b-form-input
              id="date-input"
              v-model="form.date"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="driver"
            label="Driver"
            label-for="driver-input"
            description=""
          >
            <b-form-input
              id="driver-input"
              v-model="form.driver"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="engine-hours"
            label="Engine Hours"
            label-for="engine-hours-input"
            description=""
          >
            <b-input-group>
              <b-form-input
                id="engine-hours-input"
                v-model="form.engineHours"
                type="text"
                placeholder=""
              />
              <b-input-group-append is-text>
                hrs
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="liters"
            label="Fuel"
            label-for="fuel-input"
            description=""
          >
            <b-input-group>
              <b-form-input
                id="fuel-input"
                v-model="form.fuel"
                type="text"
                placeholder=""
              />
              <b-input-group-append is-text>
                ltrs
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="cost"
            v-if="!form.isDiscounted"
            label="Cost"
            label-for="cost-input"
            description=""
          >
            <b-input-group>
              <b-form-input
                id="cost-input"
                v-model="form.cost"
                type="text"
                placeholder=""
                v-on:change="costChange"
              />
              <b-input-group-append is-text>
                {{getCurrency}}
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="cost-gross"
            v-if="form.isDiscounted"
            label="Cost (Gross)"
            label-for="cost-gross-input"
            description=""
          >
            <b-input-group>
              <b-form-input
                id="cost-gross-input"
                v-model="form.costGross"
                type="text"
                placeholder=""
                v-on:change="costGrossChange"
              />
              <b-input-group-append is-text>
                {{getCurrency}}
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="discount"
            label="Discount"
            label-for="discount-input"
          >
            <toggle-button 
              id="discount-input"
              :value="form.isDiscounted"
              @change="toggleDiscount"
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{checked: 'Discount', unchecked: 'No Discount'}"
            />
          </b-form-group>
          <b-form-group
            id="cost-net"
            v-if="form.isDiscounted"
            label="Cost (net)"
            label-for="cost-net-input"
            description=""
          >
            <b-input-group>
              <b-form-input
                id="cost-net-input"
                v-model="form.costNet"
                type="text"
                placeholder=""
              />
              <b-input-group-append is-text>
                {{getCurrency}}
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <b-button type="button" variant="outline-info" v-on:click="save">
        <b-icon-check></b-icon-check>
        Save
      </b-button>
      <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x></b-icon-x>
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import { mapActions, mapGetters } from 'vuex';
import { ToggleButton } from 'vue-js-toggle-button';
import WarningBox from '@/components/WarningBox';
import moment from 'moment';

export default Vue.extend({
  name: 'fuelEntryModal',
  props: [ 'fuelEntry' ],
  components: {
    WarningBox,
    ToggleButton
  },
  data() {
    return {
      errors: [],
      form: {
        date: null
      }
    }
  },
  computed: {
    toggleWidth: function() {
      return 100;
    },
    ...mapGetters('configuration',[
      'getCurrency'
    ])
  },
  watch: {
    fuelEntry: function(newValue){
      this.setFormContent(newValue);
    }
  },
  methods: {
    setFormContent: function(entry){
      if(entry != null){
        // apply some logic depending on the backend data to get the right values for
        // net/gross and cost
        const costGross = (entry.cost_brutto == null) ? entry.cost : entry.cost_brutto;
        const costNet   = (entry.cost_brutto == null) ? null : entry.cost;
        const cost      = (entry.cost_brutto == null) ? entry.cost : entry.cost_brutto;

        this.form = {
          id: entry.id,
          date: moment(entry.timestamp, "X").format("DD.MM.YYYY HH:mm"),
          isDiscounted: (entry.cost != null && entry.cost_brutto != null) ? true : false,
          cost: cost,
          costGross: costGross,
          costNet: costNet,
          engineHours: entry.engine_hours,
          fuel: entry.liters,
          driver: entry.user_first_name
        };
      }
    },
    toggleDiscount: function(){
      this.form.isDiscounted = !this.form.isDiscounted;
    },
    costChange: function() {
      this.form.costGross = this.form.cost;
    },
    costGrossChange: function() {
      this.form.cost = this.form.costGross;
    },
    close: function(){
      this.$bvModal.hide('fuelEntryModal');
    },
    save: function(){
      if(this.form.isDiscounted && this.form.costNet == null){
        this.errors = [ "Cost (net) cannot be empty if you select a discount."];
        return;
      }

      let cost = 0;
      let costBrutto = 0;
      if(this.form.isDiscounted){
        cost = this.form.costNet;
        costBrutto = this.form.costGross;
      }else{
        cost = this.form.cost;
        costBrutto = null;
      }

      const updatedEntry = {
        id: this.form.id,
        engine_hours: this.form.engineHours,
        liters: this.form.fuel,
        cost: cost,
        cost_brutto: costBrutto
      };

      this.updateFuelEntry(updatedEntry)
      .then(() => this.close())
      .catch((errors) => this.errors = errors);
    },
    ...mapActions('boat', [
      'updateFuelEntry'
    ]),
    ...mapActions('configuration', [
      'queryConfiguration'
    ])
  },
  created() {
    this.setFormContent(this.fuelEntry);
    this.queryConfiguration();
  }
})
</script>