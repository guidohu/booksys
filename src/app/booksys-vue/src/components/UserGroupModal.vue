<template>
  <b-modal
    id="userGroupModal"
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
          <!-- User Group -->
          <h6 class="mt-2 mb-3">User Group</h6>
          <b-form-group
            id="user-group-name"
            label="Name"
            label-for="user-group-name"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="user-group-name"
              v-model="form.user_group_name"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="user-group-description"
            label="Description"
            label-for="user-group-description"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="user-group-description"
              v-model="form.user_group_description"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <!-- User Role -->
          <h6 class="mt-5 mb-3">User Role and Permissions</h6>
          <b-form-group
            id="user-role-name"
            label="Role"
            label-for="user-role-name"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="user-role-name"
              v-model="form.user_role_name"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <b-form-group
            id="user-role-description"
            label="Description"
            label-for="user-role-description"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="user-role-description"
              v-model="form.user_role_description"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
          <!-- Pricing -->
          <h6 class="mt-5 mb-3">Pricing</h6>
          <b-form-group
            id="price"
            label="Price"
            label-for="price-input"
            description=""
            label-cols="3"
          >
            <b-input-group>
              <b-form-input
                id="price-input"
                v-model="form.price_min"
                type="text"
                placeholder=""
                disabled
              ></b-form-input>
              <b-input-group-append is-text>
                {{getCurrency}}/min
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            id="price-description"
            label="Description"
            label-for="price-description-input"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="price-description-input"
              v-model="form.price_description"
              type="text"
              placeholder=""
              disabled
            ></b-form-input>
          </b-form-group>
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <div class="text-left d-inline">
        <b-button v-if="isEditMode==true" class="mr-1" type="button" variant="outline-danger" v-on:click="remove">
          <b-icon-trash></b-icon-trash>
          Delete
        </b-button>
      </div>
      <div class="text-right d-inline">
        <b-button v-if="isEditMode==true" class="ml-4" type="button" variant="outline-info" v-on:click="save">
          <b-icon-check></b-icon-check>
          Save
        </b-button>
        <b-button v-if="isEditMode==false" class="ml-1" type="button" variant="outline-info" v-on:click="enableEditMode">
          <b-icon icon="pencil"/>
          Edit
        </b-button>
        <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
          <b-icon-x></b-icon-x>
          Cancel
        </b-button>
      </div>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import { mapGetters, mapActions } from 'vuex';

export default Vue.extend({
  name: 'UserGroupModal',
  props: [ 'userGroup', 'editMode' ],
  data() {
    return {
      errors: [],
      form: {},
      isEditMode: false,
      title: "User Group"
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getCurrency'
    ])
  },
  watch: {
    userGroup: function(newValue) {
      if(newValue != null){
        this.form = newValue;
      }else{
        this.form = {};
      }
    },
    editMode: function(newValue) {
      this.isEditMode = newValue;
      this.setTitle();
    }
  },
  methods: {
    save: function() {
      console.log("TODO implement save");
    },
    close: function() {
      this.isEditMode = false;
      this.setTitle();
      this.$bvModal.hide('userGroupModal');
    },
    remove: function() {
      console.log("TODO implement remove");
    },
    ...mapActions('configuration', [
      'queryConfiguration'
    ]),
    enableEditMode: function() {
      console.log("Enable edit mode");
      this.isEditMode = true;
      this.setTitle();
    },
    setTitle: function() {
      console.log("Set title with:", this.isEditMode);
      if(this.isEditMode == true){
        this.title = "Edit User Group";
      }else{
        this.title = "User Group";
      }
    }
  },
  created() {
    this.queryConfiguration();

    if(this.userGroup != null){
      this.form = this.userGroup;
    }

    this.isEditMode = this.editMode;
    this.setTitle();
  }
})
</script>