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
      <warning-box v-if="errors.length" :erors="errors" />
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

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import Session, {
  SESSION_TYPE_OPEN,
  SESSION_TYPE_PRIVATE,
} from "@/dataTypes/session";
import * as dayjs from "dayjs";
import * as dayjsUTC from "dayjs/plugin/utc";
import * as dayjsTimezone from "dayjs/plugin/timezone";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputNumber from "./forms/inputs/InputNumber.vue";
import InputToggle from "./forms/inputs/InputToggle.vue";
import InputTextMultiline from "./forms/inputs/InputTextMultiline.vue";
import InputDateTimeLocal from "./forms/inputs/InputDateTimeLocal.vue";

dayjs.extend(dayjsUTC);
dayjs.extend(dayjsTimezone);

export default {
  name: "SessionEditorModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalBody,
    ModalFooter,
    InputDateTimeLocal,
    InputText,
    InputTextMultiline,
    InputNumber,
    InputToggle,
  },
  props: ["defaultValues", "visible"],
  data() {
    return {
      title: "",
      errors: [],
      form: {
        id: null,
        title: null,
        description: null,
        startDate: null,
        endDate: null,
        maximumRiders: null,
        type: null,
      },
    };
  },
  computed: {
    getMaximumRidersLabel: function () {
      if (this.form.id == null) {
        return "Maximum Riders";
      } else {
        return "Additional Slots for Riders";
      }
    },
    ...mapGetters("configuration", ["getTimezone", "getMaximumNumberOfRiders"]),
  },
  methods: {
    ...mapActions("sessions", ["createSession", "editSession"]),
    save: function (event) {
      event.preventDefault();
      console.log("SessionEditorModal, save:", this.form);

      const type =
        this.form.type == true ? SESSION_TYPE_PRIVATE : SESSION_TYPE_OPEN;

      const session = new Session(
        this.form.id,
        this.form.title,
        this.form.description,
        dayjs.tz(this.form.startDate, this.getTimezone).format(),
        dayjs.tz(this.form.endDate, this.getTimezone).format(),
        this.form.maximumRiders,
        type
      );

      console.log("SessionEditorModal, save session dataType:", session);
      if (session.id == null) {
        this.createSession(session)
          .then((response) => {
            this.$emit("sessionCreatedHandler", response.session_id);
            this.close();
          })
          .catch((err) => {
            this.errors = err;
          });
      } else {
        this.editSession(session)
          .then(() => {
            this.$emit("sessionEditedHandler", session.id);
            this.close();
          })
          .catch((err) => {
            this.errors = err;
          });
      }
    },
    close: function () {
      this.errors = [];
      this.$emit("update:visible", false);
    },
    setFormContent: function () {
      console.log("SessionEditorModal: Set defaults to:", this.defaultValues);

      // set ID
      this.form.id =
        this.defaultValues != null && this.defaultValues.id != null
          ? this.defaultValues.id
          : null;

      if (this.form.id != null) {
        this.title = "Edit Session";
      } else {
        this.title = "Create Session";
      }

      // set title
      this.form.title =
        this.defaultValues != null && this.defaultValues.title != null
          ? this.defaultValues.title
          : null;

      // set description
      this.form.description =
        this.defaultValues != null && this.defaultValues.description != null
          ? this.defaultValues.description
          : null;

      // set times
      this.form.startDate =
        this.defaultValues != null && this.defaultValues.start != null
          ? dayjs(this.defaultValues.start).format("YYYY-MM-DDTHH:mm")
          : dayjs().tz(this.getTimezone).format("YYYY-MM-DDTHH:mm");
      this.form.endDate =
        this.defaultValues != null && this.defaultValues.end != null
          ? dayjs(this.defaultValues.end).format("YYYY-MM-DDTHH:mm")
          : dayjs()
              .tz(this.getTimezone)
              .add(1, "hour")
              .format("YYYY-MM-DDTHH:mm");

      // set maximum riders
      this.form.maximumRiders =
        this.defaultValues != null && this.defaultValues.maximumRiders != null
          ? this.defaultValues.maximumRiders
          : this.getMaximumNumberOfRiders;

      // set session type
      this.form.type =
        this.defaultValues != null && this.defaultValues.type != null
          ? this.defaultValues.type
          : null;

      console.log("SessionEditorModal: Form values are now:", this.form);
    },
  },
  watch: {
    defaultValues: function () {
      // set form content whenever the default props change
      this.setFormContent();
    },
  },
  created() {
    // set form content based on props
    this.setFormContent();
  },
};
</script>
