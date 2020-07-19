<template>
  <div class="row" id="status_panel">
      <div class="col-sm-3"></div>
      <div class="col-sm-6">
          <div class="panel panel-default">
              <div class="panel-heading">
                  Notification
              </div>
              <div class="panel-body" id="status_body">
                  <!-- status messages go here -->
              </div>
          </div>
      </div>
      <div class="col-sm-3"></div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'

@Component
export default class NotificationModal extends Vue {
  @Prop() private msg!: string;
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
