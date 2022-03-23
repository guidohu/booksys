<template>
  <card-module nobody>
    <heat-comment-modal
      v-model:visible="isVisibleHeatCommentModal"
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

<script>
import { mapGetters, mapActions } from "vuex";
import WarningBox from "@/components/WarningBox";
import HeatCommentModal from "@/components/HeatCommentModal";
import CardModule from "./bricks/CardModule.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";

export default {
  name: "StopWatchCard",
  components: {
    WarningBox,
    HeatCommentModal,
    InputSelect,
    CardModule,
  },
  props: ["sessionId", "visible"],
  data() {
    return {
      isVisibleHeatCommentModal: false,
      selectedRiderId: null,
      selectableRiders: [],
      comment: null,
      displayTime: "00:00",
      errors: [],
    };
  },
  computed: {
    ...mapGetters("stopwatch", [
      "getDisplayTime",
      "getIsPaused",
      "getIsRunning",
      "getUserId",
    ]),
    ...mapGetters("sessions", ["getSession"]),
  },
  watch: {
    getUserId: function (selectedUserId) {
      if (selectedUserId != this.selectedRiderId) {
        this.selectedRiderId = selectedUserId;
      }
    },
    selectedRiderId: function (newRiderId) {
      this.setUserId(newRiderId);
    },
    getSession: function (newSession) {
      console.log("session data changed to", newSession);
      if (newSession != null) {
        this.setSelectableRiders(newSession.riders);
      }
    },
  },
  methods: {
    setSelectableRiders: function (riders) {
      console.log("set selectable riders", riders);
      if (riders != null && riders.length > 0) {
        this.selectableRiders = [
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
        this.selectableRiders.push(...ridersToAdd);
        console.log("new selectableRiders", this.selectableRiders);
      }
    },
    showHeatCommentModal: function () {
      this.isVisibleHeatCommentModal = true;
    },
    changeComment: function (newComment) {
      this.comment = newComment;
      this.setComment(this.comment);
    },
    finish: function () {
      this.finishTakingTime().catch((errors) => (this.errors = errors));
    },
    ...mapActions("stopwatch", [
      "startTakingTime",
      "pauseTakingTime",
      "resumeTakingTime",
      "finishTakingTime",
      "resetTakingTime",
      "setSessionId",
      "setUserId",
      "setComment",
    ]),
    ...mapActions("sessions", ["querySession"]),
    navigateBack: function () {
      this.$router.push("/today");
    },
  },
  created() {
    if (this.sessionId == null) {
      this.errors = ["No session selected"];
    } else {
      this.setSessionId(this.sessionId);
      this.querySession(this.sessionId)
        .then(() => console.log("queried session"))
        .catch((errors) => (this.errors = errors));
    }

    // select rider (if one is stored)
    const storedUserId = this.getUserId;
    console.log("Stored UserId is:", storedUserId);
    if (storedUserId != null) {
      this.selectedRiderId = storedUserId;
    }
  },
};
</script>

<style scoped>
.stopwatch-text {
  font-size: 4rem;
}
</style>
