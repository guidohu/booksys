<template>
  <div class="modal-header">
    <h5 class="modal-title">
      {{ title }}
    </h5>
    <button
      v-if="closable"
      @click="closeHandler()"
      type="button"
      class="btn-close"
      data-bs-dismiss="modal"
      aria-label="Close"
    ></button>
  </div>
</template>

<script>
export default {
  name: "ModalHeader",
  props: ["closable", "title"],
  methods: {
    closeHandler: function () {
      console.log("closeHandler of modal-header called");
      this.$emit("close");
    },
  },
};
</script>
