<template>
  <modal-container name="confirm-dialog" :visible="visible">
    <modal-header v-if="title != null && title != ''">
      {{ title }}
    </modal-header>
    <modal-body v-if="message != null && message != ''">
      <p>
        {{ message }}
      </p>
    </modal-body>
    <modal-footer>
      <button
        class="btn btn-outline-info"
        type="button"
        data-bs-dismiss="modal"
        @click.stop="confirm"
      >
        <i class="bi bi-check"></i>
        {{ textConfirm ? textConfirm : "OK" }}
      </button>
      <button
        class="btn btn-outline-danger"
        type="button"
        data-bs-dismiss="modal"
        @click.stop="close"
      >
        <i class="bi bi-x"></i>
        {{ textDeny ? textDeny : "Cancel" }}
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import ModalBody from "@/components/bricks/ModalBody.vue";
import ModalContainer from "../ModalContainer.vue";
import ModalFooter from "../ModalFooter.vue";
import ModalHeader from "../ModalHeader.vue";

export default {
  name: "ConfirmDialog",
  props: ["title", "message", "textConfirm", "textDeny"],
  data() {
    return {
      visible: true,
    };
  },
  emits: ["confirm"],
  components: {
    ModalHeader,
    ModalBody,
    ModalContainer,
    ModalFooter,
  },
  methods: {
    close: function () {
      this.$emit("confirm", false);
      return true;
    },
    confirm: function () {
      this.$emit("confirm", true);
      return true;
    },
  },
};
</script>
