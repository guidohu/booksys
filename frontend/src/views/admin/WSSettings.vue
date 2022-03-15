<template>
  <subpage-container title="Settings">
    <warning-box v-if="errors.length > 0" :errors="errors"/>
    <show-for-desktop>
      <settings-card 
        class="mx-1"
        :show-controls="false"
        @save="save"
        @change="change"
      />
    </show-for-desktop>
    <show-for-mobile>
      <settings-card 
        class="mx-1" 
        :show-controls="true"
        @saved="navigateBack"
        @cancelled="navigateBack"
      />
    </show-for-mobile>
    <template v-slot:bottom>
      <button
        tag="button"
        class="btn btn-outline-light"
        @click="save"
      >
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

<script>
import { mapActions } from "vuex";
import SettingsCard from "@/components/SettingsCard";
import SubpageContainer from "@/components/bricks/SubpageContainer";
import ShowForDesktop from "@/components/bricks/ShowForDesktop.vue";
import ShowForMobile from "@/components/bricks/ShowForMobile.vue";
import WarningBox from '../../components/WarningBox.vue';

export default {
  name: "WSSettings",
  components: {
    SettingsCard,
    SubpageContainer,
    ShowForDesktop,
    ShowForMobile,
    WarningBox,
  },
  data() {
    return {
      settings: null,
      errors: [],
    };
  },
  methods: {
    ...mapActions("configuration", ["setConfiguration"]),
    change: function(settings) {
      this.settings = settings;
    },
    save: function() {
      if(this.settings == null) {
        this.navigateBack();
        return;
      }

      this.setConfiguration(this.settings)
        .then(() => {
          this.errors = [];
          this.navigateBack();
        })
        .catch((errors) => (this.errors = errors));
    },
    cancel: function() {
      this.navigateBack();
    },
    navigateBack: function() {
      this.$router.push("/admin");
    }
  }
};
</script>
