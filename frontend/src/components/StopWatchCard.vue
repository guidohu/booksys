<template>
  <card-module nobody>
    <heat-comment-modal
      v-model:visible="isVisibleHeatCommentModal"
      :inert="!isVisibleHeatCommentModal"
      :default-comment="comment"
      @commentChangeHandler="changeComment"
    />
    <div class="row" v-if="errors.length > 0">
      <warning-box :errors="errors" />
    </div>
    <div class="row mx-1 my-4">
      <div class="col-12">
        <card-module nobody>
          <div class="row stopwatch-text text-center">
            <div class="col-12">
              {{ getDisplayTime }}
            </div>
          </div>
        </card-module>
      </div>
    </div>
    <div class="row mx-1" v-if="comment != null">
      <div class="col-12">
        <div class="input-group">
          <input type="text" class="form-control" v-model="comment" />
          <button class="btn btn-outline-info" @click="showHeatCommentModal">
            <i class="bi bi-pencil-square" />
          </button>
        </div>
      </div>
    </div>
    <div class="row mx-1">
      <input-select
        id="rider-select"
        v-model="selectedRiderId"
        :options="selectableRiders"
        :disabled="getIsRunning || sessionId == null"
      />
    </div>
    <div class="row mx-1 mt-2 mb-2 text-center" v-if="selectedRiderId != null">
      <div class="col-12">
        <div class="btn-group" role="group" aria-label="btn group">
          <button
            v-if="!getIsRunning"
            class="btn btn-outline-info btn-lg"
            @click="navigateBack"
          >
            <i class="bi bi-arrow-left" />
            Back
          </button>
          <button
            v-if="!getIsRunning"
            class="btn btn-outline-info btn-lg"
            @click="startTakingTime"
          >
            <i class="bi bi-play" />
            Start
          </button>
          <button
            v-if="getIsRunning && !getIsPaused"
            class="btn btn-outline-info btn-lg"
            @click="pauseTakingTime"
          >
            <i class="bi bi-pause" />
            Pause
          </button>
          <button
            v-if="getIsPaused"
            class="btn btn-outline-info btn-lg"
            @click="resumeTakingTime"
          >
            <i class="bi bi-play" />
            Resume
          </button>
          <button
            v-if="getIsRunning"
            class="btn btn-outline-info btn-lg"
            @click="finish"
          >
            <i class="bi bi-check-square" />
            Finish
          </button>
          <button
            v-if="selectedRiderId != null || getIsRunning"
            class="btn btn-outline-info btn-lg dropdown-toggle"
            type="button"
            data-bs-toggle="dropdown"
            aria-expanded="false"
          ></button>
          <ul class="dropdown-menu">
            <li v-if="selectedRiderId != null || getIsRunning">
              <a
                class="dropdown-item"
                href="#"
                @click.prevent="showHeatCommentModal"
              >
                Add Comment
              </a>
            </li>
            <li v-if="getIsRunning && selectedRiderId != null">
              <hr class="dropdown-divider" />
            </li>
            <li v-if="getIsRunning">
              <a
                class="dropdown-item"
                href="#"
                @click.prevent="resetTakingTime"
              >
                Cancel Heat
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </card-module>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import WarningBox from "booksys/components/WarningBox.vue";
import HeatCommentModal from "booksys/components/HeatCommentModal.vue";
import CardModule from "./bricks/CardModule.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";

const props = defineProps(["sessionId", "visible"]);

const store = useStore();
const router = useRouter();

const isVisibleHeatCommentModal = ref(false);
const selectedRiderId = ref(null);
const selectableRiders = ref([]);
const comment = ref(null);
const displayTime = ref("00:00");
const errors = ref([]);

const getDisplayTime = computed(
  () => store.getters["stopwatch/getDisplayTime"],
);
const getIsPaused = computed(() => store.getters["stopwatch/getIsPaused"]);
const getIsRunning = computed(() => store.getters["stopwatch/getIsRunning"]);
const getUserId = computed(() => store.getters["stopwatch/getUserId"]);
const getSession = computed(() => store.getters["sessions/getSession"]);

watch(getUserId, (selectedUserId) => {
  if (selectedUserId != selectedRiderId.value) {
    selectedRiderId.value = selectedUserId;
  }
});

watch(selectedRiderId, (newRiderId) => {
  setUserId(newRiderId);
});

watch(getSession, (newSession) => {
  console.log("session data changed to", newSession);
  if (newSession != null) {
    setSelectableRiders(newSession.riders);
  }
});

function setSelectableRiders(riders) {
  console.log("set selectable riders", riders);
  if (riders != null && riders.length > 0) {
    selectableRiders.value = [
      {
        value: null,
        text: "Please select a rider",
      },
    ];
    const ridersToAdd = riders.map((r) => {
      return {
        value: r.id,
        text: r.first_name + " " + r.last_name,
      };
    });
    console.log("ridersToAdd", ridersToAdd);
    selectableRiders.value.push(...ridersToAdd);
    console.log("new selectableRiders", selectableRiders.value);
  }
}

function showHeatCommentModal() {
  isVisibleHeatCommentModal.value = true;
}

function changeComment(newComment) {
  comment.value = newComment;
  setComment(comment.value);
}

function finish() {
  finishTakingTime().catch((errs) => (errors.value = errs));
}

const startTakingTime = () => store.dispatch("stopwatch/startTakingTime");
const pauseTakingTime = () => store.dispatch("stopwatch/pauseTakingTime");
const resumeTakingTime = () => store.dispatch("stopwatch/resumeTakingTime");
const finishTakingTime = () => store.dispatch("stopwatch/finishTakingTime");
const resetTakingTime = () => store.dispatch("stopwatch/resetTakingTime");
const setSessionId = (id) => store.dispatch("stopwatch/setSessionId", id);
const setUserId = (id) => store.dispatch("stopwatch/setUserId", id);
const setComment = (comment) => store.dispatch("stopwatch/setComment", comment);
const querySession = (id) => store.dispatch("sessions/querySession", id);

function navigateBack() {
  router.push("/today");
}

if (props.sessionId == null) {
  errors.value = ["No session selected"];
} else {
  setSessionId(props.sessionId);
  querySession(props.sessionId)
    .then(() => console.log("queried session"))
    .catch((errs) => (errors.value = errs));
}

// select rider (if one is stored)
const storedUserId = getUserId.value;
console.log("Stored UserId is:", storedUserId);
if (storedUserId != null) {
  selectedRiderId.value = storedUserId;
}
</script>

<style scoped>
.stopwatch-text {
  font-size: 4rem;
}
</style>
