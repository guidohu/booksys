<template>
  <div
    class="modal fade"
    :ref="name"
    :id="name"
    data-bs-target="static"
    aria-hidden="false"
    tabindex="-1"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
import { Modal } from "bootstrap";

export default {
  name: "ModalContainer",
  props: ["visible", "name"],
  mounted() {
    this.modal = new Modal(this.$refs[this.name]);
    this.modal.show();
  },
  watch: {
    visible: function (newValue) {
      if (newValue) {
        this.modal.show();
      } else {
        this.modal.hide();
      }
    },
  },
};
</script>
