<template>
  <form @submit.stop="uploadNewLogo">
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <div class="alert alert-info">
      This logo will be displayed within the application to customize the look
      to your needs.
    </div>
    <overlay-spinner :active="isUploading">
      <!-- The existing logo -->
      <div
        class="row"
        v-if="hasExistingLogo()"
      >
        <label class="col-3 col-form-label">Current Logo</label>
        <div class="col-9">
          <img
            :src="getLogoUri"
            class="img-fluid custom-height"
            alt="The logo for the login screen"
          />
        </div>
      </div>
      <div class="row mb-3 mt-2">
        <div class="col-9 offset-3">
          <button
            v-if="logoState == States.HASLOGO"
            class="btn btn-outline-info me-1"
            @click="showReplaceLogo"
            type="button"
          >
            <i class="bi bi-arrow-down-up" />
            Replace
          </button>
          <button
            v-if="logoState == States.HASLOGO || logoState == States.REPLACE"
            class="btn btn-outline-danger"
            @click="removeLogo"
            type="button"
          >
            <i class="bi bi-trash" />
            Remove
          </button>
        </div>
      </div>
      <!-- The replacement logo -->
      <div
        class="row"
        v-if="hasReplacementLogo()"
      >
        <label class="col-3 col-form-label">New Logo</label>
        <div class="col-9">
          <img
            v-if="getLogoUri != null && newLogoUri == null"
            :src="newLogoUri"
            class="img-fluid custom-height"
            alt="The logo for the login screen"
          />
          <img
            v-if="newLogoUri != null && newLogoUri != getLogoUri"
            :src="newLogoUri"
            class="img-fluid custom-height"
            alt="The new logo for the login screen"
          />
        </div>
      </div>
      <!-- File Selection -->
      <input-file
        v-if="logoState == States.NOLOGO || logoState == States.REPLACE"
        id="logo-file"
        :label="uploadLabel()"
        v-model="form.logoFile"
        size="small"
        accept="image/jpeg, image/png, image/gif"
        description="Supported types are .png or .jpg (maximum 500kB)"
      />
      <div class="row mb-3 mt-2">
        <div class="col-9 offset-3">
          <button
            v-if="hasReplacementLogo()"
            class="btn btn-outline-danger ms-1"
            @click="removeReplacement"
            type="button"
          >
            <i class="bi bi-arrow-counterclockwise" />
            Undo
          </button>
          <button
            v-if="form.logoFile != null && newLogoUri == null"
            type="button"
            class="btn btn-outline-info ms-1"
            @click.stop="uploadNewLogo"
          >
            <i class="bi bi-upload" />
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
import InputFile from "./forms/inputs/InputFile.vue";
import OverlaySpinner from "./styling/OverlaySpinner.vue";

const States = Object.freeze({
  UNKNOWN: Symbol("UNKNOWN"),
  NOLOGO: Symbol("NOLOGO"),
  HASLOGO: Symbol("HASLOGO"),
  REPLACE: Symbol("REPLACE"),
  UPLOADREADY: Symbol("UPLOADREADY"),
});

export default {
  name: "LogoUpload",
  components: {
    WarningBox,
    InputFile,
    OverlaySpinner,
  },
  data() {
    return {
      States,
      errors: [],
      isUploading: false,
      form: {
        logoFile: null,
      },
      newLogoUri: null,
      logoRemoved: false,
      logoState: States.UNKNOWN,
    };
  },
  computed: {
    ...mapGetters("configuration", ["getLogoUri"]),
  },
  methods: {
    uploadLabel: function () {
      if (this.logoState == States.NOLOGO) {
        return "Upload Logo";
      } else if (this.logoState == States.REPLACE) {
        return "New Logo";
      }
      return "";
    },
    hasReplacementLogo: function () {
      if (this.newLogoUri != null) {
        return true;
      }
      return false;
    },
    hasExistingLogo: function () {
      if (this.getLogoUri != null && this.logoRemoved == false) {
        return true;
      }
      return false;
    },
    uploadNewLogo: function () {
      this.isUploading = true;

      uploadLogo(this.form.logoFile)
        .then((data) => {
          this.newLogoUri = data.uri;
          const filename = data.filename;
          this.isUploading = false;
          console.log("Logo uploaded successfully: ", filename);
          this.$emit("logoChanged", filename);
          this.errors = [];
          if (this.logoState == States.REPLACE) {
            this.logoState = States.REPLACEMENTREADY;
          } else if (this.logoState == States.NOLOGO) {
            this.logoState = States.UPLOADREADY;
          }
        })
        .catch((errors) => {
          this.errors = errors;
          this.isUploading = false;
        });
    },
    removeLogo: function () {
      this.form.logoFile = null;
      this.newLogoUri = null;
      this.$emit("logoChanged", null);
      this.logoState = States.NOLOGO;
      this.logoRemoved = true;
    },
    removeReplacement: function () {
      this.form.newLogoUri = null;
      this.form.logoFile = null;
      this.newLogoUri = null;
      this.$emit("logoChanged", null);
      if (this.logoState == States.REPLACEMENTREADY) {
        this.logoState = States.HASLOGO;
      } else {
        this.logoState = States.NOLOGO;
      }
    },
    showReplaceLogo: function () {
      this.form.logoFile = null;
      this.newLogoUri = null;
      this.logoState = States.REPLACE;
    },
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    this.queryConfiguration()
      .then(() => {
        if (this.getLogoUri == null) {
          this.logoState = States.NOLOGO;
        } else {
          this.logoState = States.HASLOGO;
        }
      })
      .catch((errors) => {
        console.error("failed to load configuration", errors);
        logoState = States.UNKNOWN;
      });
  },
};
</script>

<style scoped>
.custom-height {
  max-height: 50px;
}
</style>
