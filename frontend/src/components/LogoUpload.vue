<template>
  <b-form @submit="uploadLogo">
    <b-alert variant="info" show>
      This logo will be displayed in the application.
    </b-alert>
    <!-- Image Upload -->
    <b-form-group
      id="input-group-logo-file"
      label="Your Logo"
      label-for="input-logo-file"
      label-cols="3"
      description="Supported types are .png or .jpg (maximum 500kB)"
    >
      <b-input-group size="sm">
        <b-form-file
          id="input-timezone"
          v-model="form.logoFile"
          placeholder="Choose a logo or drop it here"
          drop-placeholder="Drop logo here"
          accept="image/jpeg, image/png, image/gif"
        />
      </b-input-group>
    </b-form-group>
    <b-row class="text-right">
      <b-col offset="3" cols="9">
        <b-button class="mr-1">Upload</b-button>
        <b-button>Clear</b-button>
      </b-col>
    </b-row>
  </b-form>
</template>

<script>
import {
  BForm,
  BAlert,
  BFormGroup,
  BInputGroup,
  BFormFile,
  BButton,
  BRow,
  BCol
} from 'bootstrap-vue';

export default {
  name: "LogoUpload",
  components: {
    BForm,
    BAlert,
    BFormGroup,
    BInputGroup,
    BFormFile,
    BButton,
    BRow,
    BCol
  },
  data() {
    return {
      form: {}
    }
  },
  methods: {
    uploadLogo: function(){

    }
  }
}
</script>