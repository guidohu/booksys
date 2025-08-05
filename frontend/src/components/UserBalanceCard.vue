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

<script setup>
import { ref, computed } from "vue";
import { useStore } from "vuex";
import SectionedCardModule from "./bricks/SectionedCardModule.vue";
import PaymentInfoModal from "./PaymentInfoModal.vue";
import { formatCost } from "booksys/libs/formatters";

const store = useStore();

const showPaymentInfoModal = ref(false);

const userInfo = computed(() => store.getters["login/userInfo"]);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);
const balanceRounded = computed(() => store.getters["user/balanceRounded"]);

const queryConfiguration = () =>
  store.dispatch("configuration/queryConfiguration");

function showPaymentInfo() {
  showPaymentInfoModal.value = true;
}

function formatBalance(value) {
  return formatCost(value, getCurrency.value);
}

queryConfiguration();
</script>
