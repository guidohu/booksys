<template>
  <modal-container
    name="user-details-modal"
    :visible="visible"
    :inert="!visible"
  >
    <modal-header
      :closable="true"
      :title="title"
      @close="$emit('update:visible', false)"
    />
    <modal-body>
      <warning-box v-if="errors.length" :errors="errors" />
      <div v-if="details == null" class="row">
        <div class="col-12">No user selected.</div>
      </div>
      <div v-else>
        <h6 class="mt-2 mb-3">Profile</h6>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Name</label>
          <div class="col-9 col-form-label">{{ fullName }}</div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Username</label>
          <div class="col-9 col-form-label">{{ orDash(details.username) }}</div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Email</label>
          <div class="col-9 col-form-label">{{ orDash(details.email) }}</div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Mobile</label>
          <div class="col-9 col-form-label">{{ orDash(details.mobile) }}</div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Street / Nr</label>
          <div class="col-9 col-form-label">{{ orDash(details.address) }}</div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Zip / City</label>
          <div class="col-9 col-form-label">{{ zipAndCity }}</div>
        </div>
        <div class="row mb-2" v-if="details.comment">
          <label class="col-3 col-form-label">Comment</label>
          <div class="col-9 col-form-label">{{ details.comment }}</div>
        </div>

        <h6 class="mt-4 mb-3">Membership</h6>
        <input-select
          v-if="userGroupList.length > 0"
          id="user-details-group"
          label="User Group"
          v-model="groupId"
          :options="userGroupList"
          size="small"
          @changed="groupChanged"
        />
        <div v-else class="row mb-2">
          <label class="col-3 col-form-label">User Group</label>
          <div class="col-9 col-form-label">-</div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Price</label>
          <div class="col-9 col-form-label">{{ pricePerMinute }}</div>
        </div>
        <input-toggle
          id="user-details-locked"
          label="Account"
          v-model="locked"
          offLabel="Unlocked"
          onLabel="Locked"
          @change="lockedChanged"
        />
        <div class="row mb-2">
          <label class="col-3 col-form-label">Driving License</label>
          <div class="col-9 col-form-label">
            <i v-if="details.license" class="bi bi-patch-check text-success" />
            {{ details.license ? "Authorized to drive boats" : "No license" }}
          </div>
        </div>

        <h6 class="mt-4 mb-3">Riding and Balance</h6>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Riding Time</label>
          <div class="col-9 col-form-label">{{ ridingTime }}</div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Riding Cost</label>
          <div class="col-9 col-form-label">
            {{ currency(details.total_heat_cost) }}
          </div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Payments</label>
          <div class="col-9 col-form-label">
            {{ currency(details.total_payment) }}
          </div>
        </div>
        <div class="row mb-2">
          <label class="col-3 col-form-label">Balance</label>
          <div class="col-9 col-form-label">{{ currency(balance) }}</div>
        </div>
      </div>
    </modal-body>
    <modal-footer>
      <button type="button" class="btn btn-outline-info" @click="close">
        <i class="bi bi-check"></i>
        OK
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useStore } from "vuex";
import WarningBox from "booksys/components/WarningBox.vue";
import ModalContainer from "./bricks/ModalContainer.vue";
import ModalHeader from "./bricks/ModalHeader.vue";
import ModalBody from "./bricks/ModalBody.vue";
import ModalFooter from "./bricks/ModalFooter.vue";
import InputSelect from "./forms/inputs/InputSelect.vue";
import InputToggle from "./forms/inputs/InputToggle.vue";
import { formatCurrency, formatDurationString } from "booksys/libs/formatters";

const props = defineProps(["visible", "user"]);
const emit = defineEmits(["update:visible"]);

const store = useStore();

// The user that is being shown. Changing the group or the lock state reloads
// the user list, which drops the selection in the table behind this dialog,
// so the dialog works on its own copy instead of on the selected row.
const details = ref(null);
const groupId = ref(null);
const locked = ref(false);
const errors = ref([]);

const userGroups = computed(() => store.getters["user/userGroups"]);
const getCurrency = computed(() => store.getters["configuration/getCurrency"]);

const setUserGroup = (data) => store.dispatch("user/setUserGroup", data);
const lockUser = (data) => store.dispatch("user/lockUser", data);

const userGroupList = computed(() =>
  userGroups.value.map((ug) => {
    return {
      value: ug.user_group_id,
      text: ug.user_group_name,
    };
  }),
);

const selectedGroup = computed(() =>
  userGroups.value.find((ug) => ug.user_group_id == groupId.value),
);

const pricePerMinute = computed(() => {
  if (selectedGroup.value == null) {
    return "-";
  }
  return formatCurrency(
    Number(selectedGroup.value.price_min),
    getCurrency.value + "/min",
  );
});

const title = computed(() => {
  if (details.value == null) {
    return "User";
  }
  return fullName.value;
});

const fullName = computed(() => {
  if (details.value == null) {
    return "-";
  }
  return orDash(
    [details.value.first_name, details.value.last_name]
      .filter((n) => n)
      .join(" "),
  );
});

const zipAndCity = computed(() => {
  if (details.value == null) {
    return "-";
  }
  // A zip code of 0 is what the backend sends when none is set.
  const zip = details.value.plz ? String(details.value.plz) : null;
  return orDash([zip, details.value.city].filter((p) => p).join(" "));
});

const balance = computed(() => {
  if (details.value == null) {
    return 0;
  }
  return (
    Number(details.value.total_payment) - Number(details.value.total_heat_cost)
  );
});

const ridingTime = computed(() => {
  if (details.value == null) {
    return "-";
  }
  return formatDurationString(Number(details.value.total_heat_seconds), false);
});

watch(
  () => props.visible,
  (isVisible) => {
    if (isVisible) {
      load();
    }
  },
);

function load() {
  details.value = props.user != null ? { ...props.user } : null;
  groupId.value = props.user != null ? props.user.status : null;
  locked.value = props.user != null && props.user.locked ? true : false;
  errors.value = [];
}

function groupChanged() {
  setUserGroup({
    userId: details.value.id,
    userGroupId: groupId.value,
  })
    .then(() => {
      errors.value = [];
      details.value.status = groupId.value;
    })
    .catch((errs) => {
      errors.value = errs;
      // Keep the dialog in sync with what is stored.
      groupId.value = details.value.status;
    });
}

function lockedChanged() {
  lockUser({
    user: details.value.id,
    locked: locked.value,
  })
    .then(() => {
      errors.value = [];
      details.value.locked = locked.value;
    })
    .catch((errs) => {
      errors.value = errs;
      locked.value = details.value.locked ? true : false;
    });
}

function currency(value) {
  return formatCurrency(Number(value), getCurrency.value);
}

function orDash(value) {
  return value != null && value !== "" ? value : "-";
}

function close() {
  emit("update:visible", false);
}
</script>
