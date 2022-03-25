<template>
  <modal-container name="user-heats-modal" :visible="visible">
    <modal-header
      :closable="true"
      title="Payment Information"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <div class="row">
        <div class="col-12">
          Please transfer the desired amount to the bank account below and let
          us know as soon as you did so to activate the amount.
        </div>
      </div>
      <div class="row mt-3">
        <div class="col-3">Owner</div>
        <div class="col-7">
          {{
            getConfiguration && getConfiguration.payment_account_owner
              ? getConfiguration.payment_account_owner
              : "n/a"
          }}
        </div>
      </div>
      <div class="row mt-3">
        <div class="col-3">IBAN</div>
        <div class="col-7">
          {{
            getConfiguration && getConfiguration.payment_account_iban
              ? getConfiguration.payment_account_iban
              : "n/a"
          }}
        </div>
      </div>
      <div class="row mt-3">
        <div class="col-3">BIC</div>
        <div class="col-7">
          {{
            getConfiguration && getConfiguration.payment_account_bic
              ? getConfiguration.payment_account_bic
              : "n/a"
          }}
        </div>
      </div>
      <div class="row mt-3">
        <div class="col-3">Comment</div>
        <div class="col-7">
          {{
            getConfiguration && getConfiguration.payment_account_comment
              ? getConfiguration.payment_account_comment
              : "n/a"
          }}
        </div>
      </div>
    </modal-body>
    <modal-footer>
      <button type="submit" class="btn btn-outline-info" @click="close">
        <i class="bi bi-check"></i>
        OK
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapGetters } from "vuex";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";

export default {
  name: "PaymentInfoModal",
  components: {
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
  },
  props: ["visible"],
  computed: {
    ...mapGetters("configuration", ["getConfiguration"]),
  },
  methods: {
    close: function () {
      this.$emit("update:visible", false);
    },
  },
};
</script>
