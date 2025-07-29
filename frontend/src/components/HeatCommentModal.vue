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

<script setup>
import { ref, watch, onMounted } from "vue";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";

const props = defineProps(["defaultComment", "visible"]);
const emit = defineEmits(["update:visible", "commentChangeHandler"]);

const errors = ref([]);
const comment = ref(null);

watch(() => props.defaultComment, (newComment) => {
  comment.value = newComment;
});

onMounted(() => {
  comment.value = props.defaultComment;
});

function saveComment() {
  emit("commentChangeHandler", comment.value);
  close();
}

function removeComment() {
  emit("commentChangeHandler", null);
  close();
}

function close() {
  comment.value = null;
  emit("update:visible", false);
}
</script>
