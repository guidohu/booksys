<template>
  <div>
    <ul class="nav nav-tabs" id="accounting-tab" role="tablist">
      <li class="nav-item" role="presentation">
        <a
          class="nav-link active"
          data-bs-toggle="tab"
          href="#payments"
          id="payments-tab"
          data-bs-target="#payments"
          role="tab"
          aria-controls="payments"
          aria-selected="true"
        >
          <i class="bi bi-cash-stack"></i>
          Payments
        </a>
      </li>
      <li class="nav-item" role="presentation">
        <a
          class="nav-link"
          data-bs-toggle="tab"
          href="#statistics"
          id="statistics-tab"
          data-bs-target="#statistics"
          role="tab"
          aria-controls="statistics"
          aria-selected="false"
        >
          <i class="bi bi-graph-up"></i>
          Statistics
        </a>
      </li>
    </ul>

    <!-- Responsive tab content to handle different
         height requirements -->
    <div class="tab-content">
      <div
        class="tab-pane active"
        id="payments"
        role="tabpanel"
        aria-labelledby="payments-tab"
      >
        <show-for-desktop>
          <div class="tab-limited-height">
            <payment-table />
          </div>
        </show-for-desktop>
        <show-for-mobile>
          <payment-table />
        </show-for-mobile>
      </div>
      <div
        class="tab-pane"
        id="statistics"
        role="tabpanel"
        aria-labelledby="statistics-tab"
      >
        <show-for-desktop>
          <div class="tab-limited-height">
            <payment-statistics />
          </div>
        </show-for-desktop>
        <show-for-mobile>
          <payment-statistics />
        </show-for-mobile>
      </div>
    </div>
  </div>
</template>

<script setup>
import PaymentTable from "booksys/components/PaymentTable.vue";
import ShowForMobile from "booksys/components/bricks/ShowForMobile.vue";
import ShowForDesktop from "booksys/components/bricks/ShowForDesktop.vue";
import PaymentStatistics from "booksys/components/PaymentStatistics.vue";
</script>

<style scoped>
.tab-limited-height {
  max-height: 420px;
  height: 420px;
  overflow-y: scroll;
}

.nav-link {
  color: #bdbdbd;
}

a {
  color: #bdbdbd;
}
</style>
