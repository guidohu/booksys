<template>
  <b-card no-body class="text-left">
    <b-card-header> Balance </b-card-header>
    <b-card-body>
      <b-row>
        <b-col cols="6">
          <b-row>
            <b-col cols="12">
              Current Balance: {{ balanceRounded }} {{ getCurrency }}
            </b-col>
          </b-row>
        </b-col>
        <b-col cols="6">
          <b-button
            size="sm"
            type="button"
            variant="outline-info"
            v-on:click="showPaymentInfo"
          >
            <b-icon-cash />
            Buy Sessions
          </b-button>
        </b-col>
      </b-row>
      <PaymentInfoModal :visible.sync="showPaymentInfoModal" />
    </b-card-body>
  </b-card>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import PaymentInfoModal from "./PaymentInfoModal";
import {
  BCard,
  BCardHeader,
  BCardBody,
  BRow,
  BCol,
  BButton,
  BIconCash,
} from "bootstrap-vue";

export default {
  name: "UserBalanceCard",
  components: {
    PaymentInfoModal,
    BCard,
    BCardHeader,
    BCardBody,
    BRow,
    BCol,
    BButton,
    BIconCash,
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
  },
  created() {
    this.queryConfiguration();
  }
};
</script>
