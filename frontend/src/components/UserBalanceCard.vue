<template>
  <sectioned-card-module title="Balance">
    <template v-slot:body>
      <payment-info-modal v-model:visible="showPaymentInfoModal" />
      <div class="row">
        <div class="col-6">
          <div class="row">
            <div class="col-12">
              Current Balance: {{ this.formatBalance(balanceRounded) }}
            </div>
          </div>
        </div>
        <div class="col-6">
          <button
            type="button"
            class="btn btn-outline-info btn-sm"
            @click="showPaymentInfo"
          >
            <i class="bi bi-cash-coin"></i>
            Buy Sessions
          </button>
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import SectionedCardModule from "./bricks/SectionedCardModule.vue";
import PaymentInfoModal from "./PaymentInfoModal.vue";
import { formatCost } from "booksys/libs/formatters";

export default {
  name: "UserBalanceCard",
  components: {
    PaymentInfoModal,
    SectionedCardModule,
  },
  data() {
    return {
      showPaymentInfoModal: false,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getCurrency"]),
    ...mapGetters("user", ["balanceRounded"]),
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    showPaymentInfo: function () {
      this.showPaymentInfoModal = true;
    },
    formatBalance: function(value) {
      return formatCost(value, this.getCurrency);
    }
  },
  created() {
    this.queryConfiguration();
  },
};
</script>
