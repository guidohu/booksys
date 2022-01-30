<template>
  <b-card
    no-body
    class="text-left"
  >
    <b-card-header> Statistics </b-card-header>
    <b-card-body>
      <b-row>
        <b-col cols="6">
          <b-row>
            <b-col cols="12">
              Riding Time: {{ heatTimeMinutesYTD }} min
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              Riding Time (all-time): {{ heatTimeMinutes }} min
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              Cost: {{ heatCostYTD }} {{ getCurrency }}
            </b-col>
          </b-row>
          <b-row>
            <b-col cols="12">
              Cost (all-time): {{ heatCost }} {{ getCurrency }}
            </b-col>
          </b-row>
        </b-col>
        <b-col cols="6">
          <b-button
            size="sm"
            type="button"
            variant="outline-info"
            @click="showLatestHeats"
          >
            <b-icon-list-ul />
            Latest Heats
          </b-button>
        </b-col>
      </b-row>
    </b-card-body>
    <UserHeatsModal v-model:visible="showUserHeatsModal" />
  </b-card>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import UserHeatsModal from "./UserHeatsModal";
import {
  BCard,
  BCardHeader,
  BCardBody,
  BRow,
  BCol,
  BButton,
  BIconListUl,
} from "bootstrap-vue";

export default {
  name: "UserStatisticsCard",
  components: {
    UserHeatsModal,
    BCard,
    BCardHeader,
    BCardBody,
    BRow,
    BCol,
    BButton,
    BIconListUl,
  },
  data() {
    return {
      showUserHeatsModal: false,
    };
  },
  computed: {
    ...mapGetters("login", ["userInfo"]),
    ...mapGetters("configuration", ["getCurrency"]),
    ...mapGetters("user", [
      "heatHistory",
      "heatTimeMinutes",
      "heatTimeMinutesYTD",
      "heatCost",
      "heatCostYTD",
    ]),
  },
  methods: {
    ...mapActions("configuration", ["queryConfiguration"]),
    showLatestHeats: function () {
      this.showUserHeatsModal = true;
    },
  },
  created() {
    this.queryConfiguration();
  },
};
</script>
