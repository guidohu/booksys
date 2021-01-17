<template>
  <b-form @submit="uploadNewLogo">
    <b-alert variant="info" show>
      This logo will be displayed within the application to customize the look to your needs.
    </b-alert>
    <!-- Image Upload -->
    <b-form-group
      v-if="newLogoUri == null"
      id="input-group-logo-file"
      label="Your Logo"
      label-for="input-logo-file"
      label-cols="3"
      description="Supported types are .png or .jpg (maximum 500kB)"
    >
      <b-input-group v-if="getLogoFile == null" size="sm">
        <b-form-file
          id="input-logo-file"
          v-model="form.logoFile"
          placeholder="No logo selected. Choose one or drop one here..."
          drop-placeholder="Drop logo here"
          accept="image/jpeg, image/png, image/gif"
        />
      </b-input-group>
    </b-form-group>
    <b-form-group
      v-if="newLogoUri != null"
      label="Your Logo"
      label-cols="3"
    >
      <b-img v-if="getLogoFile != null" src="getLogoFile"/>
      <b-img v-if="newLogoUri != null" :height="50" :src="newLogoUri"/>
    </b-form-group>
    <b-row class="text-right mb-3">
      <b-col offset="3" cols="9">
        <b-button 
          v-if="newLogoUri != null || form.logoFile != null"
          v-on:click="clearLogo"
          variant="outline-danger"
        >
          <b-icon-trash/>
          Clear
        </b-button>
        <b-button 
          v-if="form.logoFile != null && newLogoUri == null"
          class="ml-1" 
          v-on:click="uploadNewLogo"
          variant="outline-info"
        >
          <b-icon-upload/>
          Upload
        </b-button>
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
  BImg,
  BIconTrash,
  BIconUpload,
  BRow,
  BCol
} from 'bootstrap-vue';
import {uploadLogo} from '@/api/resources';
import { mapGetters, mapActions } from 'vuex';

export default {
  name: "LogoUpload",
  components: {
    BForm,
    BAlert,
    BFormGroup,
    BInputGroup,
    BFormFile,
    BButton,
    BImg,
    BIconTrash,
    BIconUpload,
    BRow,
    BCol
  },
  data() {
    return {
      form: {},
      logoUploaded: false,
      newLogoUri: null,
      logoProps: { blank: true, height: 75, width: 150, class: "m1", thumbnail: true }
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getLogoFile'
    ])
  },
  methods: {
    uploadNewLogo: function(){
      console.log('File to upload:', this.form.logoFile);
      uploadLogo(this.form.logoFile)
      .then((uri) => {
        this.newLogoUri = uri;
        console.log("Success")
      })
      .catch((errors) => console.error("Error", errors));
    },
    clearLogo: function(){
      this.form.logoFile = null;
      this.newLogoUri = null;
    },
    ...mapActions('configuration', [
      'queryConfiguration'
    ])
  },
  created() {
    this.queryConfiguration();
  }
}
</script>