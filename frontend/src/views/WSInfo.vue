<template>
  <subpage-container title="Location">
    <card-module :nobody="true" class="mx-1 card-height">
      <div class="row">
        <div class="col-12 mt-5 text-center">
          <div
            v-if="getLocationAddress != null && getLocationAddress.length > 0"
            cols="12"
            class="main-color"
            v-html="getLocationAddress"
          />
          <div class="main-color" v-else>
            [ no address set by the site-owner ]
          </div>
        </div>
      </div>
      <div class="row">
        <div class="col-12 text-center mt-3">
          <iframe
            v-if="getLocationMap != null && getLocationMap.length > 0"
            :src="getLocationMap"
            frameborder="0"
            style="border: 0"
            :width="mapWidth"
            :height="mapHeight"
          />
          <div class="main-color" v-else>[ no map configured ]</div>
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

<script setup>
import { useStore, mapGetters } from "vuex";
import { onMounted, computed } from "vue";
import { BooksysBrowser } from "booksys/libs/browser";
import CardModule from "booksys/components/bricks/CardModule.vue";
import SubpageContainer from "booksys/components/bricks/SubpageContainer.vue";

const store = useStore();

const getLocationAddress = computed(
  mapGetters("configuration", ["getLocationAddress"]).getLocationAddress.bind({
    $store: store,
  }),
);
const getLocationMap = computed(
  mapGetters("configuration", ["getLocationMap"]).getLocationMap.bind({
    $store: store,
  }),
);

const mapHeight = computed(() => {
  if (BooksysBrowser.isMobileResponsive()) {
    return 350;
  } else {
    return 300;
  }
});

const mapWidth = computed(() => {
  if (BooksysBrowser.isMobileResponsive()) {
    return 340;
  } else {
    return 600;
  }
});

onMounted(() => {
  console.log("Info.vue: Try to get all information required");
  store.dispatch("configuration/queryConfiguration");
});
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
