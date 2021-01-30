<template>
  <b-card no-body class="text-left">
    <b-card-header>
      <b-row>
        <b-col cols="8" class="align-middle">
          Session Details
        </b-col>
        <b-col cols="4" class="text-right">
          <b-dropdown v-if="showAddRiders || showDeleteSession" id="dropdown-1" variant="light" text="Dropdown Button" toggle-class="btn-xs" right no-caret>
            <template v-slot:button-content>
              <b-icon-list/>
            </template>
            <b-dropdown-item v-if="showAddRiders" v-on:click="addRiders">
              <b-icon-person-plus/>
              {{"  "}} Add Rider
            </b-dropdown-item>
            <b-dropdown-item v-if="showAddRiders" v-on:click="editSession">
              <b-icon-pencil-square/>
              {{"  "}} Edit Session
            </b-dropdown-item>
            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-item v-if="showDeleteSession" v-on:click="deleteSession">
              <b-icon-trash/>
              {{"  "}} Delete Session
            </b-dropdown-item>
          </b-dropdown>
        </b-col>
      </b-row>
    </b-card-header>
    <b-card-body>
      <RiderSelectionModal 
        :session="session"
        :visible.sync="showRiderSelectionModal"
      />
      <b-row v-if="session != null && session.title != null" class="font-weight-bold">
        <b-col cols="4">
          Title
        </b-col>
        <b-col cols="8">
          {{session.title}}
        </b-col>
      </b-row>
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
      <b-row v-if="showRiders" class="mt-2">
        <b-col cols="4">
          Riders
        </b-col>
        <b-col cols="8">
          <b-row v-for="rider in session.riders" :key="rider.id">
            <b-col cols="10" class="text-truncate">
              {{ rider.name }}
            </b-col>
            <b-col cols="2">
              <a href="#" v-on:click="removeRider(rider.id)">
                <b-icon-person-dash/>
              </a>
            </b-col>
          </b-row>
        </b-col>
      </b-row>
      <b-row v-if="showCreateSession" class="mt-3">
        <b-col offset="4" cols="8">
          <b-button v-on:click="createSession" size="sm" variant="success" block>Create Session</b-button>
        </b-col>
      </b-row>
      <b-row v-if="showAddRiders" class="mt-3">
        <b-col offset="4" cols="8">
          <b-button v-on:click="addRiders" size="sm" variant="success" block>Add Riders</b-button>
        </b-col>
      </b-row>
    </b-card-body>
  </b-card>
</template>

<script>
import { mapActions} from 'vuex';
import * as dayjs from 'dayjs';
import { BooksysBrowser } from '@/libs/browser';
import { UserPointer } from '@/dataTypes/user';
import RiderSelectionModal from '@/components/RiderSelectionModal';
import {
  BCard,
  BCardHeader,
  BCardBody,
  BRow,
  BCol,
  BDropdown,
  BDropdownItem,
  BDropdownDivider,
  BIconList,
  BIconPersonPlus,
  BIconPersonDash,
  BIconPencilSquare,
  BIconTrash,
  BButton
} from 'bootstrap-vue';

export default {
  name: 'SessionDetailsCard',
  components: {
    RiderSelectionModal,
    BCard,
    BCardHeader,
    BCardBody,
    BRow,
    BCol,
    BDropdown,
    BDropdownItem,
    BDropdownDivider,
    BIconList,
    BIconPersonPlus,
    BIconPersonDash,
    BIconPencilSquare,
    BIconTrash,
    BButton
  },
  props: [
    "date",
    "session"
  ],
  data() {
    return {
      showRiderSelectionModal: false,
      sessionSelected: false
    };
  },
  computed: {
    dateString: function(){
      const newDate = this.date;
      return dayjs(newDate).format('DD.MM.YYYY');
    },
    timeString: function(){
      if(this.session == null){
        return "no session selected";
      }

      const start = this.session.start;
      const end   = this.session.end;
      return dayjs(start).format('HH:mm') + " - " + dayjs(end).format('HH:mm');
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
    },
    isMobile: function () {
      return BooksysBrowser.isMobile();
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
      this.showRiderSelectionModal = true;
    },
    ...mapActions('sessions', [
      "deleteUserFromSession"
    ]),
    removeRider: function(id){
      console.log("TODO remove rider with id", id);
      this.deleteUserFromSession({
        sessionId: this.session.id,
        user: new UserPointer(id)
      })
    }
  }
}
</script>

<style scoped>
  .btn-xs {
    padding: .2rem .4rem;
    font-size: .6rem;
    line-height: 1;
    border-radius: .2rem;
  }
</style>