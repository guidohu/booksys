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

<script>
import { mapActions } from "vuex";
import * as dayjs from "dayjs";
import { UserPointer } from "@/dataTypes/user";
import RiderSelectionModal from "@/components/RiderSelectionModal";
import SectionedCardModule from "@/components/bricks/SectionedCardModule.vue";

export default {
  name: "SessionDetailsCard",
  components: {
    SectionedCardModule,
    RiderSelectionModal,
  },
  props: ["date", "session"],
  data() {
    return {
      showRiderSelectionModal: false,
      sessionSelected: false,
    };
  },
  computed: {
    dateString: function () {
      const newDate = this.date;
      return dayjs(newDate).format("DD.MM.YYYY");
    },
    timeString: function () {
      if (this.session == null) {
        return "no session selected";
      }

      const start = this.session.start;
      const end = this.session.end;
      return dayjs(start).format("HH:mm") + " - " + dayjs(end).format("HH:mm");
    },
    showCreateSession: function () {
      if (this.session != null && this.session.id == null) {
        return true;
      }
      return false;
    },
    showDeleteSession: function () {
      if (this.session != null && this.session.id != null) {
        return true;
      }
      return false;
    },
    showAddRiders: function () {
      if (this.session != null && this.session.id != null) {
        return true;
      }
      return false;
    },
    showRiders: function () {
      if (
        this.session != null &&
        this.session.riders != null &&
        this.session.riders.length > 0
      ) {
        return true;
      }
      return false;
    },
  },
  watch: {
    timeString: function () {
      if (this.session == null) {
        this.sessionSelected = false;
      } else {
        this.sessionSelected = true;
      }
    },
  },
  methods: {
    createSession: function () {
      this.$emit("createSessionHandler");
    },
    editSession: function () {
      this.$emit("editSessionHandler");
    },
    deleteSession: function () {
      this.$emit("deleteSessionHandler", { id: this.session.id });
    },
    addRiders: function () {
      this.showRiderSelectionModal = true;
    },
    ...mapActions("sessions", ["deleteUserFromSession"]),
    removeRider: function (id) {
      this.deleteUserFromSession({
        sessionId: this.session.id,
        user: new UserPointer(id),
      }).catch((errors) => {
        console.error("Cannot delete user:", errors);
      });
    },
  },
};
</script>

<style scoped>
.btn-xs {
  padding: 0.2rem 0.4rem;
  font-size: 0.6rem;
  line-height: 1;
  border-radius: 0.2rem;
}
</style>
