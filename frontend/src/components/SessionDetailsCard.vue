<template>
  <sectioned-card-module>
    <template v-slot:header>
      <div class="row">
        <div class="col-8">
          <h5 class="card-title pt-1">Session Details</h5>
        </div>
        <div class="col-4 text-end">
          <div v-if="showAddRiders || showDeleteSession" class="dropdown">
            <button
              class="btn btn-outline-info btn-sm dropdown-toggle"
              type="button"
              id="dropdownMenuButton1"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <i class="bi bi-list"></i>
            </button>
            <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1">
              <li v-if="showAddRiders">
                <a class="dropdown-item" href="#" @click.prevent="addRiders">
                  <i class="bi bi-person-plus"></i> {{ " " }} Add Rider
                </a>
              </li>
              <li v-if="showAddRiders">
                <a class="dropdown-item" href="#" @click.prevent="editSession">
                  <i class="bi bi-pencil-square"></i> {{ " " }} Edit Session
                </a>
              </li>
              <li><hr class="dropdown-divider" /></li>
              <li v-if="showDeleteSession">
                <a
                  class="dropdown-item"
                  href="#"
                  @click.prevent="deleteSession"
                >
                  <i class="bi bi-trash"></i> {{ " " }} Delete Session
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </template>
    <template v-slot:body>
      <rider-selection-modal
        v-model:visible="showRiderSelectionModal"
        :inert="!showRiderSelectionModal"
        :session="session"
      />
      <div v-if="session != null && session.title != null" class="row">
        <div class="col-4">Title</div>
        <div class="col-8">
          {{ session.title }}
        </div>
      </div>
      <div class="row">
        <div class="col-4">Date</div>
        <div class="col-8">
          {{ dateString }}
        </div>
      </div>
      <div
        v-if="
          session != null &&
          session.creator_first_name != null &&
          session.creator_last_name != null
        "
        class="row"
      >
        <div class="col-4">Host</div>
        <div class="col-8">
          {{ session.creator_first_name }} {{ session.creator_last_name }}
        </div>
      </div>
      <div class="row">
        <div class="col-4">Session</div>
        <div class="col-8">
          {{ timeString }}
        </div>
      </div>
      <div v-if="showRiders" class="row">
        <div class="col-4">Riders</div>
        <div class="col-8">
          <div v-for="rider in session.riders" :key="rider.id" class="row">
            <div class="col-10 text-truncate">
              {{ rider.name }}
            </div>
            <div class="col-2">
              <a href="#" @click.prevent="removeRider(rider.id)">
                <i class="bi bi-person-dash"></i>
              </a>
            </div>
          </div>
        </div>
      </div>
      <div v-if="showCreateSession" class="row">
        <div class="col-8 offset-4">
          <button
            type="button"
            class="btn btn-success btn-sm block"
            @click.stop="createSession"
          >
            Create Session
          </button>
        </div>
      </div>
      <div v-if="showAddRiders" class="row">
        <div class="col-8 offset-4">
          <button
            type="button"
            class="btn btn-success btn-sm block"
            @click.stop="addRiders"
          >
            Add Riders
          </button>
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import dayjs from "dayjs";
import { UserPointer } from "../dataTypes/user";
import RiderSelectionModal from "../components/RiderSelectionModal.vue";
import SectionedCardModule from "../components/bricks/SectionedCardModule.vue";

const props = defineProps(["date", "session"]);
const emit = defineEmits([
  "createSessionHandler",
  "editSessionHandler",
  "deleteSessionHandler",
]);

const store = useStore();

const showRiderSelectionModal = ref(false);
const sessionSelected = ref(false);

const dateString = computed(() => {
  const newDate = props.date;
  return dayjs(newDate).format("DD.MM.YYYY");
});

const timeString = computed(() => {
  if (props.session == null) {
    return "no session selected";
  }

  const start = props.session.start;
  const end = props.session.end;
  return dayjs(start).format("HH:mm") + " - " + dayjs(end).format("HH:mm");
});

const showCreateSession = computed(() => {
  if (props.session != null && props.session.id == null) {
    return true;
  }
  return false;
});

const showDeleteSession = computed(() => {
  if (props.session != null && props.session.id != null) {
    return true;
  }
  return false;
});

const showAddRiders = computed(() => {
  if (props.session != null && props.session.id != null) {
    return true;
  }
  return false;
});

const showRiders = computed(() => {
  if (
    props.session != null &&
    props.session.riders != null &&
    props.session.riders.length > 0
  ) {
    return true;
  }
  return false;
});

watch(timeString, () => {
  if (props.session == null) {
    sessionSelected.value = false;
  } else {
    sessionSelected.value = true;
  }
});

function createSession() {
  emit("createSessionHandler");
}

function editSession() {
  emit("editSessionHandler");
}

function deleteSession() {
  emit("deleteSessionHandler", { id: props.session.id });
}

function addRiders() {
  showRiderSelectionModal.value = true;
}

const deleteUserFromSession = (data) =>
  store.dispatch("sessions/deleteUserFromSession", data);

function removeRider(id) {
  deleteUserFromSession({
    sessionId: props.session.id,
    user: new UserPointer(id),
  }).catch((errors) => {
    console.error("Cannot delete user:", errors);
  });
}
</script>

<style scoped>
.btn-xs {
  padding: 0.2rem 0.4rem;
  font-size: 0.6rem;
  line-height: 1;
  border-radius: 0.2rem;
}
</style>
