<template>
  <modal-container
    name="rider-selection-modal"
    :visible="visible"
    @hide="$emit('update:visible', false)"
    @show="$emit('update:visible', true)"
  >
    <modal-header
      :closable="true"
      title="Add Riders"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length" :errors="errors" />
      <form @submit.prevent="add">
        <show-for-desktop>
          <input-text id="rider-search" label="Search" v-model="search" />
        </show-for-desktop>
        <input-select
          id="rider-selection"
          label="Rider"
          v-model="selected"
          :options="filteredOptions"
          size="small"
          :select-size="5"
          select-mode="multiple"
          description="use ctrl+click to select multiple users"
        />
        <show-for-desktop>
          <div class="row">
            <div class="col-9 offset-3">
              <button
                type="button"
                class="btn btn-outline-success block"
                @click="add"
              >
                <i class="bi bi-person-plus"></i>{{ " " }}Select
              </button>
            </div>
          </div>
          <div class="row my-2">
            <hr />
          </div>
          <div class="row">
            <label class="col-3 col-form-label">Selected Riders</label>
            <div class="col-9 align-middle">
              <ul v-if="usersToAdd.length > 0">
                <li v-for="u in usersToAdd" :key="u.id">
                  {{ u.firstName + " " + u.lastName + "  " }}
                  <a href="#" @click.prevent="remove(u.id)">
                    <i class="bi bi-person-dash"></i>
                  </a>
                </li>
              </ul>
              <ul v-else>
                <li>Please select riders to be added</li>
              </ul>
            </div>
          </div>
        </show-for-desktop>
      </form>
    </modal-body>
    <modal-footer>
      <show-for-desktop>
        <button type="button" class="btn btn-outline-info" @click="save">
          <i class="bi bi-check" />
          Add
        </button>
      </show-for-desktop>
      <show-for-mobile>
        <button type="button" class="btn btn-outline-info" @click="saveMobile">
          <i class="bi bi-check" />
          Add
        </button>
      </show-for-mobile>
      <button type="button" class="btn btn-outline-danger mb-1" @click="close">
        <i class="bi bi-x" />
        Cancel
      </button>
    </modal-footer>
  </modal-container>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import uniq from "lodash/uniq";
import WarningBox from "@/components/WarningBox";
import { UserPointer } from "@/dataTypes/user";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import ShowForDesktop from "./bricks/ShowForDesktop.vue";
import ShowForMobile from "./bricks/ShowForMobile.vue";

export default {
  name: "RiderSelectionModal",
  components: {
    WarningBox,
    ModalContainer,
    ModalHeader,
    ModalFooter,
    ModalBody,
    InputText,
    InputSelect,
    ShowForDesktop,
    ShowForMobile,
  },
  props: ["session", "visible"],
  data() {
    return {
      errors: [],
      userIds: [],
      form: {},
      selected: [],
      search: "",
      filteredOptions: [],
      usersToAdd: [],
    };
  },
  computed: {
    ...mapGetters("user", ["userList"]),
    userOptions: function () {
      let users = [];
      this.userList.forEach((u) => {
        users.push({
          value: u.id,
          text: u.firstName + " " + u.lastName,
        });
      });

      users = users.filter(
        (u) => !this.usersToAdd.map((uta) => uta.id).includes(u.value)
      );
      return users;
    },
  },
  watch: {
    search: function (newSearch) {
      // filteredOptions are the ones that
      this.filteredOptions = this.userOptions.filter((u) =>
        u.text.toLowerCase().includes(newSearch.toLowerCase())
      );
    },
    userOptions: function (newOptions) {
      this.filteredOptions = newOptions;
    },
    selected: function (newSelected) {
      console.log("newSelected", newSelected);
    },
    usersToAdd: function (newUsersToAdd) {
      console.log("usersToAdd", newUsersToAdd);
    },
  },
  methods: {
    ...mapActions("user", ["queryUserList"]),
    ...mapActions("sessions", ["addUsersToSession"]),
    add: function () {
      console.log("add Called");
      console.log("selected", this.selected);

      if (this.selected.length > 0) {
        console.log("Add selected");
        this.usersToAdd.push(
          ...this.userList.filter((u) => this.selected.includes(u.id))
        );
      } else if (
        this.selected.length == 0 &&
        this.filteredOptions.length == 1
      ) {
        this.usersToAdd.push(
          ...this.userList.filter((u) => u.id == this.filteredOptions[0].value)
        );
      }

      // de-duplicate selection
      this.usersToAdd = uniq(this.usersToAdd);
    },
    remove: function (id) {
      this.usersToAdd = this.usersToAdd.filter((u) => u.id != id);
      console.log("Removed user from the selection:", id);
    },
    save: function () {
      console.log(this.session);
      this.addUsersToSession({
        sessionId: this.session.id,
        users: this.usersToAdd,
      })
        .then(() => this.close())
        .catch((errs) => (this.errors = errs));
    },
    saveMobile: function () {
      console.log(this.selected);
      this.addUsersToSession({
        sessionId: this.session.id,
        users: this.selected.map((i) => new UserPointer(Number(i))),
      })
        .then(() => this.close())
        .catch((errors) => (this.errors = errors));
    },
    close: function () {
      this.usersToAdd = [];
      this.$emit("update:visible", false);
    },
  },
  created() {
    this.queryUserList()
      .then(() => {
        console.log("userList received");
      })
      .catch((error) => {
        console.log("errors");
        this.errors = error;
      });
  },
};
</script>
