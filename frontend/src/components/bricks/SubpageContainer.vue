<template>
  <div>
    <show-for-desktop id="subpage-container-show-for-desktop">
      <subpage-desktop :title="title">
        <template v-slot:content>
          <div class="content-max-height">
            <slot></slot>
          </div>
        </template>
        <template v-slot:bottom>
          <slot name="bottom">
            <router-link
              tag="button"
              class="btn btn-outline-light"
              to="/dashboard"
            >
              <i class="bi bi-house"></i>
              HOME
            </router-link>
          </slot>
        </template>
      </subpage-desktop>
    </show-for-desktop>
    <show-for-mobile id="subpage-container-show-for-mobile">
      <navbar-mobile :title="title" />
      <div class="navbar-padding">
        <div class="row">
          <div class="col-12">
            <slot></slot>
          </div>
        </div>
      </div>
    </show-for-mobile>
  </div>
</template>

<script>
import SubpageDesktop from "@/components/bricks/SubpageDesktop.vue";
import NavbarMobile from "@/components/NavbarMobile.vue";
import ShowForDesktop from "@/components/bricks/ShowForDesktop.vue";
import ShowForMobile from "@/components/bricks/ShowForMobile.vue";

export default {
  name: "SubpageContainer",
  props: ["title"],
  components: {
    ShowForDesktop,
    ShowForMobile,
    SubpageDesktop,
    NavbarMobile,
  },
};
</script>

<style scoped>
.navbar-padding {
  padding-top: 60px;
}

.content-max-height {
  min-height: 480px;
  max-height: 480px;
  overflow-y: scroll;
}
</style>
