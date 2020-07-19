<template>
  <div class="login">
    <div class="display noborder nobackground">
      <div class="vcenter_placeholder">
        <div class="vcenter_content">
          <NotificationModal/>
          <form class="form-horizontal" role="form" action="#" name="login_form" id="login_form">
            <div class="row">
              <div class="col-sm-3">
              </div>
              <div class="col-sm-6">
                <div class="panel panel-default">
                  <div class="panel-heading">
                    Please Login
                  </div>
                  <div class="panel-body">
                    <div class="form-group">
                        <div class="col-sm-3"></div>
                        <div class="col-sm-8">
                            <label for="username">e-mail</label>
                            <input type="text" class="form-control form-control-white" name="username" id="username" placeholder="e-mail" autocomplete="username"/>
                        </div>
                        <div class="col-sm-1"></div>
                    </div>
                    <div class="form-group">
                        <div class="col-sm-3"></div>
                        <div class="col-sm-8">
                            <label for="password">password</label>
                            <input type="password" class="form-control form-control-white" name="password" id="password" placeholder="password" autocomplete="current-password"/>
                        </div>
                        <div class="col-sm-1"></div>
                    </div>
                    <div class="form-group text-center">
                        <div class="col-sm-3"></div>
                        <div class="col-sm-8" id="login_label"></div>
                        <div class="col-sm-1"></div>
                    </div>
                    <div class="form-group text-right">
                        <div class="col-sm-3"></div>
                        <div class="col-sm-8">
                            <button type="submit" id="loginBtn" class="btn btn-default inline">
                                <span class="glyphicon glyphicon-log-in"></span>
                                Login
                            </button>
                        </div>
                        <div class="col-sm-1"></div>
                    </div>
                  </div>
                  <div class="panel-footer">
                    <div class="form-group text-left">
                        <div class="col-sm-3"></div>
                        <div class="col-sm-8">
                            <a href="forgot_password.html" type="button" class="btn btn-default btn-xs inline">
                                Forgot password
                            </a>
                            <a href="sign_up.html" type="button" class="btn btn-default btn-xs inline">
                                Sign-Up
                            </a>
                        </div>
                        <div class="col-sm-1"></div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="col-sm-3"></div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import NotificationModal from '@/components/NotificationModal.vue'

export default {
  name: 'LoginModal',
  components: {
    NotificationModal
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
