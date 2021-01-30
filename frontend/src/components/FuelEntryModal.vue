<template>
  <b-modal
    id="fuelEntryModal"
    title="Fuel Entry"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
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
            label-cols="3"
          >
            <b-form-input
              size="sm"
              id="date-input"
              v-model="form.date"
              type="text"
              placeholder=""
              disabled
            />
          </b-form-group>
          <b-form-group
            id="driver"
            label="Driver"
            label-for="driver-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              size="sm"
              id="driver-input"
              v-model="form.driver"
              type="text"
              placeholder=""
              disabled
            />
          </b-form-group>
          <b-form-group
            v-if="form.averageFuelPerHour != null"
            id="liters-per-hour"
            label="Consumption"
            label-for="fuel-input"
            description=""
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-input
                id="liter-per-hour-input"
                v-model="form.averageFuelPerHour"
                type="text"
                placeholder=""
                disabled
              />
              <b-input-group-append is-text>
                ltrs/h
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <engine-hours
            label="Engine Hours"
            v-model="form.engineHours"
            :display-format="getEngineHourFormat"
          />
          <b-form-group
            id="liters"
            label="Fuel"
            label-for="fuel-input"
            description=""
            label-cols="3"
          >
            <b-input-group size="sm">
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
            label-cols="3"
          >
            <b-input-group size="sm">
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
            label-cols="3"
          >
            <b-input-group size="sm">
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
            label-cols="3"
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
            label-cols="3"
          >
            <b-input-group size="sm">
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
        <b-icon-check/>
        Save
      </b-button>
      <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x/>
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { ToggleButton } from 'vue-js-toggle-button';
import WarningBox from '@/components/WarningBox';
import EngineHours from '@/components/forms/inputs/EngineHours';
import { 
  formatCurrency,
  formatFuelConsumption,
  formatFuel
} from '@/libs/formatters';
import moment from 'moment';
import {
  BModal,
  BCol,
  BRow,
  BForm,
  BFormGroup,
  BFormInput,
  BInputGroup,
  BInputGroupAppend,
  BButton,
  BIconCheck,
  BIconX
} from 'bootstrap-vue';

export default {
  name: 'fuelEntryModal',
  props: [ 'fuelEntry', 'visible' ],
  components: {
    WarningBox,
    EngineHours,
    ToggleButton,
    BModal,
    BCol,
    BRow,
    BForm,
    BFormGroup,
    BFormInput,
    BInputGroup,
    BInputGroupAppend,
    BButton,
    BIconCheck,
    BIconX
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
      'getCurrency',
      'getEngineHourFormat'
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
          cost: formatCurrency(cost, null),
          costGross: formatCurrency(costGross, null),
          costNet: formatCurrency(costNet, null),
          engineHours: entry.engine_hours,
          fuel: formatFuel(entry.liters),
          driver: entry.user_first_name,
          averageFuelPerHour: formatFuelConsumption(entry.avg_liters_per_hour)
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
      this.$emit('update:visible', false);
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
}
</script>