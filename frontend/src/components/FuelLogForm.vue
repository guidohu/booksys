<template>
  <div class="container text-left">
    <fuel-refuel-modal
      v-if="visibleRefuelModal"
      v-model:visible="visibleRefuelModal"
    />
    <form-button
      type="button"
      btn-style="info"
      @click.prevent.self="showRefuelModal"
    >
      Refuel
    </form-button>
  </div>
</template>

<script setup>
import { ref } from "vue";
import FormButton from "booksys/components/forms/FormButton.vue";
import FuelRefuelModal from "booksys/components/FuelRefuelModal.vue";

const visibleRefuelModal = ref(false);

function showRefuelModal() {
  visibleRefuelModal.value = true;
}
</script>
