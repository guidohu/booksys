<div id="view_password_reset_email">
    <!-- <div class="modal fade text-left" id="view_user_signup_created_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content" id="view_user_signup_content">
                <div class="modal-header">
                    Sign Up Successful
                </div>
                <div class="modal-body" id="view_user_signup_dialog">
                    Your account has been created successfully. We will activate it soon and let you know, when your account is ready for your first login.
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" v-on:click="done">
                        <span class="glyphicon glyphicon-ok"></span>
                        Done
                    </button>
                </div>
            </div>
        </div>
    </div> -->
    <div class="modal fade text-left" id="view_password_reset_email_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content" id="view_user_signup_content">
                <div class="modal-header">
                    <h4>
                        Password Reset
                    </h4>
                </div>
                <form>
                    <div class="modal-body" id="view_user_signup_dialog">               
                        <div v-if="errors.length" class="row">
                            <div class="alert alert-danger alert-dismissible" role="alert">
                                <b>Please correct the following error(s):</b>
                                <ul>
                                    <li v-for="error in errors">{{ error }}</li>
                                </ul>
                            </div>
                        </div>
                        <div>
                            <div class="alert alert-info">
                                Please enter your Email address to reset your password:
                            </div>
                            <!-- email -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Email</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="email" type="text" class="form-control form-control-white input-sm" placeholder=""/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- captcha -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-10 col-xs-10">
                                    <div id="captcha"></div>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                        </div>
                        <div v-else>
                            <div class="alert alert-warning">
                                Something went wrong.
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-default" v-on:click="cancel">
                            <span class="glyphicon glyphicon-remove"></span>
                            Cancel
                        </button>
                        <button type="button" class="btn btn-default" v-on:click="save">
                            <span class="glyphicon glyphicon-arrow-right"></span>
                            Continue
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>

