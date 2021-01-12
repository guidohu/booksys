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
              <b-input-group-prepend is-text>
                hrs
              </b-input-group-prepend>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="input-group-description"
            label="Description"
            label-for="input-description"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-textarea
                id="input-description"
                v-model="form.description"
                type="text"
                placeholder="Add your maintenance notes here..."
                rows="3"
                max-rows="3"
              />
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
import { mapActions, mapGetters } from 'vuex';
import WarningBox from '@/components/WarningBox';
import {
  BContainer,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BInputGroup,
  BFormInput,
  BInputGroupPrepend,
  BFormTextarea,
  BButton
} from 'bootstrap-vue';

export default {
  name: "MaintenanceLogForm",
  components: {
    WarningBox,
    BContainer,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BInputGroup,
    BFormInput,
    BInputGroupPrepend,
    BFormTextarea,
    BButton
  },
  data() {
    return {
      errors: [],
      form: {
        engineHours: null,
        description: null
      }
    }
  },
  computed: {
    ...mapGetters('login', [
      'userInfo'
    ]),
  },
  methods: {
    add: function() {
      const entry = {
        user_id: this.userInfo.id,
        engine_hours: this.form.engineHours,
        description: this.form.description
      };
      this.addMaintenanceEntry(entry)
      .then(() => {
        this.resetForm();
      })
      .catch((errors) => this.errors = errors);
    },
    resetForm: function() {
      this.form = {
        engineHours: null,
        description: null
      };
      this.errors = [];
    },
    ...mapActions('boat', [
      'addMaintenanceEntry'
    ])
  }
}
</script>