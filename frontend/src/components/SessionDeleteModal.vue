<template>
  <modal-container
    name="session-delete-modal"
    :visible="visible"
    :inert="!visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <modal-header
      :closable="true"
      title="Delete Session"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length" :errors="errors" />
      <div class="row">
        <div class="col-12">Do you really want to delete this session?</div>
      </div>
    </modal-body>
    <modal-footer>
      <button type="button" class="btn btn-outline-danger" @click="confirm">
        <i class="bi bi-trash"></i>
        Delete
      </button>
      <button type="button" class="btn btn-outline-info" @click="close">
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "booksys/components/bricks/ModalContainer.vue";
import ModalHeader from "booksys/components/bricks/ModalHeader.vue";
import ModalBody from "booksys/components/bricks/ModalBody.vue";
import ModalFooter from "booksys/components/bricks/ModalFooter.vue";

const props = defineProps(["session", "visible"]);
const emit = defineEmits(["sessionDeletedHandler", "update:visible"]);

const store = useStore();

const errors = ref([]);

const deleteSession = (session) =>
  store.dispatch("sessions/deleteSession", session);

function confirm(event) {
  if (event != null) {
    event.preventDefault();
  }

  deleteSession(props.session)
    .then(() => {
      emit("sessionDeletedHandler");
      close();
    })
    .catch((err) => {
      errors.value = err;
    });
}

function close() {
  errors.value = [];
  emit("update:visible", false);
}
</script>
