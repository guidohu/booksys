<template>
  <b-modal
    id="engineHourEntryModal"
    title="Engine Hour Entry"
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
            ></b-form-input>
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
            ></b-form-input>
          </b-form-group>
          <engine-hours
            label="Before"
            v-model="form.beforeHours"
            :display-format="displayFormat"
            :disabled="true"
          />
          <engine-hours
            label="After"
            v-model="form.afterHours"
            :display-format="displayFormat"
            :disabled="true"
          />
          <engine-hours
            label="Difference"
            v-model="form.deltaHours"
            :display-format="displayFormat"
            :disabled="true"
          />
          <b-form-group
            id="input-group-type"
            label="Type"
            label-for="input-type"
            label-cols="3"
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
import { mapActions } from 'vuex';
import { ToggleButton } from 'vue-js-toggle-button';
import WarningBox from '@/components/WarningBox';
import EngineHours from '@/components/forms/inputs/EngineHours';
import * as dayjs from 'dayjs';
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

export default {
  name: 'EngineHourEntryModal',
  props: [ 'engineHourEntry', 'visible', 'displayFormat' ],
  components: {
    WarningBox,
    EngineHours,
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
          date: dayjs(entry.time*1000).format("DD.MM.YYYY HH:mm"),
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
      this.$emit('update:visible', false);
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
}
</script>