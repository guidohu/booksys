<template>
  <div class="container text-left">
    <div class="row">
      <div class="col-12">
        <warning-box v-if="errors.length > 0" :errors="errors" />
        <form @submit.prevent="add">
          <input-text
            id="driver"
            label="Driver"
            v-model="form.driverName"
            size="small"
            :disabled="true"
            type="text"
          />
          <input-engine-hours
            id="before-hours"
            label="Before"
            :display-format="getEngineHourFormat"
            :disabled="disableBefore"
            v-model="form.beforeHours"
            :description="beforeDescription"
            placeholder="0"
            size="small"
          />
          <input-engine-hours
            v-if="showAfter"
            id="after-hours"
            label="After"
            :display-format="getEngineHourFormat"
            :disabled="!showAfter"
            v-model="form.afterHours"
            :description="afterDescription"
            placeholder="0"
            size="small"
          />
          <input-toggle
            id="type"
            label="Type"
            off-label="Private"
            on-label="Course"
            v-model="form.type"
          />
          <form-button
            type="submit"
            btn-style="info"
            btn-size="small"
            @click.prevent="add"
          >
            {{ showAfter ? "Check Out" : "Check In" }}
          </form-button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import InputEngineHours from "booksys/components/forms/inputs/InputEngineHours.vue";
import InputText from "booksys/components/forms/inputs/InputText.vue";
import InputToggle from "booksys/components/forms/inputs/InputToggle.vue";
import FormButton from "booksys/components/forms/FormButton.vue";
import WarningBox from "booksys/components/WarningBox.vue";

const store = useStore();

const form = ref({
  driverId: null,
  driverName: "",
  beforeHours: null,
  afterHours: null,
  type: false,
});

const beforeDescription = ref(null);
const afterDescription = ref(null);
const disableBefore = ref(false);
const showAfter = ref(false);
const errors = ref([]);

const getEngineHourLogLatest = computed(() => store.getters["boat/getEngineHourLogLatest"]);
const getMyNautiqueEngineHours = computed(() => store.getters["boat/getMyNautiqueEngineHours"]);
const userInfo = computed(() => store.getters["login/userInfo"]);
const getEngineHourFormat = computed(() => store.getters["configuration/getEngineHourFormat"]);
const getMyNautiqueEnabled = computed(() => store.getters["configuration/getMyNautiqueEnabled"]);
const getMyNautiqueBoatId = computed(() => store.getters["configuration/getMyNautiqueBoatId"]);

watch(userInfo, () => {
  setDriver();
});

watch(getEngineHourLogLatest, (newData) => {
  console.log("engineHourLatest just changed", newData);
  if (newData.id === 0 && parseFloat(newData.before_hours) === 0.0) {
    form.value.beforeHours = null;
    beforeDescription.value = null;
    prefillBefore();
    form.value.afterHours = null;
  } else if (parseFloat(newData.before_hours) > 0 && parseFloat(newData.after_hours) > 0) {
    form.value.beforeHours = null;
    form.value.type = false;
    prefillBefore();
    form.value.afterHours = null;
  } else {
    form.value.beforeHours = newData.before_hours;
    form.value.afterHours = null;
    prefillAfter();
  }

  setDriver();
  setDisableBefore();
  setShowAfter();
});

watch(getMyNautiqueBoatId, (boatId) => {
  queryMyNautiqueInfo(boatId);
});

watch(getMyNautiqueEngineHours, () => {
  if (disableBefore.value == false) {
    prefillBefore();
  } else {
    prefillAfter();
  }
});

const queryEngineHourLogLatest = () => store.dispatch("boat/queryEngineHourLogLatest");
const addEngineHours = (data) => store.dispatch("boat/addEngineHours", data);
const queryMyNautiqueInfo = (boatId) => store.dispatch("boat/queryMyNautiqueInfo", boatId);
const queryConfiguration = () => store.dispatch("configuration/queryConfiguration");

function prefillBefore() {
  if (getMyNautiqueEnabled.value && getMyNautiqueEngineHours.value != null) {
    form.value.beforeHours = getMyNautiqueEngineHours.value;
    beforeDescription.value = "prefilled by myNautique";
    afterDescription.value = null;
  }
}

function prefillAfter() {
  if (getMyNautiqueEnabled.value && getMyNautiqueEngineHours.value != null) {
    form.value.afterHours = getMyNautiqueEngineHours.value;
    beforeDescription.value = null;
    afterDescription.value = "prefilled by myNautique";
  }
}

function setDisableBefore() {
  if (getEngineHourLogLatest.value == null) {
    disableBefore.value = false;
    return;
  } else if (parseFloat(getEngineHourLogLatest.value.after_hours) > 0) {
    disableBefore.value = false;
    return;
  }
  disableBefore.value = true;
}

function setShowAfter() {
  if (getEngineHourLogLatest.value != null && parseFloat(getEngineHourLogLatest.value.after_hours) === 0.0) {
    showAfter.value = true;
    return;
  }
  showAfter.value = false;
}

function setDriver() {
  if (getEngineHourLogLatest.value != null && parseFloat(getEngineHourLogLatest.value.after_hours) === 0.0) {
    form.value.driverName = `${getEngineHourLogLatest.value.user_first_name} ${getEngineHourLogLatest.value.user_last_name}`;
    form.value.driverId = getEngineHourLogLatest.value.user_id;
  } else if (userInfo.value != null) {
    form.value.driverName = `${userInfo.value.first_name} ${userInfo.value.last_name}`;
    form.value.driverId = userInfo.value.id;
  } else {
    form.value.driverName = "unknown";
    form.value.driverId = null;
  }
}

function add() {
  const type = form.value.type ? 2 : 1;

  const data = {
    user_id: form.value.driverId,
    engine_hours_before: form.value.beforeHours,
    engine_hours_after: form.value.afterHours,
    type: type,
  };
  addEngineHours(data)
    .then(() => {
      errors.value = [];
    })
    .catch((errs) => (errors.value = errs));
}

queryConfiguration();
setDriver();
queryEngineHourLogLatest();
if (getMyNautiqueEnabled.value && getMyNautiqueBoatId.value) {
  queryMyNautiqueInfo(getMyNautiqueBoatId.value);
}
setDisableBefore();
setShowAfter();
</script>
