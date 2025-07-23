<template>
  <subpage-container title="Settings">
    <warning-box v-if="errors.length > 0" :errors="errors" />
    <show-for-desktop>
      <SettingsCard
        class="mx-1"
        :show-controls="false"
        @save="save"
        @change="change"
      />
    </show-for-desktop>
    <show-for-mobile>
      <SettingsCard
        class="mx-1"
        :show-controls="true"
        @saved="navigateBack"
        @cancelled="navigateBack"
      />
    </show-for-mobile>
    <template v-slot:bottom>
      <button tag="button" class="btn btn-outline-light" @click="save">
        <i class="bi bi-check"></i>
        Save
      </button>
      <button
        tag="button"
        class="btn btn-outline-light ms-1"
        @click="navigateBack"
      >
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </template>
  </subpage-container>
</template>

<script setup>
import { ref } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import SettingsCard from "@/components/SettingsCard";
import SubpageContainer from "@/components/bricks/SubpageContainer";
import ShowForDesktop from "@/components/bricks/ShowForDesktop.vue";
import ShowForMobile from "@/components/bricks/ShowForMobile.vue";
import WarningBox from "../../components/WarningBox.vue";

const store = useStore();
const router = useRouter();

const settings = ref(null);
const errors = ref([]);

const change = (newSettings) => {
  settings.value = newSettings;
};

const save = () => {
  if (settings.value == null) {
    navigateBack();
    return;
  }

  store
    .dispatch("configuration/setConfiguration", settings.value)
    .then(() => {
      errors.value = [];
      navigateBack();
    })
    .catch((err) => (errors.value = err));
};

const navigateBack = () => {
  router.push("/admin");
};
</script>
