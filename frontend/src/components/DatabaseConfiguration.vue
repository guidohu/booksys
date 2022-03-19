<template>
  <div class="container">
    <div class="row">
      <div class="col-12">
        <p>
          Setup the connection to the database. Please make sure that the
          database server is up, the database exists and there is a user with
          sufficient permissions to setup the tables.
        </p>
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <form @submit.stop="save">
          <input-text
            id="host"
            label="Host"
            size="small"
            description="Database Host (e.g. 127.0.0.1:3306)"
            v-model="conf.host"
          />
          <input-text
            id="name"
            label="Name"
            size="small"
            description="Database Name (e.g. booksys)"
            v-model="conf.name"
          />
          <input-text
            id="user"
            label="User"
            size="small"
            description="Database User (e.g. dbUser)"
            v-model="conf.user"
          />
          <input-password
            id="password"
            label="Password"
            size="small"
            description="Password for Database User"
            v-model="conf.password"
          />
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import InputText from "./forms/inputs/InputText.vue";
import InputPassword from "./forms/inputs/InputPassword.vue";

export default {
  name: "DatabaseConfiguration",
  components: {
    InputText,
    InputPassword,
  },
  emits: ["save", "config-change"],
  props: ["dbconfig"],
  data() {
    return {
      conf: {
        host: "",
        name: "",
        user: "",
        password: "",
      },
    };
  },
  watch: {
    conf: {
      handler(){
        this.$emit("config-change", this.conf);
      },
      deep: true,
    }
  },
  created() {
    this.conf.host = this.dbconfig.host;
    this.conf.name = this.dbconfig.name;
    this.conf.user = this.dbconfig.user;
    this.conf.password = this.dbconfig.password;
  },
  methods: {
    save: function () {
      this.$emit("save", this.conf);
    },
  },
};
</script>
