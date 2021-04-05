<template>
  <b-container>
    <b-row>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <p>
          Setup the connection to the database. Please make sure that the
          database server is up, the database exists and there is a user with
          sufficient permissions to setup the tables.
        </p>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-row>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <b-form v-on:submit.prevent="save">
          <b-form-group
            id="host"
            label="Host"
            label-for="host-input"
            description="Database Host (e.g. 127.0.0.1:3306)"
            label-cols="3"
          >
            <b-form-input id="host-input" v-model="dbConfig.host" type="text" />
          </b-form-group>
          <b-form-group
            id="name"
            label="Name"
            label-for="name-input"
            description="Database Name (e.g. booksys)"
            label-cols="3"
          >
            <b-form-input id="name-input" v-model="dbConfig.name" type="text" />
          </b-form-group>
          <b-form-group
            id="user"
            label="User"
            label-for="user-input"
            description="Database User (e.g. dbUser)"
            label-cols="3"
          >
            <b-form-input id="user-input" v-model="dbConfig.user" type="text" />
          </b-form-group>
          <b-form-group
            id="password"
            label="Password"
            label-for="password-input"
            description="Password for the user."
            label-cols="3"
          >
            <b-form-input
              id="password-input"
              v-model="dbConfig.password"
              type="password"
            />
          </b-form-group>
        </b-form>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
  </b-container>
</template>

<script>
import {
  BContainer,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
} from "bootstrap-vue";

export default {
  name: "DatabaseConfiguration",
  components: {
    BContainer,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
  },
  props: ["dbConfig"],
  methods: {
    save: function () {
      this.$emit("save", this.dbConfig);
    },
  },
};
</script>
