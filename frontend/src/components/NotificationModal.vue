<template>
  <div id="status_panel" class="row">
    <div class="col-sm-3" />
    <div class="col-sm-6">
      <div class="panel panel-default">
        <div class="panel-heading">Notification</div>
        <div id="status_body" class="panel-body">
          <!-- status messages go here -->
          {{ msg }}
        </div>
      </div>
    </div>
    <div class="col-sm-3" />
  </div>
</template>

<script>
export default {
  name: "NotificationModal",
  props: ["message"],
  data: function () {
    return {
      msg: this.message,
    };
  },
};
</script>
