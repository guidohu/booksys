<template>
  <b-modal
    id="engineHourEntryModal"
    title="Engine Hour Entry"
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
            id="before"
            label="Before"
            label-for="before-input"
            description=""
          >
            <b-form-input
              id="before-input"
              v-model="form.beforeHours"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="after"
            label="After"
            label-for="after-input"
            description=""
          >
            <b-form-input
              id="after-input"
              v-model="form.afterHours"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="difference"
            label="Difference"
            label-for="difference-input"
            description=""
          >
            <b-form-input
              id="difference-input"
              v-model="form.deltaHours"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="input-group-type"
            label="Type"
            label-for="input-type"
          >
            <toggle-button 
              id="input-type"
              :value="form.type"
              @change="toggleType"
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{checked: 'Course', unchecked: 'Private'}"
            />
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
import Vue from 'vue';
import { mapActions } from 'vuex';
import { ToggleButton } from 'vue-js-toggle-button';
import WarningBox from '@/components/WarningBox';
import moment from 'moment';
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BButton,
  BIconCheck,
  BIconX
} from 'bootstrap-vue';

export default Vue.extend({
  name: 'EngineHourEntryModal',
  props: [ 'engineHourEntry' ],
  components: {
    WarningBox,
    ToggleButton,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
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
      return 70;
    }
  },
  watch: {
    engineHourEntry: function(newValue){
      this.setFormContent(newValue);
    }
  },
  methods: {
    setFormContent: function(entry){
      if(entry != null){
        this.form = {
          id: entry.id,
          date: moment(entry.time, "X").format("DD.MM.YYYY HH:mm"),
          driver: entry.user_first_name,
          beforeHours: entry.before_hours,
          afterHours: entry.after_hours,
          deltaHours: entry.delta_hours,
          type: (entry.type == 0) ? false : true
        }
      }
    },
    toggleType: function(){
      this.form.type = !this.form.type;
    },
    close: function(){
      this.$bvModal.hide('engineHourEntryModal');
    },
    save: function(){
      const update = {
        id: this.form.id,
        type: (this.form.type == true) ? 1 : 0
      };
      this.updateEngineHours(update)
      .then(() => this.close())
      .catch((errors) => this.errors = errors);
    },
    ...mapActions('boat', [
      'updateEngineHours'
    ])
  },
  created() {
    this.setFormContent(this.engineHourEntry);
  }
})
</script>