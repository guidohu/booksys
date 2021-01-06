<template>
<b-modal
    id="riderSelectionModal"
    title="Add Riders"
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
    <b-form @submit="add">
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <b-form-group
            v-if="isMobile != true"
            id="search"
            label="Search"
            label-for="rider-search"
            description=""
          >
            <b-form-input
              id="rider-search"
              v-model="search"
              type="text"
            />
          </b-form-group>
          <b-form-group
            id="rider"
            label="Rider"
            label-for="rider-select"
            description="use ctrl+click to select multiple users"
          >
            <b-form-select
              id="rider-select"
              v-model="selected"
              :options="filteredOptions"
              multiple 
              :select-size="5"
            ></b-form-select>
          </b-form-group>
          <b-button v-if="isMobile != true" type="button" variant="outline-success" v-on:click="add" block>
            <b-icon-person-plus/>{{" "}}Select
          </b-button>
          <hr/>
          <b-form-group
            v-if="isMobile != true"
            id="selected"
            label="Selected Riders"
            label-for="riders-selected"
            description=""
          >
            <ul v-if="usersToAdd.length > 0">
              <li v-for="u in usersToAdd" :key="u.id">
                {{ u.firstName + " " + u.lastName + " " }}
                <a href="#" v-on:click="remove(u.id)"><b-icon icon="person-dash"></b-icon></a>
              </li>
            </ul>
            <ul v-else>
              <li>Please select riders to be added</li>
            </ul>
          </b-form-group>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <b-button v-if="isMobile != true" type="button" variant="outline-info" v-on:click="save">
        <b-icon-check/>
        Add
      </b-button>
      <b-button v-if="isMobile == true" type="button" variant="outline-info" v-on:click="saveMobile">
        <b-icon-check/>
        Add
      </b-button>
      <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x/>
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import { BooksysBrowser } from '@/libs/browser';
import { uniq } from 'lodash';
import WarningBox from '@/components/WarningBox';
import { UserPointer } from '@/dataTypes/user';
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BFormSelect,
  BButton,
  BIconPersonPlus,
  BIconCheck,
  BIconX
} from 'bootstrap-vue';

export default {
  name: 'RiderSelectionModal',
  props: [
    'session',
    'visible'
  ],
  components: {
    WarningBox,
      BModal,
      BRow,
      BCol,
      BForm,
      BFormGroup,
      BFormInput,
      BFormSelect,
      BButton,
      BIconPersonPlus,
      BIconCheck,
      BIconX
  },
  data() {
    return {
      errors: [],
      userIds: [],
      form: {},
      selected: [],
      search: "",
      filteredOptions: [],
      usersToAdd: []
    }
  },
  computed: {
    ...mapGetters('user', [
      'userList'
    ]),
    userOptions: function(){
      let users = [];
      this.userList.forEach(u => {
        users.push({
          value: u.id,
          text: u.firstName + " " + u.lastName
        })
      });

      users = users.filter(u => !this.usersToAdd.map(uta => uta.id).includes(u.value))
      return users;
    },
    isMobile: function () {
      return BooksysBrowser.isMobile()
    }
  },
  watch: {
    search: function(newSearch){
      // filteredOptions are the ones that 
      this.filteredOptions = this.userOptions
        .filter(u => u.text.toLowerCase().includes(newSearch.toLowerCase()))
    },
    userOptions: function(newOptions){
      this.filteredOptions = newOptions;
    },
    selected: function(newSelected){
      console.log("newSelected", newSelected);
    },
    usersToAdd: function(newUsersToAdd){
      console.log("usersToAdd", newUsersToAdd);
    }
  },
  methods: {
    ...mapActions('user', [
      'queryUserList'
    ]),
    ...mapActions('sessions', [
      'addUsersToSession'
    ]),
    add: function(event) {
      event.preventDefault();

      console.log("add Called");
      console.log("selected", this.selected);

      if(this.selected.length > 0){
        console.log("Add selected");
        this.usersToAdd.push(...this.userList.filter(u => this.selected.includes(u.id)))
      }else if(this.selected.length == 0 && this.filteredOptions.length == 1){
        this.usersToAdd.push(...this.userList.filter(u => u.id == this.filteredOptions[0].value))
      }

      // de-duplicate selection
      this.usersToAdd = uniq(this.usersToAdd);
    },
    remove: function(id){
      console.log(id);
      this.usersToAdd = this.usersToAdd.filter(u => u.id != id);
    },
    save: function() {
      console.log(this.session);
      this.addUsersToSession({
        sessionId: this.session.id,
        users: this.usersToAdd
      })
      .then(() => this.close())
      .catch((errs) => this.errors = errs);
    },
    saveMobile: function() {
      console.log(this.selected);
      this.addUsersToSession({
        sessionId: this.session.id,
        users: this.selected.map(i => new UserPointer(Number(i)))
      })
      .then(() => this.close())
      .catch((errors) => this.errors = errors);
    },
    close: function() {
      this.usersToAdd = [];
      this.$emit('update:visible', false);
    }
  },
  created() {
    this.queryUserList()
    .then( () => {
      console.log("userList received");
    })
    .catch( (error) => {
      console.log("errors")
      this.errors = error;
    });
  }
}
</script>