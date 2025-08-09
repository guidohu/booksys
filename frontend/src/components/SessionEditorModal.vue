<template>
  <modal-container
    name="session-editor-modal"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <modal-header
      :closable="true"
      :title="title"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length" :errors="errors" />
      <form @submit.prevent="save">
        <input-text
          id="session-title"
          label="Title"
          v-model="form.title"
          size="small"
        />
        <input-text-multiline
          id="session-description"
          label="Description"
          rows="2"
          v-model="form.description"
          size="small"
        />
        <input-date-time-local
          id="session-start"
          label="Start"
          v-model="form.startDate"
          size="small"
        />
        <input-date-time-local
          id="session-end"
          label="End"
          v-model="form.endDate"
          size="small"
        />
        <input-number
          id="session-capacity"
          :label="getMaximumRidersLabel"
          v-model="form.maximumRiders"
          size="small"
        />
        <input-toggle
          id="session-type"
          label="Type"
          off-label="Public"
          on-label="Private"
          v-model="form.type"
        />
      </form>
    </modal-body>
    <modal-footer>
      <button
        v-if="form.id == null"
        type="button"
        class="btn btn-outline-info"
        @click="save"
      >
        <i class="bi bi-check"></i>
        Create
      </button>
      <button
        v-if="form.id != null"
        type="button"
        class="btn btn-outline-info"
        @click="save"
      >
        <i class="bi bi-check"></i>
        Save
      </button>
      <button type="button" class="btn btn-outline-danger ml-1" @click="close">
        <i class="bi bi-x"></i>
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import Session, {
  SESSION_TYPE_OPEN,
  SESSION_TYPE_PRIVATE,
} from "booksys/dataTypes/session";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputNumber from "./forms/inputs/InputNumber.vue";
import InputToggle from "./forms/inputs/InputToggle.vue";
import InputTextMultiline from "./forms/inputs/InputTextMultiline.vue";
import InputDateTimeLocal from "./forms/inputs/InputDateTimeLocal.vue";

dayjs.extend(utc);
dayjs.extend(timezone);

const props = defineProps(["defaultValues", "visible"]);
const emit = defineEmits([
  "update:visible",
  "sessionCreatedHandler",
  "sessionEditedHandler",
]);

const store = useStore();

const title = ref("");
const errors = ref([]);
const form = ref({
  id: null,
  title: null,
  description: null,
  startDate: null,
  endDate: null,
  maximumRiders: null,
  type: null,
});

const getTimezone = computed(() => store.getters["configuration/getTimezone"]);
const getMaximumNumberOfRiders = computed(
  () => store.getters["configuration/getMaximumNumberOfRiders"],
);

const getMaximumRidersLabel = computed(() => {
  return form.value.id == null
    ? "Maximum Riders"
    : "Additional Slots for Riders";
});

const createSession = (session) =>
  store.dispatch("sessions/createSession", session);
const editSession = (session) =>
  store.dispatch("sessions/editSession", session);

function save(event) {
  event.preventDefault();
  console.log("SessionEditorModal, save:", form.value);

  const type = form.value.type ? SESSION_TYPE_PRIVATE : SESSION_TYPE_OPEN;

  const session = new Session(
    form.value.id,
    form.value.title,
    form.value.description,
    dayjs.tz(form.value.startDate, getTimezone.value).format(),
    dayjs.tz(form.value.endDate, getTimezone.value).format(),
    form.value.maximumRiders,
    type,
  );

  console.log("SessionEditorModal, save session dataType:", session);
  if (session.id == null) {
    createSession(session)
      .then((response) => {
        emit("sessionCreatedHandler", response.session_id);
        close();
      })
      .catch((err) => {
        console.warn("Error received:", err);
        errors.value = err;
      });
  } else {
    editSession(session)
      .then(() => {
        emit("sessionEditedHandler", session.id);
        close();
      })
      .catch((err) => {
        errors.value = err;
      });
  }
}

function close() {
  errors.value = [];
  emit("update:visible", false);
}

function setFormContent() {
  console.log("SessionEditorModal: Set defaults to:", props.defaultValues);

  form.value.id = props.defaultValues?.id || null;

  title.value = form.value.id != null ? "Edit Session" : "Create Session";

  form.value.title = props.defaultValues?.title || null;
  form.value.description = props.defaultValues?.description || null;

  form.value.startDate = props.defaultValues?.start
    ? dayjs(props.defaultValues.start).format("YYYY-MM-DDTHH:mm")
    : dayjs().tz(getTimezone.value).format("YYYY-MM-DDTHH:mm");
  form.value.endDate = props.defaultValues?.end
    ? dayjs(props.defaultValues.end).format("YYYY-MM-DDTHH:mm")
    : dayjs().tz(getTimezone.value).add(1, "hour").format("YYYY-MM-DDTHH:mm");

  form.value.maximumRiders =
    props.defaultValues?.maximumRiders || getMaximumNumberOfRiders.value;

  form.value.type = props.defaultValues?.type || null;

  console.log("SessionEditorModal: Form values are now:", form.value);
}

watch(
  () => props.defaultValues,
  () => {
    setFormContent();
  },
  { deep: true, immediate: true },
);
</script>
