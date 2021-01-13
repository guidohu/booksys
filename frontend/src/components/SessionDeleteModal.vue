<template>
  <b-modal
    id="sessionDeleteModal"
    ref="sessionDeleteModal"
    title="Delete Session"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <b-row v-if="errors.length">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <WarningBox :errors="errors"/>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-form @submit="confirm">
      <b-row>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          Do you really want to delete this session?
        </b-col>
        <b-col cols="1" class="d-none d-sm-block"></b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <b-button type="button" variant="outline-info" v-on:click="confirm">
        <b-icon-check/>
        Delete
      </b-button>
      <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
        <b-icon-x/>
        Cancel
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { mapActions } from 'vuex';
import WarningBox from '@/components/WarningBox';
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BButton,
  BIconCheck,
  BIconX
} from 'bootstrap-vue';

export default {
  name: 'SessionDeleteModal',
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BButton,
    BIconCheck,
    BIconX
  },
  props: [
    'session',
    'visible'
  ],
  data(){
    return {
      errors: []
    }
  },
  methods: {
    ...mapActions('sessions', [
      'deleteSession'
    ]),
    confirm: function(event){
      if(event != null){
        event.preventDefault();
      }

      this.deleteSession(this.session)
        .then( () => {
          this.$emit('sessionDeletedHandler');
          this.close()
        })
        .catch( err => {
          this.errors = err;
        })
    },
    close: function(){
      this.errors = [];
      this.$emit('update:visible', false);
    }
  }
}
</script>
