<template>
  <form @submit.stop="uploadNewLogo">
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <div class="alert alert-info">
      This logo will be displayed within the application to customize the look
      to your needs.
    </div>
    <overlay-spinner :active="isUploading">
      <!-- The existing logo -->
      <div class="row" v-if="hasExistingLogo()">
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
      <div class="row" v-if="hasReplacementLogo()">
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

<script setup>
import { ref, computed } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import { uploadLogo } from "booksys/api/resources";
import InputFile from "./forms/inputs/InputFile.vue";
import OverlaySpinner from "./styling/OverlaySpinner.vue";

const States = Object.freeze({
  UNKNOWN: Symbol("UNKNOWN"),
  NOLOGO: Symbol("NOLOGO"),
  HASLOGO: Symbol("HASLOGO"),
  REPLACE: Symbol("REPLACE"),
  UPLOADREADY: Symbol("UPLOADREADY"),
});

const store = useStore();
const emit = defineEmits(["logoChanged"]);

const errors = ref([]);
const isUploading = ref(false);
const form = ref({ logoFile: null });
const newLogoUri = ref(null);
const logoRemoved = ref(false);
const logoState = ref(States.UNKNOWN);

const getLogoUri = computed(() => store.getters["configuration/getLogoUri"]);

function uploadLabel() {
  if (logoState.value === States.NOLOGO) {
    return "Upload Logo";
  } else if (logoState.value === States.REPLACE) {
    return "New Logo";
  }
  return "";
}

function hasReplacementLogo() {
  return newLogoUri.value != null;
}

function hasExistingLogo() {
  return getLogoUri.value != null && !logoRemoved.value;
}

function uploadNewLogo() {
  isUploading.value = true;
  uploadLogo(form.value.logoFile)
    .then((data) => {
      newLogoUri.value = data.uri;
      const filename = data.filename;
      isUploading.value = false;
      console.log("Logo uploaded successfully: ", filename);
      emit("logoChanged", filename);
      errors.value = [];
      if (logoState.value === States.REPLACE) {
        logoState.value = States.REPLACEMENTREADY;
      } else if (logoState.value === States.NOLOGO) {
        logoState.value = States.UPLOADREADY;
      }
    })
    .catch((errs) => {
      errors.value = errs;
      isUploading.value = false;
    });
}

function removeLogo() {
  form.value.logoFile = null;
  newLogoUri.value = null;
  emit("logoChanged", null);
  logoState.value = States.NOLOGO;
  logoRemoved.value = true;
}

function removeReplacement() {
  form.value.newLogoUri = null;
  form.value.logoFile = null;
  newLogoUri.value = null;
  emit("logoChanged", null);
  if (logoState.value === States.REPLACEMENTREADY) {
    logoState.value = States.HASLOGO;
  } else {
    logoState.value = States.NOLOGO;
  }
}

function showReplaceLogo() {
  form.value.logoFile = null;
  newLogoUri.value = null;
  logoState.value = States.REPLACE;
}

store
  .dispatch("configuration/queryAdminConfiguration")
  .then(() => {
    if (getLogoUri.value == null) {
      logoState.value = States.NOLOGO;
    } else {
      logoState.value = States.HASLOGO;
    }
  })
  .catch((errs) => {
    console.error("failed to load configuration", errs);
    logoState.value = States.UNKNOWN;
  });
</script>

<style scoped>
.custom-height {
  max-height: 50px;
}
</style>
