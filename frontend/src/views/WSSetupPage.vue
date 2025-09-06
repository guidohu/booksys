<template>
  <modal-container name="setup-modal" :visible="true">
    <modal-header v-if="setupStep!=3" :title="setupSteps[setupStep].title" />
    <modal-body>
      <div class="progress mb-3" style="height: 2px">
        <div
          class="progress-bar bg-info"
          role="progressbar"
          :style="progress"
        ></div>
      </div>
      <warning-box v-if="errors.length > 0 && setupStep != 3" class="mt-4" :errors="errors" />
      <!-- Database setup -->
      <database-configuration
        v-if="setupSteps[setupStep].name == 'db'"
        :dbconfig="dbConfig"
        @config-change="dbConfigInputHandler"
        @save="setDbSettings"
      />
      <!-- Administrator Account Setup -->
      <div v-if="setupSteps[setupStep].name == 'administrator'">
        <div class="row mb-4">
          <div class="col-12">Setup the Administrator Account</div>
        </div>
        <user-sign-up
          :user-data="adminUserConfig"
          :show-disclaimer="false"
          @save="addAdminUser"
          @update:user="handleUserUpdate"
        />
      </div>
      <!-- Setup Done -->
      <div v-if="setupSteps[setupStep].name == 'done'" class="row text-center">
        <div class="col-12">
          <p class="h4 mb-2">
            <i class="bi bi-check-circle text-success" />
            Setup Done.
          </p>
          <p>Please go back to the login page and login.</p>
        </div>
      </div>
      <!-- Error -->
      <div v-if="setupSteps[setupStep].name == 'error'" class="row text-center">
        <div class="col-12">
          <p class="h4 mb-2">
            <i class="bi bi-x text-danger" />
            Error.
          </p>
          <warning-box :errors="errors" :dismissible="false" :hide-label="true"/>
        </div>
      </div>
    </modal-body>
    <modal-footer>
      <button
        v-if="setupSteps[setupStep].name != 'done' && setupStep != 3"
        class="btn btn-outline-danger me-1"
        type="button"
        @click="close"
      >
        <i class="bi bi-x" />
        Cancel
      </button>
      <button
        v-if="setupSteps[setupStep].name == 'db'"
        class="btn btn-outline-info"
        type="button"
        :disabled="isLoading"
        @click="setDbSettings"
      >
        <i class="bi bi-arrow-right" />
        Next
      </button>
      <button
        v-if="setupSteps[setupStep].name == 'administrator'"
        type="button"
        class="btn btn-outline-info"
        :disabled="isLoading"
        @click="addAdminUser"
      >
        <i class="bi bi-arrow-right" />
        Next
      </button>
      <button
        v-if="setupSteps[setupStep].name == 'done'"
        class="btn btn-outline-info me-1"
        type="button"
        variant="outline-info"
        @click="close"
      >
        <i class="bi bi-check" />
        Done
      </button>
      <button
        v-if="setupStep == 3"
        class="btn btn-outline-info me-1"
        type="button"
        variant="outline-success"
        @click="close"
      >
        <i class="bi bi-check" />
        OK
      </button>
    </modal-footer>
  </modal-container>
</template>

