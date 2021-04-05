<template>
  <b-modal
    id="paymentInfoModal"
    title="Payment Information"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <b-row class="text-left">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <b-row>
          <b-col cols="12">
            Please transfer the desired amount to the bank account below and let
            us know as soon as you did so to activate the amount.
          </b-col>
        </b-row>
        <b-row class="mt-3">
          <b-col cols="3"> Owner </b-col>
          <b-col cols="7">
            {{
              getConfiguration && getConfiguration.payment_account_owner
                ? getConfiguration.payment_account_owner
                : "n/a"
            }}
          </b-col>
        </b-row>
        <b-row>
          <b-col cols="3"> IBAN </b-col>
          <b-col cols="7">
            {{
              getConfiguration && getConfiguration.payment_account_iban
                ? getConfiguration.payment_account_iban
                : "n/a"
            }}
          </b-col>
        </b-row>
        <b-row>
          <b-col cols="3"> BIC </b-col>
          <b-col cols="7">
            {{
              getConfiguration && getConfiguration.payment_account_bic
                ? getConfiguration.payment_account_bic
                : "n/a"
            }}
          </b-col>
        </b-row>
        <b-row>
          <b-col cols="3"> Comment </b-col>
          <b-col cols="7">
            {{
              getConfiguration && getConfiguration.payment_account_comment
                ? getConfiguration.payment_account_comment
                : "n/a"
            }}
          </b-col>
        </b-row>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"> </b-col>
    </b-row>
    <div slot="modal-footer">
      <b-button type="button" variant="outline-info" v-on:click="close">
        <b-icon-check />
        Done
      </b-button>
    </div>
  </b-modal>
</template>

<script>
import { mapGetters } from "vuex";
import { BModal, BRow, BCol, BButton, BIconCheck } from "bootstrap-vue";

export default {
  name: "PaymentInfoModal",
  props: ["visible"],
  components: {
    BModal,
    BRow,
    BCol,
    BButton,
    BIconCheck,
  },
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
