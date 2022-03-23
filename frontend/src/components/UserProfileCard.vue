<template>
  <sectioned-card-module>
    <template v-slot:header>
      <div class="row">
        <div class="col-6">
          <h5 class="card-title pt-1">Profile</h5>
        </div>
        <div class="col-6 text-end">
          <div class="dropdown">
            <button
              class="btn btn-outline-info btn-sm dropdown-toggle"
              type="button"
              id="dropdownMenuButton1"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <i class="bi bi-pencil-square"></i>
            </button>
            <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1">
              <li>
                <a class="dropdown-item" href="#" @click.prevent="showUserEdit"
                  >Edit Profile</a
                >
              </li>
              <li>
                <a
                  class="dropdown-item"
                  href="#"
                  @click.prevent="showPasswordEdit"
                  >Change Password</a
                >
              </li>
            </ul>
          </div>
        </div>
      </div>
    </template>
    <template v-slot:body>
      <user-edit-modal v-model:visible="showUserEditModal" />
      <user-password-edit-modal v-model:visible="showPasswordEditModal" />
      <div v-if="userInfo">
        <div class="row">
          <div class="col-6">
            {{ userInfo.first_name }} {{ userInfo.last_name }}
          </div>
          <div class="col-6">
            {{ userInfo.email }}
          </div>
        </div>
        <div class="row">
          <div class="col-6">
            {{ userInfo.address }}
          </div>
          <div class="col-6">
            {{ userInfo.mobile }}
          </div>
        </div>
        <div class="row">
          <div class="col-6">{{ userInfo.plz }} {{ userInfo.city }}</div>
          <div class="col-6">
            Driving License {{ userInfo.license == true ? "Yes" : "No" }}
          </div>
        </div>
      </div>
    </template>
  </sectioned-card-module>
</template>

<script>
import { mapGetters } from "vuex";
import UserEditModal from "./UserEditModal";
import UserPasswordEditModal from "./UserPasswordEditModal";
import SectionedCardModule from "./bricks/SectionedCardModule.vue";

export default {
  name: "UserProfileCard",
  components: {
    UserEditModal,
    UserPasswordEditModal,
    SectionedCardModule,
  },
  data() {
    return {
      showPasswordEditModal: false,
      showUserEditModal: false,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
  },
  methods: {
    showUserEdit: function () {
      this.showUserEditModal = true;
    },
    showPasswordEdit: function () {
      this.showPasswordEditModal = true;
    },
  },
  created() {
    console.log("created: ", this.userInfo);
  },
};
</script>
