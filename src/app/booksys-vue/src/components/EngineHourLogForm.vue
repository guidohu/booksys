<template>
  <b-container float class="text-left">
    <b-row>
      <b-col cols="12">
        <b-form @submit="add">
          <b-form-group
            id="input-group-driver"
            label="Driver"
            label-for="input-driver"
            label-cols="3"
          >
            <b-form-input
              id="input-driver"
              v-model="form.driverName"
              type="text"
              size="sm"
              disabled
            />
          </b-form-group>
          <b-form-group
            id="input-group-before-hours"
            label="Before"
            label-for="input-before-hours"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-input
                id="input-before-hours"
                v-model="form.beforeHours"
                type="text"
                :disabled="disableBefore"
              />
              <b-input-group-prepend is-text>
                hrs
              </b-input-group-prepend>
            </b-input-group>
          </b-form-group>
          <b-form-group
            v-if="showAfter"
            id="input-group-after-hours"
            label="After"
            label-for="input-after-hours"
            label-cols="3"
          >
            <b-input-group size="sm">
              <b-form-input
                id="input-after-hours"
                v-model="form.afterHours"
                type="text"
              />
              <b-input-group-prepend is-text>
                hrs
              </b-input-group-prepend>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="input-group-ty[e"
            label="Type"
            label-for="input-type"
            label-cols="3"
          >
            <toggle-button 
              id="input-type"
              :value="false"
              @change="toggleType"
              color="#17a2b8"
              :width="toggleWidth"
              :labels="{checked: 'Course', unchecked: 'Private'}"
            />
          </b-form-group>
          <b-row class="text-right">
            <b-col cols="9" offset="3">
              <b-button v-if="!disableBefore" block variant="outline-info" size="sm" v-on:click="add">Start</b-button>
              <b-button v-if="disableBefore" block variant="outline-info" size="sm" v-on:click="add">Finish</b-button>
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
import { ToggleButton } from 'vue-js-toggle-button';


export default Vue.extend({
  name: "EngineHourLogForm",
  components: {
    ToggleButton
  },
  data() {
    return {
      form: {
        driverId: null,
        driverName: "",
        beforeHours: null,
        afterHours: null,
        type: false
      },
      disableBefore: false,
      showAfter: false,
    }
  },
  computed: {
    ...mapGetters('boat', [
      'getEngineHourLogLatest'
    ]),
    ...mapGetters('login', [
      'userInfo'
    ]),
    toggleWidth: function() {
      return 70;
    }
  },
  watch: {
    userInfo: function(){
      console.log("userInfo just changed");
      this.setDriver();
    },
    getEngineHourLogLatest: function(newData){
      console.log('engineHourLatest just changed', newData);

      if(newData.before_hours != null && newData.after_hours != null){
        this.form.beforeHours = null;
      }else{
        this.form.beforeHours = newData.before_hours;
      }
      this.form.afterHours = null;

      this.setDriver();
      this.setDisableBefore();
      this.setShowAfter();
    }
  },
  methods: {
    ...mapActions('boat', [
      'queryEngineHourLogLatest',
      'addEngineHours'
    ]),
    setDisableBefore: function() {
      if(this.getEngineHourLogLatest != null && this.getEngineHourLogLatest.after_hours != null){
        this.disableBefore = false;
        return;
      }
      this.disableBefore = true;
    },
    setShowAfter: function() {
      if(this.getEngineHourLogLatest != null && this.getEngineHourLogLatest.after_hours == null){
        this.showAfter = true;
        return;
      }
      this.showAfter = false;
    },
    setDriver: function() {
      if(this.getEngineHourLogLatest != null && this.getEngineHourLogLatest.after_hours == null){
        this.form.driverName = this.getEngineHourLogLatest.user_first_name + " " + this.getEngineHourLogLatest.user_last_name;
        this.form.driverId = this.getEngineHourLogLatest.user_id;
      }else if(this.userInfo != null){
        this.form.driverName = this.userInfo.first_name + " " + this.userInfo.last_name;
        this.form.driverId = this.userInfo.id;
      }else{
        this.form.driverName = "unknown";
        this.form.driverId = null;
      }
    },
    toggleType: function() {
      this.form.type = !this.form.type;
    },
    add: function() {
      const type = (this.form.type) ? 1 : 0;

      const data = {
        user_id: this.form.driverId,
        engine_hours_before: this.form.beforeHours,
        engine_hours_after: this.form.afterHours,
        type: type
      };
      this.addEngineHours(data);
    }
  },
  created() {
    this.queryEngineHourLogLatest();
    this.setDisableBefore();
    this.setShowAfter();
  }
})
</script>