<template>
  <b-modal
    id="heatEntryModal"
    title="Heat Entry"
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
              id="date-input"
              v-model="form.date"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="rider"
            label="Rider"
            label-for="rider-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="rider-input"
              v-model="form.rider"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="fare"
            label="Fare"
            label-for="fare-input"
            description=""
            label-cols="3"
            disabled
          >
            <b-input-group>
              <b-form-input
                id="fare-input"
                v-model="form.fare"
                type="text"
                placeholder=""
                disabled
              />
              <b-input-group-append is-text>
                {{getCurrency}}/min
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="duration"
            label="Duration"
            label-for="duration-input"
            description=""
            label-cols="3"
          >
            <b-input-group>
              <b-form-input
                id="duration-input"
                v-model="form.duration"
                v-on:input="durationChangeHandler"
                type="text"
                placeholder=""
              />
              <b-input-group-append is-text>
                mm:ss
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="cost"
            label="Cost"
            label-for="cost-input"
            description=""
            label-cols="3"
          >
            <b-input-group>
              <b-form-input
                id="cost-input"
                v-model="form.cost"
                type="text"
                placeholder=""
                disabled
              />
              <b-input-group-append is-text>
                {{getCurrency}}
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="comment"
            label="Comment"
            label-for="comment-input"
            description=""
            label-cols="3"
          >
            <b-form-textarea
              id="comment-input"
              v-model="form.comment"
              rows="3"
            />
          </b-form-group>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <div class="text-left d-inline">
        <b-button class="mr-1" type="button" variant="outline-danger" v-on:click="remove">
          <b-icon-trash/>
          Delete
        </b-button>
      </div>
      <div class="text-right d-inline">
        <b-button class="ml-4" type="button" variant="outline-info" v-on:click="save">
          <b-icon-check/>
          Save
        </b-button>
        <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
          <b-icon-x/>
          Cancel
        </b-button>
      </div>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import { mapGetters, mapActions } from 'vuex';
import moment from 'moment';
import { sprintf } from 'sprintf-js';
import WarningBox from '@/components/WarningBox';
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BInputGroup,
  BInputGroupAppend,
  BFormTextarea,
  BButton,
  BIconTrash,
  BIconCheck,
  BIconX
} from 'bootstrap-vue';

export default Vue.extend({
  name: 'HeatEntryModal',
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
    BInputGroup,
    BInputGroupAppend,
    BFormTextarea,
    BButton,
    BIconTrash,
    BIconCheck,
    BIconX
  },
  props: [ 'heat' ],
  data() {
    return {
      errors: [],
      form: {
        date: null
      }
    };
  },
  computed: {
    ...mapGetters('configuration', [
      'getCurrency'
    ])
  },
  watch: {
    heat: function(newHeat){
      this.setFormDefaults(newHeat);
    }
  },
  methods: {
    setFormDefaults: function(heatData) {
      this.form = {
        id: heatData.heat_id,
        date: moment(heatData.timestamp, 'X').format("DD.MM.YYYY HH:mm"),
        userId: heatData.user_id,
        rider: heatData.first_name + " " + heatData.last_name,
        fare: sprintf('%.2f', heatData.price_per_min),
        duration: this.formatDuration(heatData.duration_s),
        cost: heatData.cost,
        comment: heatData.comment
      };
    },
    formatDuration: function(durationS) {
      const seconds = durationS % 60;
      const minutes = Math.floor(((durationS - seconds)) / 60);
      return sprintf('%02d:%02d', minutes, seconds);
    },
    durationChangeHandler: function() {
      console.log("duration changed to:", this.form.duration);
      if(this.form.duration == null || !this.form.duration.match(/^\d+:\d+$/)){
        this.form.cost = "...";
        return;
      }else{
        const durationParts = this.form.duration.split(':');
        const seconds = Number(durationParts[1]);
        const minutes = Number(durationParts[0]);
        const durationSeconds = seconds + (60 * minutes);
        this.form.cost = Math.round(durationSeconds * this.form.fare / 60 * 100) / 100;
        return;
      }
    },
    close: function() {
      this.$bvModal.hide("heatEntryModal");
    },
    remove: function() {
      this.removeHeat(this.form.id)
      .then(() => this.close())
      .catch((errors) => this.errors = errors);
    },
    save: function() {
      if(this.form.duration == null){
        this.errors = ["No duration provided"];
        return;
      }
      const durations = this.form.duration.split(":");
      if(durations.length != 2){
        this.errors = ["Duration does not have a valid format such as 23:15."];
        return;
      }
      const seconds = Number(durations[1]);
      const minutes = Number(durations[0]);
      if(isNaN(seconds) || isNaN(minutes)){
        this.errors = ['The duration must use numbers such as 23:15.'];
        return;
      }

      const durationSeconds = seconds + 60 * minutes;
      const heatUpdate = {
        id: this.form.id,
        userId: this.form.userId,
        duration: durationSeconds,
        comment: this.form.comment
      };
      this.updateHeat(heatUpdate)
      .then(() => {
        this.errors = [];
        this.close();
      })
      .catch((errors) => this.errors = errors);
    },
    ...mapActions('heats', [
      'removeHeat',
      'updateHeat'
    ])
  },
  created() {
    if(this.heat != null){
      this.setFormDefaults(this.heat);
    }
  }
})
</script>