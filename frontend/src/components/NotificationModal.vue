<template>
  <div class="row" id="status_panel">
    <div class="col-sm-3"></div>
    <div class="col-sm-6">
      <div class="panel panel-default">
        <div class="panel-heading">Notification</div>
        <div class="panel-body" id="status_body">
          <!-- status messages go here -->
          {{ msg }}
        </div>
      </div>
    </div>
    <div class="col-sm-3"></div>
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
