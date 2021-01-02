<template>
  <b-modal
    id="incomeModal"
    title="Add New Income"
  >
    <b-row v-if="errors.length > 0">
      <b-col cols="1" class="d-none d-sm-block"></b-col>
      <b-col cols="12" sm="10">
        <WarningBox :errors="errors" dismissible="true"/>
      </b-col>
      <b-col cols="1" class="d-none d-sm-block"></b-col>
    </b-row>
    <b-form @submit="save">
      <b-row class="text-left">
        <b-col cols="1" class="d-none d-sm-block"></b-col>
        <b-col cols="12" sm="10">
          <b-form-group
            id="type-select-group"
            label="Type"
            label-for="type-select"
            description=""
            label-cols="3"
          >
            <b-form-select
              id="type-select"
              v-model="form.type"
              :options="incomeTypes"
            />
          </b-form-group>
          <b-form-group
            id="date-group"
            label="Date"
            label-for="date-select"
            description=""
            label-cols="3"
          >
            <b-form-input
              id="date-select"
              type="date"
              v-model="form.date"
            />
          </b-form-group>
          <b-form-group
            v-if="form.type != null"
            id="user-select-group"
            :label="userLabel"
            label-for="user-select"
            :description="userDescription"
            label-cols="3"
          >
            <b-form-select
              id="user-select"
              v-model="form.user"
              :options="users"
            />
          </b-form-group>
          <b-form-group
            v-if="form.type != null"
            id="amount"
            label="Amount"
            label-for="amount-input"
            description=""
            label-cols="3"
          >
            <b-input-group>
              <b-form-input
                id="amount-input"
                v-model="form.amount"
                type="text"
                placeholder="0.00"
              ></b-form-input>
              <b-input-group-append is-text>
                {{getCurrency}}
              </b-input-group-append>
            </b-input-group>
          </b-form-group>
          <b-form-group
            v-if="form.type != null && form.type != 4"
            id="description"
            label="Description"
            label-for="description-input"
            description=""
            label-cols="3"
          >
            <b-form-textarea
              id="description-input"
              v-model="form.description"
              type="text"
              rows="2"
              placeholder=""
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-form>
    <div slot="modal-footer">
      <div class="text-right d-inline">
        <b-button class="ml-4" type="button" variant="outline-info" v-on:click="add">
          <b-icon-check/>
          Add
        </b-button>
        <b-button class="ml-1" type="button" variant="outline-danger" v-on:click="close">
          <b-icon-x/>
          Cancel
        </b-button>
      </div>
    </div>
  </b-modal>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import WarningBox from '@/components/WarningBox';
import moment from 'moment';
import { orderBy } from 'lodash';
import {
  BModal,
  BRow,
  BCol,
  BForm,
  BFormGroup,
  BFormSelect,
  BFormInput,
  BInputGroup,
  BInputGroupAppend,
  BFormTextarea,
  BButton,
  BIconCheck,
  BIconX
} from 'bootstrap-vue';

export default {
  name: 'IncomeModal',
  components: {
    WarningBox,
    BModal,
    BRow,
    BCol,
    BForm,
    BFormGroup,
    BFormSelect,
    BFormInput,
    BInputGroup,
    BInputGroupAppend,
    BFormTextarea,
    BButton,
    BIconCheck,
    BIconX
  },
  data() {
    return {
      errors: [],
      incomeTypes: [],
      users: [],
      userLabel: "User",
      userDescription: "",
      form: {
        type: null,
        user: null,
        date: null
      }
    }
  },
  computed: {
    ...mapGetters('configuration', [
      'getCurrency'
    ]),
    ...mapGetters('user', [
      'userList'
    ]),
    ...mapGetters('accounting', [
      'getIncomeTypes'
    ])
  },
  watch: {
    getIncomeTypes: function(newValues){
      this.buildTypeSelect(newValues);
    },
    userList: function(newValues){
      // console.log("TODO, add userList:", newValues);
      this.buildUserSelect(newValues);
    },
    'form.type': function(newValue, oldValue){
      console.log(newValue, oldValue);
      if(newValue != oldValue){
          // if the type is either
          if([6, 4].includes(Number(newValue))){
            // 6: membership fee
            // 4: session
            this.userLabel = "Account";
            this.userDescription = "User that paid for sessions or membership.";
          }else if([3, 5].includes(Number(newValue))){
            // 3: invest
            // 5: other
            this.userLabel = "User";
            this.userDescription = "Payer or internal reference";
          }else if([7].includes(Number(newValue))){
            // 7: salary
            this.userLabel = "Driver";
            this.userDescription = "Driver that got paid for a session.";
          }else{
            // Default case
            this.userLabel = "User";
            this.userDescription = "";
          }
        }
    }
  },
  methods: {
    ...mapActions('configuration', [
      'queryConfiguration'
    ]),
    ...mapActions('user', [
      'queryUserList'
    ]),
    ...mapActions('accounting', [
      'queryIncomeTypes',
      'addIncome'
    ]),
    buildTypeSelect: function(types){
      this.incomeTypes = types.map((t) => {
        return {
          value: t.id,
          text:  t.name
        }
      });
      this.incomeTypes = orderBy(this.incomeTypes, ['text'], ['asc']);
      this.incomeTypes.unshift({ value: null, text: "Please select"});
    },
    buildUserSelect: function(users){
      this.users = users.map((u) => {
        return {
          value: u.id,
          text:  u.firstName + " " + u.lastName,
          lastName: u.lastName,
          firstName: u.firstName
        }
      });
      this.users = orderBy(this.users, ['text'], ['asc']);
      this.users.unshift({ value: null, text: "Please select"});
    },
    clearForm: function() {
      this.form = {
        amount: null,
        type: null,
        date: moment().format('YYYY-MM-DD'),
        user: null,
        comment: null
      };
    },
    add: function() {
      const income = {
        amount:   Number(this.form.amount),
        typeId:   Number(this.form.type),
        date:     this.form.date,
        userId:   Number(this.form.user),
        comment:  this.form.comment
      };

      this.addIncome(income)
      .then(() => {
        this.errors = [];
        this.close();
      })
      .catch((errors) => this.errors = errors);
    },
    save: function() {
      this.add();
    },
    close: function() {
      this.clearForm();
      this.$bvModal.hide('incomeModal');
    }
  },
  created() {
    this.queryConfiguration();

    this.queryUserList()
    .then(() => this.buildUserSelect(this.userList))
    .catch((errors) => this.errors.push(...errors))

    this.queryIncomeTypes()
    .then(() => this.buildTypeSelect(this.getIncomeTypes))
    .catch((errors) => this.errors.push(...errors));

    this.form.date = moment().format('YYYY-MM-DD');
  }
}
</script>