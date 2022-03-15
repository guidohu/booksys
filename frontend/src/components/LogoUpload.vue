<template>
  <form @submit.stop="uploadNewLogo">
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <div class="alert alert-info">
      This logo will be displayed within the application to customize the look
      to your needs.
    </div>
    <overlay-spinner
      :active="isUploading"
    >
      <input-file
        v-if="newLogoUri == null && (getLogoFile == null || showLogo == false)"
        id="logo-file"
        label="Your Logo"
        v-model="form.logoFile"
        size="small"
        accept="image/jpeg, image/png, image/gif"
        description="Supported types are .png or .jpg (maximum 500kB)"
      />
      <div 
        class="row" 
        v-if="newLogoUri != null || (getLogoFile != null && showLogo == true)"
      > 
        <label class="col-3 col-form-label">Your Logo</label>
        <div class="col-9">
          <img 
            v-if="getLogoFile != null && newLogoUri == null" 
            :src="getLogoFile" 
            class="img-fluid custom-height" 
            alt="The logo for the login screen"
          />
          <img 
            v-if="newLogoUri != null && newLogoUri != getLogoFile" 
            :src="newLogoUri" 
            class="img-fluid custom-height" 
            alt="The logo for the login screen"
          />
        </div>
      </div>
      <div class="row mb-3 mt-2">
        <div class="col-9 offset-3">
          <button
            v-if="
              newLogoUri != null ||
              form.logoFile != null ||
              (getLogoFile != null && showLogo == true)
            "
            class="btn btn-outline-danger"
            @click="clearLogo"
            type="button"
          >
            <i class="bi bi-trash"/>
            Remove
          </button>
          <button
            v-if="form.logoFile != null && newLogoUri == null"
            type="button"
            class="btn btn-outline-info ms-1"
            @click.stop="uploadNewLogo"
          >
            <i class="bi bi-upload"/>
            Upload
          </button>
        </div>
      </div>
    </overlay-spinner>
  </form>
</template>

<script>
import WarningBox from "@/components/WarningBox";
import { uploadLogo } from "@/api/resources";
import { mapGetters, mapActions } from "vuex";
import InputFile from './forms/inputs/InputFile.vue';
import OverlaySpinner from './styling/OverlaySpinner.vue';

export default {
  name: "LogoUpload",
  components: {
    WarningBox,
    InputFile,
    OverlaySpinner
  },
  data() {
    return {
      errors: [],
      isUploading: false,
      form: {
        logoFile: null,
      },
      newLogoUri: null,
      showLogo: true,
    };
  },
  computed: {
    ...mapGetters("configuration", ["getLogoFile"]),
  },
  methods: {
    uploadNewLogo: function () {
      this.isUploading = true;

      uploadLogo(this.form.logoFile)
        .then((uri) => {
          this.newLogoUri = uri;
          this.showLogo = true;
          this.isUploading = false;
          this.$emit("logoChanged", uri);
          this.errors = [];
        })
        .catch((errors) => {
          this.errors = errors;
          this.isUploading = false;
        });
    },
    clearLogo: function () {
      this.form.logoFile = null;
      this.newLogoUri = null;
      this.showLogo = false;
      this.$emit("logoChanged", null);
    },
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.queryConfiguration();
  },
};
</script>

<style scoped>
  .custom-height {
    max-height: 50px;
  }
</style>