<script setup>
import { defineAsyncComponent, ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { getBackendStatus } from "../api/backend";
import Configuration from "../api/configuration";
import User from "../api/user";
import ModalContainer from "../components/bricks/ModalContainer.vue";
import ModalHeader from "../components/bricks/ModalHeader.vue";
import ModalBody from "../components/bricks/ModalBody.vue";
import ModalFooter from "../components/bricks/ModalFooter.vue";

// Lazy loaded components
const DatabaseConfiguration = defineAsyncComponent(
  () => import("../components/DatabaseConfiguration.vue"),
);
const UserSignUp = defineAsyncComponent(
  () => import("../components/forms/UserSignUp.vue"),
);
const WarningBox = defineAsyncComponent(
  () => import("../components/WarningBox.vue"),
);

const router = useRouter();

const errors = ref([]);
const isLoading = ref(false);
const dbConfig = ref({});
const adminUserConfig = ref({});
const setupStep = ref(0);
const setupSteps = ref([
  {
    id: 0,
    name: "db",
    title: "Setup Database",
    progress: 33,
  },
  {
    id: 1,
    name: "administrator",
    title: "Setup User",
    progress: 66,
  },
  {
    id: 2,
    name: "done",
    title: "Setup Done",
    progress: 100,
  },
  {
    id: 3,
    name: "error",
    title: "Error",
    progress: 0,
  }
]);

const progress = computed(() => {
  const step = setupStep.value
  let value = "width: ";
  value += parseInt(setupSteps.value[step].progress);
  value += "%";
  return value;
});

const dbConfigInputHandler = (config) => {
  dbConfig.value = config;
};

const setDbSettings = () => {
  console.log("setDB called:", dbConfig.value);
  isLoading.value = true;

  Configuration.setDbConfig(dbConfig.value)
    .then(() => {
      isLoading.value = false;
      getBackendStatusInternal();
    })
    .catch((err) => {
      errors.value = err;
      isLoading.value = false;
    });
};

const handleUserUpdate = (u) => {
  adminUserConfig.value.email = u.email;
  adminUserConfig.value.password = u.password;
  adminUserConfig.value.passwordConfirm = u.passwordConfirm;
  adminUserConfig.value.firstName = u.firstName;
  adminUserConfig.value.lastName = u.lastName;
  adminUserConfig.value.street = u.street;
  adminUserConfig.value.zip = u.zip;
  adminUserConfig.value.city = u.city;
  adminUserConfig.value.phone = u.phone;
  adminUserConfig.value.ownRisk = u.ownRisk;
  adminUserConfig.value.license = u.license;
};

const addAdminUser = () => {
  // check for obvious validation errors
  const validationErrors = validateAdminUser();
  if (validationErrors.length > 0) {
    errors.value = validationErrors;
    return;
  }

  isLoading.value = true;

  User.signUp(adminUserConfig.value)
    .then((user) => {
      errors.value = [];
      makeUserAdmin(user.user_id);
    })
    .catch((err) => {
      errors.value = err;
      isLoading.value = false;
    });
};

const makeUserAdmin = (userId) => {
  User.makeAdmin(userId)
    .then(() => {
      errors.value = [];
      isLoading.value = false;
      getBackendStatusInternal();
    })
    .catch((err) => {
      errors.value = err;
      isLoading.value = false;
    });
};

const close = () => {
  router.push("/login");
};

const validateAdminUser = () => {
  const validationErrors = [];
  if (
    typeof adminUserConfig.value.password == "undefined" ||
    typeof adminUserConfig.value.email == "undefined" ||
    typeof adminUserConfig.value.passwordConfirm == "undefined"
  ) {
    validationErrors.push("Fields cannot be empty");
    return validationErrors;
  }

  if (adminUserConfig.value.password != adminUserConfig.value.passwordConfirm) {
    validationErrors.push(
      "Password and Password Confirmation are not identical.",
    );
  }
  if (adminUserConfig.value.password.length <= 8) {
    validationErrors.push("Please use a password longer than 8 characters.");
  }
  if (
    adminUserConfig.value.recaptchaResponse == null &&
    getRecaptchaKey.value
  ) {
    validationErrors.push("Please tick `I'm not a robot`.");
  }

  // password strength
  const pwUpperRegex = /[A-Z]+/;
  const pwLowerRegex = /[a-z]+/;
  const pwDigitRegex = /[0-9]+/;
  if (adminUserConfig.value.password.match(pwUpperRegex) == null) {
    validationErrors.push(
      "The password needs to contain at least one upper case letter (A-Z)",
    );
  }
  if (adminUserConfig.value.password.match(pwLowerRegex) == null) {
    validationErrors.push(
      "The password needs to contain at least one lower case letter (a-z)",
    );
  }
  if (adminUserConfig.value.password.match(pwDigitRegex) == null) {
    validationErrors.push(
      "The password needs to contain at least one digit (0-9)",
    );
  }

  return validationErrors;
};

const getBackendStatusInternal = () => {
  // reset errors
  errors.value = [];

  getBackendStatus()
    .then((status) => {
      console.log("status", status);
      if (status.configDb == false) {
        // no database configuration
        setupStep.value = 0;
        console.log("setupStep: 0");
      } else if (status.dbReachable == false) {
        // database not reachable but database configured
        // -> error
        setupStep.value = 3;
        console.log("setupStep: 4");
        errors.value = [
          "Cannot connect to the database."
        ]
      }  else if (status.usersExist == false) {
        // database is up, but there is no (admin) user yet
        setupStep.value = 1;
        console.log("setupStep: 1");
      }else {
        // setup is done
        setupStep.value = 2;
        console.log("setupStep: 2");
      }
      isLoading.value = false;
    })
    .catch((err) => {
      errors.value = err;
      isLoading.value = false;
    });
};

onMounted(() => {
  isLoading.value = true;
  getBackendStatusInternal();
});
</script>
