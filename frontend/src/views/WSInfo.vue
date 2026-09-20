<template>
  <subpage-container title="Location">
    <card-module :nobody="true" fill class="mx-1 card-height">
      <div class="info-content">
        <div class="address-section text-center mt-5">
          <div
            v-if="getLocationAddress != null && getLocationAddress.length > 0"
            class="main-color"
            v-html="getLocationAddress"
          />
          <div class="main-color" v-else>
            [ no address set by the site-owner ]
          </div>
        </div>
        <div class="map-section mt-3">
          <iframe
            v-if="getLocationMap != null && getLocationMap.length > 0"
            class="map-frame"
            :src="getLocationMap"
            frameborder="0"
          />
          <div class="main-color text-center" v-else>[ no map configured ]</div>
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
  overflow: auto;
}

.info-content {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
}

.address-section {
  flex: 0 0 auto;
}

/* Takes up whatever vertical space is left over below the address */
.map-section {
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 0.5rem;
}

.map-frame {
  border: 0;
  flex: 0 0 auto;
  width: 600px;
  max-width: 100%;
  height: 300px;
}

@media (max-width: 992px) {
  /* Fill the whole phone screen below the fixed navbar (60px padding-top) */
  .card-height {
    min-height: calc(100vh - 66px);
    max-height: calc(100vh - 66px);
    height: calc(100vh - 66px);
    min-height: calc(100dvh - 66px - env(safe-area-inset-bottom, 0px));
    max-height: calc(100dvh - 66px - env(safe-area-inset-bottom, 0px));
    height: calc(100dvh - 66px - env(safe-area-inset-bottom, 0px));
    /* Everything fits through flex sizing, no inner scrolling needed */
    overflow: hidden;
  }

  .address-section {
    margin-top: 1.5rem !important;
  }

  .map-section {
    padding: 0 0.5rem 0.5rem;
  }

  /* Grow into the leftover space instead of a fixed height */
  .map-frame {
    flex: 1 1 auto;
    min-height: 0;
    width: 100%;
    height: auto;
  }
}
</style>
