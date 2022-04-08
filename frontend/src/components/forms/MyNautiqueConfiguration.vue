<template>
  <div class="container">
    <div class="row">
      <div class="col-12">
        <p>
          This application is compatible with myNautique, which allows to
          automatically get data from your Nautique boat via myNautique App API
          to simplify many processes in this App.
        </p>
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <form @submit.stop="save">
          <input-toggle
            id="mynautique-toggle"
            label="myNautique"
            offLabel="Do not connect to myNautique app"
            onLabel="myNautique support enabled"
            v-model="form.enabled"
            @input="update()"
          />
          <input-text
            v-if="form.enabled"
            id="user"
            label="User"
            size="small"
            description="myNautique user"
            v-model="form.user"
            @input="update()"
          />
          <input-password
            v-if="form.enabled"
            id="password"
            label="Password"
            size="small"
            description="myNautique password"
            v-model="form.password"
            @input="update()"
          />
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import InputText from "./inputs/InputText.vue";
import InputPassword from "./inputs/InputPassword.vue";
import InputToggle from "./inputs/InputToggle.vue";

export default {
  name: "MyNautiqueConfiguration",
  components: {
    InputText,
    InputPassword,
    InputToggle,
  },
  emits: ["save", "update:settings"],
  props: ["settingsData"],
  data() {
    return {
      form: {
        enabled: false,
        user: "",
        password: "",
      },
    };
  },
  methods: {
    save: function () {
      this.$emit("save", this.form);
    },
    update: function () {
      this.$emit("update:settings", this.form);
    },
  },
  created() {
    this.form = this.$props.settingsData;
  },
};
</script>
