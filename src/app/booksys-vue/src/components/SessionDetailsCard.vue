<template>
  <b-card no-body class="text-left">
    <b-card-header>
      <b-row>
        <b-col cols="8">
          Session Details
        </b-col>
        <b-col cols="4" class="text-right">
          <b-dropdown v-if="showAddRiders || showDeleteSession" id="dropdown-1" size="sm" variant="light" text="Dropdown Button" class="btn-xs" right no-caret>
            <template v-slot:button-content>
              <b-icon icon="list"></b-icon>
            </template>
            <b-dropdown-item v-if="showAddRiders" v-on:click="addRiders">
              <b-icon icon="person-plus"></b-icon>
              {{"  "}} Add Rider
            </b-dropdown-item>
            <b-dropdown-item v-if="showAddRiders" v-on:click="editSession">
              <b-icon icon="pencil-square"></b-icon>
              {{"  "}} Edit Session
            </b-dropdown-item>
            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-item v-if="showDeleteSession" v-on:click="deleteSession">
              <b-icon icon="trash"></b-icon>
              {{"  "}} Delete Session
            </b-dropdown-item>
          </b-dropdown>
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body>
      <RiderSelectionModal :session="session"/>
      <b-row>
        <b-col cols="4">
          Date
        </b-col>
        <b-col cols="8">
          {{dateString}}
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="4">
          Session
        </b-col>
        <b-col cols="8">
          {{timeString}}
        </b-col>
      </b-row>
      <b-row v-if="showRiders">
        <b-col cols="4">
          Riders
        </b-col>
      </b-row>
      <b-row>
        <ul v-if="showRiders">
          <li v-for="rider in session.riders" :key="rider.id">
            {{ rider.name }} {{ " "}}
            <a href="#" v-on:click="removeRider(rider.id)"><b-icon icon="person-dash"></b-icon></a>
          </li>
        </ul>
      </b-row>
      <b-row v-if="showCreateSession">
        <b-col offset="4" cols="8">
          <b-button v-on:click="createSession" size="sm" variant="success" block>Create Session</b-button>
        </b-col>
      </b-row>
      <b-row v-if="showAddRiders">
        <b-col offset="4" cols="8">
          <b-button v-on:click="addRiders" size="sm" variant="success" block>Add Riders</b-button>
        </b-col>
      </b-row>
    </b-card-body>
  </b-card>
</template>

<script>
import Vue from 'vue';
import moment from 'moment';
import 'moment-timezone';
import RiderSelectionModal from '@/components/RiderSelectionModal.vue';

export default Vue.extend({
  name: 'SessionDetailsCard',
  components: {
    RiderSelectionModal
  },
  props: [
    "date",
    "session"
  ],
  data() {
    return {
      sessionSelected: false
    };
  },
  computed: {
    dateString: function(){
      const newDate = this.date;
      return moment(newDate).format('DD.MM.YYYY');
    },
    timeString: function(){
      if(this.session == null){
        return "no session selected";
      }

      const start = this.session.start;
      const end   = this.session.end;
      return moment(start).format('HH:mm') + " - " + moment(end).format('HH:mm');
    },
    showCreateSession: function(){
      if(this.session != null && this.session.id == null){
        return true;
      }
      return false;
    },
    showDeleteSession: function(){
      if(this.session != null && this.session.id != null){
        return true;
      }
      return false;
    },
    showAddRiders: function(){
      if(this.session != null && this.session.id != null){
        return true;
      }
      return false;
    },
    showRiders: function(){
      if(this.session != null && this.session.riders != null && this.session.riders.length > 0){
        return true;
      }
      return false;
    }
  },
  watch: {
    timeString: function(){
      if(this.session == null){
        this.sessionSelected = false;
      }else{
        this.sessionSelected = true;
      }
    }
  },
  methods: {
    createSession: function(){
      this.$emit('createSessionHandler');
    },
    editSession: function(){
      this.$emit('editSessionHandler');
    },
    deleteSession: function(){
      this.$emit('deleteSessionHandler', { id: this.session.id });
    },
    addRiders: function(){
      this.$bvModal.show('riderSelectionModal');
    },
    removeRider: function(id){
      console.log("TODO remove rider with id", id);
    }
  }
})
</script>