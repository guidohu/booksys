<template>
  <b-modal
    id="heatCommentModal"
    title="Heat Comment"
  >
    <b-row v-if="errors.length">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <WarningBox :errors="errors"/>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-form @submit="saveComment">
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <b-form-group
            id="comment"
            label=""
            label-for=""
            description=""
          >
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
      <b-button type="button" variant="outline-danger" v-on:click="removeComment" class="mr-1">
        <b-icon-x></b-icon-x>
        Delete
      </b-button>
      <b-button type="button" variant="outline-info" v-on:click="saveComment">
        <b-icon-check></b-icon-check>
        Save
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import Vue from 'vue';
import WarningBox from '@/components/WarningBox';

export default Vue.extend({
  name: "HeatCommentModal",
  components: {
    WarningBox
  },
  props: [
    'defaultComment'
  ],
  data() {
    return {
      errors: [],
      comment: null
    }
  },
  methods: {
    saveComment: function(){
      this.$emit('commentChangeHandler', this.comment);
      this.close();
    },
    removeComment: function() {
      this.$emit('commentChangeHandler', null);
      this.close();
    },
    close: function() {
      this.comment = null;
      this.$bvModal.hide('heatCommentModal');
    }
  },
  watch: {
    defaultComment: function(newComment){
      this.comment = newComment;
    }
  },
  mounted() {
    this.comment = this.defaultComment;
  }
})
</script>