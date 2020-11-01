<template>
  <div>
    <div v-if="isDesktop" class="display">
      <b-row>
        <b-col>
          <div class="main-title">
            <h3>
              Logs
            </h3>
          </div>
        </b-col>
      </b-row>
      <b-row class="ml-1 mr-1">
        <b-col cols="12">
          <b-card no-body>
            Content
          </b-card>
        </b-col>
      </b-row>
      <div class="bottom mr-2">
        <b-button variant="outline-light" to="/admin">
          <b-icon icon="gear"></b-icon>
          ADMIN
        </b-button>
      </div>
    </div>
    <div v-else>
      <NavbarAdminMobile title="Logs"/>
      <b-card no-body>
        Content
      </b-card>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { BooksysBrowser } from '@/libs/browser';
import NavbarAdminMobile from '@/components/NavbarAdminMobile';

export default Vue.extend({
  name: 'Logs',
  components: {
    NavbarAdminMobile
  },
  computed: {
    isDesktop: function () {
      return !BooksysBrowser.isMobile()
    }
  }
})
</script>