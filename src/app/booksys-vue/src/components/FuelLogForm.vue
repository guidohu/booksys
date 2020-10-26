<template>
  <b-container float class="text-left">
    <b-row>
      <b-col cols="12">
        <WarningBox v-if="errors.length > 0" :errors="errors"/>
        <b-form @submit="add">
          <b-form-group
            id="input-group-engine-hours"
            label="Engine Hrs"
            label-for="input-engine-hours"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-input
                id="input-engine-hours"
                v-model="form.engineHours"
                type="text"
              />
              <b-input-group-append is-text>
                hrs
              </b-input-group-append>
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
              <b-input-group-append is-text>
                {{getCurrency}}
              </b-input-group-append>
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
              <b-input-group-append is-text>
                ltrs
              </b-input-group-append>
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
import { mapActions, mapGetters } from 'vuex';
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
  computed: {
    ...mapGetters('login', [
      'userInfo'
    ]),
    ...mapGetters('configuration',[
      'getCurrency'
    ])
  },
  methods: {
    add: function() {
      console.log("add has been clicked");
      const entry = {
        user_id: this.userInfo.id,
        engine_hours: this.form.engineHours,
        liters: this.form.liters,
        cost: this.form.cost
      };
      this.addFuelEntry(entry)
      .then(() => {
        this.resetForm();
      })
      .catch((errors) => this.errors = errors);
    },
    resetForm: function() {
      this.form = {
        engineHours: null,
        cost: null,
        liters: null
      };
      this.errors = [];
    },
    ...mapActions('boat', [
      'addFuelEntry'
    ]),
    ...mapActions('configuration', [
      'queryConfiguration'
    ])
  },
  created() {
    this.queryConfiguration();
  }
})
</script>