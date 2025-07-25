<template>
  <modal-container
    name="session-delete-modal"
    :visible="visible"
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

<script>
import { mapActions } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "booksys/components/bricks/ModalContainer.vue";
import ModalHeader from "booksys/components/bricks/ModalHeader.vue";
import ModalBody from "booksys/components/bricks/ModalBody.vue";
import ModalFooter from "booksys/components/bricks/ModalFooter.vue";

export default {
  name: "SessionDeleteModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
  },
  props: ["session", "visible"],
  data() {
    return {
      errors: [],
    };
  },
  methods: {
    ...mapActions("sessions", ["deleteSession"]),
    confirm: function (event) {
      if (event != null) {
        event.preventDefault();
      }

      this.deleteSession(this.session)
        .then(() => {
          this.$emit("sessionDeletedHandler");
          this.close();
        })
        .catch((err) => {
          this.errors = err;
        });
    },
    close: function () {
      this.errors = [];
      this.$emit("update:visible", false);
    },
  },
};
</script>
