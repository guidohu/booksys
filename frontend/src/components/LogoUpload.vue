<template>
  <b-form @submit="uploadNewLogo">
    <b-alert variant="info" show>
      This logo will be displayed within the application to customize the look to your needs.
    </b-alert>
    <warning-box v-if="errors.length>0" :errors="errors"/>
    <!-- Image Upload -->
    <b-overlay
      id="overlay-background"
      :show="isUploading"
      spinner-type="border"
      spinner-variant="info"
      rounded="sm"
    >
      <b-form-group
        v-if="newLogoUri == null && (getLogoFile == null || showLogo == false)"
        id="input-group-logo-file"
        label="Your Logo"
        label-for="input-logo-file"
        label-cols="3"
        description="Supported types are .png or .jpg (maximum 500kB)"
      >
        <b-input-group size="sm">
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
        v-if="newLogoUri != null || (getLogoFile != null && showLogo == true)"
        label="Your Logo"
        label-cols="3"
      >
        <b-img v-if="getLogoFile != null && newLogoUri == null" :height="50" :src="getLogoFile"/>
        <b-img v-if="newLogoUri != null" :height="50" :src="newLogoUri"/>
      </b-form-group>
      <b-row class="text-right mb-3">
        <b-col offset="3" cols="9">
          <b-button 
            v-if="newLogoUri != null || form.logoFile != null || (getLogoFile != null && showLogo == true)"
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
    </b-overlay>
  </b-form>
</template>

<script>
import {
  BOverlay,
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
import WarningBox from '@/components/WarningBox';
import { uploadLogo } from '@/api/resources';
import { mapGetters, mapActions } from 'vuex';

export default {
  name: "LogoUpload",
  components: {
    WarningBox,
    BOverlay,
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
      errors: [],
      isUploading: false,
      form: {
        logoFile: null
      },
      newLogoUri: null,
      showLogo: true,
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getLogoFile'
    ])
  },
  methods: {
    uploadNewLogo: function(){
      this.isUploading = true;

      uploadLogo(this.form.logoFile)
      .then((uri) => {
        this.newLogoUri = uri;
        this.showLogo = true;
        this.isUploading = false;
        this.$emit('logoChanged', uri);
        this.errors = [];
      })
      .catch((errors) => {
        this.errors = errors;
        this.isUploading = false;
      });
    },
    clearLogo: function(){
      this.form.logoFile = null;
      this.newLogoUri = null;
      this.showLogo = false;
      this.$emit('logoChanged', null);
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