<template>
  <modal-container
    name="heatCommentModal"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <modal-header
      :closable="true"
      title="Edit Comment"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length > 0" :errors="errors" />
      <form @submit.stop="saveComment">
        <input-text
          id="comment-input"
          v-model="comment"
          placeholder="add your comment here..."
        />
      </form>
    </modal-body>
    <modal-footer>
      <button
        type="button"
        class="btn btn-outline-danger mr-1"
        @click="removeComment"
      >
        <i class="bi bi-x" />
        Delete
      </button>
      <button type="button" class="btn btn-outline-info" @click="saveComment">
        <i class="bi bi-check" />
        Save
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import WarningBox from "@/components/WarningBox";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";

export default {
  name: "HeatCommentModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputText,
  },
  props: ["defaultComment", "visible"],
  data() {
    return {
      errors: [],
      comment: null,
    };
  },
  watch: {
    defaultComment: function (newComment) {
      this.comment = newComment;
    },
  },
  mounted() {
    this.comment = this.defaultComment;
  },
  methods: {
    saveComment: function () {
      this.$emit("commentChangeHandler", this.comment);
      this.close();
    },
    removeComment: function () {
      this.$emit("commentChangeHandler", null);
      this.close();
    },
    close: function () {
      this.comment = null;
      this.$emit("update:visible", false);
    },
  },
};
</script>
