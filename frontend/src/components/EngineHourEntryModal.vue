<template>
  <modal-container name="engineHourentryModal" :visible="visible">
    <modal-header
      :closable="true"
      title="Engine Hour Entry"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <div class="row" v-if="errors.length">
        <div class="col-1 d-none d-sm-block" />
        <div class="col-12 col-sm-10">
          <warning-box :errors="errors" />
        </div>
      </div>
      <form @submit.prevent.self="save">
        <input-date-time-local
          id="date"
          label="Date"
          v-model="form.date"
          size="small"
          :disabled="true"
        />
        <input-text
          id="driver"
          label="Driver"
          v-model="form.driver"
          size="small"
          :disabled="true"
        />
        <input-engine-hours
          id="before"
          label="Before"
          v-model="form.beforeHours"
          :displayFormat="displayFormat"
          size="small"
          :disabled="true"
        />
        <input-engine-hours
          id="after"
          label="After"
          v-model="form.afterHours"
          :displayFormat="displayFormat"
          size="small"
          :disabled="true"
        />
        <input-engine-hours
          id="delta"
          label="Difference"
          v-model="form.deltaHours"
          :displayFormat="displayFormat"
          size="small"
          :disabled="true"
        />
        <input-toggle
          id="type"
          label="Type"
          v-model="form.type"
          offLabel="Private"
          onLabel="Course"
        />
      </form>
    </modal-body>
    <modal-footer>
      <button
        type="submit"
        class="btn btn-outline-info"
        @click.prevent.self="save"
      >
        <i class="bi bi-check"></i>
        Save
      </button>
      <button type="button" class="btn btn-outline-danger" @click="close">
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "booksys/components/bricks/ModalContainer.vue";
import ModalHeader from "booksys/components/bricks/ModalHeader.vue";
import ModalBody from "booksys/components/bricks/ModalBody.vue";
import ModalFooter from "booksys/components/bricks/ModalFooter.vue";
import dayjs from "dayjs";
import InputDateTimeLocal from "booksys/components/forms/inputs/InputDateTimeLocal.vue";
import InputText from "booksys/components/forms/inputs/InputText.vue";
import InputEngineHours from "booksys/components/forms/inputs/InputEngineHours.vue";
import InputToggle from "booksys/components/forms/inputs/InputToggle.vue";

const props = defineProps({
  engineHourEntry: {
    type: Object,
    required: true,
  },
  visible: {
    type: Boolean,
    required: true,
  },
  displayFormat: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["update:visible"]);
const store = useStore();

const errors = ref([]);
const form = ref({ date: null });

const toggleWidth = computed(() => 70);

watch(
  () => props.engineHourEntry,
  (newValue) => {
    console.debug("EngineHourEntryModal: new engineHourEntry", newValue);
    setFormContent(newValue);
  },
);

const setFormContent = (entry) => {
  if (entry != null) {
    form.value = {
      id: entry.id,
      date: dayjs.unix(entry.timestamp).format("YYYY-MM-DDTHH:mm"),
      driver: entry.user_first_name,
      beforeHours: entry.before_hours,
      afterHours: entry.after_hours,
      deltaHours: entry.delta_hours,
      type: entry.type == 1 ? false : true,
    };
  }
};

const toggleType = () => {
  form.value.type = !form.value.type;
};

const close = () => {
  emit("update:visible", false);
};

const save = () => {
  const update = {
    id: form.value.id,
    type: form.value.type == true ? 2 : 1,
  };
  store
    .dispatch("boat/updateEngineHours", update)
    .then(() => close())
    .catch((error) => (errors.value = error));
};

setFormContent(props.engineHourEntry);
</script>
