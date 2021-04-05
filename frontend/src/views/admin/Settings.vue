<template>
  <div>
    <div v-if="isDesktop" class="display">
      <main-title title-name="Settings" />
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <SettingsCard class="content-max-height" />
        </b-col>
      </b-row>
    </div>
    <div v-else>
      <NavbarAdminMobile title="Settings" />
      <SettingsCard />
    </div>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import NavbarAdminMobile from "@/components/NavbarAdminMobile";
import SettingsCard from "@/components/SettingsCard";
import MainTitle from "@/components/MainTitle";
import { BRow, BCol } from "bootstrap-vue";

export default {
  name: "Settings",
  components: {
    NavbarAdminMobile,
    MainTitle,
    SettingsCard,
    BRow,
    BCol,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
  },
};
</script>

<style>
.content-max-height {
  max-height: 520px;
  overflow-y: scroll;
}
</style>
