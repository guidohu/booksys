<template>
  <div>
    <div
      v-if="isDesktop"
      class="display"
    >
      <main-title title-name="Boat" />
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <b-card no-body>
            <BoatTabs />
          </b-card>
        </b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button
          variant="outline-light"
          to="/dashboard"
        >
          <b-icon-house />
          HOME
        </b-button>
      </div>
    </div>
    <div v-else>
      <NavbarMobile title="Boat" />
      <b-card no-body>
        <BoatTabs />
      </b-card>
    </div>
  </div>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import BoatTabs from "@/components/BoatTabs";
import MainTitle from "@/components/MainTitle";
import { BRow, BCol, BCard, BButton, BIconHouse } from "bootstrap-vue";

export default {
  name: "WSBoat",
  components: {
    NavbarMobile,
    MainTitle,
    BoatTabs,
    BRow,
    BCol,
    BCard,
    BButton,
    BIconHouse,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
  },
};
</script>
