<template>
  <b-modal
    id="sessionEditorModal"
    :title="title"
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
            id="title"
            label="Title"
            label-for="title-input"
            description=""
          >
            <b-form-input
              id="title-input"
              v-model="form.title"
              type="text"
              placeholder=""
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="description"
            label="Description"
            label-for="description-input"
            description=""
          >
            <b-form-textarea
              id="description-input"
              v-model="form.description"
              type="text"
              rows="2"
              placeholder=""
            />
          </b-form-group>
          <b-form-group
            id="start"
            label="Start"
            label-for="start-date-input"
            description=""
          >
            <b-row>
              <b-col cols="6">
                <b-form-input 
                  id="start-date-input"
                  type="date" 
                  v-model="form.startDate"
                  placeholder=""
                />
              </b-col>
              <b-col cols="6">
                <b-form-input 
                  id="start-time-input"
                  type="time" 
                  v-model="form.startTime"
                  placeholder=""
                />
              </b-col>
            </b-row>
          </b-form-group>
          <b-form-group
            id="end"
            label="End"
            label-for="end-date-input"
            description=""
          >
            <b-row>
              <b-col cols="6">
              <b-form-input 
                id="end-date-input"
                type="date" 
                v-model="form.endDate"
                placeholder=""
              />
              </b-col>
              <b-col cols="6">
                <b-form-input 
                  id="end-time-input"
                  type="time" 
                  v-model="form.endTime"
                  placeholder=""
                />
              </b-col>
            </b-row>
          </b-form-group>
          <b-form-group
            id="capacity"
            v-bind:label="getMaximumRidersLabel"
            label-for="capacity-input"
            description=""
          >
            <b-form-input type="number" v-model="form.maximumRiders"/>
          </b-form-group>
          <b-form-group
            id="type"
            label="Type"
            label-for="type-input"
            description="Choose whether this session is open to others or a closed group"
          >
            <b-form-checkbox 
              v-model="form.type" 
              switch
            >
              {{getTypeText}}
            </b-form-checkbox>
          </b-form-group>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <b-button v-if="form.id==null" type="button" variant="outline-info" v-on:click="save">
        <b-icon-check/>
        Create
      </b-button>
      <b-button v-if="form.id!=null" type="button" variant="outline-info" v-on:click="save">
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
import { mapGetters, mapActions } from 'vuex';
import WarningBox from '@/components/WarningBox';
import Session, { SESSION_TYPE_OPEN, SESSION_TYPE_PRIVATE } from '@/dataTypes/session';
import moment from 'moment';
import 'moment-timezone';
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BFormTextarea,
  BFormCheckbox,
  BButton,
  BIconCheck,
  BIconX
} from 'bootstrap-vue';

export default Vue.extend({
  name: 'SessionEditorModal',
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
    BFormTextarea,
    BFormCheckbox,
    BButton,
    BIconCheck,
    BIconX
  },
  props: [
    'defaultValues'
  ],
  data(){
    return {
      title: "",
      errors: [],
      form: {
        id: null,
        title: null,
        description: null,
        startDate: null,
        startTime: null,
        endDate: null,
        endTime: null,
        maximumRiders: null,
        type: null
      }
    }
  },
  computed: {
    getTypeText: function() {
      if(this.form.type == true){
        return "Private Session"
      }else{
        return "Open Session"
      }
    },
    getMaximumRidersLabel: function() {
      if(this.form.id == null){
        return "Maximum Riders"
      }else{
        return "Additional Slots for Riders"
      }
    },
    ...mapGetters('configuration', [
      'getTimezone',
      'getMaximumNumberOfRiders'
    ])
  },
  methods: {
    ...mapActions('sessions', [
      'createSession',
      'editSession'
    ]),
    save: function(event){
      event.preventDefault();
      console.log("SessionEditorModal, save:", this.form);

      const type = (this.form.type == SESSION_TYPE_PRIVATE)
        ? SESSION_TYPE_PRIVATE
        : SESSION_TYPE_OPEN;

      const startString = this.form.startDate + " " + this.form.startTime;
      // const endString = this.form.endDate + " " + this.form.endTime;
      console.log("SessionEditorModal: startString moment:", moment(startString));
      console.log("SessionEditorModal: startString moment.tz(timezone):", moment.tz(startString, this.getTimezone));
      console.log("SessionEditorModal: startString format:", moment.tz(startString, this.getTimezone).format());
      

      const session = new Session(
        this.form.id,
        this.form.title,
        this.form.description,
        moment.tz(this.form.startDate + " " + this.form.startTime, this.getTimezone).format(),
        moment.tz(this.form.endDate + " " + this.form.endTime, this.getTimezone).format(),
        this.form.maximumRiders,
        type
      );

      console.log("SessionEditorModal, save session dataType:", session);
      if(session.id == null){
        this.createSession(session)
        .then( response => {
          this.$emit("sessionCreatedHandler", response.session_id);
          this.close();
        })
        .catch( err => {
          this.errors = err;
        });
      }else{
        this.editSession(session)
        .then( () => {
          this.$emit("sessionEditedHandler", session.id);
          this.close();
        })
        .catch( err => {
          this.errors = err;
        })
      }
      
    },
    close: function(){
      this.errors = [];
      this.$bvModal.hide('sessionEditorModal');
    },
    setFormContent: function(){
      console.log("SessionEditorModal: Set defaults to:", this.defaultValues);

      // set ID
      this.form.id = (this.defaultValues != null && this.defaultValues.id != null) 
        ? this.defaultValues.id
        : null;

      if(this.form.id != null){
        this.title = "Edit Session";
      }else{
        this.title = "Create Session";
      }

      // set title
      this.form.title = (this.defaultValues != null && this.defaultValues.title != null) 
        ? this.defaultValues.title
        : null;

      // set description
      this.form.description = (this.defaultValues != null && this.defaultValues.description != null) 
        ? this.defaultValues.description
        : null;

      // set times
      this.form.startDate = (this.defaultValues != null && this.defaultValues.start != null) 
        ? moment(this.defaultValues.start).format("YYYY-MM-DD") 
        : moment().tz(this.getTimezone).format("YYYY-MM-DD");
      this.form.startTime = (this.defaultValues != null && this.defaultValues.start != null) 
        ? moment(this.defaultValues.start).format("HH:mm") 
        : moment().tz(this.getTimezone).format("HH:mm");
      this.form.endDate = (this.defaultValues != null && this.defaultValues.end != null) 
        ? moment(this.defaultValues.end).format("YYYY-MM-DD") 
        : moment().tz(this.getTimezone).add(1, 'hour').format("YYYY-MM-DD");
      this.form.endTime = (this.defaultValues != null && this.defaultValues.end != null) 
        ? moment(this.defaultValues.end).format("HH:mm") 
        : moment().tz(this.getTimezone).add(1, 'hour').format("HH:mm");

      // set maximum riders
      this.form.maximumRiders = (this.defaultValues != null && this.defaultValues.maximumRiders != null) 
        ? this.defaultValues.maximumRiders
        : this.getMaximumNumberOfRiders;

      // set session type
      this.form.type = (this.defaultValues != null && this.defaultValues.type != null) 
        ? this.defaultValues.type
        : null;

      console.log("SessionEditorModal: Form values are now:", this.form);
    }
  },
  created() {
    // set form content based on props
    this.setFormContent();
  },
  watch: {
    defaultValues: function() {
      // set form content whenever the default props change
      this.setFormContent();
    }
  }
})
</script>

<style scoped>
</style>