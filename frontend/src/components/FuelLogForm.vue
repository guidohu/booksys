<template>
  <div class="container text-left">
    <fuel-refuel-modal v-model:visible="visibleRefuelModal"/>
    <form-button
      type="button"
      btn-style="info"
      @click.prevent.self="showRefuelModal"
    >
      Refuel
    </form-button>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import FormButton from "@/components/forms/FormButton.vue";
import FuelRefuelModal from "@/components/FuelRefuelModal.vue";

export default {
  name: "FuelLogForm",
  components: {
    FuelRefuelModal,
    FormButton,
  },
  data() {
    return {
      errors: [],
      form: {
        engineHours: null,
        cost: null,
        liters: null,
      },
      visibleRefuelModal: false,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getCurrency", "getEngineHourFormat"]),
  },
  methods: {
    showRefuelModal: function(){
      this.visibleRefuelModal = true;
    },
  },
};
</script>
