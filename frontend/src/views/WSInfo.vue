<template>
  <subpage-container title="Location">
    <card-module :nobody="true" class="mx-1 card-height">
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
    </card-module>
    <template v-slot:bottom>
      <router-link tag="button" class="btn btn-outline-light" to="/dashboard">
        <i class="bi bi-house" />
        Home
      </router-link>
    </template>
  </subpage-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import { BooksysBrowser } from "@/libs/browser";
import CardModule from "@/components/bricks/CardModule.vue";
import SubpageContainer from "@/components/bricks/SubpageContainer.vue";

export default {
  name: "WSInfo",
  components: {
    CardModule,
    SubpageContainer,
  },
  computed: {
    ...mapGetters("configuration", ["getLocationAddress", "getLocationMap"]),
    isMobile: function () {
      return BooksysBrowser.isMobile();
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

<style scoped>
.card-height {
  max-height: 500px;
  min-height: 500px;
  height: 500px;
  overflow-y: scroll;
  overflow-x: hidden;
}

@media (max-width: 992px) {
  .card-height {
    min-height: 500px;
    overflow-y: scroll;
  }
}
</style>
