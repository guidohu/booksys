<template>
  <b-tabs content-class="mt-3">
    <b-tab active>
      <div slot="title">
        <p class="h6 mt-2 mb-2">
          <b-icon-cash-stack class="mr-1" />
          Payments
        </p>
      </div>
      <div :class="tabClass">
        <PaymentTable />
      </div>
    </b-tab>
    <b-tab>
      <div slot="title">
        <p class="h6 mt-2 mb-2">
          <b-icon-kanban
            rotate="180"
            class="mr-1"
          />
          Statistics
        </p>
      </div>
      <div :class="tabClass">
        <PaymentStatistics />
      </div>
    </b-tab>
  </b-tabs>
</template>

<script>
import { BooksysBrowser } from "@/libs/browser";
import PaymentTable from "@/components/PaymentTable";
import PaymentStatistics from "@/components/PaymentStatistics";
import { BTabs, BTab, BIconKanban, BIconCashStack } from "bootstrap-vue";

export default {
  name: "AccountingTabs",
  components: {
    PaymentStatistics,
    PaymentTable,
    BTabs,
    BTab,
    BIconKanban,
    BIconCashStack,
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile();
    },
    tabClass: function () {
      if (!BooksysBrowser.isMobile()) {
        return "accounting-tab-limited-height";
      } else {
        return "";
      }
    },
  },
};
</script>

<style scoped>
.accounting-tab-limited-height {
  max-height: 400px;
  overflow: scroll;
}
</style>
