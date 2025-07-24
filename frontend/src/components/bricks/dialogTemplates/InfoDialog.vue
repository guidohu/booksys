<template>
  <modal-container name="info-dialog" :visible="visible">
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
    </modal-footer>
  </modal-container>
</template>

<script setup>
import ModalBody from "@/components/bricks/ModalBody.vue";
import ModalContainer from "../ModalContainer.vue";
import ModalFooter from "../ModalFooter.vue";
import ModalHeader from "../ModalHeader.vue";
import { ref } from "vue";

defineProps(["title", "message", "textConfirm"]);
const emit = defineEmits(["confirm"]);

const visible = ref(true);

function confirm() {
  emit("confirm", true);
  return true;
}
</script>
