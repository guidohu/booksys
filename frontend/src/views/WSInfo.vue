<template>
  <div v-if="isDesktop" class="display">
    <main-title title-name="Location" />
    <div class="row">
      <div class="col-12 mt-5 text-center">
        <div
          v-if="getLocationAddress != null"
          cols="12"
          class="main-color"
          v-html="getLocationAddress"
        />
        <div v-else>no address set by the site-owner</div>
      </div>
    </div>
    <div class="row">
      <div class="col-12 text-center mt-3">
        <iframe
          v-if="getLocationMap != null"
          :src="getLocationMap"
          frameborder="0"
          style="
             {
              border: 0;
            }
          "
          :width="mapWidth"
          :height="mapHeight"
        />
        <div v-else>no map to display for this adress</div>
      </div>
    </div>
    <div class="bottom">
      <button class="btn btn-outline-light" to="/dashboard">
        <i class="bi bi-house" />
        Home
      </button>
    </div>
  </div>
  <div v-else>
    <NavbarMobile title="Location" />
    <card-module>
      <div class="row text-center">
        <div class="col-12 text-center">
          <div
            v-if="getLocationAddress != null"
            cols="12"
            v-html="getLocationAddress"
          />
          <div v-else>no address is configured</div>
        </div>
      </div>
      <div class="row">
        <div class="col-12 text-center mt-3">
          <iframe
            v-if="getLocationMap != null"
            :src="getLocationMap"
            frameborder="0"
            style="
                {
                border: 0;
              }
            "
            :width="mapWidth"
            :height="mapHeight"
          />
          <div v-else>no map is configured</div>
        </div>
      </div>
    </card-module>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import NavbarMobile from "@/components/NavbarMobile";
import MainTitle from "@/components/MainTitle";
import CardModule from "@/components/bricks/CardModule.vue";
import {
  BRow,
  BCol,
  BButton,
  BIconHouse,
  BCard,
  BCardBody,
  BCardText,
  CardModule,
} from "bootstrap-vue";

export default {
  name: "WSInfo",
  components: {
    MainTitle,
    NavbarMobile,
    BRow,
    BCol,
    BButton,
    BIconHouse,
    BCard,
    BCardBody,
    BCardText,
  },
  computed: {
    ...mapGetters("configuration", ["getLocationAddress", "getLocationMap"]),
    isMobile: function () {
      return BooksysBrowser.isMobile();
    },
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
    mapHeight: function () {
      if (this.isMobile) {
        return 350;
      } else {
        return 300;
      }
    },
    mapWidth: function () {
      if (this.isMobile) {
        return 340;
      } else {
        return 600;
      }
    },
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
  },
  created() {
    console.log("Info.vue: Try to get all information required");
    this.queryConfiguration();
  },
};
</script>
