<template>
  <modal-container
    name="rider-selection-modal"
    v-bind:visible="visible"
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

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import uniq from "lodash/uniq";
import WarningBox from "booksys/components/WarningBox.vue";
import { UserPointer } from "booksys/dataTypes/user";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputText from "./forms/inputs/InputText.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import ShowForDesktop from "./bricks/ShowForDesktop.vue";
import ShowForMobile from "./bricks/ShowForMobile.vue";

const props = defineProps(["session", "visible"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

const errors = ref([]);
const userIds = ref([]);
const form = ref({});
const selected = ref([]);
const search = ref("");
const filteredOptions = ref([]);
const usersToAdd = ref([]);

const userList = computed(() => store.getters["user/userList"]);

const userOptions = computed(() => {
  let users = [];
  userList.value.forEach((u) => {
    users.push({
      value: u.id,
      text: u.firstName + " " + u.lastName,
    });
  });

  users = users.filter(
    (u) => !usersToAdd.value.map((uta) => uta.id).includes(u.value),
  );
  return users;
});

watch(search, (newSearch) => {
  // filteredOptions are the ones that
  filteredOptions.value = userOptions.value.filter((u) =>
    u.text.toLowerCase().includes(newSearch.toLowerCase()),
  );
});

watch(userOptions, (newOptions) => {
  filteredOptions.value = newOptions;
});

watch(selected, (newSelected) => {
  console.log("newSelected", newSelected);
});

watch(usersToAdd, (newUsersToAdd) => {
  console.log("usersToAdd", newUsersToAdd);
});

// The modal stays mounted between openings, so it has to start from a clean
// state every time it is shown. Otherwise the errors and the riders of the
// previous round are still there.
watch(
  () => props.visible,
  (isVisible) => {
    if (isVisible) {
      reset();
    }
  },
);

const queryUserList = () => store.dispatch("user/queryUserList");
const addUsersToSession = (data) =>
  store.dispatch("sessions/addUsersToSession", data);

function add() {
  console.log("add Called");
  console.log("selected", selected.value);

  if (selected.value.length > 0) {
    console.log(
      "Users to add",
      userList.value.filter((u) => selected.value.includes(u.id.toString())),
    );
    usersToAdd.value.push(
      ...userList.value.filter((u) => selected.value.includes(u.id.toString())),
    );
  } else if (selected.value.length == 0 && filteredOptions.value.length == 1) {
    usersToAdd.value.push(
      ...userList.value.filter((u) => u.id == filteredOptions.value[0].value),
    );
  }

  // de-duplicate selection
  usersToAdd.value = uniq(usersToAdd.value);
}

function remove(id) {
  usersToAdd.value = usersToAdd.value.filter((u) => u.id != id);
  console.log("Removed user from the selection:", id);
}

function save() {
  console.log(props.session);
  addUsersToSession({
    sessionId: props.session.id,
    users: usersToAdd.value,
  })
    .then(() => close())
    .catch((errs) => (errors.value = errs));
}

function saveMobile() {
  console.log(selected.value);
  addUsersToSession({
    sessionId: props.session.id,
    users: selected.value.map((i) => new UserPointer(Number(i))),
  })
    .then(() => close())
    .catch((errs) => (errors.value = errs));
}

function reset() {
  errors.value = [];
  selected.value = [];
  usersToAdd.value = [];
  // The search term belongs to the selection that is being cleared. Keeping it
  // would leave the search field and the rider list out of sync.
  search.value = "";
  filteredOptions.value = userOptions.value;
}

function close() {
  reset();
  emit("update:visible", false);
}

queryUserList()
  .then(() => {
    console.log("userList received");
  })
  .catch((error) => {
    console.log("errors");
    errors.value = error;
  });
</script>
