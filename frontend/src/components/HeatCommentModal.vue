<template>
  <b-modal
    id="heatCommentModal"
    title="Heat Comment"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <b-row v-if="errors.length">
      <b-col cols="1" class="d-none d-sm-block" />
      <b-col cols="12" sm="10">
        <WarningBox :errors="errors" />
      </b-col>
      <b-col cols="1" class="d-none d-sm-block" />
    </b-row>
    <b-form @submit="saveComment">
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block" />
        <b-col cols="12" sm="10">
          <b-form-group id="comment" label="" label-for="" description="">
            <b-form-input
              id="comment-input"
              v-model="comment"
              type="text"
              placeholder="add your comment here..."
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <b-button
        type="button"
        variant="outline-danger"
        class="mr-1"
        @click="removeComment"
      >
        <b-icon-x />
        Delete
      </b-button>
      <b-button type="button" variant="outline-info" @click="saveComment">
        <b-icon-check />
        Save
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import WarningBox from "@/components/WarningBox";
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormInput,
  BButton,
  BIconX,
  BIconCheck,
} from "bootstrap-vue";

export default {
  name: "HeatCommentModal",
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormInput,
    BButton,
    BIconX,
    BIconCheck,
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
