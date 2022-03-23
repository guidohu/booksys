<template>
  <div class="display">
    <main-title :title-name="title" />
    <div class="row ml-1 mr-1">
      <div class="col-12">
        <slot name="content"></slot>
      </div>
    </div>
    <div class="bottom mr-2">
      <div class="row">
        <div class="col-12">
          <slot name="bottom"></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import MainTitle from "@/components/MainTitle.vue";

export default {
  name: "SubpageDesktop",
  props: ["title"],
  components: {
    MainTitle,
  },
};
</script>
