<template>
  <div class="container text-left">
    <div class="row">
      <div class="col-12">
        <warning-box v-if="errors.length > 0" :errors="errors" />
        <form @submit.prevent="add">
          <input-engine-hours
            id="engine-hours"
            label="Engine"
            :display-format="getEngineHourFormat"
            v-model="form.engineHours"
            placeholder="0"
            size="small"
          />
          <input-text-multiline
            id="description"
            label="Description"
            v-model="form.description"
            rows="3"
            placeholder="Add your notes here..."
            size="small"
          />
          <form-button
            type="submit"
            btn-style="info"
            btn-size="small"
            @click.prevent="add"
          >
            Add
          </form-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import InputEngineHours from "booksys/components/forms/inputs/InputEngineHours.vue";
import InputTextMultiline from "booksys/components/forms/inputs/InputTextMultiline.vue";
import FormButton from "booksys/components/forms/FormButton.vue";

const store = useStore();

const errors = ref([]);
const form = ref({
  engineHours: null,
  description: null,
});

const userInfo = computed(() => store.getters["login/userInfo"]);
const getEngineHourFormat = computed(() => store.getters["configuration/getEngineHourFormat"]);

function add() {
  if (isNaN(form.value.engineHours) || form.value.engineHours == null) {
    errors.value = ["Please add a valid number for the engine hours."];
    return;
  }

  const entry = {
    user_id: userInfo.value.id,
    engine_hours: form.value.engineHours,
    description: form.value.description,
  };
  store.dispatch("boat/addMaintenanceEntry", entry)
    .then(() => {
      resetForm();
    })
    .catch((errs) => (errors.value = errs));
}

function resetForm() {
  form.value = {
    engineHours: null,
    description: null,
  };
  errors.value = [];
}

store.dispatch("configuration/queryConfiguration");
</script>
