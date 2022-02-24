<template>
  <modal-container
    name="info-dialog"
    :visible="visible"
  >
    <modal-header v-if="title != null && title != ''">
      {{ title }}
    </modal-header>
    <modal-body v-if="message != null && message != ''">
      <p>
        {{ message}}
      </p>
    </modal-body>
    <modal-footer>
      <button class="btn btn-outline-info" type="button" data-bs-dismiss="modal" @click.stop="confirm">
        <i class="bi bi-check"></i>
        {{ textConfirm ? textConfirm : "OK" }}
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import ModalBody from "@/components/bricks/ModalBody.vue";
import ModalContainer from '../ModalContainer.vue';
import ModalFooter from '../ModalFooter.vue';
import ModalHeader from '../ModalHeader.vue';

export default {
  name: "InfoDialog",
  props: ["title", "message", "textConfirm"],
  data() {
    return {
      visible: true
    }
  },
  emits: ["confirm"], 
  components: {
    ModalHeader,
    ModalBody,
    ModalContainer,
    ModalFooter,
  }, 
  methods: {
    confirm: function() {
      this.$emit('confirm', true);
      return true;
    }
  }
}
</script>