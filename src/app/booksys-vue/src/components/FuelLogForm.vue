<template>
  <b-container float class="text-left">
    <b-row>
      <b-col cols="12">
        <WarningBox v-if="errors.length > 0" :errors="errors"/>
        <b-form @submit="add">
          <b-form-group
            id="input-group-engine-hours"
            label="Before"
            label-for="input-engine-hours"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-input
                id="input-engine-hours"
                v-model="form.engineHours"
                type="text"
              />
              <b-input-group-prepend is-text>
                hrs
              </b-input-group-prepend>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="input-group-cost"
            label="Cost"
            label-for="input-cost"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-input
                id="input-cost"
                v-model="form.cost"
                type="text"
              />
              <b-input-group-prepend is-text>
                CURRENCY
              </b-input-group-prepend>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="input-liters"
            label="Liters"
            label-for="input-liters"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-input
                id="input-liters"
                v-model="form.liters"
                type="text"
              />
              <b-input-group-prepend is-text>
                ltrs
              </b-input-group-prepend>
            </b-input-group>
          </b-form-group>
          <b-row class="text-right">
            <b-col cols="9" offset="3">
              <b-button block variant="outline-info" size="sm" v-on:click="add">Add</b-button>
            </b-col>
          </b-row>
        </b-form>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import Vue from 'vue';
import WarningBox from '@/components/WarningBox';

export default Vue.extend({
  name: "FuelLogForm",
  components: {
    WarningBox
  },
  data() {
    return {
      errors: [],
      form: {
        engineHours: null,
        cost: null,
        liters: null
      }
    }
  },
  methods: {
    add: function() {
      console.log("add has been clicked");
    }
  }
})
</script>